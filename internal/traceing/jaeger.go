package traceing

import (
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

func InitJaeger(serviceName string, agentHostPort string) (tracer opentracing.Tracer, closer io.Closer, err error) {
	// Create a new configuration object for Jaeger.
	cfg := &config.Configuration{
		ServiceName: serviceName, // Set the service name for the tracer.
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst, // Use constant sampling type.
			Param: 1,                       // Set the sampling parameter to 1 (sample all traces).
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,            // Enable logging of spans.
			BufferFlushInterval: 1 * time.Second, // Set the buffer flush interval to 1 second.
			LocalAgentHostPort:  agentHostPort,   // Set the address of the Jaeger Agent.
		},
	}

	// Create a new Jaeger tracer using the configuration.
	tracer, closer, err = cfg.NewTracer(
		config.Logger(jaeger.StdLogger),     // Use the standard logger for logging.
		config.Metrics(metrics.NullFactory), // Disable metrics.
	)
	return
}
