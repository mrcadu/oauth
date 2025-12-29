package middleware

import (
	"net/http"
	"oauth/locale"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ErrorMiddleware(c *gin.Context, err any) {
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
			var rawDoc bson.M
			if err := bson.Unmarshal(writeError.Raw, &rawDoc); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}
			keyValue, ok := rawDoc["keyValue"].(bson.M)
			var duplicatedKey string
			if ok {
				for k := range keyValue {
					duplicatedKey = k
					break
				}
			}
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
