调用 Go 的 HTTP Client 的 Get\Post 之类的方法时，默认是开启 HTTP keepalive 的，不过直接使用还是会遇到一些情况导致持久连接失效。首先，Client 构造好 HTTP 请求后，利用 Transport 来发送请求并等待结果，默认使用 DefaultTransport 来实现，大多数情况下，自定义 Client 时，配置一下自带的 Transport 即可。

transport 主要围绕着 persistConn 来实现，通过当前请求的 proxy, scheme, addr 作为 Key，对已经建立的连接进行缓存，新的请求来时，先从缓存中取一个连接，如果没有，再新发起一个连接。按照 Go 的基本法，毫无疑问会有两个 goroutine 来分别处理连接上的读和写，然后各种 channel 就开始飞来飞去，于是便让人深思这真的会比基于事件回调的实现简单吗。

好了，这里还是简单写点代码浅显的试验一下可能会导致持久连接失效的一些情况。

准备工作
在发送请求建立 TCP 连接时输出一些提示消息，这样便可以确定是否发起了新连接。自定义一个 Transport 的 Dial 方法就好了。这里输出连接的 LocalAddr，代码如下：

func PrintLocalDial(network, addr string) (net.Conn, error) {
    dial := net.Dialer{
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,
    }

    conn, err := dial.Dial(network, addr)
    if err != nil {
        return conn, err
    }

    fmt.Println("connect done, use", conn.LocalAddr().String())

    return conn, err
}
本地起了个在 8888 端口监听的 Web Server。好了，现在使用的 http.Client 如下：

client := &http.Client{
    Transport: &http.Transport{
        Dial: PrintLocalDial,
    },
}
发起请求
一个三观正确的，具有普遍意义的请求步骤如下所示：

func doGet(client *http.Client, url string, id int) {
    resp, err := client.Get(url)
    if err != nil {
        fmt.Println(err)
        return
    }

    buf, err := ioutil.ReadAll(resp.Body)
    fmt.Printf("%d: %s -- %v\n", id, string(buf), err)
    if err := resp.Body.Close(); err != nil {
        fmt.Println(err)
    }
}
调用 Get 进入 RoundTrip，首先会找到个 persistConn (从缓存中找个已经存在的或者新建一个)，再调用它的 roundTrip，这时会把这个请求发送到 writeLoop，然后等待响应的到来，当 readLoop 中读到响应数据时，便会把响应 Response 发到 roundTrip，自此，Get 方法返回，不过事情还没有结束，响应的 Body 还没有读取，readLoop 会一直阻塞等待读取数据，也就是当前这个 persistConn 一直被占用着，当读取完 resp.Body，readLoop 就会把 persistConn 放回连接缓存中，以便下个请求继续使用。

持续发几个请求试试：

const URL = "http://localhost:8888/"

for {
    go doGet(client, URL, 1)
    go doGet(client, URL, 2)
    time.Sleep(2 * time.Second)
}
每次同时发送两个请求，并等待请求完成，输出结果如下：

$ go run client.go
connect done, use [::1]:57571
connect done, use [::1]:57570
2: Hello, world -- <nil>
1: Hello, world -- <nil>
2: Hello, world -- <nil>
1: Hello, world -- <nil>
...
可见此时建立了两条 TCP 持久连接，后面的请求都复用了一开始建立好的连接。如果再加一个请求呢，每次同时发送三个请求，输出结果如下：

$ go run client.go
connect done, use [::1]:57582
connect done, use [::1]:57583
connect done, use [::1]:57584
2: Hello, world -- <nil>
1: Hello, world -- <nil>
3: Hello, world -- <nil>
connect done, use [::1]:57585
2: Hello, world -- <nil>
1: Hello, world -- <nil>
3: Hello, world -- <nil>
connect done, use [::1]:57586
1: Hello, world -- <nil>
2: Hello, world -- <nil>
3: Hello, world -- <nil>
...
可见每次都会有一个请求是新建了个 TCP 连接的，也就是说默认只保持两条持久连接，这是因为这里自定义的的 http.Transport 没有设置 MaxIdleConnsPerHost，于是便采用了默认的 DefaultMaxIdleConnsPerHost，这个值是 2，这是 RFC2616 建议的单个客户端发起的持久连接数，不过在大部分情况下，这个值有点过于保守了。如果把 MaxIdleConnsPerHost 设置为 3，结果便和第一种情况一样。

三观不正的请求
这里来一个三观不正的请求：

func doGet(client *http.Client, url string, id int) {
    resp, err := client.Get(url)
    if err != nil {
        fmt.Println(err)
        return
    }
    if err := resp.Body.Close(); err != nil {
        fmt.Println(err)
    }
    fmt.Printf("%d: done\n", id)
}
这里并不读取 resp.Body，因为读来也没用，但是也没法用 HEAD 请求(当然，这只是示例)。每次同时发送两个请求，结果如下：

$ go run client.go
connect done, use [::1]:57974
connect done, use [::1]:57973
2: done
1: done
connect done, use [::1]:57975
connect done, use [::1]:57976
2: done
1: done
connect done, use [::1]:57978
connect done, use [::1]:57977
...
额，每次都是新建的 TCP 连接，看来持久连接没用上。考虑到 TCP 接收到的数据，应用层并没有主动去读取，如果再次复用这个连接发送数据，那么上一次的数据要怎么处理？要么 Go 在库里默默的给读了，要么直接断开连接新建一条。在 Go 1.0 中就是在库里默默的读，想象一下正在 Get 的是几个 G 的东西，调用 Close 的时候... 所以最好的方法还是断开这个没有读取 Body 就直接 Close 的连接。

这里实现上利用 bodyEOFSignal 这个数据类型来包装 readLoop 生成的响应 Body，并设置回调函数 earlyCloseFn，如果 Body 并没有读取完便 Close，这个函数将执行并通知 readLoop，然后 readLoop 关闭连接并退出。所以如果想要使用持久连接，还得处理掉 Body，这就看应用的取舍了，尽量使用 HEAD 代替，也可以读取来丢弃掉。

n, err := io.Copy(ioutil.Discard, resp.Body)
后话
当然，可以设置 http.Transport 的 DisableKeepAlives 来禁用掉持久连接。前面废话说的有点多，总结一下无非就下面几条：

Web Server 得支持持久连接
如果有需要，加大 DefaultMaxIdleConnsPerHost 或者设置 MaxIdleConnsPerHost
读完 Response Body 再 Close
还记得那个默默奉献的 localhost:8888 不，感谢 Tornado。