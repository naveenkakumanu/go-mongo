package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnableSharding enables sharding for a given database using the admin command
func EnableSharding(client *mongo.Client, dbName string) error {
	adminDB := client.Database("admin")
	cmd := map[string]interface{}{
		"enableSharding": dbName,
	}
	res := adminDB.RunCommand(context.Background(), cmd)
	return res.Err()
}

// ShardCollection shards a collection on the specified key
func ShardCollection(client *mongo.Client, dbName, collName, shardKey string) error {
	adminDB := client.Database("admin")
	cmd := map[string]interface{}{
		"shardCollection": fmt.Sprintf("%s.%s", dbName, collName),
		"key":             map[string]interface{}{shardKey: 1},
	}
	res := adminDB.RunCommand(context.Background(), cmd)
	return res.Err()
}

// Example use case: Sharding the "users" collection on the "email" field as the shard key.
// This is useful if user lookups are mostly by email and emails are well-distributed.
func ShardUsersByEmail(client *mongo.Client, dbName string) error {
	// Enable sharding on the database
	if err := EnableSharding(client, dbName); err != nil {
		return fmt.Errorf("failed to enable sharding: %w", err)
	}
	// Shard the "users" collection on the "email" field
	if err := ShardCollection(client, dbName, "users", "email"); err != nil {
		return fmt.Errorf("failed to shard users collection: %w", err)
	}
	return nil
}

// MongoDBConfig holds configuration for MongoDB connection
type MongoDBConfig struct {
	URI     string
	Timeout time.Duration
}

// NewMongoClient creates and returns a new MongoDB client
func NewMongoClient(cfg MongoDBConfig) (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	clientOpts := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		cancel()
		return nil, nil, nil, fmt.Errorf("mongo connect error: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(ctx)
		cancel()
		return nil, nil, nil, fmt.Errorf("mongo ping error: %w", err)
	}
	return client, ctx, cancel, nil
}

func main() {
	cfg := MongoDBConfig{
		URI:     getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		Timeout: 10 * time.Second,
	}

	client, _, cancel, err := NewMongoClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer func() {
		cancel()
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Failed to disconnect MongoDB: %v", err)
		}
	}()

	fmt.Println("Connected to MongoDB!")

	// Graceful shutdown on SIGINT/SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("Shutting down gracefully...")
}

// getEnv returns the value of the environment variable or fallback if not set
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
