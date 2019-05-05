struct的内存布局: struct中的所有字段所在内存是连续的


type Rect1 struct{ Min, Max Point}

r1:=Rect1{Point{10,20},Point{50,60}}
|10|20|50|60| Rect1

type Rect2 struct{Min,Max *Point}
r2:=Rect2{&Point{10,20},&Point{50,60}}
|10|20|       |50|60| 


#链表定义

type Student struct{
	Name string
	Next *Student
}

每个节点包含下一个节点的地址,这样就把所有的节点串起来了。通常把链表中的第一个节点叫做链表头


# 双链表定义

type Student struct{
	Name string
	Next *Student
	Prev *Student
}

如果有两个指针分别指向前 个节点和后 个节点，我们叫做双链表


# 二叉树
type Student struct{
	Name string
	left *Student
	right *Student
}

如果每个节点有两个指针分别用来指向左子树和右子树,我们把这样的
结构叫做二叉树

map底层就是红黑树(二叉树的一种)


