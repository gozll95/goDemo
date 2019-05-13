package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type A struct {
	One int32
	Two int32
}

var a A

func main() {
	// i := uint16(1)
	// size := binary.Size(i)
	// fmt.Println(size)

	// a.One = int32(1)
	// a.Two = int32(2)

	// 将结构体序列化到一个buf中
	// 通过Size可以得到所需buffer的大小，通过Write可以将对象a的内容序列化到buffer中
	// 这里采用"小端序"的方式进行序列化(x86架构都是小端序,网络字节序都是大端序)
	// buf := new(bytes.Buffer)
	// fmt.Println("a's size is ", binary.Size(a))
	// binary.Write(buf, binary.LittleEndian, a)
	// fmt.Println("after write ，buf is:", buf.Bytes())

	// 从buf中反序列化回一个结构
	// var aa A
	// buf := new(bytes.Buffer)
	// binary.Write(buf, binary.LittleEndian, a)
	// binary.Read(buf, binary.LittleEndian, &aa)
	// fmt.Println("after aa is ", aa)

	// 将整数序列化到buf中，并从buf中反序列化出来
	/*
		我们可以通过Read/Write直接去读或者写一个uintx类型的变量来实现对xxx数的序列化和反序列化。
		由于在网络中，由于在网络中，对于×××数的序列化非常常用，因此系统库提供了type ByteOrder接口可以方便的对uint16/uint32/uint64进行序列化和反序列化：
	*/
	int16buf := new(bytes.Buffer)
	i := uint16(1)
	binary.Write(int16buf, binary.LittleEndian, i)
	fmt.Println("write buf is:", int16buf.Bytes())

	var int16buf2 [2]byte
	binary.LittleEndian.PutUint16(int16buf2[:], uint16(1))
	fmt.Println("put buffer is :", int16buf2[:])

	ii := binary.LittleEndian.Uint16(int16buf2[:])
	fmt.Println("Get buf is :", ii)

	// 一个实在的例子
	type Head struct {
		Cmd byte
		Version byte
		Magic   uint16
		Reserve byte
		HeadLen byte
		BodyLen uint16
	}
	
	func NewHead(buf []byte)*Head{
		head := new(Head)
	
		head.Cmd     = buf[0]
		head.Version = buf[1]
		head.Magic   = binary.BigEndian.Uint16(buf[2:4])
		head.Reserve = buf[4]
		head.HeadLen = buf[5]
		head.BodyLen = binary.BigEndian.Uint16(buf[6:8])
		return head
	}

}

/*
在写网络程序的时候，我们经常需要将结构体或者整数等数据类型序列化成【二进制的buffer串】。
或者从一个buffer中解析出来一个结构体出来。
最典型的就是在协议的header部分表征head length或者body length在拼包和拆包的过程中，
需要按照规定的整数类型进行解析，且涉及到【大小端序】的问题
*/

/*
通过Read接口可以将buf中得内容填充到data参数表示的数据结构中，通过Write接口可以将data参数里面包含的数据写入到buffer中。 变量BigEndian和LittleEndian是实现了ByteOrder接口的对象，通过接口中提供的方法可以直接将uintx类型序列化（uintx()）或者反序列化(putuintx())到buf中。
*/

