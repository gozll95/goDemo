package router

import (
	"fmt"

	"go-kit-gin/controller"
	"go-kit-gin/service"
)

func initUserRouter(svc service.AppService) {
	userGroup := router.Group("/users")
	fmt.Println(userGroup)
	userGroup.POST("/create", transport.CreateAccount(svc))
	userGroup.POST("/login", transport.Login(svc))
	userGroup.GET("/query", transport.Account(svc))
}
