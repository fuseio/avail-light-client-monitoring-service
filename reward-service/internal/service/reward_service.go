package service

import (
	"fmt"
	"log"
	"time"

	"reward-service/internal/database"
)

// RewardLogEntry represents a log entry for a reward
type RewardLogEntry struct {
	RewardType      string
	Address         string
	Points          int64
	NFTCount        int64
	DelegationCount int64
	Timestamp       time.Time
}

// RewardService handles the calculation and distribution of rewards
type RewardService struct {
	db            *database.Database
	logger        *log.Logger
	ticker        *time.Ticker
	done          chan bool
	rewardsLog    []RewardLogEntry
	currentCycleID string
}

// NewRewardService creates a new reward service
func NewRewardService(db *database.Database, interval time.Duration, logger *log.Logger) *RewardService {
	return &RewardService{
		db:     db,
		logger: logger,
		ticker: time.NewTicker(interval),
		done:   make(chan bool),
		rewardsLog: make([]RewardLogEntry, 0),
	}
}

// logReward adds a reward entry to the rewards log
func (s *RewardService) logReward(rewardType, address string, points int64, nftCount int64, delegationCount int64) {
	timestamp := time.Now().Format(time.RFC3339)
	logEntry := fmt.Sprintf("[%s] REWARD_CREATED: type=%s, address=%s, points=%d, nfts=%d, delegations=%d", 
		timestamp, rewardType, address, points, nftCount, delegationCount)
	s.rewardsLog = append(s.rewardsLog, RewardLogEntry{
		RewardType:      rewardType,
		Address:         address,
		Points:          points,
		NFTCount:        nftCount,
		DelegationCount: delegationCount,
		Timestamp:       time.Now(),
	})
	s.logger.Println(logEntry)
}

// printRewardsLog prints the complete rewards log for the current cycle
func (s *RewardService) printRewardsLog() {
	if len(s.rewardsLog) == 0 {
		s.logger.Println("REWARDS_LOG: No rewards were created in this cycle")
		return
	}
	
	s.logger.Println("=== REWARDS LOG BEGIN ===")
	for _, entry := range s.rewardsLog {
		s.logger.Println(entry.RewardType, entry.Address, entry.Points, entry.NFTCount, entry.DelegationCount, entry.Timestamp)
	}
	s.logger.Printf("REWARDS_LOG: Total %d rewards created in this cycle", len(s.rewardsLog))
	s.logger.Println("=== REWARDS LOG END ===")
}

