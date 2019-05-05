之前在进行微信Demo开发时曾用到过ngrok这个强大的tunnel(隧道)工具,ngrok在其github官方页面上的自我诠释是"introspected tunnels to localhost",这个诠释有两层含义:

1.可以用来建立public到localhost的tunnel,让居于内网主机上的服务可以暴露给public,俗称***内网穿透***。
2.支持对隧道中数据的introspection(内省),支持可视化的观察隧道内数据,并且replay(重放)相关请求(诸如http请求)。

因此ngrok可以很便捷的协助进行服务端程序调试,尤其在进行一些Web server开发中。ngrok更强大的一点是它支持tcp层之上的所有应用协议或者说与应用层协议无关。比如:你可以通过ngrok实现ssh登陆到内网主机,也可以通过ngrok实现远程桌面(VNC)方式访问内网主机。

今天我们就来简单分析一下这款强大工具的实现原理。ngrok本身是用go语言实现的，需要go 1.1以上版本编译。ngrok官方代码最新版为1.7，作者似乎已经完成了ngrok 2.0版本，但不知为何迟迟不放出最新代码。因此这里我们就以ngrok 1.7版本源码作为原理分析的基础。

# 一、ngrok tunnel与ngrok部署
网络tunnel（隧道）对多数人都是很”神秘“的概念，tunnel种类很多，没有标准定义，我了解的也不多（日常工作较少涉及），这里也就不 深入了。在《HTTP权威指南》中有关于HTTP tunnel（http上承载非web流量）和SSL tunnel的说明，但ngrok中的tunnel又与这些有所不同。

ngrok实现了一个tcp之上的端到端的tunnel，两端的程序在ngrok实现的Tunnel内透明的进行数据交互。

                               隧道
Server   <-----------------HTTP Request----------------- Client
(Local)  ------------------HTTP Response---------------> (public)
          | -----------------------------------------|
         Ngrok                                     Ngrokd
        
ngrok分为client端(ngrok)和服务端(ngrokd),实际使用中的部署如下:

内网主机

Server
(local)
                                                                            公网主机            Client
Server           Server On Ngrok host       内网代理        防火墙            (Ngrokd)           (Public)
(local)                  (Ngrok)

Server
(local)


内网服务程序可以与ngrok client部署在同一主机,也可以部署在内网可达的其他主机上。ngrok和ngrokd会为建立与public client间的专用通道(tunnel)。



# 二、ngrok开发调试环境搭建

生成tunnel.tonybai.com的证书

我们这里以NGROK_BASE_DOMAIN="tunnel.tonybai.com"为例，生成证书的命令如下：

$ cd ~/goproj/src/github.com/inconshreveable/ngrok
$ openssl genrsa -out rootCA.key 2048
$ openssl req -x509 -new -nodes -key rootCA.key -subj "/CN=tunnel.tonybai.com" -days 5000 -out rootCA.pem
$ openssl genrsa -out device.key 2048
$ openssl req -new -key device.key -subj "/CN=tunnel.tonybai.com" -out device.csr
$ openssl x509 -req -in device.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -out device.crt -days 5000


执行完以上命令，在ngrok目录下就会新生成6个文件：

-rw-rw-r– 1 ubuntu ubuntu 1001 Mar 14 02:22 device.crt
-rw-rw-r– 1 ubuntu ubuntu  903 Mar 14 02:22 device.csr
-rw-rw-r– 1 ubuntu ubuntu 1679 Mar 14 02:22 device.key
-rw-rw-r– 1 ubuntu ubuntu 1679 Mar 14 02:21 rootCA.key
-rw-rw-r– 1 ubuntu ubuntu 1119 Mar 14 02:21 rootCA.pem
-rw-rw-r– 1 ubuntu ubuntu   17 Mar 14 02:22 rootCA.srl

ngrok通过bindata将ngrok源码目录下的assets目录（资源文件）打包到可执行文件(ngrokd和ngrok)中 去，assets/client/tls和assets/server/tls下分别存放着用于ngrok和ngrokd的默认证书文件，我们需要将它们替换成我们自己生成的：(因此这一步务必放在编译可执行文件之前)

