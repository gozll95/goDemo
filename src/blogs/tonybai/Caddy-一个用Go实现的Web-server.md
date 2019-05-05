
这是一个Web Server的时代，apache2与nginx共舞，在追求极致性能的路上，没有最高，只有更高。但这又是一个追求个性化的时代，有些Web Server并没有去挤“Performance提升”这一独木桥，而是有着自己的定位，Caddy就是这样一个开源Web Server。
Caddy的作者Matt Holt在caddy官网以及FAQ中对caddy的目标阐释如下： 其他Web Server为Web而设计，Caddy为human设计。功能定位上，与经常充当最前端反向代理的nginx不同，caddy致力于成为一个易用的静态 文件Web Server。可以看出Caddy主打易用性，使用配置简单。并且得益于Go的跨平台特性，caddy很容易的支持了三大主流平台:Windows、 Linux、Mac。在Caddy开发者文档中，我们可以看到caddy还可以在Android(linux arm)上运行。caddy目前版本为0.7.1，还不稳定，且后续版本可能变化较大，甚至与前期版本不兼容，因此作者目前不推荐caddy在生产环境被 重度使用。
关注caddy，是因为caddy填补了go在通用web server这块的空白(也许有其他，但我还不知道)，同时Web server in go也“响应”了近期Golang去C化的趋势(Go 1.5中C is gone！)，即便caddy作者提到caddy的目标并非如nginx那样。但未来谁知道呢？一旦Go性能足够高时，一旦caddy足够稳定时，自然而 然的就会有人将其用在某些应用的生产环境中替代nginx或apache2了。一套全Go的系统，在部署、运维方面也是有优势的。

# 一、安装和运行caddy
和诸多go应用一样，我们可以直接从caddy的github.com releases页中找到最新发布版(目前是0.7.1)的二进制包。这里使用的是caddy_darwin_amd64.zip。
下载解压后，进入目录，直接执行./caddy即可将caddy运行起来。

```
$caddy
0.0.0.0:2015
```

在浏览器里访问localhost:2015，页面上没有预期显示的类似"caddy works!”之类的默认Welcome页面，而是“404 Not Found"。虽然这说明caddy已经work了，但没有一个default welcome page毕竟对于caddy beginer来说并不友好。这里已经向作者提了一个sugguestion issue。

# 二、caddy原理
Go的net/http标准库已经提供了http server的实现，大多数场合这个http server都能满足你的需要，无论是功能还是性能。Caddy实质上也是一个Go web app，它也import net/http，嵌入*http.Server，并通过handler的ServeHTTP方法为每个请求提供服务。caddy使用 http.FileServer作为处理 静态文件的基础。caddy的诱人之处在于其middleware，将诸多middleware串成一个middleware chain以提供了灵活的web服务。另外caddy中的middleware还可以独立于caddy之外使用。
caddy从当前目录的Caddyfile（默认）文件中读取配置，当然你也可以通过-conf指定配置文件路径。Caddyfile的配置格式 的确非常easy，这也符合caddy的目标。
Caddyfile总是以站点的Addr开始的。

## 单一站点的Caddyfile样例如下：

//Caddyfile
localhost:2015
gzip
log ./2015.log


Caddy也支持配置多个站点,类似virtualhost的 配置(80端口多路复用)：

//Caddyfile
foo.com:80 {
    log ./foo.log
    gzip
}
bar.com:80 {
    log ./bar.log
    gzip
}

为了实现风格上的统一，单一站点也最好配置为如下这种格式(代码内部称之为    Server Block)：


localhost:2015 {
    gzip
    log ./2015.log
}


这样Caddyfile的配置文件模板样式类似于下面这样：

host1:port {
    middleware1
    middleware2 {
        … …
    }
    … …
}
host2:port {
    middleware1
    middleware2 {
        … …
    }
    … …
}




关于***middleware***，在caddy文档中有较为详细的说明和例子。对于caddy这样一个年轻的开源项目而言，其文档还算是相对较全的，虽 然现在还不能和nginx、 apache比。
caddy中的middleware就是一个实现了middleware.Handler接口的struct，例如gzip这个 middleware:

## 重点
// middleware.go
type Middleware func(Handler) Handler
type Handler interface {
        ServeHTTP(http.ResponseWriter, *http.Request) (int, error)
}
// gzip/gzip.go
type Gzip struct {
    Next middleware.Handler
}
func (g Gzip) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
    if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
        return g.Next.ServeHTTP(w, r)
    }
    …. …
    gz := gzipResponseWriter{Writer: gzipWriter, ResponseWriter: w}
    // Any response in forward middleware will now be compressed
    status, err := g.Next.ServeHTTP(gz, r)
    … …
}
middleware.Handler的函数原型与http.Handler的不同，不能直接作为http.Server的Handler使用。caddy使用了下面这个idiomatic go pattern:
type appHandler func(http.ResponseWriter, *http.Request) (int, error)
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if status, err := fn(w, r); err != nil {
        http.Error(w, err.Error(), status)
    }
}


