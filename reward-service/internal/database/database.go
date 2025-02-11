package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client          *mongo.Client
	users           *mongo.Collection
	rewards         *mongo.Collection
	logger          *log.Logger
}

type User struct {
    Address      	 string `bson:"address"`
    DelegationsCount int64  `bson:"delegationCount"`
    GlobalPoints 	 int64  `bson:"global_points"`
    Uptime 			 string `bson:"percentage"`
    WeeklyUptime 	 string `bson:"weeklypercentage"`
    Comission 		 string `bson:"percentage"`
    Status 			 string `bson:"percentage"`
}

type RewardRecord struct {
    Address        string    `bson:"address"`
    Points         int64     `bson:"points"`
    Timestamp      time.Time `bson:"timestamp"`
    NFTCount       int64     `bson:"nft_count"`
    DelegationCount int64    `bson:"delegation_count"`
    CommissionRate float64   `bson:"commission_rate"`
}

func NewDatabase(mongoURI, dbName string, logger *log.Logger) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	db := client.Database(dbName)

	return &Database{
		client:      client,
		users:       db.Collection("users"),
		rewards:     db.Collection("rewards"),
		logger:      logger,
	}, nil
}

func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.client.Disconnect(ctx)
}

