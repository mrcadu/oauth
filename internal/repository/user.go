package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"oauth/internal/model"
	"oauth/internal/service/password_encription"
)

type UserRepositoryMongo struct {
}

func (u UserRepositoryMongo) CreateUser(user model.User) (model.User, error) {
	user.Password, _ = password_encription.HashPassword(user.Password)
	user.ID = primitive.NewObjectID()
	_, err := u.getCollection().InsertOne(context.TODO(), user)
	return user, err
}

func (u UserRepositoryMongo) UpdateUser(user model.User) (model.User, error) {
	password, err := password_encription.HashPassword(user.Password)
	user.Password = password
	updateResult, err := u.getCollection().ReplaceOne(context.TODO(), bson.D{{"_id", user.ID}}, user)
	if updateResult != nil && updateResult.ModifiedCount == 0 {
		return user, mongo.ErrNoDocuments
	}
	return user, err
}

func (u UserRepositoryMongo) GetUser(id primitive.ObjectID) (model.User, error) {
	var user model.User
	err := u.getCollection().FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&user)
	return user, err
}

func (u UserRepositoryMongo) DeleteUser(id primitive.ObjectID) (primitive.ObjectID, error) {
	deleteResult, err := u.getCollection().DeleteOne(context.TODO(), id)
	if deleteResult != nil && deleteResult.DeletedCount == 0 {
		return id, mongo.ErrNoDocuments
	}
	return id, err
}

func (u UserRepositoryMongo) getCollection() *mongo.Collection {
	return model.DB.Database("auth").Collection("User")
}
