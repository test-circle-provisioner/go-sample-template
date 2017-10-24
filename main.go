package main

import (
	"fmt"
	"github.com/apex/log"
	"github.com/segmentio/ecs-logs-go/apex"
	"github.com/segmentio/ecs-logs-go/log"
	"github.com/segmentio/stats"
	"github.com/segmentio/stats/datadog"
	"gopkg.in/alecthomas/kingpin.v2"
	stdlog "log"
	"os"
	"os/signal"
	"net/http"
	"github.com/segmentio/stats/httpstats"
	"syscall"
	"context"
)

type config struct {
	DatadogAddr string
	LogLevel    string
	BindAddr    string
	Message 	string
}

func parseConfig() *config {
	cfg := &config{}
	kingpin.Flag("bind-addr", "Address and port to listen on").Default(":3000").StringVar(&cfg.BindAddr)
	kingpin.Flag("datadog-addr", "Datadog statsd host and port").Default("127.0.0.1:8125").StringVar(&cfg.DatadogAddr)
	kingpin.Flag("log-level", "Logging level").Default("INFO").StringVar(&cfg.LogLevel)
	kingpin.Flag("message", "Message to Serve").Default("Hello, World").StringVar(&cfg.Message)
	kingpin.Parse()
	return cfg
}

const (
	program = "Hello World Service"
	version = "1.0.3"
)

func setupLogging(cfg *config) {
	handler := log_ecslogs.NewHandler(os.Stdout)
	writer := log_ecslogs.NewWriter("", stdlog.Flags(), handler)
	stdlog.SetOutput(writer)

	log.SetHandler(apex_ecslogs.NewHandler(os.Stdout))
	log.SetLevel(log.MustParseLevel(cfg.LogLevel))
	log.Log = log.WithFields(log.Fields{
		"program": program,
		"version": version,
	})
}

func setupStats(cfg *config) {
	stats.DefaultEngine = stats.NewEngine(program, stats.Discard, []stats.Tag{
		{Name: "program", Value: program},
		{Name: "version", Value: version},
	}...)
	stats.Register(datadog.NewClient(cfg.DatadogAddr))
}

// Embed the message in the title and body of a simple HTML page
func htmlBody(message string) string {
	return fmt.Sprintf("<html>\n<head><title>%s</title></head>\n<body>\n<p>%s</p>\n</body>\n</html>", message, message)
}

func getResponseWriter(cfg *config) func (http.ResponseWriter, *http.Request) {
	return func (res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(htmlBody(cfg.Message)))
	}
}


func main() {
	cfg := parseConfig()
	setupLogging(cfg)
	setupStats(cfg)
	defer stats.Flush()

	httpServer := &http.Server{Addr: cfg.BindAddr, Handler: httpstats.NewHandler(
		http.HandlerFunc(getResponseWriter(cfg)))}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	exitCode := 0
	go func() {
		select {
		case s := <-sigChan:
			log.Infof("Received %s. Shutting down http server.", s.String())
		}
		httpServer.Shutdown(context.Background())
	}()

	log.Info("Starting http server")
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithError(err).Fatal("Http server failed")
	}

	os.Exit(exitCode)

}
