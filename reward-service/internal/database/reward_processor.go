package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RewardRecord represents a reward given to a user
type RewardRecord struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserAddress string             `bson:"user_address" json:"user_address"`
	Points      int64              `bson:"points" json:"points"`
	Reason      string             `bson:"reason" json:"reason"`
	CycleID     string             `bson:"cycle_id" json:"cycle_id"`
	IsClaimed   bool               `bson:"is_claimed" json:"is_claimed"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	ClaimedAt   time.Time          `bson:"claimed_at,omitempty" json:"claimed_at,omitempty"`
}

// RewardSummary contains summary information about rewards
type RewardSummary struct {
	CycleID         string    `json:"cycle_id"`
	TotalRewards    int       `json:"total_rewards"`
	TotalPoints     int64     `json:"total_points"`
	OperatorRewards int       `json:"operator_rewards"`
	DelegatorRewards int      `json:"delegator_rewards"`
	Timestamp       time.Time `json:"timestamp"`
}

// RewardProcessor handles the reward calculation and distribution
type RewardProcessor struct {
	db           *mongo.Database
	logger       *log.Logger
	cycleID      string
	users        *mongo.Collection
	rewards      *mongo.Collection
	rewardSummary *mongo.Collection
}

func NewRewardProcessor(db *mongo.Database, logger *log.Logger) *RewardProcessor {
	// Create a unique cycleID based on the current date in Israel timezone
	israelLoc, err := time.LoadLocation("Asia/Jerusalem")
	if err != nil {
		logger.Printf("Error loading Israel timezone: %v, falling back to UTC", err)
		israelLoc = time.UTC
	}
	now := time.Now().In(israelLoc)
	cycleID := fmt.Sprintf("CYCLE_%s", now.Format("2006-01-02"))
	
	return &RewardProcessor{
		db:           db,
		logger:       logger,
		cycleID:      cycleID,
		users:        db.Collection("users"),
		rewards:      db.Collection("rewards"),
		rewardSummary: db.Collection("reward_summary"),
	}
}

// SetCycleID allows overriding the default cycle ID (useful for testing)
func (rp *RewardProcessor) SetCycleID(cycleID string) {
	rp.cycleID = cycleID
}

// ProcessRewards calculates and creates rewards for all eligible users
func (rp *RewardProcessor) ProcessRewards() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	
	rp.logger.Printf("Starting reward processing for cycle: %s", rp.cycleID)
	
	// Check if rewards for this cycle have already been processed
	count, err := rp.rewards.CountDocuments(ctx, bson.M{"cycle_id": rp.cycleID})
	if err != nil {
		return fmt.Errorf("failed to check existing rewards: %v", err)
	}
	
	if count > 0 {
		rp.logger.Printf("Rewards for cycle %s have already been processed (%d rewards found)", rp.cycleID, count)
		return fmt.Errorf("rewards for cycle %s have already been processed", rp.cycleID)
	}
	
	// 1. Find all operators
	rp.logger.Println("Finding eligible operators...")
	operatorCursor, err := rp.users.Find(ctx, bson.M{"user_type": "OPERATOR"})
	if err != nil {
		return fmt.Errorf("failed to find operators: %v", err)
	}
	defer operatorCursor.Close(ctx)
	
	var operators []User
	if err := operatorCursor.All(ctx, &operators); err != nil {
		return fmt.Errorf("failed to decode operators: %v", err)
	}
	
	// Track reward stats
	rewardSummary := RewardSummary{
		CycleID:         rp.cycleID,
		Timestamp:       time.Now(),
	}
	
	// 2. Process each operator
	eligibleOperators := make(map[string]bool)
	var rewardRecords []interface{}
	
	for _, operator := range operators {
		// Check if operator was active for at least 50% of the time
		if operator.AllUptimePercentage < 50 {
			rp.logger.Printf("Operator %s not eligible: uptime %.2f%% below 50%%", 
				operator.Address, operator.AllUptimePercentage)
			continue
		}
		
		// Operator is eligible for base points
		eligibleOperators[operator.Address] = true
		basePoints := int64(1000)
		
		// Calculate commission points
		commissionPoints := int64(0)
		
		if operator.Delegators != nil {
			for _, amount := range operator.Delegators {
				// For each delegator NFT, add commission to operator (10% of 1000 points per NFT)
				// Default to 10% commission if commission_rate is not available
				commissionRate := 0.10 // Default 10%
				
				// Ideally we'd get the actual commission rate from the client collection
				// but for simplicity we'll use the default
				commissionPoints += amount * 1000 * int64(commissionRate*100) / 100
			}
		}
		
		// Create operator reward record
		operatorReward := RewardRecord{
			UserAddress: operator.Address,
			Points:      basePoints + commissionPoints,
			Reason:      fmt.Sprintf("Operator Base: 1000, Commission: %d", commissionPoints),
			CycleID:     rp.cycleID,
			IsClaimed:   false,
			CreatedAt:   time.Now(),
		}
		
		rewardRecords = append(rewardRecords, operatorReward)
		rewardSummary.OperatorRewards++
		rewardSummary.TotalPoints += operatorReward.Points
		
		rp.logger.Printf("Calculated reward for operator %s: %d points (%d base + %d commission)", 
			operator.Address, operatorReward.Points, basePoints, commissionPoints)
	}
	
	// 3. Find and process all delegators
	rp.logger.Println("Finding eligible delegators...")
	delegatorCursor, err := rp.users.Find(ctx, bson.M{"user_type": "DELEGATOR"})
	if err != nil {
		return fmt.Errorf("failed to find delegators: %v", err)
	}
	defer delegatorCursor.Close(ctx)
	
	var delegators []User
	if err := delegatorCursor.All(ctx, &delegators); err != nil {
		return fmt.Errorf("failed to decode delegators: %v", err)
	}
	
	for _, delegator := range delegators {
		if delegator.Operators == nil || len(delegator.Operators) == 0 {
			continue
		}
		
		totalDelegatorPoints := int64(0)
		var reasonParts string
		
		// Process each operator the delegator has delegated to
		for operatorAddr, amount := range delegator.Operators {
			// Only reward if the operator is eligible
			if !eligibleOperators[operatorAddr] {
				continue
			}
			
			// Default commission rate (10%)
			commissionRate := 0.10
			
			// Calculate delegator points (1000 per NFT minus commission)
			pointsPerNFT := 1000 - int64(1000*commissionRate)
			delegatorPoints := amount * pointsPerNFT
			
			totalDelegatorPoints += delegatorPoints
			
			if reasonParts != "" {
				reasonParts += ", "
			}
			reasonParts += fmt.Sprintf("%d NFTs to %s: %d points", 
				amount, operatorAddr, delegatorPoints)
		}
		
		// Skip if no eligible delegations
		if totalDelegatorPoints == 0 {
			continue
		}
		
		// Create delegator reward record
		delegatorReward := RewardRecord{
			UserAddress: delegator.Address,
			Points:      totalDelegatorPoints,
			Reason:      fmt.Sprintf("Delegation rewards: %s", reasonParts),
			CycleID:     rp.cycleID,
			IsClaimed:   false,
			CreatedAt:   time.Now(),
		}
		
		rewardRecords = append(rewardRecords, delegatorReward)
		rewardSummary.DelegatorRewards++
		rewardSummary.TotalPoints += delegatorReward.Points
		
		rp.logger.Printf("Calculated reward for delegator %s: %d points", 
			delegator.Address, totalDelegatorPoints)
	}
	
	// 4. Insert all rewards in a bulk operation
	if len(rewardRecords) > 0 {
		rp.logger.Printf("Inserting %d reward records...", len(rewardRecords))
		_, err = rp.rewards.InsertMany(ctx, rewardRecords)
		if err != nil {
			return fmt.Errorf("failed to insert rewards: %v", err)
		}
	}
	
	// Update summary stats
	rewardSummary.TotalRewards = rewardSummary.OperatorRewards + rewardSummary.DelegatorRewards
	
	// 5. Save reward summary
	_, err = rp.rewardSummary.InsertOne(ctx, rewardSummary)
	if err != nil {
		rp.logger.Printf("Warning: Failed to save reward summary: %v", err)
	}
	
	rp.logger.Printf("Reward processing complete: %d rewards (%d operators, %d delegators) for %d total points",
		rewardSummary.TotalRewards, rewardSummary.OperatorRewards, 
		rewardSummary.DelegatorRewards, rewardSummary.TotalPoints)
	
	return nil
}

// Helper functions for the reward database operations

// GetRewardsByAddress retrieves all rewards for a specific address
func (db *Database) GetRewardsByAddress(address string) ([]RewardRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Sort by createdAt descending to get most recent first
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	
	cursor, err := db.db.Collection("rewards").Find(ctx, bson.M{
		"user_address": address,
	}, opts)
	
	if err != nil {
		return nil, fmt.Errorf("failed to find rewards: %v", err)
	}
	defer cursor.Close(ctx)
	
	var rewards []RewardRecord
	if err := cursor.All(ctx, &rewards); err != nil {
		return nil, fmt.Errorf("failed to decode rewards: %v", err)
	}
	
	return rewards, nil
}

// GetRewardByID retrieves a reward by its ID
func (db *Database) GetRewardByID(id primitive.ObjectID) (RewardRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var reward RewardRecord
	err := db.db.Collection("rewards").FindOne(ctx, bson.M{"_id": id}).Decode(&reward)
	if err != nil {
		return RewardRecord{}, fmt.Errorf("failed to find reward: %v", err)
	}
	
	return reward, nil
}

// ClaimReward marks a reward as claimed and updates the user's total claimed points
func (db *Database) ClaimReward(rewardID primitive.ObjectID, address string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Start a session for transaction
	session, err := db.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %v", err)
	}
	defer session.EndSession(ctx)
	
	// Run transaction
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		// 1. Find the reward and ensure it exists, belongs to the user, and is not claimed
		var reward RewardRecord
		err := db.db.Collection("rewards").FindOne(sc, bson.M{
			"_id": rewardID,
			"user_address": address,
			"is_claimed": false,
		}).Decode(&reward)
		
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fmt.Errorf("reward not found or already claimed")
			}
			return fmt.Errorf("failed to find reward: %v", err)
		}
		
		// 2. Update the reward to mark it as claimed
		_, err = db.db.Collection("rewards").UpdateOne(sc, bson.M{
			"_id": rewardID,
		}, bson.M{
			"$set": bson.M{
				"is_claimed": true,
				"claimed_at": time.Now(),
			},
		})
		
		if err != nil {
			return fmt.Errorf("failed to update reward: %v", err)
		}
		
		// 3. Update the user's total claimed points
		_, err = db.users.UpdateOne(sc, bson.M{
			"address": address,
		}, bson.M{
			"$inc": bson.M{
				"total_claimed_points": reward.Points,
			},
			"$set": bson.M{
				"last_claim_time": time.Now(),
			},
		})
		
		if err != nil {
			return fmt.Errorf("failed to update user points: %v", err)
		}
		
		return nil
	})
	
	return err
}

// ClaimAllRewards claims all unclaimed rewards for a user
func (db *Database) ClaimAllRewards(address string) (int, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	// Start a session for transaction
	session, err := db.client.StartSession()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to start session: %v", err)
	}
	defer session.EndSession(ctx)
	
	claimedCount := 0
	var totalPoints int64 = 0
	
	// Run transaction
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		// 1. Find all unclaimed rewards for this user
		cursor, err := db.db.Collection("rewards").Find(sc, bson.M{
			"user_address": address,
			"is_claimed": false,
		})
		
		if err != nil {
			return fmt.Errorf("failed to find rewards: %v", err)
		}
		defer cursor.Close(ctx)
		
		var rewards []RewardRecord
		if err := cursor.All(sc, &rewards); err != nil {
			return fmt.Errorf("failed to decode rewards: %v", err)
		}
		
		if len(rewards) == 0 {
			return nil // No rewards to claim
		}
		
		// 2. Update all rewards to claimed status
		var rewardIDs []primitive.ObjectID
		for _, reward := range rewards {
			rewardIDs = append(rewardIDs, reward.ID)
			totalPoints += reward.Points
		}
		
		_, err = db.db.Collection("rewards").UpdateMany(sc, bson.M{
			"_id": bson.M{"$in": rewardIDs},
		}, bson.M{
			"$set": bson.M{
				"is_claimed": true,
				"claimed_at": time.Now(),
			},
		})
		
		if err != nil {
			return fmt.Errorf("failed to update rewards: %v", err)
		}
		
		// 3. Update the user's total claimed points
		_, err = db.users.UpdateOne(sc, bson.M{
			"address": address,
		}, bson.M{
			"$inc": bson.M{
				"total_claimed_points": totalPoints,
			},
			"$set": bson.M{
				"last_claim_time": time.Now(),
			},
		})
		
		if err != nil {
			return fmt.Errorf("failed to update user points: %v", err)
		}
		
		claimedCount = len(rewards)
		return nil
	})
	
	if err != nil {
		return 0, 0, err
	}
	
	return claimedCount, totalPoints, nil
}

// GetUserClaimedPoints returns the total claimed points for a user
func (db *Database) GetUserClaimedPoints(address string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var user User
	err := db.users.FindOne(ctx, bson.M{"address": address}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil // User doesn't exist, return 0 points
		}
		return 0, fmt.Errorf("failed to find user: %v", err)
	}
	
	return user.TotalClaimedPoints, nil
} 