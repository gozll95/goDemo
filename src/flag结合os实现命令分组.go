package main

import (
	"flag"
	"fmt"
	"os"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n%s command [options]\n", os.Args[0], os.Args[0])
		fmt.Fprintln(os.Stderr, "commands:")
		fmt.Fprintln(os.Stderr, "  register_pfdstg   register pfdstg")
		fmt.Fprintln(os.Stderr, "  pfdstg_writable   switch pfdstg readonly or writable")
	}
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(-1)
	}
	flag.Parse()
}

func main() {
	switch os.Args[1] {
	case "register_pfdstg":
		RegisterPFDStgCmd(os.Args[2:])
	case "pfdstg_writable":
		SetPFDStgWritableCmd(os.Args[2:])
	}

}

func RegisterPFDStgCmd(args []string) {
	var (
		conf string
	)
	fl := flag.NewFlagSet("RegisterPFDStg", flag.ExitOnError)
	fl.StringVar(&conf, "c", "pfdstg.json", "PFDStg config file")
	fl.Parse(args)

	fmt.Println(conf)

}

func SetPFDStgWritableCmd(args []string) {
	var (
		conf string
	)
	fl := flag.NewFlagSet("RegisterPFDStg", flag.ExitOnError)
	fl.StringVar(&conf, "d", "pfdstg.json", "PFDStg config file")
	fl.Parse(args)

	fmt.Println(conf)

}

/*
flower@:~/workspace/learngo/src/myGoNotes$ go run flag结合os.go register_pfdstg -h
Usage of RegisterPFDStg:
  -c string
    	PFDStg config file (default "pfdstg.json")
exit status 2
flower@:~/workspace/learngo/src/myGoNotes$ go run flag结合os.go pfdstg_writable -h
Usage of RegisterPFDStg:
  -d string
    	PFDStg config file (default "pfdstg.json")
exit status 2
*/
