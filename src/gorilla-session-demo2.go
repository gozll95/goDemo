package main

import (
	"encoding/gob"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
)

type MPerson struct {
	FirstName string
	LastName  string
	Email     string
	Age       int
}

type M map[string]interface{}

func init() {
	fmt.Println("init")
	gob.Register(&MPerson{})
	gob.Register(&M{})
}

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func SetMHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "m")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ms := MPerson{FirstName: "king", LastName: "eastersun"}
	session.Values["ms"] = ms
	m := M{}
	m = make(map[string]interface{})
	m["name"] = "kingeasternsun"
	session.Values["m"] = m
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func GetMHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "m")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ms := session.Values["ms"]
	// var person = &MPerson{}
	if person, ok := ms.(*MPerson); ok {
		fmt.Println(person)
	} else {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := session.Values["m"]
	// var mmap = M{}
	if mmap, ok := m.(*M); ok {
		fmt.Println(mmap)
	} else {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func main() {
	// fmt.Println(string(securecookie.GenerateRandomKey(32)))
	// fmt.Println(string(securecookie.GenerateRandomKey(32)))
	routes := mux.NewRouter()
	routes.HandleFunc("/setm", SetMHandler)
	routes.HandleFunc("/getm", GetMHandler)

	http.Handle("/", routes)
	http.ListenAndServe(":8080", nil)
}