cp rootCA.pem assets/client/tls/ngrokroot.crt
cp device.crt assets/server/tls/snakeoil.crt
cp device.key assets/server/tls/snakeoil.key


make release-server
make release-client


 ./repo/ngrok/bin/ngrok -subdomain example -config=debug.yml -log=ngrok.log 8888
 ./repo/ngrok/bin/ngrokd -domain="tunnel.tonybai.com" -httpAddr=":8080" -httpsAddr=":8081"


debug.yml内容
server_addr: "tunnel.tonybai.com:4443"
trust_host_root_certs: false
tunnels:
  test:
    proto:
      http: 8888

    
cat /etc/hosts
127.0.0.1	ngrok.me
127.0.0.1	test.ngrok.me
127.0.0.1	tunnel.tonybai.com
127.0.0.1	example.tunnel.tonybai.com



ngrok                                                                                                                  (Ctrl+C to quit)

Tunnel Status                 online
Version                       1.7/1.7
Forwarding                    http://example.tunnel.tonybai.com:8080 -> 127.0.0.1:8888
Forwarding                    https://example.tunnel.tonybai.com:8080 -> 127.0.0.1:8888
Web Interface                 127.0.0.1:4040
# Conn                        0
Avg Conn Time                 0.00ms


# 三、第一阶段:Control Connection建立
在ngrokd的启动日志中我们可以看到这样一行:

[INFO] Listening for control and proxy connections on [::]:4443

ngrokd在4443端口(默认)监听control和proxy connection。Control Connection,顾名思义"控制连接",有些类似于FTP协议的控制连接(不知道ngrok作者在设计协议时是否参考了FTP协议^_^)。该连接只用于收发控制类消息。作为客户端的ngrok启动后的第一件事就是与ngrokd监理Control Connection,建立过程序列图如下:

ngrok                                         ngrokd

    -----向ngrokd的4443端口建立TCP连接--------->
    -----------Auth Msg--------------------->(进行TLS handshake过程,包括ssl证书校验等)
    <----------Auth Response----------------


