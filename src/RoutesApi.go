package main

import "github.com/kataras/iris/v12"

func ApiRoutes(app* iris.Application) {
	apiUserController := apiVPNController{25}
	apiUserControllerWithEncryption := apiUserControllerWithEncryption{25}

	v1 := app.Party("api/v1/")
	{
		v1.Any("/servers",apiUserController.servers)
		v1.Any("/status",apiUserController.status)
		//v1.Any("/get",apiUserController.getServersFiles)
	}
	v2 := app.Party("api/v2/")
	{
		v2.Get("/servers",apiUserControllerWithEncryption.servers)
		//v2.Any("/status",apiUserControllerWithEncryption.status)
		v2.Any("/read",apiUserControllerWithEncryption.readFromFile)
	}
}
