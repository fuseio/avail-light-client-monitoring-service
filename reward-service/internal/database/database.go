package database

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Database struct {
	client          *mongo.Client
	rewards         *mongo.Collection
	clients         *mongo.Collection // Reference to monitoring service clients collection
	delegations     *mongo.Collection // Reference to monitoring service delegations collection
	heartbeats      *mongo.Collection
	users           *mongo.Collection
	logger          *log.Logger
	db              *mongo.Database
	cycles          *mongo.Collection
}

type MonitoringClientInfo struct {
	Address                string    `bson:"address"`
	TotalTime              int64     `bson:"total_time"`
	LastHeartbeat          time.Time `bson:"last_heartbeat"`
	CreatedAt              time.Time `bson:"created_at"`
	NFTAmount              int64     `bson:"nft_amount"`
	CommissionRate         float64   `bson:"commission_rate"`
	Status                 string    `bson:"status"`
	AllUptimePercentage    float64   `bson:"all_uptime_percentage"`
	WeeklyUptimePercentage float64   `bson:"weekly_uptime_percentage"`
	OperatorName           string    `bson:"operator_name"`
	RewardCollectorAddress string    `bson:"reward_collector_address"`
}

type DelegationRecord struct {
	FromAddress    string    `bson:"from_address"`
	ToAddress      string    `bson:"to_address"`
	Amount         int64     `bson:"amount"`
	CommissionRate float64   `bson:"commission_rate"`
	CreatedAt      time.Time `bson:"created_at"`
}

func NewDatabase(uri, dbName string, logger *log.Logger) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	logger.Printf("Connected to MongoDB database: %s", dbName)
	
	clients := db.Collection("clients")
	heartbeats := db.Collection("heartbeats")
	delegations := db.Collection("delegations")
	rewards := db.Collection("rewards")
	users := db.Collection("users")
	cycles := db.Collection("cycles")
	
	logger.Printf("Accessing collections: clients, delegations, rewards, users, cycles")

	return &Database{
		client:      client,
		db:          db,
		clients:     clients,
		heartbeats:  heartbeats,
		delegations: delegations,
		rewards:     rewards,
		users:       users,
		cycles:      cycles,
		logger:      logger,
	}, nil
}

func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.client.Disconnect(ctx)
}

func (d *Database) GetAllClients() ([]*MonitoringClientInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	d.logger.Printf("Fetching all clients from monitoring service")

	cursor, err := d.clients.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch clients: %v", err)
	}
	defer cursor.Close(ctx)

	var clients []*MonitoringClientInfo
	for cursor.Next(ctx) {
		var client MonitoringClientInfo
		if err := cursor.Decode(&client); err != nil {
			d.logger.Printf("WARNING: Failed to decode client: %v", err)
			continue
		}
		
		now := time.Now()
		
		if client.AllUptimePercentage == 0 && client.CreatedAt.Unix() > 0 {
			allTimeTotalSeconds := float64(now.Sub(client.CreatedAt).Seconds())
			if client.TotalTime > 0 && allTimeTotalSeconds > 0 {
				client.AllUptimePercentage = math.Min(float64(client.TotalTime) / allTimeTotalSeconds * 100, 100)
			}
		}
		
		
		clients = append(clients, &client)
		d.logger.Printf("Client %d: Address=%s, NFTAmount=%d, Commission=%.2f%%, Uptime=%.2f%%, Status=%s", 
			len(clients), client.Address, client.NFTAmount, client.CommissionRate, client.AllUptimePercentage, client.Status)
	}

	d.logger.Printf("SUCCESS: Retrieved %d clients, %d with NFTs (for testing)", 
		len(clients), countClientsWithNFTs(clients))
	return clients, nil
}

func countClientsWithNFTs(clients []*MonitoringClientInfo) int {
	count := 0
	for _, client := range clients {
		if client.NFTAmount > 0 {
			count++
		}
	}
	return count
}

