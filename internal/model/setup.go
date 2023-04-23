package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
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
	_, err = client.Database("auth").Collection("Profile").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{"name", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic(err)
	}
	_, err = client.Database("auth").Collection("User").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{"username", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic(err)
	}
	DB = client
}
