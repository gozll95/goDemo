package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	levelFlag = flag.Int("level", 0, "级别")
	bnFlag    int
)

func init() {
	flag.IntVar(&bnFlag, "bn", 3, "份数")
}

func main() {

	flag.Parse()
	count := len(os.Args)
	fmt.Println("参数总个数:", count)

	fmt.Println("参数详情:")
	for i := 0; i < count; i++ {
		fmt.Println(i, ":", os.Args[i])
	}

	fmt.Println("\n参数值:")
	fmt.Println("级别:", *levelFlag)
	fmt.Println("份数:", bnFlag)
}
