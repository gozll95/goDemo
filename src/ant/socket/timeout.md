## 1.说说网络数据读取timeout的处理-以SetReadDeadline为例

Go语言天生适合于网络编程,但网络编程的复杂性也是有目共睹的、要写出稳定、高效的网络端程序,需要考虑的因素有很多。
比如其中之一:
- 从socket读取数据超时的问题。

Go语言标准网络库并没有实现epoll实现的那样"idle timeout",而是提供了Deadline机制,我们用一副图来对比一下两个机制的不同:


```
        |---------------timeout-------------------------------|
read start:wait for data                                      |
        |                                                     |
(a)----------------------------------------------------------------------------->t
        |                                                     |
        Read                                               idle timeout



                                          |--------------timeout-------------| 
read start:wait for data    data ready    |                                  |
        |                        |        |                                  |
(b)----------------------------------------------------------------------------->t
        |                        |        |                                  |
      Read                      Read  read again:wait for data             idle timeout                        



       |----------------timeout-----------------|
       |                                        |    
       |                                        |    
  first,setreaddeadline                         |
  read start:wait for data                      |
       |                                        |
(c)-------------------------------------------------------------------------------->t
       |                                        |
     Read                                     deadline





       |----------------timeout-----------------|
       |                                        |    
       |                                        |    
  first,setreaddeadline                         |
  read start:wait for data    data ready        |             data ready
       |                        |               |               |
(d)-------------------------------------------------------------------------------->t
       |                        |     |         |               |
     Read                     Read  read again: deadline       Read return timeout
                                    wait for data         
```

看上图a)和b)展示了"idle timeout"机制,所谓idle timeout就是指这个timeout是真正在没有data ready的情况的timeout(如图a),如果有数据ready可读(如图b),那么timeout机制暂停,直到数据读完后,再次进入数据等待的时候,idle timeout再次启动。

而***deadline(以read deadline为例)机制***,则是***无论是否有数据ready以及数据读取活动***,都会在到达时间(deadline)后的再次read时返回timeout error,并且***后续的所有network read operation也都会返回timeout(如图中d)***,除非重新`调用SetReadDeadline(time.Time{})取消Deadline`或`再次读取动作前重新设定deadline实现续时`的目的。Go网络编程一般是"阻塞模型",那为什么还要有SetReadDeadline呢，这是因为有时候，我们要给调用者“感知”其他“异常情况”的机会，比如是否收到了main goroutine发送过来的退出通知信息。

Deadline机制在使用起来很容易出错，这里列举两个曾经遇到的出错状况：

a) 以为SetReadDeadline后，后续每次Read都可能实现idle timeout

读取一个完整业务包的流程

-----------------------------------
defer SetReadDeadline(time.Time{})

SetReadDeadline

Read A                                  一个完整业务包由A、B、C组成

Read B

Read C
----------------------------------

在上图中,我们看到这个流程是读取一个完整业务包的过程,业务包的读取使用了三次Read调用,但是只在第一次Read前调用了SetReadDeadline。这种使用方式仅仅在Read A时实现了足额的"idle timeout",且仅当A数据始终未ready时会timeout;一旦A数据read并已经被Read,当Read B和Read C时,如果还期望足额的"idle timeout"那就误解了SetReadDeadline的真正含义了。因此要想在每次Read时实现"足额的idle timeout",需要在每次Read前都重新设定deadline。


b)一个完整"业务包"分多次读取的异常情况处理

读取一个完整业务包的流程

----------------------------------------
defer SetReadDeadline(time.Time{})

SetReadDeadline

Read A

SetReadDeadline                         一个完整业务包由A、B、C组成

Read B

SetReadDeadline

Read C
---------------------------------------

在这幅图中，每个Read前都重新设定了deadline，那么这样就一定ok了么？对于在一个过程中读取一个“完整业务包”的业务逻辑来说，我们还要考虑对每次读取异常情况的处理，尤其是timeout发生。在该例子中，有三个Read位置需要考虑异常处理。
如果Read A始终没有读到数据，deadline到期，返回timeout，这里是最容易处理的，因为此时前一个完整数据包已经被读完，新的完整数据包还没有到来，外层控制逻辑收到timeout后，重启再次启动该读流程即可。
如果Read B或Read C处没有读到数据，deadline到期，这时异常处理就棘手一些，因为一个完整数据包的部分数据（A）已经从流中被读出，剩余的数据并不是一个完整的业务数据包，不能简单地再在外层控制逻辑中重新启动该过程。我们要么在Read B或Read C处尝试多次重读，直到将完整数据包读取完整后返回；要么认为在B或C处出现timeout是不合理的，返回区别于A处的错误码给外层控制逻辑，让外层逻辑决定是否是连接存在异常。