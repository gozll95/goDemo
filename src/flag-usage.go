package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	dnssec := flag.Bool("dnssec", false, "Request DNSSEC records")
	//port := flag.String("port", "53", "Set the query port")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage:%s [OPTION] [name ...]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *dnssec {
		fmt.Println("dnssec")
	}
}
