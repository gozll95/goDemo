package godis

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/alecthomas/log4go"
	"github.com/garyburd/redigo/redis"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	//连接zk默认过期时间
	DEFAULT_ZK_CONNECT_TIMEOUT = 100

	//默认最大空间连接数
	DEFAULT_MAX_IDLE = 10

	//默认最大连接数
	DEFAULT_MAX_ACTIVE = 100
)

type proxyInfo struct {
	Addr  string `json:"addr"`
	State string `json:"state"`
}

type idleConn struct {
	c redis.Conn
	t time.Time
}

type GodisPool struct {
	//zk的接点
	ZkDir string

	//测试方法
	TestOnBorrow func(c redis.Conn, t time.Time) error

	//zk服务器连接地址
	ZkServerList []string

	//最大空闲连接
	MaxIdle int

	//最大连接数
	MaxActive int

	//超时时间
	IdleTimeout time.Duration

	//redis服务的连接池，zk接点数据里的addr，只存在线（online）的addr
	//pools []redis.Conn
	pools []string

	ZkConnTimeout time.Duration

	//zk连接的实例
	zkC zk.Conn

	//请求等待
	Wait bool

	//空闲连接池
	idle list.List

	//当前连接数
	active  int
	mutex   sync.Mutex
	cond    *sync.Cond
	nextIdx int
}

var (
	nowFunc = time.Now

	//当前zk的列表(包括在线和非在线的)，用于新服务发现后监听新接点的数据变化，
	zkMap = make(map[string]string)

	//日志开关
	godisLogHandle = false
)

//初始化方法
func (gp *GodisPool) InitPool() {

	if gp.MaxIdle <= 0 {
		gp.MaxIdle = DEFAULT_MAX_IDLE
	}
	if gp.MaxActive <= 0 {
		gp.MaxActive = DEFAULT_MAX_ACTIVE
	}

	gp.initZK()

	//设置连接池
	rsE := gp.resetPools()
	if rsE != nil {
		fmt.Println("set pools err:", rsE.Error())
		return
	}
	//启用接点监听，动态更新连接池
	go gp.poolWatcher()
}

//初始化zookeeper
func (gp *GodisPool) initZK() {
	zkConn, _, err := zk.Connect(gp.ZkServerList, gp.ZkConnTimeout)

	if err != nil {
		if godisLogHandle {
			log4go.Error("Failed to connect to zookeeper: %+v", err)
			log4go.Error("after 100 millisecond reconnect to zk...")
		}
	} else {
		gp.zkC = *zkConn
	}

}

//重设连接池
/**
 * 获取当前的zk接点，将在线的服务（State为online）创建redis连接
 * 放到连接池
 * 对于新加进的zk接点设置wather，监听此接点的数据变化
 */
func (gp *GodisPool) resetPools() error {
	gp.pools = []string{}

	//定义临时map
	tmpMap := make(map[string]string)
	for tk, tv := range zkMap {
		tmpMap[tk] = tv
	}

	proxys, _, err := gp.zkC.Children(gp.ZkDir)
	if err != nil {
		fmt.Println("connect to zookeeper get children err:", err)
		return err
	}

	//var _conn redis.Conn
	//var _err error

	for _, child := range proxys {
		connData, _, err := gp.zkC.Get(gp.ZkDir + "/" + child)

		if err != nil {
			continue
		}
		var p proxyInfo
		Uerr := json.Unmarshal(connData, &p)

		if Uerr != nil {
			fmt.Println(Uerr.Error())
		}

		/*_conn, _err = redis.Dial("tcp", p.Addr)

		if _err != nil {
			log4go.Error("Create redis connection err: %s", _err.Error())
			continue
		}*/

		if p.State == "online" {
			gp.pools = append(gp.pools, p.Addr)
		}

		_, gRs := tmpMap[child]
		if !gRs {
			//a new node
			zkMap[child] = p.Addr
			go gp.childWatcher(child)
		}

		delete(tmpMap, child)
	}

	//删除zkmap里无用的接点
	for tk, _ := range tmpMap {
		delete(zkMap, tk)
	}

	log4go.Info("new pool len----->", len(gp.pools))
	log4go.Info("zkmap------", zkMap)

	return nil
}

