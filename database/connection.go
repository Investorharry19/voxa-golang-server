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

var DB *mongo.Database

func ConnectMongoDB() {
	uri := "mongodb://localhost:27017" // <-- replace this
	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Mongo connection error:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo ping error:", err)
	}

	DB = client.Database("voxa_temp")
	fmt.Println("✅ Connected to MongoDB!")

	// Ensure unique indexes for users collection: email and username
	usersColl := DB.Collection("users")
	// Create index models

	usernameIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("username_unique"),
	}

	// Create indexes with a context timeout
	ctxIdx, cancelIdx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelIdx()
	if _, err := usersColl.Indexes().CreateMany(ctxIdx, []mongo.IndexModel{usernameIndex}); err != nil {
		// Index creation failure should not panic the app, but log it for debugging
		log.Printf("warning: could not create user indexes: %v", err)
	}
}

func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}