func (d *Database) GetClient(address string) (*MonitoringClientInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d.logger.Printf("Fetching client with address: %s", address)

	var clientInfo MonitoringClientInfo
	err := d.clients.FindOne(ctx, bson.M{"address": address}).Decode(&clientInfo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("client not found: %s", address)
		}
		return nil, fmt.Errorf("failed to fetch client: %v", err)
	}

	now := time.Now()
	
	if clientInfo.AllUptimePercentage == 0 && clientInfo.CreatedAt.Unix() > 0 {
		allTimeTotalSeconds := float64(now.Sub(clientInfo.CreatedAt).Seconds())
		if clientInfo.TotalTime > 0 && allTimeTotalSeconds > 0 {
			clientInfo.AllUptimePercentage = math.Min(float64(clientInfo.TotalTime) / allTimeTotalSeconds * 100, 100)
		}
		d.logger.Printf("Calculated all-time uptime for %s: %.2f%% (TotalTime: %d, CreatedAt: %s)", 
			clientInfo.Address, clientInfo.AllUptimePercentage, clientInfo.TotalTime, clientInfo.CreatedAt)
	}
	
	if clientInfo.WeeklyUptimePercentage == 0 {
		oneWeekAgo := now.Add(-7 * 24 * time.Hour)
		startTime := clientInfo.CreatedAt
		if startTime.Before(oneWeekAgo) {
			startTime = oneWeekAgo
		}
		
		filter := bson.M{
			"client_address": clientInfo.Address,
			"timestamp": bson.M{"$gte": startTime},
		}
		
		cursor, err := d.heartbeats.Find(ctx, filter)
		if err != nil {
			d.logger.Printf("WARNING: Failed to fetch heartbeats for %s: %v", clientInfo.Address, err)
		} else {
			defer cursor.Close(ctx)
			
			var weeklyTotalTime int64 = 0
			for cursor.Next(ctx) {
				var heartbeat struct {
					Duration int64 `bson:"duration"`
				}
				if err := cursor.Decode(&heartbeat); err != nil {
					continue
				}
				weeklyTotalTime += heartbeat.Duration
			}
			
			weeklyTotalSeconds := float64(now.Sub(startTime).Seconds())
			if weeklyTotalTime > 0 && weeklyTotalSeconds > 0 {
				clientInfo.WeeklyUptimePercentage = math.Min(float64(weeklyTotalTime) / weeklyTotalSeconds * 100, 100)
			}
			
			d.logger.Printf("Calculated weekly uptime for %s: %.2f%%", 
				clientInfo.Address, clientInfo.WeeklyUptimePercentage)
		}
	}

	

	d.logger.Printf("SUCCESS: Retrieved client %s: NFTAmount=%d, Commission=%.2f%%, Uptime=%.2f%%, Status=%s", 
		clientInfo.Address, clientInfo.NFTAmount, clientInfo.CommissionRate, clientInfo.AllUptimePercentage, clientInfo.Status)
	return &clientInfo, nil
}

func (d *Database) GetDelegationsForClient(address string) ([]*DelegationRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d.logger.Printf("Fetching delegations to address: %s", address)

	if d.delegations == nil {
		return nil, fmt.Errorf("delegations collection is not initialized")
	}

	cursor, err := d.delegations.Find(ctx, bson.M{"to_address": address})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch delegations: %v", err)
	}
	defer cursor.Close(ctx)

	var delegations []*DelegationRecord
	for cursor.Next(ctx) {
		var delegation DelegationRecord
		if err := cursor.Decode(&delegation); err != nil {
			d.logger.Printf("WARNING: Failed to decode delegation: %v", err)
			continue
		}
		delegations = append(delegations, &delegation)
	}

	d.logger.Printf("Found %d delegations to %s", len(delegations), address)


	// Get client info to check if it's eligible for test delegations
	_, err = d.GetClient(address)
	if err != nil {
		d.logger.Printf("WARNING: Failed to get client info for %s: %v", address, err)
		return delegations, nil
	}

	return delegations, nil
}

