package main

import (
    "bytes"
    "encoding/gob"
    "fmt"
    "log"
    "io/ioutil"
    "os"
)

type P struct {
    X, Y, Z int
    Name    string
}

type Q struct {
    X, Y *int32
    Name string
}

func main() {
    var fout bytes.Buffer        
    enc := gob.NewEncoder(&fout) //新的编码器
    
    err := enc.Encode(P{3, 4, 5, "Pythagoras"})//编码结构体和数据
    if err != nil {
        log.Fatal("encode error:", err)
    }
    ioutil.WriteFile("vt.dat", fout.Bytes(), 0644)//写入文件

    //读取并且解码
    fin,err := os.Open("vt.dat")//读取数据
    dec := gob.NewDecoder(fin)//解码数据
    var q Q
    err = dec.Decode(&q)
    if err != nil {
        log.Fatal("decode error:", err)
    }
    fmt.Printf("%q: {%d,%d}\n", q.Name, *q.X, *q.Y)
}