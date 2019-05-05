# Go包管理

如上图所示：app_c包直接调用lib_a包中函数，并使用了lib_b包(v0.2版本)中的类型，lib_a包vendor了lib_b包(v0.1版本)。在这样的情况下，当我们编译app_c包时，是否会出现什么问题呢？我们一起来看一下这个例子：

```
在$GOPATH/src路径下面我们查看当前示例的目录结构：

$tree
├── app_c
    ├── c.go
├── lib_a
    ├── a.go
    └── vendor
        └── lib_b
            └── b.go
├── lib_b
    ├── b.go
```


各个源文件的示例代码如下：
```
//lib_a/a.go
package lib_a

import "lib_b"

func Foo(b lib_b.B) {
    b.Do()
}

//lib_a/vendor/lib_b/b.go

package lib_b

import "fmt"

type B struct {
}

func (*B) Do() {
    fmt.Println("lib_b version:v0.1")
}

// lib_b/b.go
package lib_b

import "fmt"

type B struct {
}

func (*B) Do() {
    fmt.Println("lib_b version:v0.2")
}

// app_c/c.go
package app_c

import (
    "lib_a"
    "lib_b"
)

func main() {
    var b lib_b.B
    lib_a.Foo(b)
}
```


进入app_c目录，执行编译命令：

```
$go build c.go
# command-line-arguments
./c.go:10:11: cannot use b (type "lib_b".B) as type "lib_a/vendor/lib_b".B in argument to lib_a.Foo
```

我们看到go compiler认为
***app_c包main函数中定义的变量b的类型(lib_b.B)与lib_a.Foo的参数b的类型(lib_a/vendor/lib_b.B)是不同的类型，不能相互赋值。***

## 2.通过手工vendor解决问题
这个例子非常有代表性，那么怎么解决这个问题呢？我们需要在app_c中也使用vendor机制，即将app_c所需的lib_a和lib_b都vendor到app_c中。

```
按照上述思路解决后的示例的目录结构：

$tree
├── app_c
    ├── c.go
    └── vendor
        ├── lib_a
        │   └── a.go
        └── lib_b
            └── b.go
├── lib_a
    ├── a.go
    └── vendor
        └── lib_b
            └── b.go
├── lib_b
    ├── b.go
```


# 关于对日志level的支持以及loglevel的热更新

#  json marshal json string时的转义问题



# 四. 如何在main包之外使用flag.Parse后的命令行flag变量

我们在使用Go开发交互界面不是很复杂的command-line应用时，一般都会使用std中的flag包进行命令行flag解析，并在main包中校验和使用flag.Parse后的flag变量。常见的套路是这样的：
//testflag1.go
package main

import (
    "flag"
    "fmt"
)

var (
    endpoints string
    user      string
    password  string
)

func init() {
    flag.StringVar(&endpoints, "endpoints", "127.0.0.1:2379", "comma-separated list of etcdv3 endpoints")
    flag.StringVar(&user, "user", "", "etcdv3 client user")
    flag.StringVar(&password, "password", "", "etcdv3 client password")
}

func usage() {
    fmt.Println("flagdemo-app is a daemon application which provides xxx service.\n")
    fmt.Println("Usage of flagdemo-app:\n")
    fmt.Println("\t flagdemo-app [options]\n")
    fmt.Println("The options are:\n")

    flag.PrintDefaults()
}

func main() {
    flag.Usage = usage
    flag.Parse()

   // ... ...
   // 这里我们可以使用endpoints、user、password等flag变量了
}

在这样的一个套路中，我们可以在main包中直接使用flag.Parse后的flag变量了。但有些时候，我们需要在main包之外使用这些flag vars(比如这里的：endpoints、user、password)，怎么做呢，有几种方法，我们逐一来看看。
## 1. 全局变量法

我想大部分gopher第一个想法就是使用全局变量，即建立一个config包，包中定义全局变量，并在main中将这些全局变量绑定到flag的Parse中：
$tree globalvars
globalvars
├── config
│   └── config.go
├── etcd
│   └── etcd.go
└── main.go

// flag-demo/globalvars/config/config.go

package config

var (
    Endpoints string
    User      string
    Password  string
)

// flag-demo/globalvars/etcd/etcd.go
package etcd

import (
    "fmt"

    "../config"
)

func EtcdProxy() {
    fmt.Println(config.Endpoints, config.User, config.Password)
    //... ....
}

// flag-demo/globalvars/main.go
package main

import (
    "flag"
    "fmt"
    "time"

    "./config"
    "./etcd"
)

func init() {
    flag.StringVar(&config.Endpoints, "endpoints", "127.0.0.1:2379", "comma-separated list of etcdv3 endpoints")
    flag.StringVar(&config.User, "user", "", "etcdv3 client user")
    flag.StringVar(&config.Password, "password", "", "etcdv3 client password")
}

.... ...

func main() {
    flag.Usage = usage
    flag.Parse()

    go etcd.EtcdProxy()

    time.Sleep(5 * time.Second)
}

可以看到，我们在绑定cmdline flag时使用的是config包中定义的全局变量。并且在另外一个etcd包中，使用了这些变量。
我们运行这个程序：
./main -endpoints 192.168.10.69:2379,10.10.12.36:2378 -user tonybai -password xyz123
192.168.10.69:2379,10.10.12.36:2378 tonybai xyz123
不过这种方法要注意这些全局变量值在Go包初始化过程的顺序，比如：如果在etcd包的init函数中使用这些全局变量，那么你得到的各个变量值将为空值，因为etcd包的init函数在main.init和main.main之前执行，这个时候绑定和Parse都还未执行。
## 2. 传参法

