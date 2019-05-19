Go 1.8中增加对HTTP Server优雅退出(gracefullly exit)的支持，对应的新增方法为：
func (srv *Server) Shutdown(ctx context.Context) error
和server.Close在调用时瞬间关闭所有active的Listeners和所有状态为New、Active或idle的connections不同，server.Shutdown首先关闭所有active Listeners和所有处于idle状态的Connections，然后无限等待那些处于active状态的connection变为idle状态后，关闭它们并server退出。如果有一个connection依然处于active状态，那么server将一直block在那里。因此Shutdown接受一个context参数，调用者可以通过context传入一个Shutdown等待的超时时间。一旦超时，Shutdown将直接返回。对于仍然处理active状态的Connection，就任其自生自灭（通常是进程退出后，自动关闭）。通过Shutdown的源码我们也可以看出大致的原理：
// $GOROOT/src/net/http/server.go
... ...
func (srv *Server) Shutdown(ctx context.Context) error {
    atomic.AddInt32(&srv.inShutdown, 1)
    defer atomic.AddInt32(&srv.inShutdown, -1)

    srv.mu.Lock()
    lnerr := srv.closeListenersLocked()
    srv.closeDoneChanLocked()
    srv.mu.Unlock()

    ticker := time.NewTicker(shutdownPollInterval)
    defer ticker.Stop()
    for {
        if srv.closeIdleConns() {
            return lnerr
        }
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
        }
    }
}