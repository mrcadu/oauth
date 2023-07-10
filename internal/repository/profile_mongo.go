package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"oauth/internal/model"
	"oauth/internal/model/datasource"
)

type ProfileRepositoryMongo struct {
	datasource datasource.MongoDatasource
}

func (p ProfileRepositoryMongo) Create(profile model.Profile) (model.Profile, error) {
	profile.ID = primitive.NewObjectID()
	_, err := p.getCollection().InsertOne(context.TODO(), profile)
	return profile, err
}
func (p ProfileRepositoryMongo) Get(id primitive.ObjectID) (model.Profile, error) {
	var recoveredProfile model.Profile
	err := p.getCollection().FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&recoveredProfile)
	return recoveredProfile, err
}
func (p ProfileRepositoryMongo) GetProfileByName(name string) (model.Profile, error) {
	var recoveredProfile model.Profile
	err := p.getCollection().FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&recoveredProfile)
	return recoveredProfile, err
}
func (p ProfileRepositoryMongo) Update(profile model.Profile) (model.Profile, error) {
	updateResult, err := p.getCollection().ReplaceOne(context.TODO(), bson.D{{"_id", profile.ID}}, profile)
	if updateResult != nil && updateResult.ModifiedCount == 0 {
		return profile, mongo.ErrNoDocuments
	}
	return profile, err
}
func (p ProfileRepositoryMongo) Delete(id primitive.ObjectID) (primitive.ObjectID, error) {
	deleteResult, err := p.getCollection().DeleteOne(context.TODO(), bson.D{{"_id", id}})
	if deleteResult != nil && deleteResult.DeletedCount == 0 {
		return id, mongo.ErrNoDocuments
	}
	return id, err
}

func (p ProfileRepositoryMongo) getCollection() *mongo.Collection {
	return p.datasource.GetClient().Database("auth").Collection("Profile")
}

func NewProfileRepository() ProfileRepositoryMongo {
	return ProfileRepositoryMongo{
		datasource: datasource.GetMongoDatasource(),
	}
}
