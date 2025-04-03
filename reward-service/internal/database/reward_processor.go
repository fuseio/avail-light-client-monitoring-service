package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RewardProcessor struct {
	db            *mongo.Database
	client        *mongo.Client
	logger        *log.Logger
	cycleID       string
	users         *mongo.Collection
	rewards       *mongo.Collection
	rewardSummary *mongo.Collection
}

func NewRewardProcessor(client *mongo.Client, db *mongo.Database, logger *log.Logger, cycleID string) *RewardProcessor {
	return &RewardProcessor{
		db:            db,
		client:        client,
		logger:        logger,
		cycleID:       cycleID,
		users:         db.Collection("users"),
		rewards:       db.Collection("rewards"),
		rewardSummary: db.Collection("reward_summary"),
	}
}

func (rp *RewardProcessor) Process() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	rp.logger.Printf("Starting reward processing for cycle: %s", rp.cycleID)

	var existingSummary RewardSummary
	err := rp.rewardSummary.FindOne(ctx, bson.M{"cycle_id": rp.cycleID}).Decode(&existingSummary)
	if err == nil {
		return fmt.Errorf("rewards for cycle %s have already been processed", rp.cycleID)
	} else if err != mongo.ErrNoDocuments {
		return fmt.Errorf("error checking existing reward summary: %v", err)
	}

	rewardSummary := RewardSummary{
		CycleID:         rp.cycleID,
		TotalRewards:    0,
		TotalPoints:     0,
		OperatorRewards:  0,
		DelegatorRewards: 0,
		Timestamp:       time.Now(),
	}

	var rewardRecords []interface{}
	
	eligibleOperators := make(map[string]bool)

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

	rp.logger.Printf("Processing %d operators...", len(operators))
	
	for _, operator := range operators {
		if operator.AllUptimePercentage < 50 {
			rp.logger.Printf("Operator %s not eligible: uptime %.2f%% below 50%%", 
				operator.Address, operator.AllUptimePercentage)
			continue
		}
		
		eligibleOperators[operator.Address] = true
		basePoints := int64(1000)
		
		commissionPoints := int64(0)
		
		if operator.Delegators != nil {
			for _, amount := range operator.Delegators {
				commissionRate := 0.0
				if operator.CommissionRate > 0 {
					commissionRate = operator.CommissionRate / 100.0
				} else {
					rp.logger.Printf("WARNING: Operator %s has no commission rate set in database, using 0%%", 
						operator.Address)
				}
				
				commissionPoints += amount * 1000 * int64(commissionRate*100) / 100
			}
		}
		
		operatorReward := RewardRecord{
			UserAddress: operator.Address,
			Points:      basePoints + commissionPoints,
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
	
	if len(rewardRecords) > 0 {
		_, err = rp.rewards.InsertMany(ctx, rewardRecords)
		if err != nil {
			return fmt.Errorf("failed to insert rewards: %v", err)
		}
	}
	
	rewardSummary.TotalRewards = rewardSummary.OperatorRewards + rewardSummary.DelegatorRewards
	_, err = rp.rewardSummary.InsertOne(ctx, rewardSummary)
	if err != nil {
		return fmt.Errorf("failed to insert reward summary: %v", err)
	}
	
	rp.logger.Printf("Successfully processed rewards for cycle %s: %d operators, %d points", 
		rp.cycleID, rewardSummary.OperatorRewards, rewardSummary.TotalPoints)
	
	return nil
}

func (db *Database) ClaimAllRewards(address string) (int, int64, error) {
	if address == "" {
		return 0, 0, fmt.Errorf("empty address provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	supportsTransactions := false
	
	if db.client.NumberSessionsInProgress() >= 0 {
		serverStatus, err := db.client.Database("admin").RunCommand(
			ctx, bson.D{{Key: "serverStatus", Value: 1}},
		).DecodeBytes()
		
		if err == nil {
			process, err := serverStatus.LookupErr("process")
			if err == nil && process.StringValue() == "mongod" {
				_, err = serverStatus.LookupErr("repl")
				supportsTransactions = err == nil
			}
		}
		
		if !supportsTransactions {
			db.logger.Println("MongoDB deployment does not support transactions, using non-transactional path")
		}
	}
	
	if supportsTransactions {
		return db.claimAllRewardsWithTransaction(ctx, address)
	} else {
		return db.claimAllRewardsWithoutTransaction(ctx, address)
	}
}

func (db *Database) claimAllRewardsWithTransaction(ctx context.Context, address string) (int, int64, error) {
	session, err := db.client.StartSession()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to start session: %v", err)
	}
	defer session.EndSession(ctx)
	
	claimedCount := 0
	var totalPoints int64 = 0
	
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := sc.StartTransaction(); err != nil {
			return fmt.Errorf("failed to start transaction: %v", err)
		}

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
			return nil
		}
		
		var rewardIDs []primitive.ObjectID
		for _, reward := range rewards {
			rewardIDs = append(rewardIDs, reward.ID)
			totalPoints += reward.Points
		}
		
		result, err := db.db.Collection("rewards").UpdateMany(sc, bson.M{
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
		
		if result.ModifiedCount != int64(len(rewardIDs)) {
			db.logger.Printf("WARNING: Expected to update %d rewards but modified %d", 
				len(rewardIDs), result.ModifiedCount)
		}
		
		var user User
		err = db.users.FindOne(sc, bson.M{"address": address}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				user = User{
					Address:            address,
					TotalClaimedPoints: 0,
					CreatedAt:          time.Now(),
					UpdatedAt:          time.Now(),
				}
				
				_, err = db.users.InsertOne(sc, user)
				if err != nil {
					return fmt.Errorf("failed to create user: %v", err)
				}
			} else {
				return fmt.Errorf("failed to check if user exists: %v", err)
			}
		}
		
		updateResult, err := db.users.UpdateOne(sc, bson.M{
			"address": address,
		}, bson.M{
			"$inc": bson.M{
				"total_claimed_points": totalPoints,
			},
			"$set": bson.M{
				"last_claim_time": time.Now(),
				"updated_at":      time.Now(),
			},
		})
		
		if err != nil {
			return fmt.Errorf("failed to update user points: %v", err)
		}
		
		if updateResult.ModifiedCount != 1 {
			return fmt.Errorf("failed to update user points, modified count: %d", updateResult.ModifiedCount)
		}
		
		claimedCount = len(rewards)
		
		if err := sc.CommitTransaction(sc); err != nil {
			return fmt.Errorf("failed to commit transaction: %v", err) 
		}
		
		return nil
	})
	
	if err != nil {
		if abortErr := session.AbortTransaction(ctx); abortErr != nil {
			db.logger.Printf("ERROR: Failed to abort transaction: %v", abortErr)
		}
		return 0, 0, err
	}
	
	db.logger.Printf("Successfully claimed %d rewards with %d points for user %s", 
		claimedCount, totalPoints, address)
	
	return claimedCount, totalPoints, nil
}

