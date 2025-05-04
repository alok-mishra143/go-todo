package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Collection

func ConnectDB() *mongo.Client {
	MONGO_URI := os.Getenv("DB_URL")

	clientOptions := options.Client().ApplyURI(MONGO_URI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Connected to MongoDB")
	DB = client.Database("golang_db").Collection("todos")
	return client
}

func GetCollection() *mongo.Collection {
	return DB
}

func CloseDB(client *mongo.Client) {
	if err := client.Disconnect(context.Background()); err != nil {
		log.Fatal("Error disconnecting MongoDB:", err)
	}
}
