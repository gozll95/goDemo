//
在go当中string底层是用byte数组存的，并且是不可以改变的
例如：
s="Go编程"
fmt.Println(len(s))
输出结果应该是8因为中文字符是用3个字节存的
len(string(rune('编')))的结果是3
如果想获得我们想要的情况的话，需要先转换为rune切片
再使用内置的len函数

fmt.Println(len([]rune(s)))
结果就是4了

所以用string存储unicode的话，如果有中文，
按照下标是访问不到的，因为你只能得到一个byte.
要想访问中文的话，还是要用rune切片，这样就能按下表访问。


最直观的区别就是
rune 能操作 任何字符
byte 不支持中文的操作



 byte与rune类型有一个共性，即：它们都属于别名类型。byte是uint8的别名类型，而rune则是int32的别名类型。
  
    byte类型的值需用8个比特位表示，其表示法与uint8类型无异。因此我们就不再这里赘述了。我们下面重点说说rune类型。
  
    一个rune类型的值即可表示一个Unicode字符。Unicode是一个可以表示世界范围内的绝大部分字符的编码规范。关于它的详细信息，大家可以参看其官网（http://unicode.org/）上的文档，或在Google上搜索。用于代表Unicode字符的编码值也被称为Unicode代码点。一个Unicode代码点通常由“U+”和一个以十六进制表示法表示的整数表示。例如，英文字母“A”的Unicode代码点为“U+0041”。

    rune类型的值需要由单引号“'”包裹。例如，'A'或'郝'。这种表示方法一目了然。不过，我们还可以用另外几种形式表示rune类型值。请看下表。





   package main

import ( 
    "fmt" 
)

func main() {
    // 声明一个rune类型变量并赋值
    var char1 rune = '赞' 
    
    // 这里用到了字符串格式化函数。其中，%c用于显示rune类型值代表的字符。
    fmt.Printf("字符 '%c' 的Unicode代码点是 %s。\n", char1, (""))
}