package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	client          *mongo.Client
	db              *mongo.Database
	logger          *log.Logger
	users           *mongo.Collection
	rewards         *mongo.Collection
	rewardSummaries *mongo.Collection
}

func New(client *mongo.Client, dbName string, logger *log.Logger) *Database {
	db := client.Database(dbName)
	return &Database{
		client:          client,
		db:              db,
		logger:          logger,
		users:           db.Collection("users"),
		rewards:         db.Collection("rewards"),
		rewardSummaries: db.Collection("reward_summaries"),
	}
}

func (d *Database) GetMongoDatabase() *mongo.Database {
	return d.db
}

func (d *Database) GetMongoClient() *mongo.Client {
	return d.client
}

func (d *Database) GetDelegationsToOperator(address string) ([]DelegationRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := d.db.Collection("delegations").Find(ctx, bson.M{
		"to_address": address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch delegations: %v", err)
	}
	defer cursor.Close(ctx)

	var delegations []DelegationRecord
	if err := cursor.All(ctx, &delegations); err != nil {
		return nil, fmt.Errorf("failed to decode delegations: %v", err)
	}

	return delegations, nil
}

func (d *Database) GetAllOperators() ([]MonitoringClientInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := d.db.Collection("clients").Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch operators: %v", err)
	}
	defer cursor.Close(ctx)

	var operators []MonitoringClientInfo
	if err := cursor.All(ctx, &operators); err != nil {
		return nil, fmt.Errorf("failed to decode operators: %v", err)
	}

	return operators, nil
}

func (d *Database) GetAllDelegators() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := d.users.Find(ctx, bson.M{"user_type": "DELEGATOR"})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch delegators: %v", err)
	}
	defer cursor.Close(ctx)

	var delegators []User
	if err := cursor.All(ctx, &delegators); err != nil {
		return nil, fmt.Errorf("failed to decode delegators: %v", err)
	}

	return delegators, nil
}

func (d *Database) GetClient(operatorAddress string) (*MonitoringClientInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var user User
	err := d.users.FindOne(ctx, bson.M{
		"address":   operatorAddress,
		"user_type": "OPERATOR",
	}).Decode(&user)
	
	if err != nil {
		return nil, fmt.Errorf("failed to fetch client: %v", err)
	}
	
	d.logger.Printf("Found operator %s in database with uptime: %.2f%%", 
		operatorAddress, user.AllUptimePercentage)
	
	client := &MonitoringClientInfo{
		Address:               operatorAddress,
		LastHeartbeat:         time.Now(), 
		AllUptimePercentage:   user.AllUptimePercentage,
		Status:                user.Status,
		CommissionRate:        user.CommissionRate,
		NFTAmount:             0,
		RewardCollectorAddress: user.RewardAddress,
	}
	
	return client, nil
}

func (d *Database) StoreRewardRecord(record *RewardRecord) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := d.rewards.InsertOne(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to store reward: %v", err)
	}

	return nil
}

// Add indexes for better performance
func (db *Database) AddIndexes() error {
	_, err := db.db.Collection("rewards").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_address", Value: 1},
			{Key: "cycle_id", Value: 1},
		},
		Options: options.Index().SetUnique(false),
	})
	if err != nil {
		db.logger.Printf("Warning: Failed to create index on rewards collection: %v", err)
		return err
	}
	return nil
}

func (d *Database) GetDelegationsFromDelegator(delegatorAddress string) ([]DelegationRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	cursor, err := d.db.Collection("delegations").Find(ctx, bson.M{"from_address": delegatorAddress})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch delegations from delegator %s: %v", delegatorAddress, err)
	}
	defer cursor.Close(ctx)
	
	var delegations []DelegationRecord
	if err := cursor.All(ctx, &delegations); err != nil {
		return nil, fmt.Errorf("failed to decode delegations: %v", err)
	}
	
	return delegations, nil
}

func (d *Database) GetUserByAddress(address string, userType string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var user User
	err := d.users.FindOne(ctx, bson.M{
		"address":   address,
		"user_type": userType,
	}).Decode(&user)
	
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}
	
	return &user, nil
}

func (d *Database) StoreRewardSummary(summary *RewardSummary) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingSummary RewardSummary
	err := d.rewardSummaries.FindOne(ctx, bson.M{"cycle_id": summary.CycleID}).Decode(&existingSummary)
	if err == nil {
		return fmt.Errorf("rewards for cycle %s have already been processed", summary.CycleID)
	} else if err != mongo.ErrNoDocuments {
		return fmt.Errorf("error checking existing reward summary: %v", err)
	}

	_, err = d.rewardSummaries.InsertOne(ctx, summary)
	if err != nil {
		return fmt.Errorf("failed to store reward summary: %v", err)
	}

	return nil
}

func Connect(uri, dbName string, logger *log.Logger) (*Database, error) {
	var client *mongo.Client
	var err error
	
	maxRetries := 5
	retryDelay := 2 * time.Second
	
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		clientOptions := options.Client().
			ApplyURI(uri).
			SetRetryWrites(true).
			SetRetryReads(true).
			SetMaxPoolSize(100).
			SetMinPoolSize(10)
		
		client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer pingCancel()
			
			err = client.Ping(pingCtx, readpref.Primary())
			if err == nil {
				logger.Println("Successfully connected to MongoDB")
				break
			}
		}
		
		logger.Printf("Failed to connect to MongoDB (attempt %d/%d): %v", i+1, maxRetries, err)
		
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
			retryDelay *= 2 // Exponential backoff
		}
	}
	
	if err != nil {
		return nil, errors.New("failed to connect to MongoDB after multiple attempts")
	}
	
	db := New(client, dbName, logger)
	
	if err := db.AddIndexes(); err != nil {
		logger.Printf("Warning: Failed to create indexes: %v", err)
	}
	
	return db, nil
}

func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.client.Disconnect(ctx)
}