// Start starts the reward service
func (s *RewardService) Start() {
	s.logger.Println("Starting reward service...")
	
	// Calculate rewards immediately on start
	s.logger.Println("Running initial reward calculation...")
	go s.calculateAndDistributeRewards()
	
	// Check for latest rewards immediately
	s.logger.Println("Running initial rewards check...")
	go s.checkLatestRewards()
	
	// Then periodically
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.logger.Println("Running scheduled reward calculation...")
				go s.calculateAndDistributeRewards()
				go s.checkLatestRewards()
			case <-s.done:
				s.logger.Println("Received stop signal, shutting down reward service...")
				s.ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops the reward service
func (s *RewardService) Stop() {
	s.logger.Println("Stopping reward service...")
	s.done <- true
}

// calculateAndDistributeRewards calculates and distributes rewards to all clients
func (s *RewardService) calculateAndDistributeRewards() {
	startTime := time.Now()
	s.logger.Printf("Starting reward calculation and distribution for cycle %s", s.currentCycleID)
	
	// Get all clients from the monitoring service
	clients, err := s.db.GetAllClients()
	if err != nil {
		s.logger.Printf("ERROR: Failed to get clients: %v", err)
		return
	}
	
	s.logger.Printf("Found %d clients for reward processing", len(clients))
	
	// Process operator rewards
	var operatorRewardsProcessed int
	var totalOperatorPoints int64
	var operatorRewardsCreated []string
	
	for _, client := range clients {
		points := s.processOperatorRewards(client)
		if points > 0 {
			operatorRewardsProcessed++
			totalOperatorPoints += points
			operatorRewardsCreated = append(operatorRewardsCreated, client.Address)
		}
	}
	
	s.logger.Printf("Processed %d operator rewards with total %d points", 
		operatorRewardsProcessed, totalOperatorPoints)
	
	// Process delegator rewards
	var delegatorRewardsProcessed int
	var totalDelegatorPoints int64
	var delegatorRewardsCreated []string
	
	// Get all delegations and extract unique delegator addresses
	delegations, err := s.db.GetAllDelegations()
	if err != nil {
		s.logger.Printf("ERROR: Failed to get delegations: %v", err)
	} else {
		// Extract unique delegator addresses
		delegatorMap := make(map[string]bool)
		for _, delegation := range delegations {
			delegatorMap[delegation.FromAddress] = true
		}
		
		// Convert map keys to slice
		var delegators []string
		for delegator := range delegatorMap {
			delegators = append(delegators, delegator)
		}
		
		s.logger.Printf("Found %d unique delegators for reward processing", len(delegators))
		
		for _, delegator := range delegators {
			points := s.processDelegatorRewards(delegator)
			if points > 0 {
				delegatorRewardsProcessed++
				totalDelegatorPoints += points
				delegatorRewardsCreated = append(delegatorRewardsCreated, delegator)
			}
		}
		
		s.logger.Printf("Processed %d delegator rewards with total %d points", 
			delegatorRewardsProcessed, totalDelegatorPoints)
	}
	
	duration := time.Since(startTime)
	s.logger.Printf("Reward calculation and distribution completed in %v", duration)
	s.logger.Printf("Total rewards created: %d (Operators: %d, Delegators: %d)", 
		operatorRewardsProcessed+delegatorRewardsProcessed, 
		operatorRewardsProcessed, delegatorRewardsProcessed)
	
	// Update the last cycle time
	if err := s.db.UpdateLastCycleTime(time.Now()); err != nil {
		s.logger.Printf("ERROR: Failed to update last cycle time: %v", err)
	}
}

// processOperatorRewards processes rewards for an operator
func (s *RewardService) processOperatorRewards(client *database.MonitoringClientInfo) int64 {
	// Check for empty or invalid address (0x address)
	if client.Address == "" || client.Address == "0x" || client.Address == "0x0" || len(client.Address) <= 3 {
		s.logger.Printf("Skipping operator with invalid address: %s", client.Address)
		return 0
	}

	s.logger.Printf("Processing operator rewards for %s (status: %s, uptime: %.2f%%, NFTs: %d)", 
		client.Address, client.Status, client.AllUptimePercentage, client.NFTAmount)
	
	// Check if operator is active - if status is empty, consider it active for now
	if client.Status != "" && client.Status != "active" {
		s.logger.Printf("Skipping operator %s: not active (status: %s)", 
			client.Address, client.Status)
		return 0
	}
	
	// Check if we already processed rewards for this operator in this cycle
	exists, err := s.db.RewardExistsForCycle(client.Address, "OPERATOR", s.currentCycleID)
	if err != nil {
		s.logger.Printf("ERROR: Failed to check if reward exists: %v", err)
		return 0
	}
	
	if exists {
		s.logger.Printf("Skipping operator %s: already processed in this cycle", client.Address)
		return 0
	}
	
	// UNCOMMENT THIS: Check if the node has sufficient uptime (at least 50%)
	if client.AllUptimePercentage < 50.0 {
		s.logger.Printf("Skipping operator %s: insufficient uptime (%.2f%% < 50%%)", 
			client.Address, client.AllUptimePercentage)
		return 0
	}
	
	// Check if operator has at least 1 NFT
	if client.NFTAmount < 1 {
		s.logger.Printf("Skipping operator %s: insufficient NFTs (%d < 1)", 
			client.Address, client.NFTAmount)
		return 0
	}
	
	// Get delegations for this operator
	delegations, err := s.db.GetDelegationsForClient(client.Address)
	if err != nil {
		s.logger.Printf("ERROR: Failed to get delegations for %s: %v", client.Address, err)
		return 0
	}
	
	// Check if operator has at least 1 delegated NFT
	if len(delegations) == 0 || getTotalDelegatedNFTs(delegations) == 0 {
		s.logger.Printf("Skipping operator %s: no delegated NFTs", client.Address)
		return 0
	}
	
	s.logger.Printf("Operator %s has %d delegations and %d owned NFTs", 
		client.Address, len(delegations), client.NFTAmount)
	
	// Calculate operator points
	points := s.db.CalculateOperatorPoints(client, delegations)
	
	// Skip if no points earned
	if points <= 0 {
		s.logger.Printf("Operator %s earned 0 points", client.Address)
		return 0
	}
	
	s.logger.Printf("Operator %s earned %d points from commission", client.Address, points)
	
	// Determine the reward address - use reward_collector_address if available
	rewardAddress := client.Address
	if client.RewardCollectorAddress != "" {
		rewardAddress = client.RewardCollectorAddress
		s.logger.Printf("Using reward collector address %s instead of operator address %s", 
			rewardAddress, client.Address)
	}
	
	// Create reward record with cycle ID and claimed=false
	record := &database.RewardRecord{
		Address:         rewardAddress,
		Points:          points,
		Timestamp:       time.Now(),
		NFTCount:        client.NFTAmount,
		DelegationCount: int64(len(delegations)),
		CommissionRate:  client.CommissionRate,
		RewardType:      "OPERATOR",
		Claimed:         false,
		CycleID:         s.currentCycleID,
	}
	
	// Store reward record
	if err := s.db.StoreRewardRecord(record); err != nil {
		s.logger.Printf("ERROR: Failed to store reward record for %s: %v", client.Address, err)
		return 0
	}
	
	// Log the operator reward
	s.logReward("OPERATOR", client.Address, points, client.NFTAmount, int64(len(delegations)))
	
	return points
}

// Helper function to calculate total delegated NFTs
func getTotalDelegatedNFTs(delegations []*database.DelegationRecord) int64 {
	var total int64 = 0
	for _, delegation := range delegations {
		total += delegation.Amount
	}
	return total
}

// processDelegatorRewards processes rewards for a delegator
func (s *RewardService) processDelegatorRewards(address string) int64 {
	// Get delegations for this delegator
	delegations, err := s.db.GetDelegationsForDelegator(address)
	if err != nil {
		s.logger.Printf("ERROR: Failed to get delegations for %s: %v", address, err)
		return 0
	}
	
	// Filter delegations to only include those to active operators with sufficient uptime
	var eligibleDelegations []*database.DelegationRecord
	
	for _, delegation := range delegations {
		// Skip delegations to invalid addresses
		if delegation.ToAddress == "" || delegation.ToAddress == "0x" || delegation.ToAddress == "0x0" || len(delegation.ToAddress) <= 3 {
			s.logger.Printf("Skipping delegation to invalid address: %s", delegation.ToAddress)
			continue
		}
		
		// Get the operator's status
		operator, err := s.db.GetClient(delegation.ToAddress)
		if err != nil {
			s.logger.Printf("WARNING: Failed to get operator %s: %v", delegation.ToAddress, err)
			continue
		}
		
		// Check if operator is active
		if operator.Status != "" && operator.Status != "active" {
			s.logger.Printf("Skipping delegation to %s: operator is not active (status: %s)", 
				delegation.ToAddress, operator.Status)
			continue
		}
		
		// Check if operator has sufficient uptime
		if operator.AllUptimePercentage < 50.0 {
			s.logger.Printf("Skipping delegation to %s: operator has insufficient uptime (%.2f%% < 50%%)", 
				delegation.ToAddress, operator.AllUptimePercentage)
			continue
		}
		
		// Check if operator has at least 1 NFT
		if operator.NFTAmount < 1 {
			s.logger.Printf("Skipping delegation to %s: operator has insufficient NFTs (%d < 1)", 
				delegation.ToAddress, operator.NFTAmount)
			continue
		}
		
		// Update the delegation record with the operator's commission rate
		if delegation.CommissionRate != operator.CommissionRate {
			s.logger.Printf("Updating delegation commission rate from %.2f%% to %.2f%% (operator rate)", 
				delegation.CommissionRate, operator.CommissionRate)
			delegation.CommissionRate = operator.CommissionRate
		}
		
		eligibleDelegations = append(eligibleDelegations, delegation)
	}
	
	// Skip if no eligible delegations
	if len(eligibleDelegations) == 0 {
		s.logger.Printf("Delegator %s has no delegations to eligible operators", address)
		return 0
	}
	
	// Calculate delegator points - pass the address and eligible delegations
	points := s.db.CalculateDelegatorPoints(address, eligibleDelegations)
	
	// Skip if no points earned
	if points <= 0 {
		s.logger.Printf("Delegator %s earned 0 points", address)
		return 0
	}
	
	// Calculate commission deducted
	var commissionDeducted int64 = 1000 * eligibleDelegations[0].Amount - points
	s.logger.Printf("Commission deducted: %d points (%.2f%%)", 
		commissionDeducted, float64(commissionDeducted)*100.0/1000.0)
	
	// Create reward record
	record := &database.RewardRecord{
		Address:         address,
		Points:          points,
		Timestamp:       time.Now(),
		NFTCount:        eligibleDelegations[0].Amount,
		DelegationCount: int64(len(eligibleDelegations)),
		CommissionRate:  0, // Not applicable for delegator rewards
		RewardType:      "DELEGATOR",
		Claimed:         false,
		CycleID:         s.currentCycleID,
	}
	
	// Store reward record
	if err := s.db.StoreRewardRecord(record); err != nil {
		s.logger.Printf("ERROR: Failed to store reward record for delegator %s: %v", address, err)
		return 0
	}
	
	// Log the delegator reward
	s.logReward("DELEGATOR", address, points, eligibleDelegations[0].Amount, int64(len(eligibleDelegations)))
	
	return points
}

// checkLatestRewards checks and logs the latest rewards
func (s *RewardService) checkLatestRewards() {
	s.logger.Println("Checking latest rewards...")
	
	// Get the 10 most recent rewards
	rewards, err := s.db.GetLatestRewards(10)
	if err != nil {
		s.logger.Printf("ERROR: Failed to check latest rewards: %v", err)
		return
	}
	
	if len(rewards) == 0 {
		s.logger.Println("REWARDS_CHECK: No rewards found in the database")
		return
	}
	
	// Log the latest rewards
	s.logger.Println("=== LATEST REWARDS BEGIN ===")
	for i, reward := range rewards {
		s.logger.Printf("REWARD #%d: type=%s, address=%s, points=%d, timestamp=%v, NFTs=%d, delegations=%d", 
			i+1, reward.RewardType, reward.Address, reward.Points, reward.Timestamp, reward.NFTCount, reward.DelegationCount)
	}
	s.logger.Println("=== LATEST REWARDS END ===")
	
	// Calculate total points
	var totalPoints int64
	for _, reward := range rewards {
		totalPoints += reward.Points
	}
	
	s.logger.Printf("REWARDS_CHECK_SUMMARY: Found %d rewards, total points: %d, latest timestamp: %v", 
		len(rewards), totalPoints, rewards[0].Timestamp)
}

// RunRewardCycle runs a reward cycle
func (s *RewardService) RunRewardCycle() {
	// Generate a unique cycle ID based on the current time
	s.currentCycleID = time.Now().Format("2006-01-02-15-04-05")
	s.logger.Printf("Starting reward cycle with ID: %s", s.currentCycleID)
	
	// Clear previous reward logs
	s.rewardsLog = []RewardLogEntry{}
	
	// Check if it's time for a monthly reward cycle
	now := time.Now()
	
	// Get the last cycle time from database or use a default
	lastCycleTime, err := s.db.GetLastCycleTime()
	if err != nil {
		s.logger.Printf("WARNING: Failed to get last cycle time: %v, assuming first run", err)
		// If error, assume it's the first run
		lastCycleTime = time.Time{}
	}
	
	// Check if a month has passed since the last cycle
	if !lastCycleTime.IsZero() && now.Sub(lastCycleTime) < 30*24*time.Hour {
		s.logger.Printf("Skipping reward cycle: not yet time for monthly rewards. Last cycle: %s", 
			lastCycleTime.Format(time.RFC3339))
		return
	}
	
	// Run the reward calculation
	s.calculateAndDistributeRewards()
	
	// Update the last cycle time
	if err := s.db.UpdateLastCycleTime(now); err != nil {
		s.logger.Printf("ERROR: Failed to update last cycle time: %v", err)
	}
}