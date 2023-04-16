package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"oauth/internal/config"
	"oauth/internal/model"
	"strings"
	"time"
)

var secretKey = []byte(config.GetProperty("JWT_SECRET_KEY"))

func CreateToken(user model.User) (string, error) {
	claims := jwt.MapClaims{
		"exp":        float64(time.Now().Add(10 * time.Minute).Unix()),
		"authorized": true,
		"username":   user.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GetToken(context *gin.Context) (jwt.MapClaims, error) {
	token := context.Request.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	claims := &jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	return *claims, err
}
