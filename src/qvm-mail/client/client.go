package client

import (
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golib/aws/service/awserr"
	"github.com/golib/aws/service/client"
	"github.com/golib/aws/service/client/metadata"
	"github.com/golib/aws/service/credentials"
	"github.com/golib/aws/service/request"
	"github.com/golib/aws/service/session"
	"github.com/golib/aws/service/signer/v4"
)

var (
	ignoreHeaders = map[string]bool{
		"X-Forward-For":   true,
		"X-Forwarded-For": true,
	}
)

type Client struct {
	region  string
	session *session.Session
}

func New(accessKeyId, accessKeySecret string) *Client {
	provider := &credentials.StaticProvider{
		Value: credentials.Value{
			AccessKeyID:     accessKeyId,
			SecretAccessKey: accessKeySecret,
		},
	}

	cred := credentials.NewCredentials(provider)

	sess := session.New()
	sess.Config.WithCredentials(cred)
	sess.Config.WithMaxRetries(1)
	sess.Config.WithForcePathStyle(true)
	// sess.Config.WithLogLevel(service.LogDebug)

	return &Client{
		session: sess,
	}
}

func NewWithRegion(accessKeyId, accessKeySecret, region string) *Client {
	proxy := New(accessKeyId, accessKeySecret)
	proxy.region = region

	return proxy
}

func (proxy *Client) WithRegion(region string) *Client {
	proxy.region = region

	return proxy
}

func (proxy *Client) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	serviceName := "qvm"
	serviceEndpoint := r.URL.Scheme + "://" + r.Host

	serviceRegion := r.Header.Get("X-Aws-Region")
	if serviceRegion == "" {
		serviceRegion = proxy.region
	}
	if serviceRegion == "" {
		serviceRegion = "cn-proxy-1"
	}

	c := proxy.session.ClientConfig(serviceName)
	c.Config.WithEndpoint(serviceEndpoint)
	c.Config.WithRegion(serviceRegion)

	svc := client.New(
		*c.Config,
		metadata.ClientInfo{
			ServiceName:   serviceName,
			Endpoint:      serviceEndpoint,
			SigningRegion: serviceRegion,
			APIVersion:    "2006-03-01",
		},
		c.Handlers,
	)

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(request.NamedHandler{
		Name: "evm.proxy.Header",
		Fn: func(req *request.Request) {
			// NOTE: inject region header
			req.HTTPRequest.Header.Set("X-Aws-Region", serviceRegion)

			// tampers auth token
			token := r.Header.Get("X-Auth-Token")
			if token != "" {
				req.HTTPRequest.Header.Add("X-Aws-Token", token)

				r.Header.Del("X-Auth-Token")
			}

			// injects all extra headers
			for name, headers := range r.Header {
				// NOTE: avoid duplicated headers when retrying!
				req.HTTPRequest.Header.Del(name)

				// ignore X-Forward-For specs
				if ignoreHeaders[http.CanonicalHeaderKey(name)] {
					continue
				}

				for _, value := range headers {
					req.HTTPRequest.Header.Add(name, value)
				}
			}
		},
	})
	svc.Handlers.UnmarshalError.PushBackNamed(request.NamedHandler{
		Name: "evm.proxy.UnmarshalError",
		Fn: func(req *request.Request) {
			defer req.HTTPResponse.Body.Close()

			buf, err := ioutil.ReadAll(req.HTTPResponse.Body)
			if err != nil {
				req.Error = awserr.New("SerializationError", "failed to read from HTTP response body", err)
				return
			}

			// return raw error messages
			req.Error = awserr.New("UnmarshalError", string(buf), nil)
		},
	})

	op := &request.Operation{
		Name:       "Proxy",
		HTTPMethod: r.Method,
		HTTPPath:   r.URL.Path,
	}

	req := svc.NewRequest(op, nil, nil)

	// NOTE: fix query string
	req.HTTPRequest.URL.RawQuery = r.URL.Query().Encode()

	switch r.Body.(type) {
	case io.ReadSeeker:
		rser, _ := r.Body.(io.ReadSeeker)

		// sign request body
		req.HTTPRequest.Header.Set("X-Aws-Content-Sha256", hex.EncodeToString(Helpers.MakeSha256Reader(rser)))

		req.SetReaderBody(rser)

	case nil:
		// sign request body
		req.HTTPRequest.Header.Set("X-Aws-Content-Sha256", hex.EncodeToString(Helpers.MakeSha256([]byte(""))))

	default:
		// NOTE: performance issue!
		buf, _ := ioutil.ReadAll(r.Body)

		// sign request body
		req.HTTPRequest.Header.Set("X-Aws-Content-Sha256", hex.EncodeToString(Helpers.MakeSha256(buf)))

		r.Body.Close()

		req.SetBufferBody(buf)
	}

	err = req.Send()
	if err == nil {
		resp = req.HTTPResponse
	}

	return
}
