package tracing

import (
	"context"
	"fmt"
	"net"

	"github.com/smallnest/rpcx/share"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitJaeger(serviceName string, agentHostPort string) (*trace.TracerProvider, context.Context, error) {
	ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, make(map[string]string))
	// Create a new configuration object for Jaeger.
	host, port, err := net.SplitHostPort(agentHostPort)
	if err != nil {
		return nil, ctx, err
	}
	exporter, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(host), jaeger.WithAgentPort(port)))
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to create exporter: %v", err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to create resource: %v", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resources),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tracerProvider, ctx, nil
}
