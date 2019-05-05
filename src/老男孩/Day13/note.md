# for循环里使用 
for k,v:=range xx{
    func(){
        xxx
        defer xxx
        xxx
    }()
}

# 
type AppConfig struct{
    host string
    ip string
}

type AppConfigMgr struct{
    atomic.Value
}


var config Value // holds current server configuration
// Create initial config value and store into config.
config.Store(loadConfig())
go func(){
    // Reload config every 10s
    // and update config value with the new version.
    for{
        time.Sleep(10*time.Second)
        config.Store(loadConfig())
    }
}()
// Create worker goroutines that handle incoming requests
// using the latest config value
for i:=0;i<10;i++{
    go func(){
        for r:=range requests(){
            c:=config.Load()
            _,_=r.c
        }
    }()
}