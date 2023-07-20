package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"oauth/internal/model"
	"oauth/internal/model/datasource"
	"oauth/internal/utils"
)

type UserRepositoryMongo struct {
	datasource                datasource.MongoDatasource
	passwordEncryptionService utils.PasswordEncryption
	profileRepositoryMongo    ProfileRepositoryMongo
}

func (u UserRepositoryMongo) Create(user model.User) (model.User, error) {
	user.Password, _ = u.passwordEncryptionService.HashPassword(user.Password)
	user.ID = primitive.NewObjectID()
	_, err := u.getCollection().InsertOne(context.TODO(), user)
	return user, err
}

func (u UserRepositoryMongo) Update(user model.User) (model.User, error) {
	password, err := u.passwordEncryptionService.HashPassword(user.Password)
	user.Password = password
	updateResult, err := u.getCollection().ReplaceOne(context.TODO(), bson.D{{"_id", user.ID}}, user)
	if updateResult != nil && updateResult.ModifiedCount == 0 {
		return user, mongo.ErrNoDocuments
	}
	return user, err
}

func (u UserRepositoryMongo) Get(id primitive.ObjectID) (model.User, error) {
	var user model.User
	err := u.getCollection().FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&user)
	return user, err
}

func (u UserRepositoryMongo) GetByUsername(username string) (model.User, error) {
	var user model.User
	err := u.getCollection().FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&user)
	return user, err
}

func (u UserRepositoryMongo) Delete(id primitive.ObjectID) (primitive.ObjectID, error) {
	deleteResult, err := u.getCollection().DeleteOne(context.TODO(), bson.D{{"_id", id}})
	if deleteResult != nil && deleteResult.DeletedCount == 0 {
		return id, mongo.ErrNoDocuments
	}
	return id, err
}

func (u UserRepositoryMongo) getCollection() *mongo.Collection {
	return u.datasource.GetClient().Database("auth").Collection("User")
}

func NewUserRepository() UserRepositoryMongo {
	return UserRepositoryMongo{
		datasource:                datasource.GetMongoDatasource(),
		passwordEncryptionService: utils.PasswordEncryptionImpl{},
		profileRepositoryMongo:    NewProfileRepository(),
	}
}
