func httpDo(ctx context.Context,req *http.Request,f func(*http.Response,error)error)error{
	tr:=&http.Transport{}
	client:=&http.Client(Transport:tr)
	c:=make(chan error,1)
	go func(){
		c<-f(client.Do(req))
	}()
	select{
	case <- ctx.Done():
		tr.CancelRequest(req)
		<-c 
		return ctx.Err()
	case err:=<-c:
		return err
	}
}

要么ctx被取消，要么request请求出错。