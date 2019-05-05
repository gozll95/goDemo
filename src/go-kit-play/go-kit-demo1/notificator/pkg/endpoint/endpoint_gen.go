// THIS FILE IS AUTO GENERATED BY GK-CLI DO NOT EDIT!!
package endpoint

import (
	endpoint "github.com/go-kit/kit/endpoint"
	service "go-kit-demo1/notificator/pkg/service"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	SendEmailEndpoint endpoint.Endpoint
}

// New returns a Endpoints struct that wraps the provided service, and wires in all of the
// expected endpoint middlewares
func New(s service.NotificatorService, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{SendEmailEndpoint: MakeSendEmailEndpoint(s)}
	for _, m := range mdw["SendEmail"] {
		eps.SendEmailEndpoint = m(eps.SendEmailEndpoint)
	}
	return eps
}
