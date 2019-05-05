import (
	"os"
	"time"
 
	"golang.org/x/time/rate"
 
	"github.com/op/go-logging"
 )
 
 var log = logging.MustGetLogger("example")
 
 // Example format string. Everything except the message has a custom color
 // which is dependent on the log level. Many fields have a custom output
 // formatting too, eg. the time returns the hour down to the milli second.
 var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
 )
 
 func main() {
 
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled, backend2Formatter)
 
	r := rate.Every(1)
	limit := rate.NewLimiter(r, 10)
	for {
		if limit.AllowN(time.Now(), 8) {
			log.Info("log:event happen")
		} else {
			log.Info("log:event not allow")
		}
 
	}
 
 }
 

 网络请求频率限制大多用的Allow，如tollbooth，服务器在对每一个请求相应之前，先从令牌池中获取令牌，如果没有令牌可用，则忽略或丢弃当前请求。
tollbooth 可以基于方法，IP等进行限制，基本实现方法就是把方法、ip作为一个key，然后对每一个key关联一个Limiter。

    // Map of limiters without TTL
    tokenBucketsNoTTL map[string]*rate.Limiter

    // Map of limiters with TTL
    tokenBucketsWithTTL *gocache.Cache
gocache.Cache 也是一个map，不过还实现了”有效期“功能，如果某个key超过了有效期就会从map中清除，这个机制实现的很巧妙，会在以后的文章中介绍。

然后对于每一个网络请求，提取出key,然后判断key对应的Limiter是否有可用的令牌,如下：

func (l *Limiter) limitReachedNoTokenBucketTTL(key string) bool {
    l.Lock()
    defer l.Unlock()

    if _, found := l.tokenBucketsNoTTL[key]; !found {
        l.tokenBucketsNoTTL[key] = rate.NewLimiter(rate.Every(l.TTL), int(l.Max))
    }

    return !l.tokenBucketsNoTTL[key].AllowN(time.Now(), 1)
}

...

func (l *Limiter) limitReachedWithCustomTokenBucketTTL(key string, tokenBucketTTL time.Duration) bool {
    l.Lock()
    defer l.Unlock()

    if _, found := l.tokenBucketsWithTTL.Get(key); !found {
        l.tokenBucketsWithTTL.Set(
            key,
            rate.NewLimiter(rate.Every(l.TTL), int(l.Max)),
            tokenBucketTTL,
        )
    }

    expiringMap, found := l.tokenBucketsWithTTL.Get(key)
    if !found {
        return false
    }

    return !expiringMap.(*rate.Limiter).AllowN(time.Now(), 1)
}

作者：kingeasternsun
链接：https://www.jianshu.com/p/4ce68a31a71d
來源：简书
简书著作权归作者所有，任何形式的转载都请联系作者获得授权并注明出处。