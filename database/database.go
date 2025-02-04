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

type OwnershipStatus int

const (
	NoNFT         OwnershipStatus = iota // 0
	OwnsNFT                              // 1
	HasDelegation                        // 2
)

type Database struct {
	client     *mongo.Client
	collection *mongo.Collection
	logger     *log.Logger
}

type OperationPointRecord struct {
	Amount         int64     `bson:"amount"`
	Timestamp     time.Time  `bson:"timestamp"`
	CommissionRate float64   `bson:"commission_rate"`
	Time           int64     `bson:"time"`
}

type DelegationPointRecord struct {
	Address        string    `bson:"address"`
	Amount         int64     `bson:"amount"`
	Timestamp     time.Time  `bson:"timestamp"`
	CommissionRate float64   `bson:"commission_rate"`
	Time           int64     `bson:"time"`
}

type ClientInfo struct {
	Address          string        `bson:"address"`
	OperationPoints  []OperationPointRecord `bson:"operation_points"`
	DelegationPoints []DelegationPointRecord `bson:"delegation_points"`
	TotalTime        int64         `bson:"total_time"`
	LastHeartbeat    time.Time     `bson:"last_heartbeat"`
	CreatedAt        time.Time     `bson:"created_at"`
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

	// Get collection using the provided database name
	collection := client.Database(dbName).Collection("user")

	// Create unique index on address if it doesn't exist
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "address", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %v", err)
	}

	return &Database{
		client:     client,
		collection: collection,
		logger:     logger,
	}, nil
}

func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.client.Disconnect(ctx)
}

func (d *Database) RegisterClient(address string, operationPoints OperationPointRecord, delegationPoints DelegationPointRecord, totalTime int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{"address": address}
	
	update := bson.M{
		"$set": bson.M{
			"address":         address,
		},
		"$setOnInsert": bson.M{
			"created_at": now,
		},
	}
	// check if operationPoints is not empty
	if operationPoints.Amount > 0 {
		update["$set"] = bson.M{
			"total_time":     totalTime,
			"last_heartbeat": now,
		}
		update["$push"] = bson.M{
			"operation_points": operationPoints,
		}
	}
	// check if delegationPoints is not empty
	if delegationPoints.Amount > 0 {
		update["$push"] = bson.M{
			"delegation_points": delegationPoints,
		}
	}

	opts := options.Update().SetUpsert(true)
	_, err := d.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		d.logger.Printf("Error registering client: %v", err)
		return err
	}
	return nil
}

func (d *Database) ClientExists(address string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := d.collection.CountDocuments(ctx, bson.M{"address": address})
	if err != nil {
		d.logger.Printf("Error checking client existence: %v", err)
		return false, err
	}
	return count > 0, nil
}

func (d *Database) GetClient(address string) (ClientInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var raw bson.M
	err := d.collection.FindOne(ctx, bson.M{"address": address}).Decode(&raw)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ClientInfo{}, nil
		}
		d.logger.Printf("Error getting client: %v", err)
		return ClientInfo{}, err
	}

	client := ClientInfo{
		Address: raw["address"].(string),
	}

	// Handle optional fields with proper type conversion
	if raw["last_heartbeat"] != nil {
		if timestamp, ok := raw["last_heartbeat"].(primitive.DateTime); ok {
			client.LastHeartbeat = timestamp.Time()
		}
	}
	if raw["created_at"] != nil {
		if timestamp, ok := raw["created_at"].(primitive.DateTime); ok {
			client.CreatedAt = timestamp.Time()
		}
	}
	if raw["total_time"] != nil {
		client.TotalTime = raw["total_time"].(int64)
	}

	// Handle operation points array
	if opPoints, ok := raw["operation_points"].(primitive.A); ok {
		client.OperationPoints = make([]OperationPointRecord, 0, len(opPoints))
		for _, point := range opPoints {
			if pointMap, ok := point.(bson.M); ok {
				timestamp := pointMap["timestamp"].(primitive.DateTime).Time()
				record := OperationPointRecord{
					Amount:         pointMap["amount"].(int64),
					Timestamp:      timestamp,
					CommissionRate: pointMap["commission_rate"].(float64),
					Time:           pointMap["time"].(int64),
				}
				client.OperationPoints = append(client.OperationPoints, record)
			}
		}
	}

	// Handle delegation points array
	if delPoints, ok := raw["delegation_points"].(primitive.A); ok {
		client.DelegationPoints = make([]DelegationPointRecord, 0, len(delPoints))
		for _, point := range delPoints {
			if pointMap, ok := point.(bson.M); ok {
				timestamp := pointMap["timestamp"].(primitive.DateTime).Time()
				record := DelegationPointRecord{
					Address:        pointMap["address"].(string),
					Amount:         pointMap["amount"].(int64),
					Timestamp:      timestamp,
					CommissionRate: pointMap["commission_rate"].(float64),
					Time:           pointMap["time"].(int64),
				}
				client.DelegationPoints = append(client.DelegationPoints, record)
			}
		}
	}

	return client, nil
}

// Add this new function at the end of the file
func (d *Database) GetAllClients() ([]ClientInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := d.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		d.logger.Printf("Error querying clients: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var clients []ClientInfo
	if err = cursor.All(ctx, &clients); err != nil {
		d.logger.Printf("Error decoding clients: %v", err)
		return nil, err
	}

	return clients, nil
}
