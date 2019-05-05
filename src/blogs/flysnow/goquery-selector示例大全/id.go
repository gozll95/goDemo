/*
这个是使用频次最多的，类似于上面的例子，有两个div元素，其实我们只需要其中的一个，那么我们只需要给这个标记一个唯一的id即可，这样我们就可以使用id选择器，精确定位了。
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
			<div id="div1">DIV1</div>
			<div>DIV2</div>
			<span>SPAN</span>
		</body>
		`
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}
	dom.Find("#div1").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
}

/*
DIV1
*/
