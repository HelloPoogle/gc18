package main

import (
	"io"
	"log"

	"github.com/gobuffalo/envy"

	"github.com/gophercon/gc18/gophercon/actions"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

// ServiceName is the string name of the service, since
// it is used in multiple places, it's an exported Constant
const ServiceName = "gophercon.web"

func main() {

	tracer, closer, err := initTracer()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	port := envy.Get("PORT", "3000")
	actions.Tracer = tracer
	app := actions.App()

	log.Fatal(app.Start(port))
}
func initTracer() (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	closer, err := cfg.InitGlobalTracer(
		"serviceName",
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return nil, closer, err
	}
	defer closer.Close()
	tracer, closer, err := cfg.New(ServiceName, jaegercfg.Logger(jLogger), jaegercfg.Metrics(jMetricsFactory))
	return tracer, closer, err
}
