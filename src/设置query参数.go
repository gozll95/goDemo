// Results is an ordered list of search results.
type Results []Result

// A Result contains the title and URL of a search result.
type Result struct {
	Title, URL string
}


func handleSearch(w http.ResponseWriter, req *http.Request) 
	- //从query params里获取timeout参数
	- timeout, err := time.ParseDuration(req.FormValue("timeout"))
	- //生成带有timeout的ctx
	- ctx, cancel = context.WithTimeout(context.Background(), timeout)
	- defer cancel() //这个要注意
	- //从query params里获取q参数
	- query := req.FormValue("q")
	- //获取userIP
	- userIP, err := userip.FromRequest(req) 
	- ctx = userip.NewContext(ctx, userIP)
			- func NewContext(ctx context.Context, userIP net.IP) context.Context 
				- return context.WithValue(ctx, userIPKey, userIP)
	- results, err := google.Search(ctx, query)
			- func Search(ctx context.Context, query string) (Results, error) 
				- //+学会这种request的方式:先req好再set query!
				- req, err := http.NewRequest("GET", "https://developers.google.com/custom-search/", nil)
				- q := req.URL.Query()
				- q.Set("q", query)
				- //-!
				- //set userIp
				- userIP, ok := userip.FromContext(ctx)
						- userIP, ok := ctx.Value(userIPKey).(net.IP)
				- q.Set("userip", userIP.String())
				- //+!
				- req.URL.RawQuery = q.Encode()
				- //-!
				- httpDo(
					ctx,
					req,
					func(resp *http.Response,err error)error{
						if err!=nil{
							return err
						}	
						defer resp.Body.Close()
					
						//parse the JSON search result
						...

						result=... 

						return nil
					}
				)
						- func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error
							- tr := &http.Transport{}
							- client := &http.Client{Transport: tr}
							- c := make(chan error, 1)
							- go func
								- c<-f(client.Do(req))
							- select case <-ctx.Done:
								- tr.CancelRequest(req)
								- <-c 
								- return ctx.Err()
							- select case err:=<-c:
								- return err



	- 模板...
		
