package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	//"strings"
)

// func sayhelloName(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
// 	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
// 	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
// 	fmt.Println("path", r.URL.Path)
// 	fmt.Println("scheme", r.URL.Scheme)
// 	fmt.Println(r.Form["url_long"])
// 	for k, v := range r.Form {
// 		fmt.Println("key:", k)
// 		fmt.Println("val:", strings.Join(v, ""))
// 	}
// 	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
// }

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}

func elsepage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("else.html")
		t.Execute(w, nil)
	} else {
	}
}

func selectpage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("select.html")
		t.Execute(w, nil)
	} else {
	}
}

func twoselectpage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("twoselect.html")
		t.Execute(w, nil)
	} else {
	}
}

func dongtaipage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("dongtai2.html")
		t.Execute(w, nil)
	} else {
	}
}

func modal(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("modal.html")
		t.Execute(w, nil)
	} else {
	}
}

func dynamic(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("dynamic.html")
		t.Execute(w, nil)
	} else {
	}
}

func sort(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("sortable3.html")
		t.Execute(w, nil)
	} else {
	}
}

func merge(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("mergetable.html")
		t.Execute(w, nil)
	} else {
	}
}

func timerange(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("timerange.html")
		t.Execute(w, nil)
	} else {
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("template")))

	//http.HandleFunc("/", sayhelloName) //设置访问的路由
	http.HandleFunc("/login", login)
	http.HandleFunc("/else", elsepage)           //设置访问的路由
	http.HandleFunc("/select", selectpage)       //设置访问的路由
	http.HandleFunc("/twoselect", twoselectpage) //二级联动
	http.HandleFunc("/dongtai", dongtaipage)     //二级联动
	http.HandleFunc("/modal", modal)             //二级联动
	http.HandleFunc("/dynamic", dynamic)
	http.HandleFunc("/sort", sort)
	http.HandleFunc("/merge", merge)         //合并单元格
	http.HandleFunc("/timerange", timerange) //timerange
	err := http.ListenAndServe(":8888", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
