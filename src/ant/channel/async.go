

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


func (client *Client) AddAsyncTask(task func()) (err error) {
	if client.asyncTaskQueue != nil {
		client.asyncChanLock.RLock()
		defer client.asyncChanLock.RUnlock()
		if client.isRunning {
			client.asyncTaskQueue <- task
		}
	} else {
		err = errors.NewClientError(errors.AsyncFunctionNotEnabledCode, errors.AsyncFunctionNotEnabledMessage, nil)
	}
	return
}


func (client *Client) Shutdown() {
	// lock the addAsync()
	client.asyncChanLock.Lock()
	defer client.asyncChanLock.Unlock()
	if client.asyncTaskQueue != nil {
		close(client.asyncTaskQueue)
	}
	client.isRunning = false
}