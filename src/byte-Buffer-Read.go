package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	Read()
	ReadByte()
	ReadRune()
}

func Read() {
	bufs := bytes.NewBufferString("Learning swift.")
	fmt.Println(bufs.String())

	//声明一个空的slice,容量为8
	l := make([]byte, 8)
	//把bufs的内容读入到l内,因为l容量为8,所以只读了8个过来
	bufs.Read(l)
	fmt.Println("::bufs缓冲器内容::")
	fmt.Println(bufs.String())
	//空的l被写入了8个字符,所以为 Learning
	fmt.Println("::l的slice内容::")
	fmt.Println(string(l))
	//把bufs的内容读入到l内,原来的l的内容被覆盖了
	bufs.Read(l)
	fmt.Println("::bufs缓冲器被第二次读取后剩余的内容::")
	fmt.Println(bufs.String())
	fmt.Println("::l的slice内容被覆盖,由于bufs只有7个了,因此最后一个g被留下来了::")
	fmt.Println(string(l))

}

// =======Read=======
// Learning swift.
// ::bufs缓冲器内容::
// swift.
// ::l的slice内容::
// Learning
// ::bufs缓冲器被第二次读取后剩余的内容::
// ::l的slice内容被覆盖::
// swift.g

func ReadByte() {
	bufs := bytes.NewBufferString("Learning swift.")
	fmt.Println(bufs.String())
	//读取第一个byte,赋值给b
	b, _ := bufs.ReadByte()
	fmt.Println(bufs.String())
	fmt.Println(string(b))
}

// =======ReadByte===
// Learning swift.
// earning swift.
// L

func ReadRune() {
	bufs := bytes.NewBufferString("学swift.")
	fmt.Println(bufs.String())

	//读取第一个rune,赋值给r
	r, z, _ := bufs.ReadRune()
	//打印中文"学",缓冲器头部第一个被拿走
	fmt.Println(bufs.String())
	//打印"学","学"作为utf8储存占3个byte
	fmt.Println("r=", string(r), ",z=", z)

}

//ReadBytes需要一个byte作为分隔符，读的时候从缓冲器里找第一个出现的分隔符（delim），找到后，把从缓冲器头部开始到分隔符之间的所有byte进行返回，作为byte类型的slice，返回后，缓冲器也会空掉一部分

func ReadBytes() {
	bufs := bytes.NewBufferString("现在开始 Learning swift.")
	fmt.Println(bufs.String())

	var delim byte = 'L'
	line, _ := bufs.ReadBytes(delim)
	fmt.Println(bufs.String())
	fmt.Println(string(line))
}

// =======ReadBytes==
// 现在开始 Learning swift.
// earning swift.
// 现在开始 L

//从一个实现io.Reader接口的r，把r里的内容读到缓冲器里，n返回读的数量
func ReadFrom() {
	//test.txt 内容是 "未来"
	file, _ := os.Open("learngo/bytes/text.txt")
	buf := bytes.NewBufferString("Learning swift.")
	buf.ReadFrom(file) //将text.txt内容追加到缓冲器的尾部
	fmt.Println(buf.String())
}

// =======ReadFrom===
// Learning swift.未来

func Reset() {
	bufs := bytes.NewBufferString("现在开始 Learning swift.")
	fmt.Println(bufs.String())

	bufs.Reset()
	fmt.Println("::已经清空了bufs的缓冲内容::")
	fmt.Println(bufs.String())
}

// =======Reset======
// 现在开始 Learning swift.
// ::已经清空了bufs的缓冲内容::

//string
//将未读取的数据返回成 string
