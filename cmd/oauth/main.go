package main

import (
	"oauth/api/v1/router"
	"oauth/internal/model"
	"oauth/locale"
)

func main() {
	locale.Setup()
	model.ConnectDatabase()
	router.CreateRouter()
}
