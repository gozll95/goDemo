棋牌类游戏
支持玩家进行下面的基本操作
- 登录游戏
- 查看房间列表
- 创建房间
- 加入房间
- 进行游戏
- 房间内聊天
- 游戏完成,退出房间
- 退出登录

因为goroutine可创建的个数不受系统资源的限制,原则上一台服务器可以创建上百万个goroutine,也就是可能可以支撑上百万个房间。当然,考虑到每个房间都需要耗费计算和内存资源，实际上不可能达到这么高的数字，但我们可以预测与使用系统线程和系统进程来对应一个房间相比，显然使用goroutine可以支持得最多很多。

接下来我们开始进行系统设计。先简化登录流程:用户只需要输入用户名就可以直接登录，无需验证过程。因此，对于用户管理，就是一个回话的管理流程。每个玩家对应的信息如下:
    - 用户唯一ID
    - 用户名，用于显示
    - 玩家等级
    - 经验值


总体上，我们可以将该实例划分成以下子系统:
- 玩家会话管理系统，用于管理每一位登录的玩家，包括玩家信息和玩家状态
- 大厅管理
- 房间管理系统，创建、管理和销毁每一个房间
- 游戏会话管理系统，管理房间内的所有动作，包括游戏进程和房间内聊天
- 聊天管理系统，用于接受管理员的广播信息

为了 免 出太多源代码，这里我们只实现了最基础的会话管理系统和聊 管理系统。因为 它们足以展示以下的技术问题:



1.简单IPC框架
简单IPC(进程间通信)框架的目的很简单，就是封装通信包的编码细节，让使用者可以专注于业务。我们这里使用channel作为模块之间的通信方式。虽然channel可以传递任何数据类型，甚至包含另外一个channel,但是为了让我们的架构更容易分拆，我们还是严格限制了只能用于传递JSON格式的字符串类型数据。这样如果之后像将这样的单进程示例修改为多进程的分布式架构，也不需要全部重写，只需要替换通信层即可。

2.中央服务器
中央服务器作为全局唯一实例，从原则上需要承担以下责任:
- 在线玩家的状态管理
- 服务器管理
- 聊天系统

我们想现在因为没有实现其他服务器，所以服务器管理这一块先空着，目前聊天系统也先只实现了广播，要实现房间内聊天或者私聊，其实都可以根据当前的实现进行扩展。


整个流程已经串联完 ，现在可以进行我们的这个半成品游戏服务器程序了:
    $ go run cgss.go
    Casual Game Server Solution
    A new session has been created successfully.
    Commands:
        login <username><level><exp>
        logout <username>
        send <message>
        listplayer
        quit(q)
        help(h)
    Command> login Tom 1 101
    Command> login Jerry 2 321
    Command> listplayer
    1 : &{Tom 1 101 0 <nil>}
    2 : &{Jerry 2 321 0 <nil>}
    Command> send Hello everybody.
    Tom received message: Hello everybody.
    Jerry received message: Hello everybody.
    Command> logout Tom
    Command> listplayer
    1 : &{Jerry 2 321 0 <nil>}
    Command> send Hello the people online.
    Jerry received message: Hello the people online.
    Command> logout Jerry
    Command> listplayer
    Failed.  No player online.
    Command> q
    $
