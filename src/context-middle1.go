package main 

func timeout(h CustomHandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		respChannel:=make(chan Response)

		ctx,cancel:=context.WithTimeout(r.Context(),time.Second*time.Duration(rtr.opt.Timeout.Timeout))
		defer cancel()

		r=r.WithContext(ctx)
		go func(){
			handlerResp,err:=h(w,r)
			if err!=nil{
				//do something with error
			}
			respChannel<-handlerResp
		}()
		select {
			// do something when context done happened
			case <-ctx.Done():
				w.WriteHeader(http.StatusRequestTimeout)
				w.Write([]byte(“Request time out”))
				Return
			case resp := <- respChannel:
				// do something with response
		}
	}
}

//client
req, err := http.NewRequest(“GET”, “someurlhere”, nil)
if err != nil {
	// do something with error
}

ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
defer cancel() // don’t forget to cancel the context, otherwise it will leaking

req = req.WithContext(ctx)
resp, err := http.DefaultClient.Do(req)


//db
db, err := sqlt.Open(“dsnmaster”, “dsnslave”)
if err != nil {
	// do something with error
}
 
var result []int
query := “SELECT id, number FROM something WHERE id = ?”

// cancel query that is running more than 2 secs
ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
defer cancel()

err = db.SelectContext(ctx, query, &result, 10)
if err != nil {
	// do something with error
}