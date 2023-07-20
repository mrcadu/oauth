package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"oauth/api/http_error"
	"oauth/internal/model"
	"oauth/internal/repository"
	"oauth/internal/utils"
)

type UserHandler interface {
	CreateUser(context *gin.Context)
	DeleteUser(context *gin.Context)
	UpdateUser(context *gin.Context)
	GetUser(context *gin.Context)
}

type UserHandlerImpl struct {
	userRepository            repository.UserRepositoryMongo
	passwordEncryptionService utils.PasswordEncryptionImpl
}

func (u UserHandlerImpl) CreateUser(context *gin.Context) {
	var user model.User
	err := context.ShouldBind(&user)
	if err != nil {
		panic(err)
	}
	createdUser, err := u.userRepository.Create(user)
	if err != nil {
		panic(err)
	}
	createdUser.Password = ""
	context.JSON(http.StatusCreated, createdUser)
}

func (u UserHandlerImpl) DeleteUser(context *gin.Context) {
	id := context.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	deletedUser, err := u.userRepository.Delete(hex)
	if err != nil {
		panic(err)
	}
	context.JSON(http.StatusOK, deletedUser)
}

func (u UserHandlerImpl) UpdateUser(context *gin.Context) {
	var user model.UserChangePassword
	err := context.ShouldBind(&user)
	id := context.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	user.ID = hex
	recoveredUser, err := u.userRepository.Get(hex)
	if err != nil {
		panic(err)
	}
	isPasswordCorrect := u.passwordEncryptionService.CheckPasswordHash(user.Password, recoveredUser.Password)
	if !isPasswordCorrect {
		context.JSON(http.StatusUnauthorized, http_error.Unauthorized("user", recoveredUser.Username))
		return
	}
	userToBeUpdated := model.User{
		ID:         user.ID,
		Username:   user.Username,
		ProfileIds: user.ProfileIds,
	}
	if user.NewPassword != "" {
		userToBeUpdated.Password = user.NewPassword
	} else {
		userToBeUpdated.Password = user.Password
	}
	updatedUser, err := u.userRepository.Update(userToBeUpdated)
	if err != nil {
		panic(err)
	}
	updatedUser.Password = ""
	context.JSON(http.StatusOK, updatedUser)
}

func (u UserHandlerImpl) GetUser(context *gin.Context) {
	id := context.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		context.JSON(401, http_error.Unauthorized("user", ""))
		return
	}
	user, err := u.userRepository.Get(hex)
	if err != nil {
		context.JSON(err.(http_error.HttpError).Status, err)
	} else {
		user.Password = ""
		context.JSON(http.StatusOK, user)
	}
}

func NewUserHandler() UserHandlerImpl {
	return UserHandlerImpl{
		userRepository:            repository.NewUserRepository(),
		passwordEncryptionService: utils.PasswordEncryptionImpl{},
	}
}
