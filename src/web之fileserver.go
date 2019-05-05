package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func someHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Normal Handler")
}

func main() {
	r := mux.NewRouter()
	http.Handle("/doc/", http.StripPrefix("/doc/", http.FileServer(http.Dir("/Users/flower/workspace/learngo/src/myGoNotes"))))

	r.HandleFunc("/", someHandler)
	http.Handle("/", r)

	if err := http.ListenAndServe(":8090", nil); err != nil {
		panic(err)
	}
}
