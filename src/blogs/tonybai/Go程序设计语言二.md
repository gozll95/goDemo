# 切片容量

package main

import (
	"fmt"
)

func main() {
	var ar = [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	var a = ar[5:7]

	fmt.Println(a)
	fmt.Println(len(a), cap(a))

	a = a[0:4]

	fmt.Println(a)

	fmt.Println(len(a), cap(a))
}

[5 6]
2 5
[5 6 7 8]
4 5

# 调整切片大小
切片可被当作可增长的数组用。使用make分配一个切片，并指定其长度和容量。当要增长时，我们可以做重新切片：
 
var sl = make([]int, 0, 100) // 长度 0, 容量 100
func appendToSlice(i int, sl []int) []int {
    if len(sl) == cap(sl) { error(…) }
    n := len(sl)
    sl = sl[0:n+1] // 长度增加1
    sl[n] = i
    return sl
}
 
因此，sl的长度总是元素的个数，但其容量可根据需要增加。
 
这种手法代价很小，并且是Go语言中的惯用法。

# 切片使用的代价很小
你可以根据需要自由地分配和调整切片大小,它们的传递仅需要很小的代价;不必分配。

记住它们是引用,因此下层的存储可以被修改。

例如,I/O使用切片,而不是计数:

func Read(fd int,b []byte)int

var buffer [100]byte 

for i:=0;i<100;i++{
    // 每次向Buffer中填充一个字节
    Read(fd,buffer[i:i+1]) // no allocation here
}

拆分一个Buffer:
header,data:=buf[:n],buf[n:]

字符串也可以被切片,而且效率相似。