起始代码:

func test(){
    i := 0
    for {
        select {
        case <-time.After(time.Second * time.Duration(2)):
            i++
            if i == 5{
                fmt.Println("break now")
                break 
            }
            fmt.Println("inside the select: ")
        }
        fmt.Println("inside the for: ")
    }
}




解决方法一：使用golang中break的特性，在外层for加一个标签

func test(){
	i := 0
	ForEnd:
	for {
		select {
		case <-time.After(time.Second * time.Duration(2)):
			i++
			if i == 5{
				fmt.Println("break now")
				break ForEnd
			}
			fmt.Println("inside the select: ")
		}
		fmt.Println("inside the for: ")
	}
}
解决方法二： 使用goto直接跳出循环

func test(){
	i := 0

	for {
		select {
		case <-time.After(time.Second * time.Duration(2)):
			i++
			if i == 5{
				fmt.Println("break now")
				goto ForEnd
			}
			fmt.Println("inside the select: ")
		}
		fmt.Println("inside the for: ")
	}
	ForEnd：
}
