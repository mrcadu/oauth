package model

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"oauth/internal/config"
	"time"
)

var DB *mongo.Client

func ConnectDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.GetProperty("DB_URL")))
	if err != nil {
		log.Fatal(err)
	}
	DB = client
}
