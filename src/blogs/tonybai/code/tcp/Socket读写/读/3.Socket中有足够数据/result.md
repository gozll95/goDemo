3、Socket中有足够数据

如果socket中有数据，且长度大于等于一次Read操作所期望读出的数据长度，那么Read将会成功读出这部分数据并返回。这个情景是最符合我们对Read的期待的了：Read将用Socket中的数据将我们传入的slice填满后返回：n = 10, err = nil。
我们通过client2.go向Server2发送如下内容：abcdefghij12345，执行结果如下：
$go run client2.go abcdefghij12345
2015/11/17 13:38:00 begin dial...
2015/11/17 13:38:00 dial ok

$go run server2.go
2015/11/17 13:38:00 accept a new connection
2015/11/17 13:38:00 start to read from conn
2015/11/17 13:38:02 read 10 bytes, content is abcdefghij
2015/11/17 13:38:02 start to read from conn
2015/11/17 13:38:02 read 5 bytes, content is 12345
client端发送的内容长度为15个字节，Server端Read buffer的长度为10，因此Server Read第一次返回时只会读取10个字节；Socket中还剩余5个字节数据，Server再次Read时会把剩余数据读出（如：情形2）。