package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"

	"sourcegraph.com/sourcegraph/appdash"
	"sourcegraph.com/sourcegraph/appdash/traceapp"

	lightstepot "github.com/lightstep/lightstep-tracer-go"
	"github.com/opentracing/opentracing-go"

	appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"
)

var (
	port           = flag.Int("port", 8080, "Example app port.")
	appdashPort    = flag.Int("appdash.port", 8700, "Run appdash locally on this port.")
	lightstepToken = flag.String("lightstep.token", "", "Lightstep access token.")
)

func main() {
	flag.Parse()

	var tracer opentracing.Tracer

	// Would it make sense to embed Appdash?
	if len(*lightstepToken) > 0 {
		tracer = lightstepot.NewTracer(lightstepot.Options{AccessToken: *lightstepToken})
	} else {
		addr := startAppdashServer(*appdashPort)
		tracer = appdashot.NewTracer(appdash.NewRemoteCollector(addr))
	}

	opentracing.InitGlobalTracer(tracer)

	fmt.Printf("Go to http://localhost:%d/home to start a request!\n", *port)

	createSpan()
	select {}
}

// Returns the remote collector address.
func startAppdashServer(appdashPort int) string {
	store := appdash.NewMemoryStore()

	// Listen on any available TCP port locally.
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		log.Fatal(err)
	}
	collectorPort := l.Addr().(*net.TCPAddr).Port

	// Start an Appdash collection server that will listen for spans and
	// annotations and add them to the local collector (stored in-memory).
	cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
	go cs.Start()

	// Print the URL at which the web UI will be running.
	appdashURLStr := fmt.Sprintf("http://localhost:%d", appdashPort)
	appdashURL, err := url.Parse(appdashURLStr)
	if err != nil {
		log.Fatalf("Error parsing %s: %s", appdashURLStr, err)
	}
	fmt.Printf("To see your traces, go to %s/traces\n", appdashURL)

	// Start the web UI in a separate goroutine.
	tapp, err := traceapp.New(nil, appdashURL)
	if err != nil {
		log.Fatalf("Error creating traceapp: %v", err)
	}
	tapp.Store = store
	tapp.Queryer = store
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appdashPort), tapp))
	}()
	return fmt.Sprintf(":%d", collectorPort)
}
