
import "net/http"

type CustomHttpClient interface {
	Client() *http.Client
}

type CustomTransport struct {
	Transport http.RoundTripper
}

func NewCustomTransport(t http.RoundTripper) CustomHttpClient {
	return &CustomTransport{
		Transport: t,
	}
}

func (t *CustomTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

func (t *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.transport().RoundTrip(req)
}

func (t *CustomTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}