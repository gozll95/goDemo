package main


import(
	"fmt"
)

func EnterBookInsertPage() {
	
}

func Menu() {
	fmt.Println("1. 书籍录入")
	fmt.Println("2. 书籍查询")
	fmt.Println("3. 学生信息录入")
	fmt.Println("4. 借书")
	fmt.Println("5. 书籍借阅状态")
	fmt.Println("6. 我的借阅")
	fmt.Println("7. 退出")

	var sel int
	fmt.Scanf("%d", &sel)

	switch sel {
	case 1:
		EnterBookInsertPage()
	case 2:
		//EnterBookQuery()
	case 7:
		return
	}
}

func main() {
	Menu()
}