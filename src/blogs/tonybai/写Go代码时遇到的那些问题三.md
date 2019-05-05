# 一、语言篇

## 1.len(channel)的使用
* 当channel为unbuffered channel时，len(channel)总是返回0；
* 当channel为buffered channel时，len(channel)返回当前channel中尚未被读取的元素个数。

这样一来，所谓len(channel)中的channel就是针对buffered channel。len(channel)从语义上来说一般会被用来做“判满”、”判有”和”判空”逻辑：

```
// 判空

if len(channel) == 0 {
    // 这时：channel 空了 ?
}

// 判有

if len(channel) > 0 {
    // 这时：channel 有数据了 ?
}

// 判满
if len(channel) == cap(channel) {
    // 这时:   channel 满了 ?
}
```


大家看到了，我在上面代码中注释：“空了”、“有数据了”和“满了”的后面打上了问号！channel多用于多个goroutine间的通讯，一旦多个goroutine共同读写channel，len(channel)就会在多个goroutine间形成”竞态条件”，单存的依靠len(channel)来判断队列状态，不能保证在后续真正读写channel的时候channel状态是不变的。以判空为例：


                        |                                           |
                        |                                           |
                        | goroutine 1                               |  goroutine2
                        |                                           |

                  if len(channel)==0
                          No
                        |                                         read from channel
len(channel)=1          |                                           |
                        |                                           |
                        |                                           |
len(channel)=0      read from channel 




从上图可以看到，当goroutine1使用len(channel)判空后，便尝试从channel中读取数据。但在真正从Channel读数据前，另外一个goroutine2已经将数据读了出去，goroutine1后面的读取将阻塞在channel上，导致后面逻辑的失效。因此，为了不阻塞在channel上，常见的方法是将***“判空与读取”***放在一起做、将****”判满与写入”***一起做，通过select实现操作的“事务性”:

func readFromChan(ch<-chan int)(int,bool){
    select{
        case i:=<-ch:
            return i,true
        default:
            return 0,false // channel is empty
    }
}


func writeToChan(ch chan <-int,i int) bool{
    select{
        case ch<-i:
            return true
        default:
            return false // channel is full
    }
}

我们看到由于用到了Select-default的trick，当channel空的时候，readFromChan不会阻塞；当channel满的时候，writeToChan也不会阻塞。这种方法也许适合大多数的场合，但是这种方法有一个“问题”，那就是“改变了channel的状态”：读出了一个元素或写入了一个元素。有些时候，我们不想这么做，我们想在不改变channel状态下单纯地侦测channel状态！很遗憾，目前没有哪种方法可以适用于所有场合。但是在特定的场景下，我们可以用len(channel)实现。比如下面这个场景：


producer -----
             |
             |---->
producer ---------> chan ------------->controller     for{
             |---->   |                  | 创建            if len(channel) > 0{
             |        |                  |                      create a new consumer
producer -----        |--------------->consumer            }
                                                        ...
                                                        wait consumer
                                                    }   

这是一个***"多producer+1 consumer"***的场景。***controller是一个总控协程***,初始情况下,它会判断channel中是否有消息。如果有消息,它本身不消费"消息",而是创建一个consumer来消费消息,直到consumer因某种情况退出,控制权再回到controller,controller不会立即创建new consumer,而是等待channel下一次有消息时才创建。这样一个场景中,我们就可以使用len(channel)来判断是否有消息。


# 三、实践篇
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

而***deadline(以read deadline为例)机制***,则是***无论是否有数据ready以及数据读取活动***,都会在到达时间(deadline)后的再次read时返回timeout error,并且***后续的所有network read operation也都会返回timeout(如图中d)***,除非重新调用SetReadDeadline(time.Time{})取消Deadline或再次读取动作前重新设定deadline实现续时的目的。Go网络编程一般是"阻塞模型",那为什么还要有SetReadDeadline呢，这是因为有时候，我们要给调用者“感知”其他“异常情况”的机会，比如是否收到了main goroutine发送过来的退出通知信息。

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