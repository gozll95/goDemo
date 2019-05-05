4、写入超时

如果非要给Write增加一个期限，那我们可以调用SetWriteDeadline方法。我们copy一份client5.go，形成client6.go，在client6.go的Write之前增加一行timeout设置代码：

```
conn.SetWriteDeadline(time.Now().Add(time.Microsecond * 10))
```

启动server6.go，启动client6.go，我们可以看到写入超时的情况下，Write的返回结果：
$go run client6.go
2015/11/17 15:26:34 begin dial...
2015/11/17 15:26:34 dial ok
2015/11/17 15:26:34 write 65536 bytes this time, 65536 bytes in total
... ...
2015/11/17 15:26:34 write 65536 bytes this time, 655360 bytes in total
2015/11/17 15:26:34 write 24108 bytes, error:write tcp 127.0.0.1:62325->127.0.0.1:8888: i/o timeout
2015/11/17 15:26:34 write 679468 bytes in total
可以看到在写入超时时，依旧存在部分数据写入的情况。