func main() {
	html := `<body>
				<div id="div1">DIV1</div>
				<div class="name">DIV2</div>
				<span>SPAN</span>
			</body>
			`
	dom,err:=goquery.NewDocumentFromReader(strings.NewReader(html))
	if err!=nil{
		log.Fatalln(err)
	}
	dom.Find(".name").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
}

/*
class也是HTML中常用的属性，我们可以通过class选择器来快速的筛选需要的HTML元素，它的用法和ID选择器类似，为Find(".class")。
以上示例中，就筛选出来class为name的这个div元素。
*/

/*
Element Class 选择器
class选择器和id选择器一样，也可以结合着HTML元素使用，他们的语法也类似Find(element.class)，这样就可以筛选特定element、并且指定class的元素。
*/