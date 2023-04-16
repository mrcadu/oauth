package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"oauth/api/http_error"
	"oauth/internal/model"
	"oauth/internal/repository"
	"oauth/internal/service/auth"
)

var userRepository = repository.UserRepositoryMongo{}

func CreateUser(context *gin.Context) {
	var user model.User
	err := context.ShouldBind(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}

	createdUser, err := userRepository.CreateUser(user)
	if err != nil {
		context.JSON(err.(http_error.HttpError).Status, err)
	} else {
		context.JSON(http.StatusCreated, createdUser)
	}
}

func DeleteUser(context *gin.Context) {
	var user model.User
	err := context.ShouldBind(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	deletedUser, err := userRepository.DeleteUser(user.Username, user.Password)
	if err != nil {
		context.JSON(err.(http_error.HttpError).Status, err)
	} else {
		context.JSON(http.StatusOK, deletedUser)
	}
}

func UpdateUser(context *gin.Context) {
	var user model.User
	err := context.ShouldBind(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	claims, err := auth.GetToken(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	username := fmt.Sprint(claims["username"])
	updatedUser, err := userRepository.UpdateUser(username, user)
	if err != nil {
		context.JSON(err.(http_error.HttpError).Status, err)
		return
	} else {
		context.JSON(http.StatusOK, updatedUser)
		return
	}
}

func GetUser(context *gin.Context) {
	claims, err := auth.GetToken(context)
	if err != nil {
		context.JSON(401, http_error.Unauthorized("user", ""))
		return
	}
	username := fmt.Sprint(claims["username"])
	user, err := repository.UserRepositoryMongo.GetUser(repository.UserRepositoryMongo{}, username)
	if err != nil {
		context.JSON(err.(http_error.HttpError).Status, err)
	} else {
		context.JSON(http.StatusOK, user)
	}
}
