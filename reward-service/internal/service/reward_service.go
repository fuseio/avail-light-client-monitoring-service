package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"reward-service/internal/database"
)

type RewardService struct {
	db             *database.Database
	logger         *log.Logger
	currentCycleID string
	mu             sync.RWMutex // Protects currentCycleID
	done           chan struct{}
	wg             sync.WaitGroup
}

func New(db *database.Database, logger *log.Logger) *RewardService {
	now := time.Now()
	cycleID := fmt.Sprintf("CYCLE_%s", now.Format("2006-01-02"))
	
	return &RewardService{
		db:            db,
		logger:        logger,
		currentCycleID: cycleID,
		done:          make(chan struct{}),
	}
}

func NewRewardService(db *database.Database, logger *log.Logger) *RewardService {
	return New(db, logger)
}

func (s *RewardService) SetCycleID(cycleID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentCycleID = cycleID
}

func (s *RewardService) GetCycleID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentCycleID
}

func (s *RewardService) logReward(rewardType string, address string, points int64, nftCount int64, delegationCount int64) {
	s.logger.Printf("Reward: %s, Address: %s, Points: %d, NFTs: %d, Delegations: %d",
		rewardType, address, points, nftCount, delegationCount)
}

func (s *RewardService) ScheduleRewardsAt(hour, minute int) {
	s.logger.Printf("Scheduling rewards to run daily at %02d:%02d", hour, minute)
	
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
			
			if next.Before(now) {
				next = next.Add(24 * time.Hour)
			}
			
			delay := next.Sub(now)
			s.logger.Printf("Next reward run scheduled in %v at %s", delay, next.Format("2006-01-02 15:04:05"))
			
			select {
			case <-time.After(delay):
				now = time.Now() // Update current time
				cycleID := fmt.Sprintf("CYCLE_%s", now.Format("2006-01-02"))
				s.SetCycleID(cycleID)
				
				maxRetries := 3
				backoff := 1 * time.Minute
				var err error
				
				for attempt := 0; attempt <= maxRetries; attempt++ {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
					err = s.ProcessRewardsWithContext(ctx)
					cancel()
					
					if err == nil {
						s.logger.Printf("Successfully processed rewards at scheduled time %s", 
							now.Format("2006-01-02 15:04:05"))
						break
					}
					
					if attempt == maxRetries {
						s.logger.Printf("ERROR: Failed to process rewards after %d attempts: %v", 
							maxRetries+1, err)
						break
					}
					
					retryDelay := backoff * time.Duration(1<<uint(attempt))
					s.logger.Printf("ERROR: Failed to process rewards (attempt %d/%d): %v. Retrying in %v...", 
						attempt+1, maxRetries+1, err, retryDelay)
					
					select {
					case <-time.After(retryDelay):
						continue
					case <-s.done:
						return
					}
				}
			case <-s.done:
				return
			}
		}
	}()
}

func (s *RewardService) Stop() {
	s.logger.Println("Stopping reward service...")
	close(s.done)
	
	waitCh := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(waitCh)
	}()
	
	select {
	case <-waitCh:
		s.logger.Println("Reward service stopped gracefully")
	case <-time.After(10 * time.Second):
		s.logger.Println("Reward service stop timed out, some tasks may not have completed")
	}
}

func (s *RewardService) ProcessRewardsWithContext(ctx context.Context) error {
	s.logger.Println("Processing rewards...")
	
	cycleID := s.GetCycleID()
	
	rewardSummary := &database.RewardSummary{
		CycleID:          cycleID,
		TotalRewards:     0,
		TotalPoints:      0,
		OperatorRewards:  0,
		DelegatorRewards: 0,
		Timestamp:        time.Now(),
	}
	
	operatorCommissions := make(map[string]int64)
	
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	
	s.logger.Println("Processing delegator rewards...")
	delegators, err := s.db.GetAllDelegators()
	if err != nil {
		return fmt.Errorf("failed to get delegators: %v", err)
	}
	
	for _, delegator := range delegators {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		points := s.processDelegatorRewards(delegator.Address, operatorCommissions)
		if points > 0 {
			rewardSummary.DelegatorRewards++
			rewardSummary.TotalPoints += points
		}
	}
	
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	
	s.logger.Println("Processing operator rewards...")
	operators, err := s.db.GetAllOperators()
	if err != nil {
		return fmt.Errorf("failed to get operators: %v", err)
	}
	
	for _, operator := range operators {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		commission := operatorCommissions[operator.Address]
		points := s.processOperatorRewards(operator.Address, commission)
		if points > 0 {
			rewardSummary.OperatorRewards++
			rewardSummary.TotalPoints += points
		}
	}
	
	rewardSummary.TotalRewards = rewardSummary.OperatorRewards + rewardSummary.DelegatorRewards
	
	if err := s.db.StoreRewardSummary(rewardSummary); err != nil {
		return fmt.Errorf("failed to store reward summary: %v", err)
	}
	
	s.logger.Printf("Reward processing completed. Total rewards: %d, Total points: %d",
		rewardSummary.TotalRewards, rewardSummary.TotalPoints)
	return nil
}

func (s *RewardService) ProcessRewards() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	return s.ProcessRewardsWithContext(ctx)
}

