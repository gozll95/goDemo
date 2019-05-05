// http://blog.csdn.net/defonds/article/details/51984683

//自定义session结构体
type Session struct {
	SessionID  string        `json:"sessionId" bson:"sessionId"`
	User       *User         `json:"-" bson:"user"`
	UserType   string        `json:"userType" bson:"userType"`
	NickName   string        `json:"nickName" bson:"nickName"`
	CreateTime time.Time     `json:"-" bson:"createTime"`
	UpdateTime time.Time     `json:"-" bson:"updateTime"`
	Expires    time.Time     `json:"-" bson:"expires"`
	Locale     string        `json:"-" bson:"locale"` // default is zh_CN
	Menus      []wmodel.Menu `json:"menus" bson:"menus"`
}

//session的保存
//使用 json.Marshal 将结构体 json 化之后保存到 redis：
/*
	【增】
	描述：插入一个 session 对象
	session 顶级 key，顶级 key 可以设置过期时间
	<[session]： 要插入的 session 对象
	>[error]：插入失败相关信息
*/
func (s *sessionService) SetSession(session *model.Session) error {
	// 从池里获取连接
	conn := pool.Get()
	if conn == nil {
		log.Errorf("redis connnection is nil")
		return errors.New("redis connnection is nil")
	}
	// 用完后将连接放回连接池
	defer conn.Close()
	// 将session转换成json数据，注意：转换后的value是一个byte数组
	value, err := json.Marshal(session)
	if err != nil {
		log.Errorf("json marshal err,%s", err)
		return err
	}
	log.Infof("send data[%s]", session.SessionID, value)
	_, err = conn.Do("SET", session.SessionID, value, "EX", sessionTimeOutInSeconds)
	if err != nil {
		return err
	}
	return nil
}



//golang的测试验证
func TestSetSession(t *testing.T) {
	s := &model.Session{
		SessionID: "20150421120000",
		UserType:  "admin",
		NickName:  "df",
	}
	err := SessionService.SetSession(s)
	if err != nil {
		t.Errorf("fail to add one session(%+v): %s", s, err)
		t.FailNow()
	}
}


//session 的删除

/*
	【删】
	描述： 删除一个 session 对象
	session 顶级 key，一般情况下 session 会在用户无操作 30 分钟后自行过期删除
	但用户登出操作可以提前对 session 进行删除，这就是本方法被调用的地方
	<[sessionID]： 要删除的 session 对象的 id
	>[error]：删除失败相关信息
*/
func (s *sessionService) DelSession(sessionID string) (err error) {
	// 从池里获取连接
	conn := pool.Get()
	if conn == nil {
		log.Errorf("redis connnection is nil")
		return errors.New("redis connnection is nil")
	}
	// 用完后将连接放回连接池
	defer conn.Close()
	log.Infof("move data[%s]", sessionID)
	_, err = conn.Do("DEL", sessionID)
	if err != nil {
		return err
	}
	return nil
}


//session 的获取

/*
	【查】
	描述： 查看并返回一个 session 实体
	session 顶级 key
	<[sessionID]： 要查看的 session 对象的 id
	>[error]：查看失败相关信息
*/
func (s *sessionService) GetSession(sessionID string) (session *model.Session, err error) {
	// 从池里获取连接
	conn := pool.Get()
	if conn == nil {
		log.Errorf("redis connnection is nil")
		return nil, errors.New("redis connnection is nil")
	}
	// 用完后将连接放回连接池
	defer conn.Close()
	log.Infof("exists data[%s]", sessionID)
	// 先查看该session是否存在
	var ifExists bool
	ifExists, err = SessionService.ExistsSession(sessionID)
	if err != nil {
		log.Errorf("fail to exists one session(%s): %s", sessionID, err)
		return nil, errors.New("session not exists, sessionID: " + sessionID)
	}
	if ifExists {
		// json数据在go中是[]byte类型，所以此处用redis.Bytes转换
		valueBytes, err2 := redis.Bytes(conn.Do("GET", sessionID))
		if err2 != nil {
			return nil, err2
		}
		//log.Infof("receive data[%s]:%s", sessionID, string(valueBytes))
		session = &model.Session{}
		err = json.Unmarshal(valueBytes, session)
		if err != nil {
			return nil, err
		}
		return session, nil
	} else {
		return nil, errors.New("session not exists, sessionID: " + sessionID)
	}

}

/*
redis 的作者为了保持简单的架构只允许我们对 top level 的 key 设置超时时间，次级 key(比如 Hash 里边的每个子 key)是不能设置超时时间的，所以我们单独使用了一个 index 为 3 的 redis 库专门存放 session，就是方便管理。另外，我们将每个 session 的 key(即 top level 的 key)有效期设置为半小时，有效期断定及处理托管 redis，避开了程序里对 session 超时机制管理的复杂性——特别是分布式环境。
另外，本文只提供了 session 的增、删、查操作，对于 session 的修改以及有效期推延操作笔者建议可以先删除再增加。
*/