package repository

import (
	"oauth/internal/model"
)

type ProfileRepository interface {
	Repository[model.Profile]
	GetProfileByName(name string) (model.Profile, error)
}
