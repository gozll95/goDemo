字符串比较相等

func main() {
    fmt.Println("ADN" == "ADN")
    fmt.Println("ADN" == "adn")
    fmt.Println(strings.Compare("ADN", "ADN"))
    fmt.Println(strings.Compare("ADN", "adn"))
    fmt.Println(strings.EqualFold("ADN", "ADN"))
    fmt.Println(strings.EqualFold("ADN", "adn"))
}
Compare比"=="快，两种方法都区分大小写

EqualFold比较UTF-8编码在小写的条件下是否相等，不区分大小写。