func (d *Database) GetDelegationsFromClient(address string) ([]DelegationRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d.logger.Printf("Fetching delegations to address: %s", address)

	filter := bson.M{"to_address": bson.M{"$regex": primitive.Regex{Pattern: "^" + address + "$", Options: "i"}}}
	
	cursor, err := d.delegations.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find delegations: %v", err)
	}
	defer cursor.Close(ctx)

	var delegations []DelegationRecord
	if err = cursor.All(ctx, &delegations); err != nil {
		return nil, fmt.Errorf("failed to decode delegations: %v", err)
	}

	d.logger.Printf("Found %d delegations to %s", len(delegations), address)

	for i, delegation := range delegations {
		if i < 5 { // Log only first 5 delegations to avoid log spam
			d.logger.Printf("Delegation to %s: from=%s, amount=%d, commission=%.2f%%", 
				address, delegation.FromAddress, delegation.Amount, delegation.CommissionRate)
		}
	}

	if len(delegations) > 0 {
		d.logger.Printf("SUCCESS: Retrieved %d delegations to client %s", len(delegations), address)
	}
	
	return delegations, nil
}

func (d *Database) CalculateOperatorPoints(client *MonitoringClientInfo, delegations []*DelegationRecord) int64 {
	if client.NFTAmount < 1 {
		d.logger.Printf("Operator %s has no delegated NFTs, skipping", client.Address)
		return 0
	}
	
	// Check if uptime is at least 50%
	if client.AllUptimePercentage < 50.0 {
		d.logger.Printf("Operator %s has uptime below 50%% (%.2f%%), skipping", 
			client.Address, client.AllUptimePercentage)
		return 0
	}
	
	// Calculate total delegated NFTs from the delegations
	var delegatedNFTs int64 = 0
	for _, delegation := range delegations {
		delegatedNFTs += delegation.Amount
	}
	
	// Calculate operator points using the new formula
	basePoints := int64(1000)
	commissionPoints := int64(1000 * float64(delegatedNFTs) * (client.CommissionRate / 100.0))
	totalPoints := basePoints + commissionPoints
	
	d.logger.Printf("Operator %s points calculation: %d base + %d commission = %d total", 
		client.Address, basePoints, commissionPoints, totalPoints)
	
	return totalPoints
}

func (d *Database) CalculateDelegatorPoints(address string, delegations []*DelegationRecord) int64 {
	var totalPoints int64 = 0
	
	d.logger.Printf("Calculating delegator points for %d delegations", len(delegations))
	
	for i, delegation := range delegations {
		d.logger.Printf("Processing delegation %d/%d: from=%s, to=%s, amount=%d, delegation_commission=%.2f%%", 
			i+1, len(delegations), delegation.FromAddress, delegation.ToAddress, delegation.Amount, delegation.CommissionRate)
		
		operator, err := d.GetClient(delegation.ToAddress)
		
		// Base points for this delegation (1000 points per delegation, not per NFT)
		basePoints := int64(1000)
		d.logger.Printf("Base points before commission: %d (1000 per delegation)", basePoints)
		
		var commissionRate float64
		if err != nil {
			d.logger.Printf("Warning: Failed to get operator info for %s: %v", delegation.ToAddress, err)
			commissionRate = delegation.CommissionRate / 100.0
			d.logger.Printf("Using delegation record's commission rate: %.2f%%", commissionRate * 100)
		} else {
			d.logger.Printf("Using operator's commission rate from client record: %.2f%%", operator.CommissionRate)
			commissionRate = operator.CommissionRate / 100.0
		}
		
		// Calculate commission deduction
		commissionDeduction := int64(float64(basePoints) * commissionRate)
		d.logger.Printf("Commission deduction: %d (%.2f%% of %d)", 
			commissionDeduction, commissionRate * 100, basePoints)
		
		// Calculate points after commission
		pointsAfterCommission := basePoints - commissionDeduction
		d.logger.Printf("Points after commission: %d (%d - %d)", 
			pointsAfterCommission, basePoints, commissionDeduction)
		
		// Add to total points
		totalPoints += pointsAfterCommission
	}
	
	d.logger.Printf("SUCCESS: Total delegator points: %d", totalPoints)
	
	return totalPoints
}

