package main

import (
	"oauth/api/v1/router"
	"oauth/cmd/datasource"
	"oauth/locale"
)

func main() {
	locale.Setup()
	datasource.Setup()
	router.NewGin().CreateRouter()
}
