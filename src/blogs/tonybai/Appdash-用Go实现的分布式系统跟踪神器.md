中文版google dapper: http://bigbully.github.io/Dapper-translation/
wiki分布式系统: https://en.wikipedia.org/wiki/Distributed_computing

学会使用sourcegraph



分布式系统三要素:
- concurrency of components
- lack of a global clock
- independent failure of components


分布式系统的例子 
- SOA-based systems 
- massively multiplayer online games 
- peer-to-peer applications.

分布式系统传递信息: 
- pure HTTP
- RPC-like
- message queues


CAP理论


各种硬件和软件体系结构用于分布式计算。在较低级别上，需要将多个CPU与某种网络互连，而不管该网络是印刷在电路板上还是由松散耦合的设备和电缆组成。在更高层次上，有必要将这些CPU上运行的进程与某种通信系统互连。[ 引用需要 ]

分布式编程通常分为以下几个基本架构中的一种：客户端-服务器，三层，Ñ -tier，或对等网络 ; 或类别：松耦合，或紧耦合。[23]

- 客户端 - 服务器：智能客户端联系服务器获取数据然后格式化并显示给用户的体系结构。客户端的输入在代表永久性更改时被提交回服务器。
- 三层：将客户智能转移到中间层以便可以使用无状态客户端的体系结构。这简化了应用程序部署。大多数Web应用程序都是三层的。
- n- tier：通常指代Web应用程序的体系结构，它将请求进一步转发到其他企业服务。这种类型的应用程序是应用程序服务器成功最主要的应用程序。
- 点对点：没有专门的机器来提供服务或管理网络资源的体系结构。[24]：227相反，所有责任在所有机器中统一分配，即所谓的同伴。同行可以作为客户和服务器[25]。
分布式计算体系结构的另一个基本方面是在并发进程之间进行通信和协调工作的方法。通过各种消息传递协议，进程可以彼此直接通信，通常以主/从关系进行通信。或者，“以数据库为中心”的体系结构可以通过利用共享数据库，在没有任何形式的直接进程间通信的情况下实现分布式计算。[26]