# 服务注册和发现

### 1.服务注册:
- 通过Consul的服务注册HTTP API
- 通过在配置文件中定义服务的方式进行注册(建议)

```
//web3.json
{
  "service": {
    "name": "web3",
    "tags": ["master"],
    "address": "127.0.0.1",
    "port": 10000,
    "checks": [
      {
        "http": "http://localhost:10000/health",
        "interval": "10s"
      }
    ]
  }
}
```

对应程序:

```
//web3.go
package main
import (
    "fmt"
    "net/http"
)
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("hello Web3! This is n3")
    fmt.Fprintf(w, "Hello Web3! This is n3")
}
func healthHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("health check!")
}
func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/health", healthHandler)
    http.ListenAndServe(":10000", nil)
}
```



### 2.服务发现
- 通过HTTP API查看存在哪些服务
- 通过consul agent内置的DNS服务来做(可以根据服务check的实时状态动态调整available服务节点列表)

在配置和部署完web3服务后，我们就可以通过DNS命令来查询服务的具体信息了。consul为服务编排的内置域名为 ***“NAME.service.consul"***，这样我们的web3的域名为:web3.service.consul。我们在n1通过dig工具来查看一 下，注意是在n1上，n1上并未定义和部署web3服务，但集群中服务的信息已经被同步到n1上了，信息是一致的：


```
$ dig @127.0.0.1 -p 8600 web3.service.consul SRV
; <<>> DiG 9.9.5-3-Ubuntu <<>> @127.0.0.1 -p 8600 web3.service.consul SRV
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 6713
;; flags: qr aa rd; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 2
;; WARNING: recursion requested but not available
;; QUESTION SECTION:
;web3.service.consul.        IN    SRV
;; ANSWER SECTION:
web3.service.consul.    0    IN    SRV    1 1 10000 n2.node.dc1.consul.
web3.service.consul.    0    IN    SRV    1 1 10000 n3.node.dc1.consul.
;; ADDITIONAL SECTION:
n2.node.dc1.consul.    0    IN    A    127.0.0.1
n3.node.dc1.consul.    0    IN    A    127.0.0.1
;; Query time: 2 msec
;; SERVER: 127.0.0.1#8600(127.0.0.1)
;; WHEN: Mon Jul 06 12:12:53 CST 2015
;; MSG SIZE  rcvd: 219
```


我们可以看到consul agent将health check失败的web3从结果列表中剔除了，这样web3服务的客户端在服务发现过程中就只能获取到当前可用的web3服务节点了，这个好处是在实际应 用中大大降低了客户端实现”服务发现“时的难度。另外consul agent DNS在返回查询结果时也支持DNS Server常见的策略，至少是支持轮询。你可以多次执行dig命令，可以看到n2和n3的排列顺序是不同的。还有一点值得注意的是：由于考虑DNS cache对consul agent查询结果的影响，默认情况下所有由consul agent返回的结果TTL值均设为0，也就是说不支持dns结果缓存。


- 结果支持轮询
- 无 DNS cache

### 3.demo级别的服务发现客户端

// servicediscovery.go
package main
import (
    "fmt"
    "log"
    "github.com/miekg/dns"
)
const (
        srvName = "web3.service.consul"
        agentAddr = "127.0.0.1:8600"
)
func main() {
    c := new(dns.Client)
    m := new(dns.Msg)
    m.SetQuestion(dns.Fqdn(srvName), dns.TypeSRV)
    m.RecursionDesired = true
    r, _, err := c.Exchange(m, agentAddr)
    if r == nil {
        log.Fatalf("dns query error: %s\n", err.Error())
    }
    if r.Rcode != dns.RcodeSuccess {
        log.Fatalf("dns query error: %v\n", r.Rcode)
    }
   
    for _, a := range r.Answer {
        b, ok := a.(*dns.SRV)
        if ok {
            m.SetQuestion(dns.Fqdn(b.Target), dns.TypeA)
            r1, _, err := c.Exchange(m, agentAddr)
            if r1 == nil {
                log.Fatalf("dns query error: %v, %v\n", r1.Rcode, err)
            }
            for _, a1 := range r1.Answer {
                c, ok := a1.(*dns.A)
                if ok {
                   fmt.Printf("%s – %s:%d\n", b.Target, c.A, b.Port)
                }
            }
        }
    }
}

我们执行该程序：
$ go run servicediscovery.go
n2.node.dc1.consul. – 10.10.126.101:10000
n3.node.dc1.consul. – 10.10.126.187:10000
注意各个node上的服务check是由其node上的agent上进行的，一旦那个node上的agent出现问题，则位于那个node上的所有 service也将会被置为unavailable状态。比如我们停掉n3上的agent，那么我们在进行web3服务节点查询时，就只能获取到n2这一 个节点上有可用的web3服务了。


在真实的程序中，我们可以像上面demo中那样，每Request都做一次DNS查询，不过这样的代价也很高。稍复杂些，我们可以结合dns结果本地缓存+定期查询+每遇到Failed查询的方式来综合考量服务的发现方法或利用Consul提供的watch命令等。
以上仅仅是Consul的一个入门。真实场景中，理想的方案需要考虑的事情还有很多。Consul自身目前演进到0.5.2版本，还有不完善之处，但它已 经被很多公司用于production环境。Consul不是孤立的，要充分发挥出Consul的优势，在真实方案中，我们还要考虑与 Docker，HAProxy，Mesos等工具的结合。