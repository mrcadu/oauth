package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"oauth/internal/model"
	"oauth/internal/repository"
)

type ProfileHandler interface {
	CreateProfile(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)
	GetProfile(ctx *gin.Context)
	DeleteProfile(ctx *gin.Context)
	GetProfileByName(ctx *gin.Context)
}
type ProfileHandlerImpl struct {
	profileRepository repository.ProfileRepository
}

func (p ProfileHandlerImpl) CreateProfile(ctx *gin.Context) {
	var profile model.Profile
	err := ctx.ShouldBind(&profile)
	if err != nil {
		panic(err)
	}
	createdProfile, err := p.profileRepository.Create(profile)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusCreated, createdProfile)
}

func (p ProfileHandlerImpl) UpdateProfile(ctx *gin.Context) {
	var profile model.Profile
	err := ctx.ShouldBind(&profile)
	id := ctx.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	profile.ID = hex
	if err != nil {
		panic(err)
	}
	updatedProfile, err := p.profileRepository.Update(profile)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, updatedProfile)
}

func (p ProfileHandlerImpl) GetProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	profile, err := p.profileRepository.Get(hex)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, profile)
}

func (p ProfileHandlerImpl) GetProfileByName(ctx *gin.Context) {
	name := ctx.Param("name")
	profile, err := p.profileRepository.GetProfileByName(name)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, profile)
}

func (p ProfileHandlerImpl) DeleteProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	_, err = p.profileRepository.Delete(hex)
	if err != nil {
		panic(err)
	}
	ctx.Status(http.StatusNoContent)
}
func NewProfileHandler() ProfileHandlerImpl {
	return ProfileHandlerImpl{
		profileRepository: repository.NewProfileRepository(),
	}
}
