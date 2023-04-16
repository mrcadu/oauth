package main

import (
	"oauth/api/v1/router"
	"oauth/internal/model"
)

func main() {
	model.ConnectDatabase()
	router.CreateRouter()
}
