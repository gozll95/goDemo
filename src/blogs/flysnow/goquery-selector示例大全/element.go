/*
这个比较简单，就是基于a,p等这些HTML的基本元素进行选择，这种直接使用Element名称作为选择器即可。比如dom.Find("div")。

*/
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	html := `<body>
			<div>DIV1</div>
			<div>DIV2</div>
			<span>SPAN</span>
		</body>
		`
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}
	dom.Find("div").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
}

/*
DIV1
DIV2
*/

//以上示例，可以把div元素筛选出来，而body,span并不会被筛选。
