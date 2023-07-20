package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oauth/api/http_error"
	"oauth/internal/service"
	"strings"
)

var tokenService = service.NewTokenService()

func CheckPermission(permission string) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			c.JSON(http.StatusUnauthorized, http_error.Unauthorized("authorization", "token"))
			c.Abort()
			return
		}
		claims, err := tokenService.GetToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, http_error.Unauthorized("authorization", "token"))
			c.Abort()
			return
		}
		var permissions []interface{}
		if claims["permissions"] != nil {
			permissions = (claims["permissions"]).([]interface{})
		}
		hasUserPermission := false
		for _, userPermission := range permissions {
			if userPermission == permission {
				hasUserPermission = true
			}
		}
		if !hasUserPermission {
			c.JSON(http.StatusUnauthorized, http_error.Unauthorized("authorization", "token"))
			c.Abort()
			return
		}
	}
}