func (s *RewardService) processDelegatorRewards(delegatorAddress string, operatorCommissions map[string]int64) int64 {
	delegations, err := s.db.GetDelegationsFromDelegator(delegatorAddress)
	if err != nil {
		s.logger.Printf("ERROR: Failed to get delegations for delegator %s: %v", delegatorAddress, err)
		return 0
	}
	
	if len(delegations) == 0 {
		s.logger.Printf("No dedicated delegation records for %s, checking user mapping", delegatorAddress)
		user, err := s.db.GetUserByAddress(delegatorAddress, "DELEGATOR")
		if err != nil || user.Operators == nil || len(user.Operators) == 0 {
			s.logger.Printf("Skipping delegator %s: no delegations found", delegatorAddress)
			return 0
		}
		
		for operatorAddr, amount := range user.Operators {
			delegations = append(delegations, database.DelegationRecord{
				FromAddress: delegatorAddress,
				ToAddress:   operatorAddr,
				Amount:      amount,
				Timestamp:   time.Now(),
			})
		}
	}
	
	s.logger.Printf("Delegator %s has %d delegations", delegatorAddress, len(delegations))
	
	totalPoints := int64(0)
	
	for _, delegation := range delegations {
		client, err := s.db.GetClient(delegation.ToAddress)
		if err != nil {
			s.logger.Printf("ERROR: Failed to get operator %s info: %v", delegation.ToAddress, err)
			continue
		}
		
		if client.AllUptimePercentage < 50.0 {
			s.logger.Printf("Skipping delegation to operator %s: insufficient operator uptime (%.2f%% < 50%%)",
				delegation.ToAddress, client.AllUptimePercentage)
			continue
		}
		
		commissionRate := 0.0
		if client.CommissionRate > 0 {
			commissionRate = client.CommissionRate / 100.0
		} else {
			s.logger.Printf("WARNING: Operator %s has no commission rate set in database, using 0%%", 
				delegation.ToAddress)
		}
		
		basePointsPerNFT := int64(1000)
		
		commissionPerNFT := int64(float64(basePointsPerNFT) * commissionRate)
		
		delegatorPointsPerNFT := basePointsPerNFT - commissionPerNFT
		points := delegation.Amount * delegatorPointsPerNFT
		
		commission := delegation.Amount * commissionPerNFT
		
		record := &database.RewardRecord{
			UserAddress:     delegatorAddress,
			OperatorAddress: delegation.ToAddress,
			Points:          points,
			CreatedAt:       time.Now(),
			IsClaimed:       false,
			CycleID:         s.GetCycleID(),
			Type:            "DELEGATOR",
		}
		
		if err := s.db.StoreRewardRecord(record); err != nil {
			s.logger.Printf("ERROR: Failed to store reward record for %s: %v", delegatorAddress, err)
			continue
		}
		
		operatorCommissions[delegation.ToAddress] += commission
		
		s.logger.Printf("Delegator %s gets %d points (%d per NFT after %.2f%% commission) for %d NFTs delegated to %s", 
			delegatorAddress, points, delegatorPointsPerNFT, client.CommissionRate, 
			delegation.Amount, delegation.ToAddress)
		
		s.logger.Printf("Added commission of %d points (%d per NFT) for operator %s", 
			commission, commissionPerNFT, delegation.ToAddress)
		
		s.logReward("DELEGATOR", delegatorAddress, points, delegation.Amount, int64(len(delegations)))
		
		totalPoints += points
	}
	
	return totalPoints
}

func (s *RewardService) processOperatorRewards(operatorAddress string, commission int64) int64 {
	client, err := s.db.GetClient(operatorAddress)
	if err != nil {
		s.logger.Printf("ERROR: Failed to get client info for operator %s: %v", operatorAddress, err)
		return 0
	}
	
	s.logger.Printf("Operator %s has uptime %.2f%%", operatorAddress, client.AllUptimePercentage)
	
	if client.AllUptimePercentage < 50.0 {
		s.logger.Printf("Skipping operator %s: insufficient uptime (%.2f%% < 50%%)",
			operatorAddress, client.AllUptimePercentage)
		return 0
	}
	
	delegations, err := s.db.GetDelegationsToOperator(operatorAddress)
	if err != nil {
		s.logger.Printf("ERROR: Failed to get delegations for operator %s: %v", operatorAddress, err)
		return 0
	}
	
	if len(delegations) == 0 {
		s.logger.Printf("Skipping operator %s: no delegations", operatorAddress)
		return 0
	}
	
	s.logger.Printf("Operator %s has %d delegations to them", 
		operatorAddress, len(delegations))
	
	basePoints := int64(1000)
	totalPoints := basePoints + commission
	
	rewardAddress := client.Address
	if client.RewardCollectorAddress != "" {
		rewardAddress = client.RewardCollectorAddress
		s.logger.Printf("Using reward collector address %s instead of operator address %s", 
			rewardAddress, client.Address)
	}
	
	record := &database.RewardRecord{
		UserAddress:     client.Address,
		Points:          totalPoints,
		CreatedAt:       time.Now(),
		IsClaimed:       false,
		CycleID:         s.GetCycleID(),
		Type:            "OPERATOR",        
	}
	
	if err := s.db.StoreRewardRecord(record); err != nil {
		s.logger.Printf("ERROR: Failed to store reward record for %s: %v", client.Address, err)
		return 0
	}
	
	s.logReward("OPERATOR", client.Address, totalPoints, client.NFTAmount, int64(len(delegations)))
	
	return totalPoints
}