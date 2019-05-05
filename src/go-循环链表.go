package main

import (
	"fmt"
	"os"
)

type Node struct {
	data  int
	pNext *Node
}

func initList() *Node {
	pHead := new(Node)
	pHead.pNext = pHead
	return pHead //返回头指针
}

//创建尾指针的单循环链表
func createList(list **Node) {
	if !isempty(*list) {
		cleanList(*list)
	}
	var val int
	p, q := *list, *list
	fmt.Println("请输入结点数据,输入0结束输入")
	fmt.Scanf("%d", &val)
	for val != 0 {
		pnew := new(Node)
		fmt.Printf("pnew: point is %p, data is %v\t,pnew next is %v\n", pnew, pnew.data, pnew.pNext)
		fmt.Println("========================")
		pnew.data = val
		pnew.pNext = p
		fmt.Printf("pnew point is %p\t,data is %v\t,pnew next is %v\n", pnew, pnew.data, pnew.pNext)
		fmt.Println("========================")
		q.pNext = pnew
		fmt.Printf("q point is %p\t,data is %v\t,q next is %v\n", q, q.data, q.pNext)
		fmt.Printf("pnew point is %p\t,data is %v\t,pnew next is %v\n", pnew, pnew.data, pnew.pNext)
		fmt.Println("========================")
		q = pnew

		fmt.Printf("q point is %p\t,data is %v\t,q next is %v\n", q, q.data, q.pNext)
		fmt.Printf("p point is %p\t,data is %v\t,q next is %v\n", p, p.data, p.pNext)
		fmt.Printf("pnew: point is %p, data is %v\t,pnew next is %v\n", pnew, pnew.data, pnew.pNext)
		fmt.Println("========================")

		fmt.Scanf("%d", &val)
	}
	*list = q
}

//清空循环链表
func cleanList(list *Node) {
	if isempty(list) {
		return
	}
	phead := list.pNext   //头结点
	p := list.pNext.pNext //第一个结点
	q := p
	for p != list.pNext {
		q = p.pNext
		p = nil
		p = q
	}
	phead.pNext = phead
}

//插入结点
func insertList(list **Node) {
	var index, val int
	fmt.Printf("请输入要插入的位置：（值范围：1-%d）\n", listLength(*list)+1)
	fmt.Scanf("%d", &index)
	if index < 1 || index > listLength(*list)+1 {
		fmt.Println("位置值越界")
		return
	}
	fmt.Println("请输入要插入的值：")
	fmt.Scanf("%d", &val)
	j := 1
	p, q := (*list).pNext, (*list).pNext //头结点
	for j < index {
		p = p.pNext
		j++
	}
	pnew := new(Node)
	pnew.data = val
	pnew.pNext = p.pNext
	p.pNext = pnew
	if pnew.pNext == q {
		*list = pnew
	}
}

func deleList(list **Node) {
	var index int
	fmt.Printf("请输入要删除的位置：（值范围：1-%d）\n", listLength(*list))
	fmt.Scanf("%d", &index)
	if index < 1 || index > listLength(*list) {
		fmt.Println("位置值越界")
		return
	}
	j := 1
	p, q := (*list).pNext, (*list).pNext //头结点
	//查找index-1结点
	for j < index {
		p = p.pNext
		j++
	}
	cur := p.pNext
	p.pNext = cur.pNext
	if p.pNext == q {
		*list = p
	}
	cur = nil
}

func locateList(list *Node) {
	fmt.Println("请输入要查找的值：")
	var val int
	fmt.Scanf("%d", &val)
	q := list.pNext.pNext //第一个结点
	var loc int = 0
	for q != list.pNext {
		loc++
		if q.data == val {
			break
		}
		q = q.pNext
	}
	if loc == 0 {
		fmt.Println("链表中未找到你要的值")
	} else {
		fmt.Printf("你查找的值的位置为：%d\n", loc)
	}
}

func traverse(list *Node) {
	if isempty(list) {
		fmt.Println("空链表")
		return
	}
	fmt.Println("链表内容如下：")
	p := list.pNext.pNext //第一个结点
	for p != list.pNext {
		fmt.Printf("%5d", p.data)
		p = p.pNext
	}
	fmt.Println()
}

func isempty(list *Node) bool {
	if list.pNext == list {
		return true
	} else {
		return false
	}
}

func listLength(list *Node) int {
	fmt.Println("list data is ", list.data)
	if isempty(list) {
		return 0
	}
	var len int = 0
	fmt.Println(" data is ", list.pNext.data)
	p := list.pNext.pNext //第一个结点
	fmt.Println("p data is ", p.data)
	for p != list.pNext {
		len++
		p = p.pNext
	}
	return len
}

func main() {
	list := initList()
	var flag int
	fmt.Println("1.初始化链表")
	fmt.Println("2.插入结点")
	fmt.Println("3.删除结点")
	fmt.Println("4.返回结点位置")
	fmt.Println("5.遍历链表")
	fmt.Println("0.退出")
	fmt.Println("请选择你的操作:")
	fmt.Scanf("%d", &flag)
	for flag != 0 {
		switch flag {
		case 1:
			createList(&list)
		case 2:
			insertList(&list)
		case 3:
			deleList(&list)
		case 4:
			locateList(list)
		case 5:
			traverse(list)
		case 0:
			os.Exit(0)
		default:
			fmt.Println("无效操作")
		}
		fmt.Println("请选择你的操作:")
		fmt.Scanf("%d", &flag)
	}
}
