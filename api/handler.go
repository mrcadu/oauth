package handler

import (
	"net/http"
	"oauth/api/v1/router"
	"oauth/cmd/datasource"
	"oauth/locale"
)

var initialized bool

func Handler(w http.ResponseWriter, r *http.Request) {
	if !initialized {
		locale.Setup()
		datasource.Setup()
		initialized = true
	}

	ginRouter := router.NewGin()
	engine := ginRouter.CreateRouter()
	engine.ServeHTTP(w, r)
}
