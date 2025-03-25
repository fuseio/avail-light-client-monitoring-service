package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RewardRecord struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Address         string             `bson:"address"`
	Points          int64              `bson:"points"`
	Timestamp       time.Time          `bson:"timestamp"`
	NFTCount        int64              `bson:"nft_count"`
	DelegationCount int64              `bson:"delegation_count"`
	CommissionRate  float64            `bson:"commission_rate"`
	RewardType      string             `bson:"reward_type"`
	Claimed         bool               `bson:"claimed"`
	ClaimedAt       time.Time          `bson:"claimed_at,omitempty"`
	CycleID         string             `bson:"cycle_id"` 
}


func (d *Database) StoreRewardRecord(record *RewardRecord) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := d.rewards.InsertOne(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to insert reward record: %v", err)
	}

	return nil
}


func (d *Database) ClaimReward(rewardID primitive.ObjectID, address string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	
	var reward RewardRecord
	err := d.rewards.FindOne(ctx, bson.M{
		"_id":      rewardID,
		"address":  address,
		"claimed":  false,
	}).Decode(&reward)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("reward not found or already claimed")
		}
		return fmt.Errorf("failed to find reward: %v", err)
	}

	
	session, err := d.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %v", err)
	}
	defer session.EndSession(ctx)

	
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		
		now := time.Now()
		_, err := d.rewards.UpdateOne(
			sc,
			bson.M{"_id": rewardID},
			bson.M{
				"$set": bson.M{
					"claimed":    true,
					"claimed_at": now,
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed to update reward: %v", err)
		}

		
		userType := UserTypeOperator
		if reward.RewardType == "DELEGATOR" {
			userType = UserTypeDelegate
		}

		
		_, err = d.GetUserByType(address, userType)
		if err != nil {
			return fmt.Errorf("user not found: %v", err)
		}

		
		err = d.UpdateUserClaimedPoints(address, userType, reward.Points)
		if err != nil {
			return fmt.Errorf("failed to update user claimed points: %v", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}

	return nil
}


func (d *Database) RewardExistsForCycle(address string, rewardType string, cycleID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	count, err := d.rewards.CountDocuments(ctx, bson.M{
		"address":     address,
		"reward_type": rewardType,
		"cycle_id":    cycleID,
	})
	
	if err != nil {
		return false, fmt.Errorf("failed to check if reward exists: %v", err)
	}
	
	return count > 0, nil
}


func (d *Database) GetRewardByID(id primitive.ObjectID) (RewardRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var reward RewardRecord
	err := d.rewards.FindOne(ctx, bson.M{"_id": id}).Decode(&reward)
	if err != nil {
		return RewardRecord{}, err
	}

	return reward, nil
}


func (d *Database) GetRewardsByAddress(address string) ([]RewardRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := d.rewards.Find(ctx, bson.M{"address": address})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rewards []RewardRecord
	if err = cursor.All(ctx, &rewards); err != nil {
		return nil, err
	}

	return rewards, nil
}


func (d *Database) ClaimAllRewards(address string) (int, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	
	cursor, err := d.rewards.Find(ctx, bson.M{
		"address": address,
		"claimed": false,
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
		return 0, 0, nil 
	}

	
	session, err := d.client.StartSession()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to start session: %v", err)
	}
	defer session.EndSession(ctx)

	var totalPoints int64
	var claimedCount int

	
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		now := time.Now()
		
		
		userType := UserTypeOperator
		if len(rewards) > 0 && rewards[0].RewardType == "DELEGATOR" {
			userType = UserTypeDelegate
		}

		
		_, err = d.GetUserByType(address, userType)
		if err != nil {
			return fmt.Errorf("user not found: %v", err)
		}

		
		for _, reward := range rewards {
			_, err := d.rewards.UpdateOne(
				sc,
				bson.M{"_id": reward.ID},
				bson.M{
					"$set": bson.M{
						"claimed":    true,
						"claimed_at": now,
					},
				},
			)
			if err != nil {
				return fmt.Errorf("failed to update reward %s: %v", reward.ID.Hex(), err)
			}
			
			totalPoints += reward.Points
			claimedCount++
		}

		
		err = d.UpdateUserClaimedPoints(address, userType, totalPoints)
		if err != nil {
			return fmt.Errorf("failed to update user claimed points: %v", err)
		}

		return nil
	})

	if err != nil {
		return 0, 0, fmt.Errorf("transaction failed: %v", err)
	}

	return claimedCount, totalPoints, nil
}


func (d *Database) GetUserByType(address string, userType UserType) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	err := d.users.FindOne(ctx, bson.M{
		"address": address,
		"user_type": userType,
	}).Decode(&user)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found with address %s and type %s", address, userType)
		}
		return nil, fmt.Errorf("failed to find user: %v", err)
	}
	
	return &user, nil
} 