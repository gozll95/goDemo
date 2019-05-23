https://blog.gopheracademy.com/advent-2016/go-timers

go runtime 实际上仅仅是启动了一个单独的goroutine,运行timerproc函数,维护了一个"最小堆",定期wake up后,读取堆顶的timer,
执行timer对应的f函数,并移除该timer element。
创建一个Timer实则就是在这个最小堆中添加一个element;
Stop一个timer,则是从堆中删除对应的element.


作为Timer的使用者，我们要做的就是尽量减少在使用Timer时对最小堆管理goroutine和GC的压力即可，即：及时调用timer的Stop方法从最小堆删除timer element(如果timer 没有expire)以及reuse active timer。


https://ggaaooppeenngg.github.io/zh-CN/2016/02/09/timer%E5%9C%A8go%E5%8F%AF%E4%BB%A5%E6%9C%89%E5%A4%9A%E7%B2%BE%E7%A1%AE/