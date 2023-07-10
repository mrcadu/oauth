package datasource

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"oauth/internal/config"
	"time"
)

type MongoDatasource interface {
	SetupMongo()
	CreateMongoIndexes()
	GetClient() *mongo.Client
}

type MongoDatasourceImpl struct {
	client *mongo.Client
}

var client *mongo.Client

func (m MongoDatasourceImpl) GetClient() *mongo.Client {
	return client
}

func (m MongoDatasourceImpl) Setup() {
	m.SetupMongo()
}

func (m MongoDatasourceImpl) SetupMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.GetProperty("DB_URL")))
	if err != nil {
		panic(err)
	}
	client = mongoClient
	m.CreateMongoIndexes()
}

func (m MongoDatasourceImpl) CreateMongoIndexes() {
	_, err := client.Database("auth").Collection("Profile").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
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
}
