package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
)

func main() {
	traceCtx := httptrace.WithClientTrace(context.Background(), &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			fmt.Printf("Prepare to get a connection for %s.\n", hostPort)
		},
		GotConn: func(info httptrace.GotConnInfo) {
			fmt.Printf("Got a connection: reused: %v, from the idle pool: %v.\n",
				info.Reused, info.WasIdle)
		},
		PutIdleConn: func(err error) {
			if err == nil {
				fmt.Println("Put a connection to the idle pool: ok.")
			} else {
				fmt.Println("Put a connection to the idle pool:", err.Error())
			}
		},
		ConnectStart: func(network, addr string) {
			fmt.Printf("Dialing... (%s:%s).\n", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			if err == nil {
				fmt.Printf("Dial is done. (%s:%s)\n", network, addr)
			} else {
				fmt.Printf("Dial is done with error: %s. (%s:%s)\n", err, network, addr)
			}
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			if info.Err == nil {
				fmt.Println("Wrote a request: ok.")
			} else {
				fmt.Println("Wrote a request:", info.Err.Error())
			}
		},
		GotFirstResponseByte: func() {
			fmt.Println("Got the first response byte.")
		},
	})
	req, err := http.NewRequest("GET", "http://www.golang.org/", nil)
	if err != nil {
		log.Fatal("Fatal error:", err)
	}
	req = req.WithContext(traceCtx)
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Request error: %v\n", err)
		os.Exit(1)
	}
}

/*
package httptrace_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
)

func Example() {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	_, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}
}
*/
