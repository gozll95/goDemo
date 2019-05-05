package main

import "fmt"

//go:generate echo hello
//go:generate go run main.go
//go:generate  echo file=$GOFILE pkg=$GOPACKAGE
func main() {
	fmt.Println("main func")
}

/*
➜ /Users/flower/workspace/goDemo/src/test1 git:(master) ✗ >go generate
hello
main func
file=main.go pkg=main
*/

/*
介绍
go generate命令是go 1.4版本里面新添加的一个命令，当运行go generate时，它将扫描与当前包相关的源代码文件，找出所有包含"//go:generate"的特殊注释，提取并执行该特殊注释后面的命令，命令为可执行程序，形同shell下面执行。

有几点需要注意：

该特殊注释必须在.go源码文件中。
每个源码文件可以包含多个generate特殊注释时。
显示运行go generate命令时，才会执行特殊注释后面的命令。
命令串行执行的，如果出错，就终止后面的执行。
特殊注释必须以"//go:generate"开头，双斜线后面没有空格。
应用
在有些场景下，我们会使用go generate：

yacc：从 .y 文件生成 .go 文件。
protobufs：从 protocol buffer 定义文件（.proto）生成 .pb.go 文件。
Unicode：从 UnicodeData.txt 生成 Unicode 表。
HTML：将 HTML 文件嵌入到 go 源码 。
bindata：将形如 JPEG 这样的文件转成 go 代码中的字节数组。
再比如：

string方法：为类似枚举常量这样的类型生成String()方法。
宏：为既定的泛型包生成特定的实现，比如用于ints的sort.Ints。

*/

// https://www.jianshu.com/p/a866147021da