func (d *Database) GetUser(address string) (*MonitoringClientInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var clientInfo MonitoringClientInfo
	err := d.clients.FindOne(ctx, bson.M{"address": address}).Decode(&clientInfo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("client not found: %s", address)
		}
		return nil, fmt.Errorf("failed to fetch client: %v", err)
	}

	return &clientInfo, nil
}

func (d *Database) GetAllUsers() ([]MonitoringClientInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := d.clients.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch clients: %v", err)
	}
	defer cursor.Close(ctx)

	var clients []MonitoringClientInfo
	if err := cursor.All(ctx, &clients); err != nil {
		return nil, fmt.Errorf("failed to decode clients: %v", err)
	}

	return clients, nil
}

func (d *Database) GetRewardsSummary() (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d.logger.Printf("Generating rewards summary")

	rewardsCount, err := d.rewards.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to count rewards: %v", err)
	}

	usersCount, err := d.clients.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %v", err)
	}

	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":         nil,
				"totalPoints": bson.M{"$sum": "$points"},
			},
		},
	}

	cursor, err := d.rewards.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate rewards: %v", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		TotalPoints int64 `bson:"totalPoints"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode aggregation results: %v", err)
	}

	var totalPoints int64 = 0
	if len(results) > 0 {
		totalPoints = results[0].TotalPoints
	}

	opts := options.FindOne().SetSort(bson.M{"timestamp": -1})
	var latestReward RewardRecord
	err = d.rewards.FindOne(ctx, bson.M{}, opts).Decode(&latestReward)
	
	var latestRewardTime time.Time
	if err == nil {
		latestRewardTime = latestReward.Timestamp
	} else if err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("failed to get latest reward: %v", err)
	}
	
	operatorCount, _ := d.rewards.CountDocuments(ctx, bson.M{"reward_type": "OPERATOR"})
	delegatorCount, _ := d.rewards.CountDocuments(ctx, bson.M{"reward_type": "DELEGATOR"})
	testCount, _ := d.rewards.CountDocuments(ctx, bson.M{"reward_type": "TEST"})

	summary := map[string]interface{}{
		"total_rewards":        rewardsCount,
		"total_users":          usersCount,
		"total_points":         totalPoints,
		"latest_reward_time":   latestRewardTime,
		"average_points_per_user": float64(0),
		"operator_rewards":     operatorCount,
		"delegator_rewards":    delegatorCount,
		"test_rewards":         testCount,
	}

	if usersCount > 0 {
		summary["average_points_per_user"] = float64(totalPoints) / float64(usersCount)
	}

	d.logger.Printf("Rewards summary: %d rewards (Operator=%d, Delegator=%d, Test=%d), %d users, %d total points", 
		rewardsCount, operatorCount, delegatorCount, testCount, usersCount, totalPoints)
	return summary, nil
}

func (d *Database) GetLatestRewards(limit int) ([]RewardRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d.logger.Printf("Fetching latest %d rewards", limit)

	options := options.Find().SetSort(bson.M{"timestamp": -1}).SetLimit(int64(limit))
	cursor, err := d.rewards.Find(ctx, bson.M{}, options)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest rewards: %v", err)
	}
	defer cursor.Close(ctx)

	var rewards []RewardRecord
	if err = cursor.All(ctx, &rewards); err != nil {
		return nil, fmt.Errorf("failed to decode rewards: %v", err)
	}

	d.logger.Printf("Found %d recent rewards", len(rewards))
	
	for i, reward := range rewards {
		if i < 5 { // Log only first 5 rewards to avoid log spam
			d.logger.Printf("Recent reward %d: Type=%s, Address=%s, Points=%d, Timestamp=%v, NFTs=%d, Delegations=%d", 
				i+1, reward.RewardType, reward.Address, reward.Points, reward.Timestamp, reward.NFTCount, reward.DelegationCount)
		}
	}

	if len(rewards) > 0 {
		d.logger.Printf("REWARDS_MONITOR: Retrieved %d latest rewards, most recent from %v", 
			len(rewards), rewards[0].Timestamp)
	}
	
	return rewards, nil
}

