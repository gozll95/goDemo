package main 


import(
    "path"
    "fmt"
)

func main(){
    m:="http://"
    a:="xxx"
    b:="add"
    fmt.Println(path.Join(m,a, b))
}