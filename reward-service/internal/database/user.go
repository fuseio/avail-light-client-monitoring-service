package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


type UserType string

const (
	UserTypeOperator  UserType = "OPERATOR"
	UserTypeDelegate  UserType = "DELEGATOR"
)


type User struct {
	Address           string    `bson:"address"`
	UserType          UserType  `bson:"user_type"`
	TotalClaimedPoints int64     `bson:"total_claimed_points"`
	LastClaimTime     time.Time `bson:"last_claim_time"`
	CreatedAt         time.Time `bson:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at"`
	Delegators        map[string]int64   `bson:"delegators,omitempty"` 
}


func (d *Database) UpdateUserClaimedPoints(address string, userType UserType, points int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	_, err := d.users.UpdateOne(
		ctx,
		bson.M{
			"address": address,
			"user_type": userType,
		},
		bson.M{
			"$inc": bson.M{"total_claimed_points": points},
			"$set": bson.M{
				"last_claim_time": now,
				"updated_at": now,
			},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to update user claimed points: %v", err)
	}

	return nil
}


func (d *Database) GetUserClaimedPoints(address string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	err := d.users.FindOne(ctx, bson.M{"address": address}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil 
		}
		return 0, fmt.Errorf("failed to find user: %v", err)
	}

	return user.TotalClaimedPoints, nil
} 