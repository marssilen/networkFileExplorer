package main

import (
	"github.com/kataras/iris/v12"
)

func WebRoutes(app* iris.Application) {
	webController:= WebController{}
	app.Get("/", webController.Home)
}
