package config

import (
	"log"
	"oauth/internal/config/env"
)

var SysEnv env.Env = env.SystemEnv{}

func GetProperty(key string) string {
	err := SysEnv.LoadEnv(".env")
	if err != nil {
		log.Print("Error loading .env file")
	}
	return SysEnv.GetEnv(key)
}
