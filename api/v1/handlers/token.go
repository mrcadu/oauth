package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oauth/api/dto"
	"oauth/api/http_error"
	"oauth/internal/model"
	"oauth/internal/service/auth"
	"oauth/internal/service/password_encription"
)

func CreateToken(context *gin.Context) {
	var user model.User
	err := context.ShouldBind(&user)
	if err != nil {
		panic(err)
	}
	recoveredUser, err := userRepository.Get(user.ID)
	if err != nil {
		panic(err)
	}
	isPasswordCorrect := password_encription.CheckPasswordHash(user.Password, recoveredUser.Password)
	if !isPasswordCorrect {
		context.JSON(http.StatusUnauthorized, http_error.Unauthorized("user", recoveredUser.Username))
		return
	}
	token, err := auth.CreateToken(recoveredUser)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	context.JSON(http.StatusOK, dto.Auth{AccessToken: token})
}

func GetToken(context *gin.Context) {
	claims, err := auth.GetToken(context)
	if err != nil {
		context.JSON(401, http_error.Unauthorized("user", ""))
		return
	} else {
		context.JSON(200, claims)
		return
	}
}
