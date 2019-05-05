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
	//与http.ServerMux不同的是mux.Route是完全的正则匹配，设置路由路径/index/,如果访问路径/inde/hello会返回404
	//设置路由路径为/index/访问路径/index也会是报404的，需要设置r.StrictSlash(true),/index/与/index才能匹配

	r.HandleFunc("/index/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root path"))
	})

	//mux.Vars(r)会返回该请求所解析的所有参数（map[string]string）
	//访问/hello/ghbai 会输出hello ghbai

	r.HandleFunc("/hello/{name:[a-zA-Z]+}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("hello %s", mux.Vars(r)["name"])))
	})

	//访问http://127.0.0.1:8091/hello/hzxxxxxxxxxx
	//hello hzxxxxxxxxxx

	http.Handle("/", r)

	if err := http.ListenAndServe(":8091", nil); err != nil {
		panic(err)
	}

}
