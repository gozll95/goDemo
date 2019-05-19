func (client *Client) Add(request *AddRequest) (response *AddResponse, err error) {
	response = CreateAddResponse()
	err = client.DoAction(request, response)
	return
}


## 异步
func (client *Client) AddWithChan(request *AddRequest) (<-chan *AddResponse, <-chan error) {
	responseChan := make(chan *AddResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.Add(request)
		responseChan <- response
		errChan <- err
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

GET     = "GET"
	PUT     = "PUT"
	POST    = "POST"
	DELETE  = "DELETE"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"

	Json = "application/json"
	Xml  = "application/xml"
	Raw  = "application/octet-stream"
	Form = "application/x-www-form-urlencoded"

	Header = "Header"
	Query  = "Query"
	Body   = "Body"
	Path   = "Path"






 var hookDo = func(fn func(req *http.Request) (*http.Response, error)) func(req *http.Request) (*http.Response, error) {
	 return fn
 }
 
httpResponse, err = hookDo(client.httpClient.Do)(httpRequest)