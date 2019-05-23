//$GOROOT/src/time/sleep.go

func AfterFunc(d Duration, f func()) *Timer {
	t := &Timer{
		r: runtimeTimer{
			when: when(d),
			f:    goFunc,
			arg:  f,
		},
	}
	startTimer(&t.r)
	return t
}

func goFunc(arg interface{}, seq uintptr) {
	go arg.(func())()
}

/*
注意：从AfterFunc源码可以看到，外面传入的f参数并非直接赋值给了内部的f，而是作为wrapper function：goFunc的arg传入的。而goFunc则是启动了一个新的goroutine来执行那个外部传入的f。
这是因为timer expire对应的事件处理函数的执行是在go runtime内唯一的timer events maintenance goroutine: timerproc中。为了不block timerproc的执行，必须启动一个新的goroutine。
*/