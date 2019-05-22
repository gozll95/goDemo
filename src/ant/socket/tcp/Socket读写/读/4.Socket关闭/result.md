如果***client端主动关闭了socket***，那么Server的Read将会读到什么呢？这里分为“有数据关闭”和“无数据关闭”。
“有数据关闭”是指在client关闭时，socket中还有server端未读取的数据，我们在go-tcpsock/read_write/client3.go和server3.go中模拟这种情况：
$go run client3.go hello
2015/11/17 13:50:57 begin dial...
2015/11/17 13:50:57 dial ok

$go run server3.go
2015/11/17 13:50:57 accept a new connection
2015/11/17 13:51:07 start to read from conn
2015/11/17 13:51:07 read 5 bytes, content is hello
2015/11/17 13:51:17 start to read from conn
2015/11/17 13:51:17 conn read error: EOF
从输出结果来看，当client端close socket退出后，server3依旧没有开始Read，10s后第一次Read成功读出了5个字节的数据，当第二次Read时，由于client端 socket关闭，Read返回EOF error。
通过上面这个例子，我们也可以猜测出“无数据关闭”情形下的结果，那就是Read直接返回***EOF error***。