package Helper  
  
import (  
    “crypto/rand”  
    “encoding/base64″  
    “io”  
    “net/http”  
    “net/url”  
    “strconv”  
    “sync”  
    “time”  
)  
  
/*Session会话管理*/  
type SessionMgr struct {  
    mCookieName  string       //客户端cookie名称  
    mLock        sync.RWMutex //互斥(保证线程安全)  
    mMaxLifeTime int64        //垃圾回收时间  
  
    mSessions map[string]*Session //保存session的指针[sessionID] = session  
}  
  
//创建会话管理器(cookieName:在浏览器中cookie的名字;maxLifeTime:最长生命周期)  
func NewSessionMgr(cookieName string, maxLifeTime int64) *SessionMgr {  
    mgr := &SessionMgr{mCookieName: cookieName, mMaxLifeTime: maxLifeTime, mSessions: make(map[string]*Session)}  
  
    //启动定时回收  
    go mgr.GC()  
  
    return mgr  
}  
  
//在开始页面登陆页面，开始Session  
func (mgr *SessionMgr) StartSession(w http.ResponseWriter, r *http.Request) string {  
    mgr.mLock.Lock()  
    defer mgr.mLock.Unlock()  
  
    //无论原来有没有，都重新创建一个新的session  
    newSessionID := url.QueryEscape(mgr.NewSessionID())  
  
    //存指针  
    var session *Session = &amp;Session{mSessionID: newSessionID, mLastTimeAccessed: time.Now(), mValues: make(map[interface{}]interface{})}  
    mgr.mSessions[newSessionID] = session  
    //让浏览器cookie设置过期时间  
    cookie := http.Cookie{Name: mgr.mCookieName, Value: newSessionID, Path: “/”, HttpOnly: true, MaxAge: int(mgr.mMaxLifeTime)}  
    http.SetCookie(w, &amp;cookie)  
  
    return newSessionID  
}  
  
//结束Session  
func (mgr *SessionMgr) EndSession(w http.ResponseWriter, r *http.Request) {  
    cookie, err := r.Cookie(mgr.mCookieName)  
    if err != nil || cookie.Value == “” {  
        return  
    } else {  
        mgr.mLock.Lock()  
        defer mgr.mLock.Unlock()  
  
        delete(mgr.mSessions, cookie.Value)  
  
        //让浏览器cookie立刻过期  
        expiration := time.Now()  
        cookie := http.Cookie{Name: mgr.mCookieName, Path: “/”, HttpOnly: true, Expires: expiration, MaxAge: -1}  
        http.SetCookie(w, &amp;cookie)  
    }  
}  
  
//结束session  
func (mgr *SessionMgr) EndSessionBy(sessionID string) {  
    mgr.mLock.Lock()  
    defer mgr.mLock.Unlock()  
  
    delete(mgr.mSessions, sessionID)  
}  
  
//设置session里面的值  
func (mgr *SessionMgr) SetSessionVal(sessionID string, key interface{}, value interface{}) {  
    mgr.mLock.Lock()  
    defer mgr.mLock.Unlock()  
  
    if session, ok := mgr.mSessions[sessionID]; ok {  
        session.mValues[key] = value  
    }  
}  
  
//得到session里面的值  
func (mgr *SessionMgr) GetSessionVal(sessionID string, key interface{}) (interface{}, bool) {  
    mgr.mLock.RLock()  
    defer mgr.mLock.RUnlock()  
  
    if session, ok := mgr.mSessions[sessionID]; ok {  
        if val, ok := session.mValues[key]; ok {  
            return val, ok  
        }  
    }  
  
    return nil, false  
}  
  
//得到sessionID列表  
func (mgr *SessionMgr) GetSessionIDList() []string {  
    mgr.mLock.RLock()  
    defer mgr.mLock.RUnlock()  
  
    sessionIDList := make([]string, 0)  
  
    for k, _ := range mgr.mSessions {  
        sessionIDList = append(sessionIDList, k)  
    }  
  
    return sessionIDList[0:len(sessionIDList)]  
}  
  
//判断Cookie的合法性（每进入一个页面都需要判断合法性）  
func (mgr *SessionMgr) CheckCookieValid(w http.ResponseWriter, r *http.Request) string {  
    var cookie, err = r.Cookie(mgr.mCookieName)  
  
    if cookie == nil ||  
        err != nil {  
        return “”  
    }  
  
    mgr.mLock.Lock()  
    defer mgr.mLock.Unlock()  
  
    sessionID := cookie.Value  
  
    if session, ok := mgr.mSessions[sessionID]; ok {  
        session.mLastTimeAccessed = time.Now() //判断合法性的同时，更新最后的访问时间  
        return sessionID  
    }  
  
    return “”  
}  
  
