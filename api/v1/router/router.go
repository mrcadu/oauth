package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"net/http"
	"oauth/api/v1/handlers"
	"oauth/internal/config"
	"oauth/locale"
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
	router.Use(gin.CustomRecovery(r.ErrorHandler))
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

func (r GinImpl) ErrorHandler(c *gin.Context, err any) {
	switch err.(type) {
	case validator.ValidationErrors:
		errorsAmount := len(err.(validator.ValidationErrors))
		fieldErrors := make([]gin.H, errorsAmount)
		for i, fieldError := range err.(validator.ValidationErrors) {
			fieldErrors[i] = gin.H{"namespace": fieldError.Namespace(), "field": fieldError.Field(), "tag": fieldError.Tag()}
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": fieldErrors})
	case mongo.WriteException:
		fieldErrors := make([]gin.H, len(err.(mongo.WriteException).WriteErrors))
		for i, writeError := range err.(mongo.WriteException).WriteErrors {
			doc, err := bsonx.ReadDoc(writeError.Raw)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.(error).Error()})
			}
			duplicatedKey := doc[doc.IndexOf("keyValue")].Value.Document()[0].Key
			fieldErrors[i] = gin.H{"message": locale.GetMessageLocaleFromRequest("Duplicated Key", c, map[string]string{
				"key": duplicatedKey,
			})}
		}
		status := http.StatusBadRequest
		if err.(mongo.WriteException).HasErrorCode(11000) {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"errors": fieldErrors})
	case error:
		if err.(error).Error() == mongo.ErrNoDocuments.Error() {
			id := c.Param("id")
			c.JSON(http.StatusNotFound, gin.H{"message": locale.GetMessageLocaleFromRequest("Not Found", c, map[string]string{
				"Name": id,
			})})
		} else if err.(error).Error() == primitive.ErrInvalidHex.(error).Error() {
			id := c.Param("id")
			c.JSON(http.StatusBadRequest, gin.H{"message": locale.GetMessageLocaleFromRequest("Invalid Hex", c, map[string]string{
				"Hex": id,
			})})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.(error).Error()})
		}
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.(error).Error()})
	}
}
func NewGin() GinImpl {
	return GinImpl{
		userHandler:    handlers.NewUserHandler(),
		profileHandler: handlers.NewProfileHandler(),
		tokenHandler:   handlers.NewTokenHandler(),
	}
}
