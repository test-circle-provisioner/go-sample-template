package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/justinas/alice"
	"github.com/segmentio/conf"

	"github.com/segmentio/events"
	_ "github.com/segmentio/events/ecslogs"
	"github.com/segmentio/events/httpevents"
	_ "github.com/segmentio/events/log"
	_ "github.com/segmentio/events/sigevents"
	_ "github.com/segmentio/events/text"

	"github.com/segmentio/go-hello-world"

	"github.com/segmentio/rpc"
	"github.com/segmentio/rpc/rpcevents"
	"github.com/segmentio/rpc/rpcstats"

	"github.com/segmentio/stats"
	"github.com/segmentio/stats/datadog"
	"github.com/segmentio/stats/httpstats"
	"github.com/segmentio/stats/procstats"
)

var (
	// Version represents the version of the program at runtime. The default value
	// is set as a placeholder. Our build script will inject a real value during
	// compilation based on the Git tags.
	Version = "x.x.x"
)

func main() {
	// Define and load the configuration options provided by this program.
	config := struct {
		DatadogAddress string `conf:"datadog_address" help:"datadog adress"`
		Debug          bool   `conf:"debug" help:"enables debug logging"`
		Address        string `conf:"address" help:"address on which the server should listen"`
	}{
		DatadogAddress: ":8125",
		Debug:          true,
		Address:        ":3000",
	}
	conf.Load(&config)

	// Configure debug logging. The events package implicitly handles setting up
	// debugging and ECS compatible logging using initialization only imports.
	events.DefaultLogger.EnableDebug = config.Debug

	// Log the starting
	events.Log("starting service", events.Args{
		{Name: "version", Value: Version},
		{Name: "config", Value: config},
	})

	// Initialize stats client.
	dd := datadog.NewClientWith(datadog.ClientConfig{
		Address: config.DatadogAddress,
	})
	stats.Register(dd)
	defer dd.Close()
	defer stats.Flush()

	// Collect Go runtime stats.
	goMetrics := procstats.StartCollector(procstats.NewGoMetrics())
	defer goMetrics.Close()
	// Collect process stats.
	procMetrics := procstats.StartCollector(procstats.NewProcMetrics())
	defer procMetrics.Close()

	// Create our service implementation.
	helloWorldService := helloWorld.New()

	// Bind our service as an JSON RPC service.
	rpcHandler := rpc.NewHandler()
	rpcHandler.Register("HelloWorld", newRPCService(helloWorldService))

	// Install HTTP level logging and metrics handlers. Do not use this for health
	// checks to reduce noise in production.
	chain := alice.New(httpstats.NewHandler, httpevents.NewHandler)

	// Create an HTTP handler for our service. This sets up 3 components:
	// 1. Our RPC service. By convention, this is setup at /rpc.
	// 2. A healthcheck. This is required for ELB health checks.
	// 3. http.DefaultServeMux. The initialization only import from the pprof package installs it's server on http.DefaultServeMux.
	mux := http.NewServeMux()
	mux.Handle("/rpc", chain.Then(rpcHandler))
	mux.HandleFunc("/internal/health", health)
	mux.Handle("/", chain.Then(http.DefaultServeMux))

	// Create a HTTP server that binds our handler to the configured address. This
	// server handles graceful shutdowns. It starts on a background goroutine and
	// listens for a signal to start the shutdown logic. In production, a SIGTERM
	// signal is sent by EC2. See https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_StopTask.html for details.
	server := &http.Server{Addr: config.Address, Handler: mux}
	defer server.Shutdown(context.Background())

	go func() {
		server.ListenAndServe()
	}()

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigchan
	events.Log("stopping in response to signal %{signal}s.", sig)
}

// newRPCService wraps an RPC receiver with RPC level logging and metrics.
func newRPCService(rcvr interface{}) rpc.Invoker {
	s := rpcstats.NewInvokerWith(stats.DefaultEngine, rpc.NewService(rcvr))
	s = rpcevents.NewInvokerWith(events.DefaultLogger, s)
	return s
}

// health implements an ELB health check and responds with an HTTP 200.
func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(200)
}
