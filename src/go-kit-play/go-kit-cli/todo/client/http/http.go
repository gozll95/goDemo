package http

import (
	"bytes"
	"context"
	"encoding/json"
	endpoint "github.com/go-kit/kit/endpoint"
	http "github.com/go-kit/kit/transport/http"
	endpoint1 "go-kit-cli/todo/pkg/endpoint"
	http2 "go-kit-cli/todo/pkg/http"
	service "go-kit-cli/todo/pkg/service"
	"io/ioutil"
	http1 "net/http"
	"net/url"
	"strings"
)

// New returns an AddService backed by an HTTP server living at the remote
// instance. We expect instance to come from a service discovery system, so
// likely of the form "host:port".
func New(instance string, options map[string][]http.ClientOption) (service.TodoService, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}
	var getEndpoint endpoint.Endpoint
	{
		getEndpoint = http.NewClient("POST", copyURL(u, "/get"), encodeHTTPGenericRequest, decodeGetResponse, options["Get"]...).Endpoint()
	}

	var addEndpoint endpoint.Endpoint
	{
		addEndpoint = http.NewClient("POST", copyURL(u, "/add"), encodeHTTPGenericRequest, decodeAddResponse, options["Add"]...).Endpoint()
	}

	var setCompleteEndpoint endpoint.Endpoint
	{
		setCompleteEndpoint = http.NewClient("POST", copyURL(u, "/set-complete"), encodeHTTPGenericRequest, decodeSetCompleteResponse, options["SetComplete"]...).Endpoint()
	}

	var removeCompleteEndpoint endpoint.Endpoint
	{
		removeCompleteEndpoint = http.NewClient("POST", copyURL(u, "/remove-complete"), encodeHTTPGenericRequest, decodeRemoveCompleteResponse, options["RemoveComplete"]...).Endpoint()
	}

	var deleteEndpoint endpoint.Endpoint
	{
		deleteEndpoint = http.NewClient("POST", copyURL(u, "/delete"), encodeHTTPGenericRequest, decodeDeleteResponse, options["Delete"]...).Endpoint()
	}

	var getByIdEndpoint endpoint.Endpoint
	{
		getByIdEndpoint = http.NewClient("POST", copyURL(u, "/get-by-id"), encodeHTTPGenericRequest, decodeGetByIdResponse, options["GetById"]...).Endpoint()
	}

	return endpoint1.Endpoints{
		AddEndpoint:            addEndpoint,
		DeleteEndpoint:         deleteEndpoint,
		GetByIdEndpoint:        getByIdEndpoint,
		GetEndpoint:            getEndpoint,
		RemoveCompleteEndpoint: removeCompleteEndpoint,
		SetCompleteEndpoint:    setCompleteEndpoint,
	}, nil
}

// EncodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// SON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPGenericRequest(_ context.Context, r *http1.Request, request interface{}) error {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeGetResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded concat response from the HTTP response body. If the response
// as a non-200 status code, we will interpret that as an error and attempt to
//  decode the specific error message from the response body.
func decodeGetResponse(_ context.Context, r *http1.Response) (interface{}, error) {
	if r.StatusCode != http1.StatusOK {
		return nil, http2.ErrorDecoder(r)
	}
	var resp endpoint1.GetResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeAddResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded concat response from the HTTP response body. If the response
// as a non-200 status code, we will interpret that as an error and attempt to
//  decode the specific error message from the response body.
func decodeAddResponse(_ context.Context, r *http1.Response) (interface{}, error) {
	if r.StatusCode != http1.StatusOK {
		return nil, http2.ErrorDecoder(r)
	}
	var resp endpoint1.AddResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeSetCompleteResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded concat response from the HTTP response body. If the response
// as a non-200 status code, we will interpret that as an error and attempt to
//  decode the specific error message from the response body.
func decodeSetCompleteResponse(_ context.Context, r *http1.Response) (interface{}, error) {
	if r.StatusCode != http1.StatusOK {
		return nil, http2.ErrorDecoder(r)
	}
	var resp endpoint1.SetCompleteResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeRemoveCompleteResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded concat response from the HTTP response body. If the response
// as a non-200 status code, we will interpret that as an error and attempt to
//  decode the specific error message from the response body.
func decodeRemoveCompleteResponse(_ context.Context, r *http1.Response) (interface{}, error) {
	if r.StatusCode != http1.StatusOK {
		return nil, http2.ErrorDecoder(r)
	}
	var resp endpoint1.RemoveCompleteResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeDeleteResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded concat response from the HTTP response body. If the response
// as a non-200 status code, we will interpret that as an error and attempt to
//  decode the specific error message from the response body.
func decodeDeleteResponse(_ context.Context, r *http1.Response) (interface{}, error) {
	if r.StatusCode != http1.StatusOK {
		return nil, http2.ErrorDecoder(r)
	}
	var resp endpoint1.DeleteResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeGetByIdResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded concat response from the HTTP response body. If the response
// as a non-200 status code, we will interpret that as an error and attempt to
//  decode the specific error message from the response body.
func decodeGetByIdResponse(_ context.Context, r *http1.Response) (interface{}, error) {
	if r.StatusCode != http1.StatusOK {
		return nil, http2.ErrorDecoder(r)
	}
	var resp endpoint1.GetByIdResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
func copyURL(base *url.URL, path string) (next *url.URL) {
	n := *base
	n.Path = path
	next = &n
	return
}
