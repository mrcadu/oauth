package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"oauth/internal/model"
	"oauth/internal/repository"
)

var profileRepository = repository.ProfileRepositoryMongo{}

func CreateProfile(ctx *gin.Context) {
	var profile model.Profile
	err := ctx.ShouldBind(&profile)
	if err != nil {
		panic(err)
	}
	createdProfile, err := profileRepository.CreateProfile(profile)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusCreated, createdProfile)
}

func UpdateProfile(ctx *gin.Context) {
	var profile model.Profile
	err := ctx.ShouldBind(&profile)
	id := ctx.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	profile.ID = hex
	if err != nil {
		panic(err)
	}
	updatedProfile, err := profileRepository.UpdateProfile(profile)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, updatedProfile)
}

func GetProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	profile, err := profileRepository.GetProfile(hex)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, profile)
}

func DeleteProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	_, err = profileRepository.DeleteProfile(hex)
	if err != nil {
		panic(err)
	}
	ctx.Status(http.StatusNoContent)
}
