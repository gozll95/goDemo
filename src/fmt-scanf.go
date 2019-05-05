// Scanf 用于扫描 os.Stdin 中的数据，并根据 format 指定的格式
// 将扫描出的数据填写到参数列表 a 中
// 当 r 中的数据被全部扫描完毕或者扫描长度超出 format 指定的长度时
// 则停止扫描（换行符会被当作空格处理）
func Scanf(format string, a ...interface{}) (n int, err error)

func main() {
var name string
var age int
// 注意：这里必须传递指针 &name, &age
// 要获取的数据前后必须有空格
fmt.Scanf("%s %d", &name, &age)
// 在控制台输入：Golang 4
fmt.Printf("我的名字叫 %s ，今年 %d 岁", name, age)
// 我的名字叫 Golang ，今年 4 岁
}