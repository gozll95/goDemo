// 是有问题的
// https://studygolang.com/articles/6176


package main

import (
    "fmt"
    "net"
    "os"
    "runtime"
    "strconv"
    "time"
)

func loop(startport, endport int, inport chan int) {
    for i := startport; i <= endport; i++ {
        inport <- i
    }
    close(inport)
}

func scanner(inport, outport, out chan int, ip string, endport int) {
    for {
        in, ok := <-inport
        if !ok {
            out <- 1
        }
        host := fmt.Sprintf("%s:%d", ip, in)
        tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
        if err != nil {
            outport <- 0
        } else {
            conn, err := net.DialTimeout("tcp", tcpAddr.String(), 200*time.Millisecond)
            if err != nil {
                outport <- 0
            } else {
                outport <- in
                conn.Close()
            }
        }
    }
}

func main() {
    runtime.GOMAXPROCS(4)
    inport := make(chan int)
    starttime := time.Now().Unix()
    outport := make(chan int)
    out := make(chan int)
    collect := []int{}
    if len(os.Args) != 4 {
        fmt.Println("Usage: scanner.exe IP startport endport")
        fmt.Println("Endport must be larger than startport")
        os.Exit(0)
    }
    ip := string(os.Args[1])
    if os.Args[3] < os.Args[2] {
        fmt.Println("Usage: scanner IP startport endport")
        fmt.Println("Endport must be larger than startport")
        os.Exit(0)
    }
    fmt.Printf("the ip is %s \r\n", ip)
    startport, _ := strconv.Atoi(os.Args[2])
    endport, _ := strconv.Atoi(os.Args[3])
    fmt.Printf("%d---------%d\r\n", startport, endport)
    go loop(startport, endport, inport)
    for {
        select {
        case <-out:
            fmt.Println(collect)
            endtime := time.Now().Unix()
            fmt.Println("The scan process has spent ", endtime-starttime)
            os.Exit(0)
        default:
            go scanner(inport, outport, out, ip, endport)
            port := <-outport
            if port != 0 {
                collect = append(collect, port)
            }
        }
    }
}