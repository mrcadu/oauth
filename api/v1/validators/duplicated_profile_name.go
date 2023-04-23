package validators

import (
	"github.com/go-playground/validator/v10"
	"oauth/internal/repository"
)

var ProfileRepository = repository.ProfileRepositoryMongo{}

var DuplicatedName validator.Func = func(fl validator.FieldLevel) bool {
	name, _ := fl.Field().Interface().(string)
	_, err := ProfileRepository.GetProfileByName(name)
	return err != nil
}
