package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"oauth/api/v1/handlers"
	"oauth/api/v1/validators"
	"oauth/internal/config"
	"oauth/locale"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.CustomRecovery(ErrorHandler))
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("duplicated_profile_name", validators.DuplicatedName)
		if err != nil {
			log.Fatal(err)
		}
	}
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
		profileRoutes := v1Routes.Group("/profile")
		{
			profileRoutes.POST("", handlers.CreateProfile)
			profileRoutes.PUT("", handlers.UpdateProfile)
			profileRoutes.GET("/:id", handlers.GetProfile)
			profileRoutes.DELETE("/:id", handlers.DeleteProfile)
		}
	}
	err := router.Run("localhost:" + config.GetProperty("SERVER_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	return router
}

func ErrorHandler(c *gin.Context, err any) {
	switch err.(type) {
	case validator.ValidationErrors:
		errorsAmount := len(err.(validator.ValidationErrors))
		fieldErrors := make([]gin.H, errorsAmount)
		for i, fieldError := range err.(validator.ValidationErrors) {
			fieldErrors[i] = gin.H{"namespace": fieldError.Namespace(), "field": fieldError.Field(), "tag": fieldError.Tag()}
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": fieldErrors})
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
	}
}
