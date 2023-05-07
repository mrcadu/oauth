package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"oauth/internal/model"
)

type ProfileRepositoryMongo struct {
}

func (p ProfileRepositoryMongo) CreateProfile(profile model.Profile) (model.Profile, error) {
	profile.ID = primitive.NewObjectID()
	_, err := p.getCollection().InsertOne(context.TODO(), profile)
	return profile, err
}
func (p ProfileRepositoryMongo) GetProfile(id primitive.ObjectID) (model.Profile, error) {
	var recoveredProfile model.Profile
	err := p.getCollection().FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&recoveredProfile)
	return recoveredProfile, err
}
func (p ProfileRepositoryMongo) GetProfileByName(name string) (model.Profile, error) {
	var recoveredProfile model.Profile
	err := p.getCollection().FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&recoveredProfile)
	return recoveredProfile, err
}
func (p ProfileRepositoryMongo) UpdateProfile(profile model.Profile) (model.Profile, error) {
	updateResult, err := p.getCollection().ReplaceOne(context.TODO(), bson.D{{"_id", profile.ID}}, profile)
	if updateResult != nil && updateResult.ModifiedCount == 0 {
		return profile, mongo.ErrNoDocuments
	}
	return profile, err
}
func (p ProfileRepositoryMongo) DeleteProfile(id primitive.ObjectID) (primitive.ObjectID, error) {
	deleteResult, err := p.getCollection().DeleteOne(context.TODO(), bson.D{{"_id", id}})
	if deleteResult != nil && deleteResult.DeletedCount == 0 {
		return id, mongo.ErrNoDocuments
	}
	return id, err
}

func (p ProfileRepositoryMongo) getCollection() *mongo.Collection {
	return model.DB.Database("auth").Collection("Profile")
}
