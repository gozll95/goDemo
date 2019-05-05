# 一、使用姿势
beego的使用姿势如下:

### First you must import it
```
import (
	"github.com/astaxie/beego/session"
)
```

### Then in you web app init the global session manager

```
var globalSessions *session.Manager
```

### Use memory as provider:

```
  func init() {
  	globalSessions, _ = session.NewManager("memory", `{"cookieName":"gosessionid","gclifetime":3600}`)
  	go globalSessions.GC()
  }
```

### Use file as provider, the last param is the path where you want file to be stored:
```
  func init() {
  	globalSessions, _ = session.NewManager("file",`{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"./tmp"}`)
  	go globalSessions.GC()
  }
```

### Use Redis as provider, the last param is the Redis conn address,poolsize,password:

```
  func init() {
  	globalSessions, _ = session.NewManager("redis", `{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"127.0.0.1:6379,100,astaxie"}`)
  	go globalSessions.GC()
  }
```

### Use MySQL as provider, the last param is the DSN, learn more from mysql:
```
  func init() {
  	globalSessions, _ = session.NewManager(
  		"mysql", `{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"username:password@protocol(address)/dbname?param=value"}`)
  	go globalSessions.GC()
  }
```


### Use Cookie as provider:
```
  func init() {
  	globalSessions, _ = session.NewManager(
  		"cookie", `{"cookieName":"gosessionid","enableSetCookie":false,"gclifetime":3600,"ProviderConfig":"{\"cookieName\":\"gosessionid\",\"securityKey\":\"beegocookiehashkey\"}"}`)
  	go globalSessions.GC()
  }
```

### Finally in the handlerfunc you can use it like this

```
func login(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	username := sess.Get("username")
	fmt.Println(username)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		fmt.Println("username:", r.Form["username"])
		sess.Set("username", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}
```


# 二、分析
beego的设计分为

| (w,r)     |   
| :-------- | 
| Manager[负责对上暴露`Set`/`Get`接口以及从`request里获取cookie`以及将最后的`session落地`(即落到session store对应的`存储介质`中)] |
| Provider[负责对上暴露`Set`/`Get`/`GC`等方法以及对下进行`按照更新时间排序Session`] |   
| SessionStore[负责`Set`/`Get`/`GC`/`落地`等]|  


# 三、接口
### 1.Store

```
// Store contains all data for one session process with specific id.
type Store interface {
	Set(key, value interface{}) error     //set session value
	Get(key interface{}) interface{}      //get session value
	Delete(key interface{}) error         //delete session value
	SessionID() string                    //back current sessionID
	SessionRelease(w http.ResponseWriter) // release the resource & save data to provider & return the data
	Flush() error                         //delete all data
}
```

### 2.Privoder

var provides = make(map[string]Provider)


// Provider contains global session methods and saved SessionStores.
// it can operate a SessionStore by its id.
```
type Provider interface {
	SessionInit(gclifetime int64, config string) error
	SessionRead(sid string) (Store, error)
	SessionExist(sid string) bool
	SessionRegenerate(oldsid, sid string) (Store, error)
	SessionDestroy(sid string) error
	SessionAll() int //get all active session
	SessionGC()
}
```

### 3.Manager

```
// ManagerConfig define the session config
type ManagerConfig struct {
	CookieName              string `json:"cookieName"`
	EnableSetCookie         bool   `json:"enableSetCookie,omitempty"`
	Gclifetime              int64  `json:"gclifetime"`
	Maxlifetime             int64  `json:"maxLifetime"`
	DisableHTTPOnly         bool   `json:"disableHTTPOnly"`
	Secure                  bool   `json:"secure"`
	CookieLifeTime          int    `json:"cookieLifeTime"`
	ProviderConfig          string `json:"providerConfig"`
	Domain                  string `json:"domain"`
	SessionIDLength         int64  `json:"sessionIDLength"`
	EnableSidInHTTPHeader   bool   `json:"EnableSidInHTTPHeader"`
	SessionNameInHTTPHeader string `json:"SessionNameInHTTPHeader"`
	EnableSidInURLQuery     bool   `json:"EnableSidInURLQuery"`
}

// Manager contains Provider and its configuration.
type Manager struct {
	provider Provider
	config   *ManagerConfig
}
```


