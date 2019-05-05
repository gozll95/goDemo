package main

import (
	"github.com/gorilla/sessions"
	"io"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func pageHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "get_name_session")
	name, ok := session.Values["name"].(string)
	session.Values["name"] = "Dean"
	session.Save(r, w)
	if ok {
		io.WriteString(w, "hello, "+string(name))
	}
}

func main() {
	http.HandleFunc("/", pageHandler)
	http.ListenAndServe("localhost:8000", nil)

}

//第一次是空
//第二次是hello Dean