第二种比较直接的想法就是将Parse后的flag变量以参数的形式、以某种init的方式传给其他要使用这些变量的包。
$tree parampass
parampass
├── etcd
│   └── etcd.go
└── main.go

// flag-demo/parampass/etcd/etcd.go
package etcd
... ...

func EtcdProxy(endpoints, user, password string) {
    fmt.Println(endpoints, user, password)
}

// flag-demo/parampass/main.go
package main

import (
    "flag"
    "fmt"
    "time"

    "./etcd"
)

var (
    endpoints string
    user      string
    password  string
)

func init() {
    flag.StringVar(&endpoints, "endpoints", "127.0.0.1:2379", "comma-separated list of etcdv3 endpoints")
    flag.StringVar(&user, "user", "", "etcdv3 client user")
    flag.StringVar(&password, "password", "", "etcdv3 client password")
}

... ...

func main() {
    flag.Usage = usage
    flag.Parse()

    go etcd.EtcdProxy(endpoints, user, password)

    time.Sleep(5 * time.Second)
}

这种方法非常直观，这里就不解释了。但注意：一旦使用这种方式，一定需要在main包与另外的包之间建立某种依赖关系，至少main包会import那些使用flag变量的包。
## 3. 配置中心法

全局变量法直观，而且一定程度上解除了其他包与main包的耦合。但是有一个问题，那就是一旦flag变量发生增减，config包就得相应添加或删除变量定义。是否有一种方案可以在flag变量发生变化时，config包不受影响呢？我们可以用配置中心法。所谓的配置中心法，就是实现一个与flag变量类型和值无关的通过配置存储结构，我们在main包中向该结构注入parse后的flag变量，在其他需要flag变量的包中，我们使用该结构得到flag变量的值。
$tree configcenter
configcenter
├── config
│   └── config.go
└── main.go

//flag-demo/configcenter/config/config.go
package config

import (
    "log"
    "sync"
)

var (
    m  map[string]interface{}
    mu sync.RWMutex
)

func init() {
    m = make(map[string]interface{}, 10)
}

func SetString(k, v string) {
    mu.Lock()
    m[k] = v
    mu.Unlock()
}

func SetInt(k string, i int) {
    mu.Lock()
    m[k] = i
    mu.Unlock()
}

func GetString(key string) string {
    mu.RLock()
    defer mu.RUnlock()
    v, ok := m[key]
    if !ok {
        return ""
    }
    return v.(string)
}

func GetInt(key string) int {
    mu.RLock()
    defer mu.RUnlock()
    v, ok := m[key]
    if !ok {
        return 0
    }
    return v.(int)
}

func Dump() {
    log.Println(m)
}

// flag-demo/configcenter/main.go

package main

import (
    "flag"
    "fmt"
    "time"

    "./config"
)

var (
    endpoints string
    user      string
    password  string
)

func init() {
    flag.StringVar(&endpoints, "endpoints", "127.0.0.1:2379", "comma-separated list of etcdv3 endpoints")
    flag.StringVar(&user, "user", "", "etcdv3 client user")
    flag.StringVar(&password, "password", "", "etcdv3 client password")
}
... ...
func main() {
    flag.Usage = usage
    flag.Parse()

    // inject flag vars to config center
    config.SetString("endpoints", endpoints)
    config.SetString("user", user)
    config.SetString("password", password)

    time.Sleep(5 * time.Second)
}

我们在main中使用config的SetString将flag vars注入配置中心。之后，我们在其他包中就可以使用：GetString、GetInt获取这些变量值了，这里就不举例了。
4、“黑魔法”: flag.Lookup

flag包中提供了一种类似上述的”配置中心”的机制，但这种机制不需要我们显示注入“flag vars”了，我们只需按照flag提供的方法在其他package中读取对应flag变量的值即可。
$tree flaglookup
flaglookup
├── etcd
│   └── etcd.go
└── main.go

// flag-demo/flaglookup/main.go
package main

import (
    "flag"
    "fmt"
    "time"

    "./etcd"
)

var (
    endpoints string
    user      string
    password  string
)

func init() {
    flag.StringVar(&endpoints, "endpoints", "127.0.0.1:2379", "comma-separated list of etcdv3 endpoints")
    flag.StringVar(&user, "user", "", "etcdv3 client user")
    flag.StringVar(&password, "password", "", "etcdv3 client password")
}

......

func main() {
    flag.Usage = usage
    flag.Parse()

    go etcd.EtcdProxy()

    time.Sleep(5 * time.Second)
}

// flag-demo/flaglookup/etcd/etcd.go
package etcd

import (
    "flag"
    "fmt"
)

func EtcdProxy() {
    endpoints := flag.Lookup("endpoints").Value.(flag.Getter).Get().(string)
    user := flag.Lookup("user").Value.(flag.Getter).Get().(string)
    password := flag.Lookup("password").Value.(flag.Getter).Get().(string)

    fmt.Println(endpoints, user, password)
}

运行该程序：
$go run main.go -endpoints 192.168.10.69:2379,10.10.12.36:2378 -user tonybai -password xyz123
192.168.10.69:2379,10.10.12.36:2378 tonybai xyz123
输出与我们的预期是一致的。
## 5、对比

我们用一幅图来对上述几种方法进行对比：
```
全局变量法
main----> config ----->other_packages

参数传递法:

main<------------other_packages


配置中心法:
main--->config<-----other_packages


“黑魔法”falg lookup
main----->falg<-------other_packages
```


很显然，经过简单包装后，“黑魔法”flaglookup应该是比较优异的方案。main包、other packages只需import flag即可。
注意：在main包中定义exported的全局flag变量并被其他package import的方法是错误的，很容易造成import cycle问题。并且任何其他package import main包都是不合理的。