func (db *Database) claimAllRewardsWithoutTransaction(ctx context.Context, address string) (int, int64, error) {
	var totalPoints int64 = 0
	claimedCount := 0
	
	cursor, err := db.db.Collection("rewards").Find(ctx, bson.M{
		"user_address": address,
		"is_claimed": false,
	})
	
	if err != nil {
		return 0, 0, fmt.Errorf("failed to find rewards: %v", err)
	}
	defer cursor.Close(ctx)
	
	var rewards []RewardRecord
	if err := cursor.All(ctx, &rewards); err != nil {
		return 0, 0, fmt.Errorf("failed to decode rewards: %v", err)
	}
	
	if len(rewards) == 0 {
		db.logger.Printf("No unclaimed rewards found for user %s", address)
		return 0, 0, nil
	}
	
	var rewardIDs []primitive.ObjectID
	for _, reward := range rewards {
		rewardIDs = append(rewardIDs, reward.ID)
		totalPoints += reward.Points
	}
	
	result, err := db.db.Collection("rewards").UpdateMany(ctx, bson.M{
		"_id": bson.M{"$in": rewardIDs},
	}, bson.M{
		"$set": bson.M{
			"is_claimed": true,
			"claimed_at": time.Now(),
		},
	})
	
	if err != nil {
		return 0, 0, fmt.Errorf("failed to update rewards: %v", err)
	}
	
	if result.ModifiedCount != int64(len(rewardIDs)) {
		db.logger.Printf("WARNING: Expected to update %d rewards but modified %d", 
			len(rewardIDs), result.ModifiedCount)
	}
	
	var user User
	err = db.users.FindOne(ctx, bson.M{"address": address}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			user = User{
				Address:            address,
				TotalClaimedPoints: 0,
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			}
			
			_, err = db.users.InsertOne(ctx, user)
			if err != nil {
				return 0, 0, fmt.Errorf("failed to create user: %v", err)
			}
		} else {
			return 0, 0, fmt.Errorf("failed to check if user exists: %v", err)
		}
	}
	
	updateResult, err := db.users.UpdateOne(ctx, bson.M{
		"address": address,
	}, bson.M{
		"$inc": bson.M{
			"total_claimed_points": totalPoints,
		},
		"$set": bson.M{
			"last_claim_time": time.Now(),
			"updated_at":      time.Now(),
		},
	})
	
	if err != nil {
		return 0, 0, fmt.Errorf("failed to update user points: %v", err)
	}
	
	if updateResult.ModifiedCount != 1 {
		db.logger.Printf("WARNING: Failed to update user points or user was already up to date, modified count: %d", 
			updateResult.ModifiedCount)
	}
	
	claimedCount = len(rewards)
	
	db.logger.Printf("Successfully claimed %d rewards with %d points for user %s", 
		claimedCount, totalPoints, address)
	
	return claimedCount, totalPoints, nil
}
