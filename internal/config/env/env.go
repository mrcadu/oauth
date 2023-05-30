package env

import (
	"github.com/joho/godotenv"
	"os"
)

type Env interface {
	LoadEnv(key string) error
	GetEnv(key string) string
}

type SystemEnv struct {
}

func (s SystemEnv) LoadEnv(key string) error {
	return godotenv.Load(key)
}

func (s SystemEnv) GetEnv(key string) string {
	return os.Getenv(key)
}
