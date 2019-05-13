func (client *Client) Add(request *AddRequest) (response *AddResponse, err error) {
	response = CreateAddResponse()
	err = client.DoAction(request, response)
	return
}

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

func (client *Client) AddWithCallback(request *AddRequest, callback func(response *AddResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *AddResponse
		var err error
		defer close(result)
		response, err = client.Add(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}