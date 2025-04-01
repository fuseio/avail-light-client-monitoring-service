package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"monitoring-service/internal/uptime"
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
	Address                string    `bson:"address"`
	TotalTime              int64     `bson:"total_time"`
	LastHeartbeat          time.Time `bson:"last_heartbeat"`
	CreatedAt              time.Time `bson:"created_at"`
	NFTAmount              int64     `bson:"nft_amount"`
	CommissionRate         float64   `bson:"commission_rate"`
	Status                 string    `bson:"status"`
	AllUptimePercentage     float64   `bson:"all_uptime_percentage"`
	WeeklyUptimePercentage  float64   `bson:"weekly_uptime_percentage"`
	OperatorName           string    `bson:"operator_name"`
	RewardCollectorAddress  string    `bson:"reward_collector_address"`
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
	
	// Check if heartbeats collection exists before creating it
	collections, err := db.ListCollectionNames(ctx, bson.M{"name": "heartbeats"})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %v", err)
	}
	
	// Only create collection if it doesn't exist
	if len(collections) == 0 {
		err = db.CreateCollection(ctx, "heartbeats", options.CreateCollection().SetTimeSeriesOptions(
			options.TimeSeries().
				SetTimeField("timestamp").
				SetMetaField("client_address").
				SetGranularity("minutes"),
		))
		if err != nil {
			return nil, fmt.Errorf("failed to create heartbeats collection: %v", err)
		}
		logger.Println("Created time series collection: heartbeats")
	} else {
		logger.Println("Heartbeats collection already exists, skipping creation")
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

func (d *Database) RegisterClient(address string, operationPoints OperationPointRecord, totalTime int64, operatorName string, rewardCollectorAddress string) error {
	ctx := context.Background()
	now := time.Now()

	
	address = strings.ToLower(address)
	rewardCollectorAddress = strings.ToLower(rewardCollectorAddress)

	// Update all relevant fields in ClientInfo
	clientUpdate := bson.M{
		"$set": bson.M{
			"address":                address,
			"total_time":             totalTime,
			"last_heartbeat":         now,
			"nft_amount":             operationPoints.Amount,
			"commission_rate":        operationPoints.CommissionRate,
			"operator_name":          operatorName,
			"reward_collector_address": rewardCollectorAddress,
		},
		"$setOnInsert": bson.M{
			"created_at": now,
		},
	}

	// Record heartbeat only if amount > 0 and time > 0.
	if operationPoints.Amount > 0 && operationPoints.Time > 0 {
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

	_, err := d.clients.UpdateOne(
		ctx,
		bson.M{"address": address},
		clientUpdate,
		options.Update().SetUpsert(true),
	)
	return err
}

func (d *Database) RegisterDelegation(address string, delegationPoints DelegationPointRecord) error {
	ctx := context.Background()
	now := time.Now()
	
	
	address = strings.ToLower(address)
	delegationPoints.Address = strings.ToLower(delegationPoints.Address)

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

	return err
}

func (d *Database) ClientExists(address string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	
	address = strings.ToLower(address)
	
	count, err := d.clients.CountDocuments(ctx, bson.M{"address": address})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *Database) GetClient(address string) (*ClientInfo, error) {
	ctx := context.Background()
	var client ClientInfo
	
	
	address = strings.ToLower(address)
	
	err := d.clients.FindOne(ctx, bson.M{"address": address}).Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &client, nil
}

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

	// Create uptime calculator
	uptimeCalc := uptime.NewCalculator(d.heartbeats)

	// Calculate uptime percentages and status for each client
	for i, client := range clients {
		// Calculate uptime percentages (all-time and weekly)
		allUptimePercentage, weeklyUptimePercentage, err := uptimeCalc.GetUptimePercentages(
			ctx, 
			client.Address,
			client.CreatedAt,
		)
		
		if err != nil {
			d.logger.Printf("Error calculating uptime for client %s: %v", client.Address, err)
			// Continue with next client instead of failing entirely
			continue
		}
		
		// Update client with calculated percentages
		clients[i].AllUptimePercentage = allUptimePercentage
		clients[i].WeeklyUptimePercentage = weeklyUptimePercentage

		// Set status based on last heartbeat
		if time.Since(client.LastHeartbeat) > 10 * time.Minute {
			clients[i].Status = "Offline"
		} else if time.Since(client.LastHeartbeat) > 5 * time.Minute {
			clients[i].Status = "Inactive"
		} else {
			clients[i].Status = "Active"
		}
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
	
	
	address = strings.ToLower(address)
	
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

func (d *Database) GetFromDelegationsByAddress(address string) ([]DelegationRecord, error) {
	ctx := context.Background()
	
	
	address = strings.ToLower(address)

	cursor, err := d.delegations.Find(ctx, bson.M{
		"from_address": address,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var delegations []DelegationRecord
	if err = cursor.All(ctx, &delegations); err != nil {
		return nil, err
	}

	return delegations, nil
}

// ClearDelegationsForAddress removes all delegation records for a specific address
// that are no longer valid based on the current blockchain state
func (d *Database) ClearDelegationsForAddress(address string, validFromAddresses []string) error {
	ctx := context.Background()
	
	
	address = strings.ToLower(address)
	normalizedValidAddrs := make([]string, len(validFromAddresses))
	for i, addr := range validFromAddresses {
		normalizedValidAddrs[i] = strings.ToLower(addr)
	}
	
	// If we have valid delegations, only remove the ones not in the list
	if len(normalizedValidAddrs) > 0 {
		filter := bson.M{
			"to_address": address,
			"from_address": bson.M{
				"$nin": normalizedValidAddrs,
			},
		}
		
		_, err := d.delegations.DeleteMany(ctx, filter)
		return err
	}
	
	// If no valid delegations, remove all delegations to this address
	_, err := d.delegations.DeleteMany(ctx, bson.M{
		"to_address": address,
	})
	
	return err
}