//zk接点监听，在监听到子接点的数量变化后，更新连接池
func (gp *GodisPool) poolWatcher() {
	log4go.Info("start to listen children node change ...")
	for {
		_, _, evtC, err := gp.zkC.ChildrenW(gp.ZkDir)

		if err != nil {
			log4go.Error("watch zkNode %s err: %s", gp.ZkDir, err.Error())
			return
		}

		evt := <-evtC

		if evt.Type == zk.EventSession {
			if evt.State == zk.StateConnected {
			}

			if evt.State == zk.StateExpired {
				gp.zkC.Close()
				log4go.Info("Zookeeper session expired, reconnecting...")
				gp.initZK()

			}
		}

		switch evt.Type {
		case zk.EventNodeChildrenChanged, zk.EventNodeDataChanged:
			log4go.Info("zk children node change reset pools....")
			time.Sleep(100 * time.Millisecond)
			gp.resetPools()

		}
	}
}

//监听子接点的数据变化，
// 子接点服务的状态变化后（比如某个codis下线，状态由online改为offline），更新连接池
func (gp *GodisPool) childWatcher(childPath string) {
	for {
		_, _, evtC, err := gp.zkC.GetW(gp.ZkDir + "/" + childPath)

		//在接点删除时，会报此错误
		if err != nil {
			log4go.Error("wath zkChildNode err: ", err.Error())
			return
		}

		evt := <-evtC

		switch evt.Type {
		case zk.EventNodeDataChanged:
			gp.resetPools()
		case zk.EventNodeDeleted:
			delete(zkMap, childPath)
			return
		}
	}
}

//overWrite get method
func (gp *GodisPool) Get() redis.Conn {
	c, err := gp.get()
	if err != nil {
		return errorConnection{err}
	}
	return &pooledConnection{p: gp, c: c}
}

/**
 *在请求连接池里拿一个连接
 * 超时判断，将超时的连接删除
 * 在idle里拿连接资源，如果没有空间资源 ， 从redis连接池拿一条（此处加平均原则），使用后放入idle
 *
 */
func (gp *GodisPool) get() (redis.Conn, error) {
	gp.mutex.Lock()
	if timeout := gp.IdleTimeout; timeout > 0 {
		for i, n := 0, gp.idle.Len(); i < n; i++ {
			e := gp.idle.Back()

			if e == nil {
				break
			}
			ic := e.Value.(idleConn)

			if ic.t.Add(timeout).After(nowFunc()) {
				break
			}

			gp.idle.Remove(e)
			gp.release()
			gp.mutex.Unlock()
			ic.c.Close()
			gp.mutex.Lock()

		}
	}
	log4go.Info("active=-----------------", gp.active)
	log4go.Info("idle=---------------->", gp.idle.Len())
	for {
		//// Get idle connection.
		for i, n := 0, gp.idle.Len(); i < n; i++ {
			e := gp.idle.Front()
			if e == nil {
				break
			}

			ic := e.Value.(idleConn)
			gp.idle.Remove(e)
			test := gp.TestOnBorrow
			gp.mutex.Unlock()
			if test == nil || test(ic.c, ic.t) == nil {
				fmt.Println("---get----from----idle---")
				return ic.c, nil
			}

			ic.c.Close()
			gp.mutex.Lock()
			gp.release()
		}

		// Dial new connection if under limit.
		if gp.MaxActive == 0 || gp.active < gp.MaxActive {
			if len(gp.pools) == 0 {
				rsE := gp.resetPools()
				if rsE != nil {
					gp.mutex.Unlock()
					return nil, rsE
				}
			}

			gp.nextIdx += 1
			if gp.nextIdx >= len(gp.pools) {
				gp.nextIdx = 0
			}
			if len(gp.pools) == 0 {
				gp.mutex.Unlock()
				err := errors.New("Proxy list empty")
				log4go.Error(err)
				return nil, err
			} else {
				fmt.Println("---get----from----new---")
				c := gp.pools[gp.nextIdx]
				gp.active += 1
				_conn, _err := redis.Dial("tcp", c)
				gp.mutex.Unlock()
				if _err != nil {
					log4go.Error("Create redis connection err: %s", _err.Error())
					return nil, _err
				}
				test := gp.TestOnBorrow
				if test == nil || test(_conn, nowFunc()) == nil {
					return _conn, nil
				}
				_conn = nil
				gp.mutex.Lock()
				gp.release()
				gp.mutex.Unlock()

				return _conn, errors.New("Create redis connection err")
			}
		}

		if !gp.Wait {
			gp.mutex.Unlock()
			return nil, errors.New("connect pool exhausted")
		}

		if gp.cond == nil {
			gp.cond = sync.NewCond(&gp.mutex)
		}

		gp.cond.Wait()
	}
}

