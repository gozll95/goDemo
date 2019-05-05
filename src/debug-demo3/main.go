//go-debug-profile-optimization/step0/demo.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

var add = make(chan struct{})
var total = make(chan int)

func getVistors() int {
	return <-total
}

func addVistors() {
	add <- struct{}{}
}

func teller() {
	var visitors int
	for {
		select {
		case <-add:
			visitors += 1
		case total <- visitors:
		}
	}
}

func handleHi(w http.ResponseWriter, r *http.Request) {
	if match, _ := regexp.MatchString(`^\w*$`, r.FormValue("color")); !match {
		http.Error(w, "Optional color is invalid", http.StatusBadRequest)
		return
	}
	addVistors()
	visitNum := getVistors()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<h1 style='color: " + r.FormValue("color") +
		"'>Welcome!</h1>You are visitor number " + fmt.Sprint(visitNum) + "!"))
}
func main() {
	log.Printf("Starting on port 8080")
	go teller()
	http.HandleFunc("/hi", handleHi)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
