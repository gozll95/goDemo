/*
 * request里包含reply的channel
 */ 
package main 

import(
	"fmt"
)
type request struct{
	a,b int 
	replyc chan int 
}

func(r *request)String()string{
    return fmt.Sprintf("%d+%d=%d",r.a,r.b,<-r.replyc)
}

type binOp func(a,b int)int 

func run(op binOp,req *request){
	req.replyc<-op(req.a,req.b)
}

func server(op binOp,service<-chan *request,quit<-chan bool){
    for{
        select{
            case req:=<-service:
                go run(op,req) // don`t wait for it
            case <- quit:
                return
        }
    }
}



func startSerever(op binOp)(service chan<-*request,quit chan<-bool){
    service=make(chan *request)
    quit=make(chan bool)
    go server(op,service,quit)
    return service,quit
}



func main(){
	adderChan,quitChan:=startServer(
		func(a,b int)int{return a+b }
	)

	req1:=&request{7,8,make(chan int)}
	req2:=&request{17,18,make(chan int)}

	adderChan<-req1
	adderChan<-req2

	fmt.Println(req2,req1)

	quitChan<-true
}