//释放一个当前存活数
func (gp *GodisPool) release() {
	gp.active -= 1
	if gp.cond != nil {
		gp.cond.Signal()
	}
}

//将一个连接放回idle
//如果空闲数已到设置的值，将此连接关闭
func (gp *GodisPool) put(c redis.Conn) error {
	err := c.Err()
	gp.mutex.Lock()

	if err == nil {
		if gp.idle.Len() < gp.MaxIdle {
			gp.idle.PushFront(idleConn{c: c, t: nowFunc()})

			fmt.Println("add-to-idle=---------------->", gp.idle.Len())
			if gp.cond != nil {
				gp.cond.Signal()
			}
			gp.release()
			gp.mutex.Unlock()
			return nil
		}
	} else {
		fmt.Println("----errr===", err)
	}
	gp.release()
	gp.mutex.Unlock()

	return c.Close()
}

//Get方法返回的结构体重写
type pooledConnection struct {
	p     *GodisPool
	c     redis.Conn
	state int
}

func (pc *pooledConnection) Close() error {
	c := pc.c

	pc.p.put(c)

	return nil
}

func (pc *pooledConnection) Err() error {
	return pc.c.Err()
}

func (pc *pooledConnection) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	ci := LookupCommandInfo(commandName)
	pc.state = (pc.state | ci.Set) &^ ci.Clear
	return pc.c.Do(commandName, args...)
}

func (pc *pooledConnection) Send(commandName string, args ...interface{}) error {
	ci := LookupCommandInfo(commandName)
	pc.state = (pc.state | ci.Set) &^ ci.Clear
	return pc.c.Send(commandName, args...)
}

func (pc *pooledConnection) Flush() error {
	return pc.c.Flush()
}

func (pc *pooledConnection) Receive() (reply interface{}, err error) {
	return pc.c.Receive()
}

const (
	WatchState = 1 << iota
	MultiState
	SubscribeState
	MonitorState
)

type CommandInfo struct {
	Set, Clear int
}

var commandInfos = map[string]CommandInfo{
	"WATCH":      {Set: WatchState},
	"UNWATCH":    {Clear: WatchState},
	"MULTI":      {Set: MultiState},
	"EXEC":       {Clear: WatchState | MultiState},
	"DISCARD":    {Clear: WatchState | MultiState},
	"PSUBSCRIBE": {Set: SubscribeState},
	"SUBSCRIBE":  {Set: SubscribeState},
	"MONITOR":    {Set: MonitorState},
}

//初始化配置
func init() {
	for n, ci := range commandInfos {
		commandInfos[strings.ToLower(n)] = ci
	}

	log4go.LoadConfiguration("log4g.xml")

	godisLogHandle = true

}

func LookupCommandInfo(commandName string) CommandInfo {
	if ci, ok := commandInfos[commandName]; ok {
		return ci
	}
	return commandInfos[strings.ToUpper(commandName)]
}

type errorConnection struct{ err error }

func (ec errorConnection) Do(string, ...interface{}) (interface{}, error) { return nil, ec.err }
func (ec errorConnection) Send(string, ...interface{}) error              { return ec.err }
func (ec errorConnection) Err() error                                     { return ec.err }
func (ec errorConnection) Close() error                                   { return ec.err }
func (ec errorConnection) Flush() error                                   { return ec.err }
func (ec errorConnection) Receive() (interface{}, error)                  { return nil, ec.err }
