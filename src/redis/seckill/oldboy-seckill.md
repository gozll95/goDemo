# 基本
- redis
- etcd
- beego

# 基本思路
- 利用etcd的watcher机制更新配置
- 利用redis的list的BRPOP/LUSH等来存储request/result
- 很好的利用channel来达到request和result的统一

> request struct{
    chan *result
}

# 大致流程
## product方面
product的信息存储在`etcd`里,由etcd的`watcher机制`来及时更新product信息

product
```
type Activity struct {
	ProductId int    `json:"product_id"`
	StartTime int64  `json:"start_time"`
	EndTime   int64	 `json:"end_time"`
	Count     int    `json:"count"`
	Status    int	 `json:"status"`
}
```

product的状态由于
```
activityMap map[int]*common.Activity
```
来保存

etcd中product对应的key的更新`activityMap`


##关于Seckill(API)
Recv-Q中存的resp,
Send-Q中存的req,

```
Seckill request(包含product id + user id)
--->
通过activityMap查看对应的product是否状态合法(在售卖/不为空/等等)
--->
map[v1_uid_productId]=request
--->
request chan <- request (这里有个timeout)
--->
<-request.result chan(这里有个timeout)
```

## 关于request chan
```

有个SendThread
range request chan
---->
get request
---->
LPUSH request to Send-Q
```

## 线程X
```
LPOP Send-Q
--->
get request
--->
通过各种算法
--->
generate result
--->
result chan <-result
```


## 关于result chan
```
range result chan
--->
get result
--->
LPUSH result to RECV-Q
```

```
BRPOP RECV-Q
--->
get result
--->
get request by result.uid+result.productId in map[v1_uid_productId]request
--->
request.result chan <- result
```

# 优秀代码
etc.go
```
package main


import (
	"context"
	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/clientv3"
	"time"
	"github.com/sherlockhua/goproject/seckill/common"
)

var etcdClient *clientv3.Client
var productChan chan string

func initEtcd(conf *common.SkillConf) (err error) {

	etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{conf.EtcdAddr},
		DialTimeout: 3 * time.Second,
	})

	if err != nil {
		logs.Warn("init etcd client failed, err:%v", err)
		return
	}

	productChan = make(chan string, 16)
	ctx, cancel := context.WithTimeout(context.Background(), 2 *time.Second)
	//   etcd_key /seckill/product/conf
	resp, err := etcdClient.Get(ctx, conf.EtcdProductKey)
	cancel()
	if err != nil {
		logs.Warn("get key %s failed, err:%v", conf.EtcdProductKey, err)
		return
	}

	for _, ev := range resp.Kvs {
		logs.Debug(" %q : %q\n",  ev.Key, ev.Value)
		productChan <- string(ev.Value)
	}
	
	go WatchEtcd(conf.EtcdProductKey)
	return
}

func WatchEtcd(key string) {
	
	for {
		rch := etcdClient.Watch(context.Background(), key)
		
		for resp := range rch {
			for _, ev := range resp.Events {
				logs.Debug("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				if ev.Type == clientv3.EventTypePut {
					productChan <- string(ev.Kv.Value)
				}
			}
		}
	}

}

func GetProductChan() chan string {
	return productChan
}
```
