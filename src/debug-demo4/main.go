//go-debug-profile-optimization/step0/demo.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

var vistors int
var c = make(chan struct{}, 1)

func handleHi(w http.ResponseWriter, r *http.Request) {
	c <- struct{}{}
	if match, _ := regexp.MatchString(`^\w*$`, r.FormValue("color")); !match {
		http.Error(w, "Optional color is invalid", http.StatusBadRequest)
		return
	}
	vistors++
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<h1 style='color: " + r.FormValue("color") +
		"'>Welcome!</h1>You are visitor number " + fmt.Sprint(vistors) + "!"))
	<-c
}
func main() {
	log.Printf("Starting on port 8080")
	http.HandleFunc("/hi", handleHi)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
