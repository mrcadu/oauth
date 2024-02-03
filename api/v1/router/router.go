package router

import (
	"github.com/gin-gonic/gin"
	"log"
	"oauth/api/v1/handlers"
	"oauth/api/v1/middleware"
	"oauth/internal/config"
)

type Gin interface {
	CreateRouter() *gin.Engine
	ErrorHandler(c *gin.Context, err any)
}

type GinImpl struct {
	userHandler    handlers.UserHandler
	profileHandler handlers.ProfileHandler
	tokenHandler   handlers.TokenHandler
}

func (r GinImpl) CreateRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.CustomRecovery(middleware.ErrorMiddleware))
	v1Routes := router.Group("/api/v1")
	{
		userRoutes := v1Routes.Group("/user")
		{
			userRoutes.POST("", handlers.CheckPermission("CREATE_USERS"), r.userHandler.CreateUser)
			userRoutes.DELETE("/:id", handlers.CheckPermission("DELETE_USERS"), r.userHandler.DeleteUser)
			userRoutes.PUT("/:id", handlers.CheckPermission("UPDATE_USERS"), r.userHandler.UpdateUser)
			userRoutes.GET("/:id", handlers.CheckPermission("GET_USERS"), r.userHandler.GetUser)
		}
		tokenRoutes := v1Routes.Group("/token")
		{
			tokenRoutes.POST("", r.tokenHandler.CreateToken)
			tokenRoutes.GET("", r.tokenHandler.GetToken)
		}
		profileRoutes := v1Routes.Group("/profile")
		{
			profileRoutes.POST("", handlers.CheckPermission("CREATE_PROFILES"), r.profileHandler.CreateProfile)
			profileRoutes.PUT("/:id", handlers.CheckPermission("UPDATE_PROFILES"), r.profileHandler.UpdateProfile)
			profileRoutes.GET("/:id", handlers.CheckPermission("GET_PROFILES"), r.profileHandler.GetProfile)
			profileRoutes.GET("/name/:name", handlers.CheckPermission("GET_PROFILES"), r.profileHandler.GetProfileByName)
			profileRoutes.DELETE("/:id", handlers.CheckPermission("DELETE_PROFILES"), r.profileHandler.DeleteProfile)
		}
	}
	err := router.Run(":" + config.GetProperty("SERVER_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	return router
}
func NewGin() GinImpl {
	return GinImpl{
		userHandler:    handlers.NewUserHandler(),
		profileHandler: handlers.NewProfileHandler(),
		tokenHandler:   handlers.NewTokenHandler(),
	}
}
