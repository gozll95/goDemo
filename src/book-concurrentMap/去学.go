***对于unsafe.Pointer类型的值是可以实施原子操作的。***

atomic.Value

分段锁

单链表

包级私有


// findSegment 会根据给定参数寻找并返回对应散列段。
func (cmap *myConcurrentMap) findSegment(keyHash uint64) Segment {
	if cmap.concurrency == 1 {
		return cmap.segments[0]
	}
	var keyHash32 uint32
	if keyHash > math.MaxUint32 {
		keyHash32 = uint32(keyHash >> 32)
	} else {
		keyHash32 = uint32(keyHash)
	}
	return cmap.segments[int(keyHash32>>16)%(cmap.concurrency-1)]
}

可以看到,该算法的核心思想就是使用高位的几个字节来决定散列段的索引。这样可以让键-元素对在segments中分布得更广、更均匀一些。