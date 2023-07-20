package service

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"oauth/internal/config"
	"oauth/internal/model"
	"oauth/internal/repository"
	"oauth/internal/utils/set"
	"strings"
	"time"
)

var secretKey = []byte(config.GetProperty("JWT_SECRET_KEY"))

type TokenService interface {
	CreateToken(user model.User) (string, error)
	GetToken(context *gin.Context) (jwt.MapClaims, error)
}

type TokenServiceImpl struct {
	profileRepositoryMongo repository.ProfileRepositoryMongo
}

func (a TokenServiceImpl) CreateToken(user model.User) (string, error) {
	permissions := set.StringSet{}
	for _, id := range user.ProfileIds {
		profileId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return "", err
		}
		profile, err := a.profileRepositoryMongo.Get(profileId)
		if err != nil {
			return "", err
		}
		for _, permission := range profile.Permissions {
			permissions.Add(permission)
		}
	}
	claims := jwt.MapClaims{
		"exp":         float64(time.Now().Add(10 * time.Minute).Unix()),
		"authorized":  true,
		"username":    user.Username,
		"permissions": permissions.Strings(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func (TokenServiceImpl) GetToken(context *gin.Context) (jwt.MapClaims, error) {
	token := context.Request.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	claims := &jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	return *claims, err
}

func NewTokenService() TokenServiceImpl {
	return TokenServiceImpl{
		profileRepositoryMongo: repository.NewProfileRepository(),
	}
}
