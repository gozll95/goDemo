闭包
闭包可以是一个函数里边返回的另一个匿名函数,该匿名函数包含了定义在它外面的值。

func safehandler(fn http.HandlerFunc)http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Requset){
		defer func(){
			if e,ok:=recover().(error);ok{
				http.Error(w,err.Error(),http.StatusInternaleServerError)
				//或者输出自定义的50x错误页面
				//w.WriteHeader(http.StatusInternaleServerError)
				//renderHtml(w,"error",e)
				//logging
				log.Println("WARN: panic in %v - %v",fn,e)
				log.Println(string(debug.Stack()))
			}
		}()
		fn(w,r)
	}
}