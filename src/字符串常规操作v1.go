//https://my.oschina.net/ivandongqifan/blog/486354
package main

import s "strings"
import "fmt"

var p = fmt.Println

func main() {
	p("Contains: ", s.Contains("test", "es"))      //是否包含 true
	p("Count: ", s.Count("test", "t"))             //字符串出现字符的次数 2
	p("HasPrefix: ", s.HasPrefix("test", "te"))    //判断字符串首部 true
	p("HasSuffix: ", s.HasSuffix("test", "st"))    //判断字符串结尾 true
	p("Index: ", s.Index("test", "e"))             //查询字符串位置 1
	p("Join: ", s.Join([]string{"a", "b"}, "-"))   //字符串数组 连接 a-b
	p("Repeat: ", s.Repeat("a", 5))                //重复一个字符串 aaaaa
	p("Replace: ", s.Replace("foo", "o", "0", -1)) //字符串替换 指定起始位置为小于0,则全部替换 f00
	p("Replace: ", s.Replace("foo", "o", "0", 1))  //字符串替换 指定起始位置1 f0o
	p("Split: ", s.Split("a-b-c-d-e", "-"))        //字符串切割 [a b c d e]
	p("ToLower: ", s.ToLower("TEST"))              //字符串 小写转换 test
	p("ToUpper: ", s.ToUpper("test"))              //字符串 大写转换 TEST
	p("Len: ", len("hello"))                       //字符串长度
	p("Char:", "hello"[1])                         //标取字符串中的字符，类型为byte
}

//*********************修剪字符串***
    你可以使用 strings.TrimSpace(s) 来剔除字符串开头和结尾的空白符号；如果你想要剔除指定字符，则可以使用strings.Trim(s, "cut") 来将开头和结尾的 cut 去除掉。该函数的第二个参数可以包含任何字符，如果你只想剔除开头或者结尾的字符串，则可以使用 TrimLeft 或者 TrimRight 来实现
    *********/



//strings.Join 性能好一点