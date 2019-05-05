package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Static("/assets", "./views")
	router.StaticFS("/more_static", http.Dir("views"))
	router.StaticFile("/favicon.ico", "./views/form.html")

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
