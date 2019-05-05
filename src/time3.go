
Gaemma
 
首页
 
 
归档
 
 
标签
Golang 中时间与时区
发表于 2017-04-16   |   分类于 Go   |   0 Comments
Golang 中每个 time.Time 实例都有一个相关联的 time.Location，也就是相应的时区数据。如果没有指定 time.Location，当会以当前系统的时区为默认值。如果没有注意到这一点，在解析或者格式化时间时可能会遇到一些问题。

CN 在 +8000 时区，我一般使用 Asia/Shanghai 来表示。假设我们运行 Go 程序的环境在 UTC 时区，当前时间是北京时间 2017-04-16 21:05:08。

格式化时间
fmt.Println(time.Now().Format("2006-01-02 15:04:05")) 会输出 2017-04-16 13:05:08，相差了 8 个小时，因为默认是 UTC 时区。

当然，如果这样写 fmt.Println(time.Now().Format("2006-01-02 15:04:05 -0700 MST"))，你会看到带时区信息的时间： 2017-04-16 13:05:08 +0000 UTC。

这个时候，我们需要用到 time.Time 的 In 方法来设置时区。

1
2
loc, _ := time.LoadLocation("Asia/Shanghai")
fmt.Println(time.Now().In(loc).Format("2006-01-02 15:04:05"))
以上会输出 2017-04-16 21:05:08。

解析时间
格式化时间更多是为了展示时间，但是解析时间如果没有设置正确的时区，则可能会导致业务错误。

假如有一个这样的业务需求，在 +8000 时区每天早晨 8 点后开始执行一个任务。如果我们解析时间时简单使用 time.Parse：

date := time.Now().Format("2006-01-02")
expectedTime, _ := time.Parse("2006-01-02 15:04", date + " 08:00")
if expectedTime.After(time.Now()) {
    // do task.
}
这样的话任务会延后 8 个小时才运行。这个时候可以使用 time.ParseInLocation 来解析时间：

loc, _ := time.LoadLocation("Asia/Shanghai")
date := time.Now().Format("2006-01-02")
expectedTime, _ := time.ParseInLocation("2006-01-02 15:04", date + " 08:00", loc)
if expectedTime.After(time.Now()) {
    // do task.
}
以上就不会有问题了。

#时间 #时区
 在 golang 中使用 redis cluster  在 Debian Stretch 上安装 nginx 1.12.1 + LuaJIT 

文章目录  站点概览
1. 格式化时间
2. 解析时间
© 2018  gaemma
由 Hexo 强力驱动  主题 - NexT.Muse
