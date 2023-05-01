package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"oauth/api/http_error"
	"oauth/internal/model"
	"oauth/internal/repository"
)

var userRepository = repository.UserRepositoryMongo{}

func CreateUser(context *gin.Context) {
	var user model.User
	err := context.ShouldBind(&user)
	if err != nil {
		panic(err)
	}
	createdUser, err := userRepository.CreateUser(user)
	if err != nil {
		panic(err)
	}
	context.JSON(http.StatusCreated, createdUser)
}

func DeleteUser(context *gin.Context) {
	id := context.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	deletedUser, err := userRepository.DeleteUser(hex)
	if err != nil {
		panic(err)
	}
	context.JSON(http.StatusOK, deletedUser)
}

func UpdateUser(context *gin.Context) {
	var user model.User
	err := context.ShouldBind(&user)
	id := context.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	user.ID = hex
	updatedUser, err := userRepository.UpdateUser(user)
	if err != nil {
		panic(err)
	}
	context.JSON(http.StatusOK, updatedUser)
}

func GetUser(context *gin.Context) {
	id := context.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		context.JSON(401, http_error.Unauthorized("user", ""))
		return
	}
	user, err := userRepository.GetUser(hex)
	if err != nil {
		context.JSON(err.(http_error.HttpError).Status, err)
	} else {
		context.JSON(http.StatusOK, user)
	}
}
