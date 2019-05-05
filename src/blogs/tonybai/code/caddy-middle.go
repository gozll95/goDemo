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

// 可以看flysnow 接口型函数
// 接口型函数 转变x函数 实现y接口 前提:拥有相似的请求参数