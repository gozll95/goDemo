client 与 server 之间通过conn通信

client call 方法 里
conn <- request  这个conn是server提供的方法xx返回的,这个xx方法里对conn进行了读并且写结果的操作
这里再等结果
response<-conn


server对request的处理可以用interface来实现
