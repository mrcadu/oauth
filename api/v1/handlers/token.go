package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oauth/api/dto"
	"oauth/api/http_error"
	"oauth/internal/model"
	"oauth/internal/repository"
	"oauth/internal/service"
	"oauth/internal/utils"
)

type TokenHandler interface {
	CreateToken(context *gin.Context)
	GetToken(context *gin.Context)
}

type TokenHandlerImpl struct {
	authService               service.TokenService
	passwordEncryptionService utils.PasswordEncryption
	userRepository            repository.UserRepository
}

func (t TokenHandlerImpl) CreateToken(context *gin.Context) {
	var user model.User
	err := context.ShouldBind(&user)
	if err != nil {
		panic(err)
	}
	recoveredUser, err := t.userRepository.GetByUsername(user.Username)
	if err != nil {
		panic(err)
	}
	isPasswordCorrect := t.passwordEncryptionService.CheckPasswordHash(user.Password, recoveredUser.Password)
	if !isPasswordCorrect {
		context.JSON(http.StatusUnauthorized, http_error.Unauthorized("user", recoveredUser.Username))
		return
	}
	token, err := t.authService.CreateToken(recoveredUser)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	context.JSON(http.StatusOK, dto.Auth{AccessToken: token})
}

func (t TokenHandlerImpl) GetToken(context *gin.Context) {
	claims, err := t.authService.GetToken(context)
	if err != nil {
		context.JSON(401, http_error.Unauthorized("user", ""))
		return
	} else {
		context.JSON(200, claims)
		return
	}
}

func NewTokenHandler() TokenHandlerImpl {
	return TokenHandlerImpl{
		authService:               service.NewTokenService(),
		passwordEncryptionService: utils.PasswordEncryptionImpl{},
		userRepository:            repository.NewUserRepository(),
	}
}