当然这个pattern有很多变种，但思路大致类似。一个middleware chain大致就是handler1(handler2(handler3))的调用传递。
前面说过caddy是基于http.FileServer的静态文件Web Server，FileServer总会作为middleware chain的最后一环，如果没有配置任何middleware，那你的server就是一个静态文件server。


# 三、caddy典型应用
##【静态文件Server】
caddy的最基础应用实际就是一个静态文件Server，底层由http.FileServer承载，当然caddy封装了http.FileServer，做了一些拦截处理，最后将w, r传递给http.ServeContent去处理文件数据。
第一次执行./caddy，实际上就启动了一个静态文件Server。但这个server不默认支持你navigate directory。如果你知道website root目录(如果没有指定root，则caddy执行的当前路径会作为website的root路径)下的文件名，比如foo.txt，你可以在浏览器 中输入：localhost:2015/foo.txt，caddy会执行正确的服务，浏览器也会显示foo.txt的全文。
对于静态文件Server，caddy支持在website的root路径下首先查找是否有如下四个文件：
//caddy/middleware/browse/browse.go
var IndexPages = []string{
    "index.html",
    "index.htm",
    "default.html",
    "default.htm",
}
如果查到有其中一个，则优先返回这个文件内容，这就是静态站点的首页。
如果要支持目录文件列表浏览，则需要为website配置browse middleware，这样对于无index file的目录，我们可以看到目录文件列表。
localhost:2015 {
    browse
}    
## 【反向代理】
caddy支持基本的反向代理功能。反向代理配置通过proxy middleware实现。
localhost:2015 {
    log ./2015.log
    proxy /foo localhost:9001
    proxy /bar localhost:9002
}
当你访问localhost:2015/foo时，实际上访问的是9001端口的服务程序；
当你访问localhost:2015/bar时，实际上访问的是9002端口的服务程序。
##【负载均衡】
Caddy支持负载均衡配置，并支持三种负载均衡算法：random（随机）、least_conn（最少连接）以及round_robin(轮询调度)。
负载均衡同样是通过proxy middleware实现的。
localhost:2015 {
    log ./2015.log
    proxy / localhost:9001 localhost:9003 {
        policy round_robin
    }
    proxy /bar localhost:9002 localhost:9004 {
        policy least_conn
    }
}
##【支持fastcgi代理】
caddy同样支持fastcgi代理，可以将请求通过fastcgi接口发送给后端的实现fastcgi的server。我们以一个"hello world"的php server为例。
mac os上自带了php-fpm，一个实现了fastcgi的php cgi进程管理器。caddy将请求转发给php-fpm监听的端口，后者会启动php-cgi解释器，解释index.php，并将结果返回给caddy。
mac os上的php-fpm默认没有随机启动。我们需要简单配置一下：
$mkdir phptest
$mkdir -p phptest/etc
$mkdir -p phptest/log
$cd phptest
$sudo cp /private/etc/php-fpm.conf.default ./etc
$cd ./etc
$sudo chown tony php-fpm.conf.default
$mv php-fpm.conf.default php-fpm.conf
编辑php-fpm.conf，保证下面两项是非注释状态的：
error_log = log/php-fpm.log
listen = 127.0.0.1:9000 
我们通过network socket进行fastcgi通信。
回到phptest目录下，执行:
php-fpm -p ~/test/go/caddy/phptest
执行后，php-fpm就会转入后台执行了。
接下来我们来配置Caddyfile：
localhost:2015 {
    fastcgi / 127.0.0.1:9000 php
    log ./2015.log
}
这里配置的含义是：将全部请求转发到9000端口，这里的php是一个preset（预配置集合），相当于：
ext   .php
split .php
index index.php
我们在phptest目录下创建一个index.php文件，内容如下：
<?php echo "Hello World\n"; ?>
好了，现在启动caddy，并使用浏览器访问localhost:2015试试。你会看到"Hello World"呈现在浏览器中。
##【git push发布】 ？？ 这个蛮感兴趣的
对于一些静态站点，caddy支持git directive，实现在server启动以及运行时定期git pull你的项目库，将最新更新pull到server上。
caddy文档中给出两个例子：
第一个是一个php站点，定期pull项目库，实现server更新：
git git@github.com:user/myphpsite {
    key /home/user/.ssh/id_rsa
}
fastcgi / 127.0.0.1:9000 php
第二个是一个hugo支撑的静态站点，每次pull后，执行hugo命令生成新的静态页面：
git github.com/user/site {
    path  ../
    then  hugo –destination=/home/user/hugosite/public
}
注意：git directive并非middleware，而是一个单独的goroutine实现的。

#四、小结
caddy的功能不局限于上面的几个例子，上面只是几个最为常见的场景而已。caddy目前还很年轻，应用不多，但知名golang网站 gopheracademy.com（GopherCon组织方）是由Caddy support的。caddy还在积极进化，有兴趣的Gopher可持续关注。