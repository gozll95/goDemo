Write操作存在写入部分数据的情况，比如上面例子中，当client端输出日志停留在“write 65536 bytes this time, 655360 bytes in total”时，我们杀掉server5，这时我们会看到client5输出以下日志：
...
2015/11/17 15:19:14 write 65536 bytes this time, 655360 bytes in total
2015/11/17 15:19:16 write 24108 bytes, error:write tcp 127.0.0.1:62245->127.0.0.1:8888: write: broken pipe
2015/11/17 15:19:16 write 679468 bytes in total
显然Write并非在655360这个地方阻塞的，而是后续又写入24108后发生了阻塞，server端socket关闭后，我们看到Wrote返回er != nil且n = 24108，程序需要对这部分写入的24108字节做特定处理。