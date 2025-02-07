package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	client          *mongo.Client
	clients         *mongo.Collection
	heartbeats      *mongo.Collection
	delegations     *mongo.Collection
	logger          *log.Logger
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
	Address        string    `bson:"address"`
	TotalTime      int64     `bson:"total_time"`
	LastHeartbeat  time.Time `bson:"last_heartbeat"`
	CreatedAt      time.Time `bson:"created_at"`
	NFTAmount      int64     `bson:"nft_amount"`
	CommissionRate float64   `bson:"commission_rate"`
}

type HeartbeatRecord struct {
	ClientAddress  string    `bson:"client_address"`
	Timestamp      time.Time `bson:"timestamp"`
	Duration       int64     `bson:"duration"`
	Amount         int64     `bson:"amount"`
}

type DelegationRecord struct {
	FromAddress    string    `bson:"from_address"`
	ToAddress      string    `bson:"to_address"`
	Amount         int64     `bson:"amount"`
	CommissionRate float64   `bson:"commission_rate"`
	Timestamp      time.Time `bson:"timestamp"`
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
	
	// Create time series collection for heartbeats
	err = db.CreateCollection(ctx, "heartbeats", options.CreateCollection().SetTimeSeriesOptions(
		options.TimeSeries().
			SetTimeField("timestamp").
			SetMetaField("client_address").
			SetGranularity("minutes"),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create heartbeats collection: %v", err)
	}

	// Get collection using the provided database name
	collection := db.Collection("clients")

	// Create unique index on address if it doesn't exist
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "address", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %v", err)
	}

	return &Database{
		client:      client,
		clients:     collection,
		heartbeats:  db.Collection("heartbeats"),
		delegations: db.Collection("delegations"),
		logger:      logger,
	}, nil
}

func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.client.Disconnect(ctx)
}

func (d *Database) RegisterClient(address string, operationPoints OperationPointRecord, delegationPoints *DelegationPointRecord, totalTime int64) error {
	ctx := context.Background()
	now := time.Now()

	// Update client info
	clientUpdate := bson.M{
		"$set": bson.M{
			"address": address,
			"total_time": totalTime,
			"last_heartbeat": now,
		},
		"$setOnInsert": bson.M{
			"created_at": now,
		},
	}

	// Add operation-specific fields if present
	if operationPoints.Amount > 0 {
		clientUpdate["$set"].(bson.M)["nft_amount"] = operationPoints.Amount
		clientUpdate["$set"].(bson.M)["commission_rate"] = operationPoints.CommissionRate

		// Record heartbeat if time > 0
		if operationPoints.Time > 0 {
			heartbeat := HeartbeatRecord{
				ClientAddress: address,
				Timestamp:     now,
				Duration:      operationPoints.Time,
				Amount:        operationPoints.Amount,
			}
			if _, err := d.heartbeats.InsertOne(ctx, heartbeat); err != nil {
				return err
			}
		}
	}

	// Record delegation ONLY if delegation record is provided
	if delegationPoints != nil {
		delegation := DelegationRecord{
			FromAddress:    delegationPoints.Address,
			ToAddress:      address,
			Amount:         delegationPoints.Amount,
			CommissionRate: delegationPoints.CommissionRate,
			Timestamp:      now,
		}
		
		filter := bson.M{
			"from_address": delegationPoints.Address,
			"to_address":   address,
		}
		update := bson.M{"$set": delegation}
		
		_, err := d.delegations.UpdateOne(
			ctx,
			filter,
			update,
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return err
		}
	}

	_, err := d.clients.UpdateOne(
		ctx,
		bson.M{"address": address},
		clientUpdate,
		options.Update().SetUpsert(true),
	)
	return err
}

func (d *Database) ClientExists(address string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := d.clients.CountDocuments(ctx, bson.M{"address": address})
	if err != nil {
		d.logger.Printf("Error checking client existence: %v", err)
		return false, err
	}
	return count > 0, nil
}

func (d *Database) GetClient(address string) (*ClientInfo, error) {
	ctx := context.Background()
	var client ClientInfo
	
	err := d.clients.FindOne(ctx, bson.M{"address": address}).Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		d.logger.Printf("Error getting client: %v", err)
		return nil, err
	}

	return &client, nil
}

// Add this new function at the end of the file
func (d *Database) GetAllClients() ([]ClientInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := d.clients.Find(ctx, bson.M{}, opts)
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

func (d *Database) setupIndexes(ctx context.Context) error {
	_, err := d.heartbeats.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "timestamp", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(60 * 60 * 24 * 30),
	})
	return err
}

func (d *Database) GetClientWithHistory(address string) (*ClientInfo, []HeartbeatRecord, []DelegationRecord, error) {
	ctx := context.Background()
	
	// Get client info
	client, err := d.GetClient(address)
	if err != nil {
		return nil, nil, nil, err
	}

	// Get recent heartbeats (last 24h)
	cutoff := time.Now().Add(-24 * time.Hour)
	heartbeatCursor, err := d.heartbeats.Find(ctx, bson.M{
		"client_address": address,
		"timestamp": bson.M{"$gte": cutoff},
	})
	if err != nil {
		return nil, nil, nil, err
	}
	var heartbeats []HeartbeatRecord
	if err = heartbeatCursor.All(ctx, &heartbeats); err != nil {
		return nil, nil, nil, err
	}

	// Get delegations
	delegationCursor, err := d.delegations.Find(ctx, bson.M{
		"$or": []bson.M{
			{"from_address": address},
			{"to_address": address},
		},
	})
	if err != nil {
		return nil, nil, nil, err
	}
	var delegations []DelegationRecord
	if err = delegationCursor.All(ctx, &delegations); err != nil {
		return nil, nil, nil, err
	}

	return client, heartbeats, delegations, nil
}