func (d *Database) GetUserRewards(address string) ([]RewardRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d.logger.Printf("Fetching reward records for %s", address)

	cursor, err := d.rewards.Find(ctx, bson.M{"address": address})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rewards: %v", err)
	}
	defer cursor.Close(ctx)

	var rewards []RewardRecord
	if err := cursor.All(ctx, &rewards); err != nil {
		return nil, fmt.Errorf("failed to decode rewards: %v", err)
	}

	d.logger.Printf("Found %d reward records for %s", len(rewards), address)
	
	var operatorCount, delegatorCount, testCount int
	var totalPoints int64
	
	for _, reward := range rewards {
		totalPoints += reward.Points
		switch reward.RewardType {
		case "OPERATOR":
			operatorCount++
		case "DELEGATOR":
			delegatorCount++
		case "TEST":
			testCount++
		}
	}
	
	if len(rewards) > 0 {
		d.logger.Printf("REWARDS_SUMMARY for %s: Total=%d (Operator=%d, Delegator=%d, Test=%d), Points=%d", 
			address, len(rewards), operatorCount, delegatorCount, testCount, totalPoints)
	}
	
	return rewards, nil
}

func (d *Database) GetLastCycleTime() (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var result struct {
		LastCycleTime time.Time `bson:"last_cycle_time"`
	}
	
	err := d.client.Database("reward_service").Collection("system_settings").FindOne(ctx, bson.M{"_id": "reward_cycle"}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return time.Time{}, nil // Return zero time if no record exists
		}
		return time.Time{}, fmt.Errorf("failed to get last cycle time: %v", err)
	}
	
	return result.LastCycleTime, nil
}

func (d *Database) UpdateLastCycleTime(t time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err := d.client.Database("reward_service").Collection("system_settings").UpdateOne(
		ctx,
		bson.M{"_id": "reward_cycle"},
		bson.M{"$set": bson.M{"last_cycle_time": t}},
		options.Update().SetUpsert(true),
	)
	
	if err != nil {
		return fmt.Errorf("failed to update last cycle time: %v", err)
	}
	
	return nil
}

func (d *Database) Initialize(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	d.client = client
	d.db = client.Database("reward_service")
	
	d.clients = d.db.Collection("clients")
	d.delegations = d.db.Collection("delegations")
	d.rewards = d.db.Collection("rewards")
	d.users = d.db.Collection("users")
	d.cycles = d.db.Collection("cycles")
	
	// Ensure users collection exists by creating it if it doesn't
	collections, err := d.db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %v", err)
	}
	
	hasUsers := false
	for _, name := range collections {
		if name == "users" {
			hasUsers = true
			break
		}
	}
	
	if !hasUsers {
		err = d.db.CreateCollection(ctx, "users")
		if err != nil {
			return fmt.Errorf("failed to create users collection: %v", err)
		}
		d.logger.Printf("Created users collection")
	}
	
	d.logger.Printf("Connected to MongoDB successfully")
	return nil
}

// GetAllDelegations gets all delegations from the database
func (d *Database) GetAllDelegations() ([]DelegationRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := d.delegations.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get delegations: %v", err)
	}
	defer cursor.Close(ctx)

	var delegations []DelegationRecord
	if err := cursor.All(ctx, &delegations); err != nil {
		return nil, fmt.Errorf("failed to decode delegations: %v", err)
	}

	return delegations, nil
}

func (d *Database) GetDelegationsForDelegator(address string) ([]*DelegationRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d.logger.Printf("Fetching delegations from address: %s", address)

	if d.delegations == nil {
		return nil, fmt.Errorf("delegations collection is not initialized")
	}

	cursor, err := d.delegations.Find(ctx, bson.M{"from_address": address})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch delegations: %v", err)
	}
	defer cursor.Close(ctx)

	var delegations []*DelegationRecord
	for cursor.Next(ctx) {
		var delegation DelegationRecord
		if err := cursor.Decode(&delegation); err != nil {
			d.logger.Printf("WARNING: Failed to decode delegation: %v", err)
			continue
		}
		delegations = append(delegations, &delegation)
	}

	d.logger.Printf("Found %d delegations from %s", len(delegations), address)

	return delegations, nil
}

