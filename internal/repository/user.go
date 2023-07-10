package repository

import "oauth/internal/model"

type UserRepository interface {
	Repository[model.User]
}
