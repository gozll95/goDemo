下面,我会带领你编写一个并发安全的整数数组类型,其中的无锁化方案会使用原子值实现。

在这里,依然先定义接口。

//并发安全的整数数组接口
type ConcurrentArray interface{
    //用于设置指定索引上的元素值
    Set(index uint32,elem int)(err error)
    //用于获取指定索引上的元素值
    Get(index uint32)(elem int,err error)
    // 用于获取数组的长度
    Len() uint32
}

https://book.douban.com/reading/43015239/