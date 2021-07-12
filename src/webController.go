package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"log"
)

type WebController struct{

}
func (self*WebController) Home(ctx iris.Context){

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	ctx.View("index.html")
}

func (self *WebController) Login(ctx iris.Context) {
	ctx.View("auth/login.html")
}