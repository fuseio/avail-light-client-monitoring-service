package database

import (
	"context"
	"log"
		
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client *mongo.Client
	db     *mongo.Database
	users  *mongo.Collection
	logger *log.Logger
}

// NewDatabase initializes a new database connection
func NewDatabase(client *mongo.Client, db *mongo.Database, logger *log.Logger) *Database {
	return &Database{
		client: client,
		db:     db,
		users:  db.Collection("users"),
		logger: logger,
	}
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

// Add this method to the Database struct
func (db *Database) GetMongoDatabase() *mongo.Database {
	return db.db
}