客户端 
func Main()
// parse options
// set up logging
// read configuration file
... 
... 
NewController().Run(config)
    - func (ctl *Controller) Run(config *Configuration) {
        - var model *ClientModel
        - if ctl.model == nil 
          - model = ctl.SetupModel(config)
        - else
          - model = ctl.model.(*ClientModel)
        - // init the model
        - // init web ui
        - // init term ui
        ... 
        - ***ctl.Go(ctl.model.Run)***
            - func (c *ClientModel) Run() //***Run函数调用c.control来运行Control Connection的主逻辑，并在control connection断开后，尝试重连。***
              - ... 
              - for 
                - // ***run the control channel***
                - ***c.control() c.control是ClientModel的一个method,用来真正建立ngrok到ngrokd的control connection,并完成ngrok的鉴权(用户名、密码配置在配置文件中)***
                    - func (c *ClientModel) control()
                      - ... 
                      - var( ctlConn conn.Conn )
                      - var( err error)
                      - if c.proxyUrl==""
                        - // simple non-proxied case,just connect to the server
                        - ctlConn,err=***conn.Dial(c.serverAddr,"ctl",c.tlsConfig)***
                             - ngrok封装了connection相关操作，代码在ngrok/src/ngrok/conn下面，包名conn。
                             - func Dial(addr, typ string, tlsCfg *tls.Config) (conn *loggedConn, err error)
                                - var rawConn net.Conn
                                - if rawConn,err=net.Dial("tcp",addr);err!=nil{
                                  return
                                }
                                - conn=wrapConn(rawConn,typ)
                                - conn.Debug("New connection to: %v",rawConn.RemoteAddr())
                                - if tlsCfg!=nil
                                  - ***conn.StartTLS(tlsCfg) //ngrok首先创建一条TCP了连接,并基于该连接创建了TLS client,不过此时并未进行TLS的初始化,即handshake。handshake发生在ngrok首次向ngrokd发送auth消息(msg.WriteMsg,ngrok/src/ngrok/msg/msg.go)时,go标准库的TLS相关函数默默的完成了这一handshanke过程。我们经常遇到的ngrok证书验证失败等问题,就发生在该过程中。*** 
                                        - func (c *loggedConn) StartTLS(tlsCfg *tls.Config)
                                          - c.Conn=tls.Client(c.conn,tlsCfg)
                                - return

                      - else
                        - ... 
                      -// ***authenticate with the server***
                      - auth:=&msg.Auth{
                        ClientId: c.id,
                        OS: runtime.GOOS,
                        Arch: runtime.GOARCH,
                        Version: version.Proto,
                        MmVersion: version.MajorMinor(),
                        User: c.authToken,
                      }
                      - if err=***msg.WriteMsg(ctlConn,auth)***;err!=nil{
                        panic(err)
                      }
                      - ***// wait for ther server to authenticate us***
                      - var authResp msg.AuthResp
                      - if err=***msg.ReadMsgInfo(ctlConn,&authResp)***;err!=nil{
                        panic(err)
                      }
                      - ...
                      -c.id=authResp.ClientId
                      - ...
                - ... 
                - if c.connStatus==mvc.ConnOnline
                  - wait=1*time.Second
                - ... 
                - c.connStatus=mvc.ConnReconnecting
                - c.update()
        ... 

在AuthResp中,ngrokd为该Control Connection分配一个ClientID,该ClientID在后续Proxy Connection建立时使用,用于关联和校验之用。



前面的逻辑和代码都是ngrok客户端的,现在我们再从ngorkd server端代码review一遍Control Connection的建立过程。

ngrokd的代码放在ngrok/src/ngrok/server下面,entrypoint如下:

func Main()
- // parse options
- opts=parseArgs()
- // init logging
- // init tunnel/control registry
- ... 
- // start listeners
- listeners=make(map[string]*conn.Listener)
- // load tls configuration
- tlsConfig,err:=LoadTLSConfig(opts.tlsCrt,opts.tlsKey)
- // listen for http
- // listen for https
- ... 
- // ngrok clients
- ***tunnelListener(opts.tunnelAddr,tlsConfig) // ngrokd启动了三个监听,其中最后一个tunnelListenner用于监听ngrok发起的Control Connection或者后续的proxy connection，作者意图通过一个端口，监听两种类型连接，旨在于方便部署。***
    - func tunnelListener(addr string, tlsConfig *tls.Config)
      - // listen for incoming connections
      - listener,err:=conn.Listen(addr,"tun",tlsConfig)
      - ... 
      - for c:=range listener.Conns
        - go func(tunnelConn conn.Conn){
            - ... 
            - var rawMsg msg.Message
            - if rawMsg,err=msg.ReadMsg(tunnelConn)
            - switch m:=rawMsg.(type)
            - case *msg.Auth:
              - ***NewControl(tunnelConn,m) //可以看到,当ngrokd在新建立的Control Connection上收到Auth消息后,ngrokd执行NewControl来处理该Control Connection上后续的事情***
                  - func NewControl(ctlConn conn.Conn, authMsg *msg.Auth)
                    - var err error
                    - // create the object
                    - c:=&Control{...}
                    - // register the clientid
                    - ... 
                    - // register the control
                    - ... 
                    - start the writer first so that the following messages get sent
                    - go c.writer()
                    - // Respond to authentication
                    - c.out <- &msg.AuthResp{
                      Verion: version.Proto,
                      MmVersion: version.MajorMinor(),
                      ClientId: c.id,
                    }
                    - // As a performance optimization,ask for a proxy connection up front
                    - c.out<-&msg.ReqProxy{}
                    - // manage the connection
                    - go c.manager()
                    - go c.reader()
                    - go c.stopper()
          }(c)


在NewControl中，ngrokd返回了AuthResp。到这里，一条新的Control Connection建立完毕。

我们最后再来看一下Control Connection建立过程时ngrok和ngrokd的输出日志，增强一下感性认知：

ngrok Server:

[INFO] [tun:d866234] ***New connection*** from 127.0.0.1:59949
[DEBG] [tun:d866234] Waiting to read message
[DEBG] [tun:d866234] Reading message with length: 126
[DEBG] [tun:d866234] Read message {"Type":***"Auth***",
"Payload":{"Version":"2","MmVersion":"1.7","User":"","Password":"","OS":"darwin","Arch":"amd64","ClientId":""}}
[INFO] [ctl:d866234] Renamed connection tun:d866234
[INFO] [registry] [ctl] Registered control with id ac1d14e0634f243f8a0cc2306bb466af
[DEBG] [ctl:d866234] [ac1d14e0634f243f8a0cc2306bb466af] Writing message: {"Type":"***AuthResp***","Payload":{"Version":"2","MmVersion":"1.7","ClientId":"ac1d14e0634f243f8a0cc2306bb466af","Error":""}}

Client:

[INFO] (ngrok/log.Info:112) Reading configuration file debug.yml
[INFO] (ngrok/log.(*PrefixLogger).Info:83) [client] Trusting root CAs: [assets/client/tls/ngrokroot.crt]
[INFO] (ngrok/log.(*PrefixLogger).Info:83) [view] [web] Serving web interface on 127.0.0.1:4040
[INFO] (ngrok/log.Info:112) Checking for update
[DEBG] (ngrok/log.(*PrefixLogger).Debug:79) [view] [term] Waiting for update
[DEBG] (ngrok/log.(*PrefixLogger).Debug:79) [ctl:31deb681] ***New connection to***: 127.0.0.1:4443
[DEBG] (ngrok/log.(*PrefixLogger).Debug:79) [ctl:31deb681] Writing message: {"Type":"***Auth***","Payload":{"Version":"2","MmVersion":"1.7","User":"","Password":"","OS":"darwin","Arch":"amd64","ClientId":""}}
[DEBG] (ngrok/log.(*PrefixLogger).Debug:79) [ctl:31deb681] Waiting to read message
(ngrok/log.(*PrefixLogger).Debug:79) [ctl:31deb681] Reading message with length: 120
(ngrok/log.(*PrefixLogger).Debug:79) [ctl:31deb681] Read message {"Type":"***AuthResp***","Payload":{"Version":"2","MmVersion":"1.7","ClientId":"ac1d14e0634f243f8a0cc2306bb466af","Error":""}}
[INFO] (ngrok/log.(*PrefixLogger).Info:83) [client] Authenticated with server, client id: ac1d14e0634f243f8a0cc2306bb466af


           

# 四、Tunnel Creation

ngrok                   ngrokd
    -----ReqTunnel-------->
    <----NewTunnel---------


Tunnel Creation是ngrok将配置文件中的tunnel信息通过刚刚建立的Contro Connection传输给ngrokd,ngrokd登记、启动相应端口监听(如果配置了remoete_port或多路复用ngrokd默认监听的http和https端口)并返回相应应答。ngrok和ngrokd之间并未真正建立新连接。

我们回到ngrok的model.go，继续看ClientModel的control方法。在收到AuthResp后，ngrok还做了如下事情：


reqIdToTunnelConfig := make(map[string]*TunnelConfiguration)
for _, config := range c.tunnelConfig {
  // create the protocol list to ask for
  var protocols []string
  for proto, _ := range config.Protocols {
      protocols = append(protocols, proto)
  }

  reqTunnel := ***&msg.ReqTunnel***{
      … …
  }

  // send the tunnel request
  if err = msg.WriteMsg(ctlConn, reqTunnel); err != nil {
      panic(err)
  }

  // save request id association so we know which local address
  // to proxy to later
  reqIdToTunnelConfig[reqTunnel.ReqId] = config
}

// main control loop
for {
  var rawMsg msg.Message
  
  switch m := rawMsg.(type) {
  … …
  case *msg.***NewTunnel***:
      … …

      tunnel := mvc.Tunnel{
          … …
      }

      c.tunnels[tunnel.PublicUrl] = tunnel
      c.connStatus = mvc.ConnOnline
      
      c.update()
  … …
  }
}


ngrok将配置的Tunnel信息逐一以ReqTunnel消息发送ngrokd以注册登记Tunnel,并在随后的main control loop中处理ngrokd回送的NewTunnel消息,完
成一些登记索引工作。



ngrokd Server端对tunnel creation的处理是在NewControl的结尾处：

//ngrok/src/ngrok/server/control.go
func NewControl(ctlConn conn.Conn, authMsg *msg.Auth) {
    … …
    // manage the connection
    ***go c.manager()***
    … …
}

func (c *Control) manager() {
//… …
for {
    select {
    case <-reap.C:
        … …

    case mRaw, ok := <-c.in:
        // c.in closes to indicate shutdown
        if !ok {
            return
        }

        switch m := mRaw.(type) {
        case *msg.ReqTunnel:
  ***c.registerTunnel(m)***

        .. …
        }
    }
}
}



Control的manager在收到ngrok发来的ReqTunnel消息后，调用registerTunnel进行处理。

// ngrok/src/ngrok/server/control.go
// Register a new tunnel on this control connection
func (c *Control) registerTunnel(rawTunnelReq *msg.ReqTunnel)
  - for _, proto := range strings.Split(rawTunnelReq.Protocol, "+")
    - tunnelReq := *rawTunnelReq
    - tunnelReq.Protocol = proto
    - c.conn.Debug("Registering new tunnel")
    - t, err := ***NewTunnel***(&tunnelReq, c)
    - if err != nil {
        - c.out <- ***&msg.NewTunnel***{Error: err.Error()}
        - if len(c.tunnels) == 0 
                c.shutdown.Begin()
        - // we're done
        - return

        // add it to the list of tunnels
        - c.tunnels = append(c.tunnels, t)

        // acknowledge success
        - ***c.out <- &msg.NewTunnel***{
            Url:      t.url,
            Protocol: proto,
            ReqId:    rawTunnelReq.ReqId,
        }
        - rawTunnelReq.Hostname = strings.Replace(t.url, proto+"://", "", 1)
    }
}

