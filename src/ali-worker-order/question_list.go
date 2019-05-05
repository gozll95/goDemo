package workorder

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) QuestionList(request *QuestionListRequest) (response *QuestionListResponse, err error) {
	response = CreateQuestionListResponse()
	err = client.DoAction(request, response)
	return
}

func (client *Client) QuestionListWithChan(request *QuestionListRequest) (<-chan *QuestionListResponse, <-chan error) {
	responseChan := make(chan *QuestionListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QuestionList(request)
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

func (client *Client) QuestionListWithCallback(request *QuestionListRequest, callback func(response *QuestionListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QuestionListResponse
		var err error
		defer close(result)
		response, err = client.QuestionList(request)
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

type QuestionListRequest struct {
	*requests.RpcRequest
	AddEndTime     requests.Integer `position:"Query" name:"AddEndTime"`
	ProductIds     string           `position:"Query" name:"ProductIds"`
	AddStartTime   requests.Integer `position:"Query" name:"AddStartTime"`
	PageSize       requests.Integer `position:"Query" name:"PageSize"`
	Ids            string           `position:"Query" name:"Ids"`
	QuestionStatus string           `position:"Query" name:"QuestionStatus"`
	PageStart      requests.Integer `position:"Query" name:"PageStart"`
}

type QuestionListResponse struct {
	*responses.BaseResponse
	Success    bool   `json:"Success" xml:"Success"`
	Code       string `json:"Code" xml:"Code"`
	Message    string `json:"Message" xml:"Message"`
	Count      int    `json:"Count" xml:"Count"`
	ListResult struct {
		QuestionDetail []struct {
			Id             string `json:"Id" xml:"Id"`
			AddTime        int    `json:"AddTime" xml:"AddTime"`
			QuestionStatus string `json:"QuestionStatus" xml:"QuestionStatus"`
			Title          string `json:"Title" xml:"Title"`
			ProductId      int    `json:"ProductId" xml:"ProductId"`
			Uid            int    `json:"Uid" xml:"Uid"`
		} `json:"QuestionDetail" xml:"QuestionDetail"`
	} `json:"ListResult" xml:"ListResult"`
}

func CreateQuestionListRequest() (request *QuestionListRequest) {
	request = &QuestionListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Workorder", "2016-09-19", "QuestionList", "", "")
	return
}

func CreateQuestionListResponse() (response *QuestionListResponse) {
	response = &QuestionListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
