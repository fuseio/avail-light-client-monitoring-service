package service

import (
	"fmt"
	"log"
	"time"

	"reward-service/internal/database"
)

type RewardService struct {
	db             *database.Database
	logger         *log.Logger
	currentCycleID string
	done           chan struct{}
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
	s.currentCycleID = cycleID
}

func (s *RewardService) logReward(rewardType string, address string, points int64, nftCount int64, delegationCount int64) {
	s.logger.Printf("Reward: %s, Address: %s, Points: %d, NFTs: %d, Delegations: %d",
		rewardType, address, points, nftCount, delegationCount)
}

func (s *RewardService) ScheduleRewardsAt(hour, minute int) {
	s.logger.Printf("Scheduling rewards to run daily at %02d:%02d", hour, minute)
	
	go func() {
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
				
				if err := s.ProcessRewards(); err != nil {
					s.logger.Printf("ERROR: Failed to process rewards: %v", err)
				} else {
					s.logger.Printf("Successfully processed rewards at scheduled time %s", 
						now.Format("2006-01-02 15:04:05"))
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
}

func (s *RewardService) ProcessRewards() error {
	s.logger.Println("Processing rewards...")
	
	rewardSummary := &database.RewardSummary{
		CycleID:          s.currentCycleID,
		TotalRewards:     0,
		TotalPoints:      0,
		OperatorRewards:  0,
		DelegatorRewards: 0,
		Timestamp:        time.Now(),
	}
	
	operatorCommissions := make(map[string]int64)
	
	s.logger.Println("Processing delegator rewards...")
	delegators, err := s.db.GetAllDelegators()
	if err != nil {
		return fmt.Errorf("failed to get delegators: %v", err)
	}
	
	for _, delegator := range delegators {
		points := s.processDelegatorRewards(delegator.Address, operatorCommissions)
		if points > 0 {
			rewardSummary.DelegatorRewards++
			rewardSummary.TotalPoints += points
		}
	}
	
	s.logger.Println("Processing operator rewards...")
	operators, err := s.db.GetAllOperators()
	if err != nil {
		return fmt.Errorf("failed to get operators: %v", err)
	}
	
	for _, operator := range operators {
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
			CycleID:         s.currentCycleID,
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
		CycleID:         s.currentCycleID,
		Type:            "OPERATOR",        
	}
	
	if err := s.db.StoreRewardRecord(record); err != nil {
		s.logger.Printf("ERROR: Failed to store reward record for %s: %v", client.Address, err)
		return 0
	}
	
	s.logReward("OPERATOR", client.Address, totalPoints, client.NFTAmount, int64(len(delegations)))
	
	return totalPoints
}