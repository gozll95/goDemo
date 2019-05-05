# 并发

虽然channel是golang一个处理并发很好的东西,但是并非所有场合都需要。比如标准库中就很少在API中使用channel的。
- 将使用channel的位置向上层移动。
- 可以使用回调函数。
- 不要混合使用mutex和channel。

# 什么时候发起goroutine
- 有一些库的New()会发起他们的goroutine,这是不好的。
- 标准库使用的是***Serve()函数***,以及对应的***Close()函数***。
- 将goroutine***向上推***。