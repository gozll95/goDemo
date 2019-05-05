package main

import (
	"net/http"

	"video-crawler/crawler-concurrent-show/frontend/controller"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("crawler-concurrent-show/frontend/view")))
	http.Handle("/search", controller.CreateSearchReusltHandler("crawler-concurrent-show/frontend/view/template.html"))
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
