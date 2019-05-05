session:
全局session管理器
	保证sessionId的全局唯一性
	为每个用户关联一个session
session的存储
session过期处理

type Manger struct{
	cookieName string 
	lock sync.Mutex
	provider Provider
	maxlifetime int64
}

func NewManager(provideName, cookieName string, maxlifetime int64) (*Manager, error) {
    provider, ok := provides[provideName]
    if !ok {
        return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
    }
    return &Manager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}

//生成全局唯一的sessionId
func(manager *Manager)sessionId()string{
	 b := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, b); err != nil {
        return ""
    }
    return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
    manager.lock.Lock()
    defer manager.lock.Unlock()
    cookie, err := r.Cookie(manager.cookieName)
    if err != nil || cookie.Value == "" {
		//生成全局唯一的sessionId
        sid := manager.sessionId()
		//调用provider初始化一个sessionId
        session, _ = manager.provider.SessionInit(sid)
        cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
        http.SetCookie(w, &cookie)
    } else {
        sid, _ := url.QueryUnescape(cookie.Value)
		//调用provider去read这个session
        session, _ = manager.provider.SessionRead(sid)
    }
    return
}

session的操作
createtime := sess.Get("createtime")
sess.Set("createtime", time.Now().Unix())
globalSessions.SessionDestroy(w, r)


//用户退出,销毁session
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request){
	... 
	//调用provide进行
	manager.provider.SessionDestroy(cookie.Value)
}

//GC根据过期时间
func(manager *Manager)GC(){
	manager.lock.Lock()
	defer manager.lock.Unlock()
	//调用provider的SessionGC
	 manager.provider.SessionGC(manager.maxlifetime)
}
/*
SessionInit函数实现Session的初始化，操作成功则返回此新的Session变量
SSessionRead函数返回sid所代表的Session变量，如果不存在，那么将以sid为参数调用SessionInit函数创建并返回一个新的Session变量
SessionDestroy函数用来销毁sid对应的Session变量
SessionGC根据maxLifeTime来删除过期的数据
*/
type Provider interface{
	SessionInit(sid string)(Session,error) //new
	SessionRead(sid string)(Session,error)
	SessionDestroy(sid string)
	SessionGC(maxLifeTime int64) 
}

type Session interface {
    Set(key, value interface{}) error //set session value
    Get(key interface{}) interface{}  //get session value
    Delete(key interface{}) error     //delete session value
    SessionID() string                //back current sessionID
}


var provides=make(map[string]Provide)

//注册provider
func Register(name string, provide Provide) {
    if driver == nil {
        panic("session: Register provide is nil")
    }
    if _, dup := provides[name]; dup {
        panic("session: Register called twice for provide " + name)
    }
    provides[name] = provide
}



/*
session存储
*/
//session model
type SessionStore struct {
	sid          string                      //session id唯一标示
	timeAccessed time.Time                   //最后访问时间
	value        map[interface{}]interface{} //session里面存储的值
}


//provider model
type Provider struct{
	lock sync.Mutex
	sessions map[string]*list.Element //存的是list的index
	list *list.List //维护一个List列表,用来维护一个时间顺序
}

//初始化一个session
func(pder *Provider)SessionInit(sid string)(session.Session,error){
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	//压入list
	element:=pder.list.PushBack(newsess)
	pder.sessions[sid] = element
	return newsess,nil

}

//读取一个session
func (pder *Provider) SessionRead(sid string) (session.Session, error) {
	if element, ok := pder.sessions[sid]; ok {
		//这里用了list,element.Value
		return element.Value.(*SessionStore), nil
	} else {
		sess, err := pder.SessionInit(sid)
		return sess, err
	}
	return nil, nil
}

//销毁一个session
func (pder *Provider) SessionDestroy(sid string) error {
	if element, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		//使用list.Remove
		pder.list.Remove(element)
		return nil
	}
	return nil
}
//GC
func (pder *Provider) SessionGC(maxlifetime int64) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	for {
		//拿最后一个element
		element := pder.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

//update就把这个session挪到list的最前面
func (pder *Provider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		//挪到最前面
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}


func init(){
	pder.sessions=make(map[string]*list.Element,0)
	session.Register("memory",pder)
}

func (st *SessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) Get(key interface{}) interface{} {
	pder.SessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (st *SessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) SessionID() string {
	return st.sid
}

/*
有一个全局manager
有一个provider提供接口
有一个session

在memory例子里
provider里有list.list的数据类型

type Provider struct{
	lock sync.Mutex
	sessions map[string]*list.Element //存的是list的index
	list *list.List //维护一个List列表,用来维护一个时间顺序
}
*/
