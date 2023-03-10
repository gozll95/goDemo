# 前言:

书里把Map翻译为映射，我觉得太硬，所以这篇文章里，我还是用英文Map。

Map是一种数据结构，是一个集合，用于存储一系列无序的键值对。它基于键存储的，键就像一个索引一样，这也是Map强大的地方，可以快速快速检索数据，键指向与该键关联的值。

# 内部实现:

Map是基于***散列表***来实现，就是我们常说的Hash表，所以我们每次迭代Map的时候，打印的Key和Value是无序的，每次迭代的都不一样，即使我们按照一定的顺序存在也不行。

Map的***散列表包含一组桶，每次存储和查找键值对的时候，都要先选择一个桶***。如何选择桶呢？就是把指定的键传给散列函数，就可以索引到相应的桶了，进而找到对应的键值。

这种方式的好处在于，***存储的数据越多，索引分布越均匀，***所以我们访问键值对的速度也就越快，当然存储的细节还有很多，大家可以参考***Hash相关的知识***，这里我们只要记住Map存储的是***无序的键值对***集合。


# 在函数间传递Map
***传递本身***
函数间传递Map是不会拷贝一个该Map的副本的，也就是说如果一个Map传递给一个函数，该函数对这个Map做了修改，那么这个Map的所有引用，都会感知到这个修改。

```
func main() {
	dict := map[string]int{"王五": 60, "张三": 43}
	modify(dict)
	fmt.Println(dict["张三"])
}
func modify(dict map[string]int) {
	dict["张三"] = 10
}
```

上面这个例子输出的结果是10,也就是说已经被函数给修改了，可以证明传递的并不是一个Map的副本。这个特性和切片是类似的，这样就会更高，因为复制整个Map的代价太大了。

# 重要
***struct里如果有map,那么经过函数传递的时候,map会被修改***

非常明显的，age的值已经被改变。如果结构体里有引用类型的值，比如map，那么我们即使传递的是结构体的值副本，如果修改这个map的话，原结构的对应的map值也会被修改，这里不再写例子，大家可以验证下。