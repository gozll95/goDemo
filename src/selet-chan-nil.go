//https://studygolang.com/articles/1675
package main

func main() {
	maybe := func(flag bool, ch chan int) <-chan int {
		if flag {
			return ch
		}
		return nil
	}
	select {
	case <-maybe(cond_a, chan_a):
		// do something
	case <-maybe(cond_b, chan_b):
		// do something
	case <-chan_def:
		// do something
	}

	//这里实际是利用了nil channel永远阻塞的特性，但是如果我们创建一个channel，但是不向它写数据也不关闭它，而是只从它读数据，那么也是可以实现永远阻塞的。以下代码实现了同样的效果：

	// var _BLOCK = make(<-chan int)
	// maybe := func(flag bool, ch chan int) <-chan int {
	// 	if flag {
	// 		return ch
	// 	}
	// 	return _BLOCK
	// }
	// select {
	// case <-maybe(cond_a, chan_a):
	// 	// do something
	// case <-maybe(cond_b, chan_b):
	// 	// do something
	// case <-chan_def:
	// 	// do something
	// }

	//这也就意味着：对于实现一个"guarded selective wating"模式来说，nil channel的永久阻塞的特性并不是必须的，因为有其他替代实现方式。但是显然用nil channel更方便，也不需要额外浪费资源去创建一个用来永久阻塞的channel。
}

//go语言的channel有一个看上去很奇怪的特性，就是如果向一个为空值（nil）的channel写入或者读取数据，当前goroutine将永远阻塞
