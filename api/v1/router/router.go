package router

import (
	"github.com/gin-gonic/gin"
	"log"
	"oauth/api/v1/handlers"
	"oauth/internal/config"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()
	v1Routes := router.Group("/api/v1")
	{
		userRoutes := v1Routes.Group("/user")
		{
			userRoutes.POST("", handlers.CreateUser)
			userRoutes.DELETE("", handlers.DeleteUser)
			userRoutes.PUT("", handlers.UpdateUser)
			userRoutes.GET("", handlers.GetUser)
		}
		tokenRoutes := v1Routes.Group("/token")
		{
			tokenRoutes.POST("", handlers.CreateToken)
			tokenRoutes.GET("", handlers.GetToken)
		}
	}
	err := router.Run("localhost:" + config.GetProperty("SERVER_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	return router
}