//更新最后访问时间  
func (mgr *SessionMgr) GetLastAccessTime(sessionID string) time.Time {  
    mgr.mLock.RLock()  
    defer mgr.mLock.RUnlock()  
  
    if session, ok := mgr.mSessions[sessionID]; ok {  
        return session.mLastTimeAccessed  
    }  
  
    return time.Now()  
}  
  
//GC回收  
func (mgr *SessionMgr) GC() {  
    mgr.mLock.Lock()  
    defer mgr.mLock.Unlock()  
  
    for sessionID, session := range mgr.mSessions {  
        //删除超过时限的session  
        if session.mLastTimeAccessed.Unix()+mgr.mMaxLifeTime &lt; time.Now().Unix() {  
            delete(mgr.mSessions, sessionID)  
        }  
    }  
  
    //定时回收  
    time.AfterFunc(time.Duration(mgr.mMaxLifeTime)*time.Second, func() { mgr.GC() })  
}  
  
//创建唯一ID  
func (mgr *SessionMgr) NewSessionID() string {  
    b := make([]byte, 32)  
    if _, err := io.ReadFull(rand.Reader, b); err != nil {  
        nano := time.Now().UnixNano() //微秒  
        return strconv.FormatInt(nano, 10)  
    }  
    return base64.URLEncoding.EncodeToString(b)  
}  
  
//——————————————————————————  
/*会话*/  
type Session struct {  
    mSessionID        string                      //唯一id  
    mLastTimeAccessed time.Time                   //最后访问时间  
    mValues           map[interface{}]interface{} //其它对应值(保存用户所对应的一些值，比如用户权限之类)  
}



var sessionMgr *Helper.SessionMgr = nil //session管理器


1 //创建session管理器,”TestCookieName”是浏览器中cookie的名字，3600是浏览器cookie的有效时间（秒）  
2 sessionMgr = Helper.NewSessionMgr(“TestCookieName”, 3600)



//处理登录  
func login(w http.ResponseWriter, r *http.Request) {  
    if r.Method == “GET” {  
        t, _ := template.ParseFiles(“web/MgrSvr_login.html”)  
        t.Execute(w, nil)  
  
    } else if r.Method == “POST” {  
        //请求的是登陆数据，那么执行登陆的逻辑判断  
        r.ParseForm()  
  
        //可以使用template.HTMLEscapeString()来避免用户进行js注入  
        username := r.FormValue(“username”)  
        password := r.FormValue(“password”)  
  
        //在数据库中得到对应数据  
        var userID int = 0  
  
        userRow := db.QueryRow(loginUserQuery, username, password)  
        userRow.Scan(&amp;userID)  
  
        //TODO:判断用户名和密码  
        if userID != 0 {  
            //创建客户端对应cookie以及在服务器中进行记录  
            var sessionID = sessionMgr.StartSession(w, r)  
  
            var loginUserInfo = UserInfo{ID: userID, UserName: username, Password: password, Alias: alias,  
                Desc: desc, ChannelAuth: channel_authority, IsSuperAdmin: is_super_admin, IsNewClientAuth: is_newclient_authority,  
                IsPayAuth: is_pay_authority, IsItemsAuth: is_itmes_authority, IsRealtimeAuth: is_realtime_authority,  
                IsPayCodeAuth: is_paycode_authority, IsUserAuth: is_user_authority, IsBgOpAuth: is_bgop_authority, IsHZRaidenNMMWeak: is_hz_raidenn_mmweak,  
                IsManualDataMgr: is_manual_data_mgr, IsManualDataQuery: is_manual_data_query}  
  
            //踢除重复登录的  
            var onlineSessionIDList = sessionMgr.GetSessionIDList()  
  
            for _, onlineSessionID := range onlineSessionIDList {  
                if userInfo, ok := sessionMgr.GetSessionVal(onlineSessionID, “UserInfo”); ok {  
                    if value, ok := userInfo.(UserInfo); ok {  
                        if value.ID == userID {  
                            sessionMgr.EndSessionBy(onlineSessionID)  
                        }  
                    }  
                }  
            }  
  
            //设置变量值  
            sessionMgr.SetSessionVal(sessionID, “UserInfo”, loginUserInfo)  
  
            //TODO 设置其它数据  
  
            //TODO 转向成功页面  
  
            return  
        }  
    }  
}



//处理退出  
func logout(w http.ResponseWriter, r *http.Request) {  
    sessionMgr.EndSession(w, r) //用户退出时删除对应session  
    http.Redirect(w, r, “/login”, http.StatusFound)  
    return  
}


func test_session_valid(w http.ResponseWriter, r *http.Request) {  
    var sessionID = sessionMgr.CheckCookieValid(w, r)  
  
    if sessionID == “” {  
        http.Redirect(w, r, “/login”, http.StatusFound)  
        return  
    }  
}