Server端创建tunnel的实际工作由NewTunnel完成：
// ngrok/src/ngrok/server/tunnel.go
func NewTunnel(m *msg.ReqTunnel, ctl *Control) (t *Tunnel, err error)
- t = &Tunnel{...}
- proto := t.req.Protocol
- switch proto {
- case "tcp":
    - bindTcp := func(port int) error
        - if t.listener, err = net.ListenTCP("tcp",&net.TCPAddr{IP: net.ParseIP("0.0.0.0"),Port: port}); err != nil
            - ...
            - return err

        - // create the url
        - addr := t.listener.Addr().(*net.TCPAddr)
        - t.url = fmt.Sprintf("tcp://%s:%d", opts.domain, addr.Port)

        - // register it
        - if err = tunnelRegistry.RegisterAndCache(t.url, t);err != nil {
            - ...
            - return err

        - go t.listenTcp(t.listener)
        - return nil

        - // use the custom remote port you asked for
        - if t.req.RemotePort != 0
            - bindTcp(int(t.req.RemotePort))
            - return

        - // try to return to you the same port you had before
        - cachedUrl := tunnelRegistry.GetCachedRegistration(t)
        - if cachedUrl != ""
            - ...

        - // Bind for TCP connections
        - bindTcp(0)
        - return

- case "http", "https":
    - l, ok := listeners[proto]
    - if !ok
        - ... 
        - return

    - if err = registerVhost(t, proto, l.Addr.(*net.TCPAddr).Port);err != nil
        - return

    - default:
        - err = fmt.Errorf("Protocol %s is not supported", proto)
        - return

- ...
- metrics.OpenTunnel(t)
- return

可以看出,NewTunnel区别对待tcp和http/https隧道:
- 对于Tcp隧道,NewTunnel先要看是否配置了remote_port,如果remote_port不为空,则启动监听这个remote_port。否则尝试从cache里找出你之前创建tunnel时使用的端口号,如果可用,则监听这个端口号,否则bindTcp(0)，即随即选择一个端口作为该tcp tunnel的remote_port。
- 对于http/https隧道,ngrokd启动时就默认监听了80和443,如果ngrok请求建立http/https隧道(目前不支持设置remote_port),则ngrokd通过一种自实现的vhost的机制实现所有http/https请求多路复用到80和443端口上。ngrokd不会新增监听端口。


# 五、Proxy Connection和Private Connection
到目前为止,我们知道了Control Connection:用于ngrok和ngrokd之间传输命令;Public Connection:外部发起的,尝试向内网服务建立的链接。

这节当中，我们要接触到Proxy Connection和Private Connection。

ngrok                          ngrokd                             public
    <---------ReqProxy------------- 
          (On Control Connection)

    --------Proxy Connection------>

    ----------RegProxy------------>
          (On Proxy Connection)
                                      <------Public Connection----->
    <----------Start Proxy---------
          (On Proxy Connection)

    --
      | (Private Connection)
    <--
                                      <------Data Transporting-----
    <--------Data Transfering------
          (On Proxy Connection)

    --
      | (Data Transfering 
      | on Private Connection )
    <--


  前面ngrok和ngrokd的交互进行到了NewTunnel,这些数据都是通过之前已经建立的Control Connection上传输的。

  ngrokd侧,NewControl方法的结尾有这样一行代码:
    // As a performance optimization, ask for a proxy connection up front
    c.out <- &msg.ReqProxy{}

  服务端ngrokd在Control Connection上向ngrok发送了"ReqProxy"的消息,意为请求ngrok向ngrokd建立一条Proxy Connection,该链接将作为隧道数据流的承载者。

  客户端ngrok在ClientModel control方法的main control loop中收到ReqProxy并处理该消息:
  case *msg.ReqProxy:
    - ***c.ctl.Go(c.proxy)***

// Establishes and manages a tunnel proxy connection with the server
func (c *ClientModel) proxy() {
    if c.proxyUrl == "" {
        remoteConn, err = conn.Dial(c.serverAddr, "pxy", c.tlsConfig)
    }……

    err = msg.WriteMsg(remoteConn, &msg.RegProxy{ClientId: c.id})
    if err != nil {
        remoteConn.Error("Failed to write RegProxy: %v", err)
        return
    }
    … …
}

ngrok客户端收到ReqProxy后，创建一条新连接到ngrokd，该连接即为Proxy Connection。并且ngrok将RegProxy消息通过该新建立的Proxy Connection发到ngrokd，以便ngrokd将该Proxy Connection与对应的Control Connection以及tunnel关联在一起。

// ngrok服务端
func tunnelListener(addr string, tlsConfig *tls.Config) {
    …. …
    case *msg.RegProxy:
                NewProxy(tunnelConn, m)
    … …
}


#【到目前为止, tunnel、Proxy Connection都已经建立了，万事俱备，就等待Public发起Public connection到ngrokd了。】

下面我们以Public发起一个http连接到ngrokd为例,比如我们通过curl命令,向test.ngrok.me发起一次http请求。

前面说过，ngrokd在启动时默认启动了80和443端口的监听，并且与其他http/https隧道共同多路复用该端口（通过vhost机制)。ngrokd server对80端口的处理代码如下：


// ngrok/src/ngrok/server/main.go
func Main() {
    … …
 // listen for http
    if opts.httpAddr != "" {
        listeners["http"] =
          ***startHttpListener(opts.httpAddr, nil)***
    }

    … …
}

startHttpListener针对每个连接，启动一个goroutine专门处理：

//ngrok/src/ngrok/server/http.go
func startHttpListener(addr string,
    tlsCfg *tls.Config) (listener *conn.Listener) {
    // bind/listen for incoming connections
    var err error
    if listener, err = conn.Listen(addr, "pub", tlsCfg);
        err != nil {
        panic(err)
    }

    proto := "http"
    if tlsCfg != nil {
        proto = "https"
    }

   … …
    go func() {
        for conn := range listener.Conns {
            ***go httpHandler(conn, proto)***
        }
    }()

    return
}

// Handles a new http connection from the public internet
func httpHandler(c conn.Conn, proto string) {
    … …
    // let the tunnel handle the connection now
    ***tunnel.HandlePublicConnection(c)***
}

我们终于看到server端处理public connection的真正方法了:

//ngrok/src/ngrok/server/tunnel.go
func (t *Tunnel) HandlePublicConnection(publicConn conn.Conn) {
    … …
    var proxyConn conn.Conn
    var err error
    for i := 0; i < (2 * proxyMaxPoolSize); i++ {
        // get a proxy connection
        if proxyConn, err = ***t.ctl.GetProxy()***;
           err != nil {
            … …
        }
        defer proxyConn.Close()
       … …

        // tell the client we're going to
        // start using this proxy connection
        startPxyMsg := ***&msg.StartProxy***{
            Url:        t.url,
            ClientAddr: publicConn.RemoteAddr().String(),
        }

        if err = msg.WriteMsg(proxyConn, startPxyMsg);
            err != nil {
           … …
        }
    }

  … …
  // join the public and proxy connections
  bytesIn, bytesOut := ***conn.Join***(publicConn, proxyConn)
  …. …
}

HandlePublicConnection通过选出的Proxy connection向ngrok client发送StartProxy信息，告知ngrok proxy启动。然后通过conn.Join方法将publicConn和proxyConn关联到一起。

// ngrok/src/ngrok/conn/conn.go
// ***核心核心核心核心***
***func Join(c Conn, c2 Conn) (int64, int64)*** {
    var wait sync.WaitGroup

    pipe := func(to Conn, from Conn, bytesCopied *int64) {
        defer to.Close()
        defer from.Close()
        defer wait.Done()

        var err error
        *bytesCopied, err = io.Copy(to, from)
        if err != nil {
            from.Warn("Copied %d bytes to %s before failing with error %v", *bytesCopied, to.Id(), err)
        } else {
            from.Debug("Copied %d bytes to %s", *bytesCopied, to.Id())
        }
    }

    wait.Add(2)
    var fromBytes, toBytes int64
    go pipe(c, c2, &fromBytes)
    go pipe(c2, c, &toBytes)
    c.Info("Joined with connection %s", c2.Id())
    wait.Wait()
    return fromBytes, toBytes
}


***Join通过io.Copy实现public conn和proxy conn数据流的转发，单向被称作一个pipe，Join建立了两个Pipe，实现了双向转发，每个Pipe直到一方返回EOF或异常失败才会退出。后续在ngrok端，proxy conn和private conn也是通过conn.Join关联到一起的。***

我们现在就来看看ngrok在收到StartProxy消息后是如何处理的。我们回到ClientModel的proxy方法中。在向ngrokd成功建立proxy connection后，ngrok等待ngrokd的StartProxy指令。

    // wait for the server to ack our register
    var startPxy msg.StartProxy
    if err = msg.ReadMsgInto(remoteConn, &startPxy);
             err != nil {
        remoteConn.Error("Server failed to write StartProxy: %v",
                   err)
        return
    }

一旦收到StartProxy，ngrok将建立一条private connection：
    // start up the private connection
    start := time.Now()
    localConn, err := conn.Dial(tunnel.LocalAddr, "prv", nil)
    if err != nil {
       … …
        return
    }
并将private connection和proxy connection通过conn.Join关联在一起，实现数据透明转发。

    m.connTimer.Time(func() {
        localConn := tunnel.Protocol.WrapConn(localConn,
             mvc.ConnectionContext{Tunnel: tunnel,
              ClientAddr: startPxy.ClientAddr})
        bytesIn, bytesOut := conn.Join(localConn, remoteConn)
        m.bytesIn.Update(bytesIn)
        m.bytesOut.Update(bytesOut)
        m.bytesInCount.Inc(bytesIn)
        m.bytesOutCount.Inc(bytesOut)
    })

这样一来，public connection上的数据通过proxy connection到达ngrok，ngrok再通过private connection将数据转发给本地启动的服务程序，从而实现所谓的内网穿透。从public视角来看，就像是与内网中的那个服务直接交互一样。