# 四、分析Manager
Manager=Provider+configuration

Providers 是一组以不同name为标识的Provider,比如"file" Provider+"redis" Provider + "mem" Provoder


```
NewManager("file",xxx-config)--> Manager  
--->go Manager.GC()[
    call Manager的Provider.GC

    time.AfterFunc(time.Duration(manager.config.Gclifetime)*time.Second, func() { manager.GC() })
]
--->SesstionStart(w,r) (sessionStore)[
    如果request里可以拿到session Id,则通过对应的session Id返回对应的session Store
    如果request里拿不到session Id,则随机生成一个session Id并且返回这个session id对应的session Store
]
```
---> defer sess.SessionRelease(w)[
    落地落地落地,比如file就刷到文件里,redis就刷到redis里,mem就不用管了哦
]
---> sessionStore.Get("xx")
---> sessionStore.Set("xx","tt")
> `**备注：**这里从request里拿cookie是从Cookie里拿->传入参数里-->header里***`。


## 五、以mem为例

```
// MemSessionStore memory session store.
// it saved sessions in a map in memory.
type MemSessionStore struct {
	sid          string                      //session id
	timeAccessed time.Time                   //last access time
	value        map[interface{}]interface{} //session store
	lock         sync.RWMutex
}
```

SessionStore的Get/Set/Del等同于对map操作
> `**备注：这里记timeAccessed time是为了后续GC,这里类似Mgo的update time**`。

```
type MemProvider struct {
	lock        sync.RWMutex             // locker
	sessions    map[string]*list.Element // map in memory
	list        *list.List               // for gc
	maxlifetime int64
	savePath    string
}
```
> `**备注：这里&list.Element对应session Store, sessions是按照sid区分的, list.List是按照时间排序的`。


```
func (pder *MemProvider) SessionRead(sid string) (Store, error) 
---> 如果sid对应的list.Element有[
    则取list.Value,
    // 更新时间
    go pder.SessionUpdate(sid)[
        list.Value的timeAccessed=time.Now()
        pder.list.MoveToFront(element)
    ]
]
```


> `GC这段直接看吧,list.List用的蛮好的`。
```
// SessionGC clean expired session stores in memory session
func (pder *MemProvider) SessionGC() {
	pder.lock.RLock()
	for {
		element := pder.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*MemSessionStore).timeAccessed.Unix() + pder.maxlifetime) < time.Now().Unix() {
			pder.lock.RUnlock()
			pder.lock.Lock()
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*MemSessionStore).sid)
			pder.lock.Unlock()
			pder.lock.RLock()
		} else {
			break
		}
	}
	pder.lock.RUnlock()
}


```

# 优质代码:

```
//!--学会灵活使用time.AfterFunc()
// GC Start session gc process.
// it can do gc in times after gc lifetime.
func (manager *Manager) GC() {
	manager.provider.SessionGC()
	time.AfterFunc(time.Duration(manager.config.Gclifetime)*time.Second, func() { manager.GC() })
}


// SessionGC clean expired session stores in memory session
func (pder *MemProvider) SessionGC() {
	pder.lock.RLock()
	for {
		element := pder.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*MemSessionStore).timeAccessed.Unix() + pder.maxlifetime) < time.Now().Unix() {
			pder.lock.RUnlock()
			pder.lock.Lock()
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*MemSessionStore).sid)
			pder.lock.Unlock()
			pder.lock.RLock()
		} else {
			break
		}
	}
	pder.lock.RUnlock()
}

```



