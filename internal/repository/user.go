package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"oauth/api/http_error"
	"oauth/internal/model"
	"oauth/internal/service/password_encription"
)

type UserRepositoryMongo struct {
}

func (u UserRepositoryMongo) CreateUser(user model.User) (model.User, error) {
	var existingUser model.User
	err := getCollection().FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&existingUser)
	if err == nil {
		return existingUser, http_error.ConflictError("user", user.Username)
	}
	user.Password, _ = password_encription.HashPassword(user.Password)
	user.ID = primitive.NewObjectID()
	result, err := getCollection().InsertOne(context.TODO(), user)
	if err != nil {
		return model.User{}, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, err
}

func (u UserRepositoryMongo) UpdateUser(oldUsername string, user model.User) (model.User, error) {
	existingUser, err := u.GetUser(oldUsername)
	if err != nil {
		return existingUser, http_error.NotFoundError("user", oldUsername)
	}
	password, err := password_encription.HashPassword(user.Password)
	if err != nil {
		return existingUser, http_error.BadRequest()
	}
	userToUpdate := model.User{
		ID:       existingUser.ID,
		Password: password,
		Username: user.Username,
	}
	_, err = getCollection().ReplaceOne(context.TODO(), bson.D{{"username", oldUsername}}, userToUpdate)
	return userToUpdate, err
}

func (u UserRepositoryMongo) GetUser(username string) (model.User, error) {
	var user model.User
	err := getCollection().FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&user)
	if err != nil {
		return user, http_error.NotFoundError("user", username)
	}
	return user, err
}

func (u UserRepositoryMongo) DeleteUser(username string, password string) (string, error) {
	var existingUser model.User
	err := getCollection().FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&existingUser)
	if err != nil {
		return "", http_error.NotFoundError("user", username)
	}
	isPasswordCorrect := password_encription.CheckPasswordHash(password, existingUser.Password)
	if isPasswordCorrect {
		_, err := getCollection().DeleteOne(context.TODO(), bson.D{{"username", username}, {"password", existingUser.Password}})
		return username, err
	} else {
		return "", http_error.NotFoundError("user", username)
	}
}

func getCollection() *mongo.Collection {
	return model.DB.Database("auth").Collection("User")
}
