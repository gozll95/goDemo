package main

import (
	"fmt"
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

func (this UserController) List() {

}

func main() {
	userController := UserController{}

	//注册路由
	beego.Router("/user/list/:name/:age", &userController, "*:List")
	beego.Router("/user/list", &userController, "*:List")

	//创建url
	//{{urlfor "UserController.List" ":name" "astaxie" ":age" "25"}}
	url := userController.UrlFor("UserController.List", ":name", "astaxie", ":age", "25")
	//输出 /user/list/astaxie/25
	fmt.Println(url)

	//{{urlfor "UserController.List" "name" "astaxie" "age" "25"}}
	url = userController.UrlFor("UserController.List", "name", "astaxie", "age", "25")
	//输出 /user/list?name=astaxie&age=25
	fmt.Println(url)
}
