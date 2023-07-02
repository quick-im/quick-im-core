package plugin

import (
	"context"
	"fmt"

	"github.com/quick-im/quick-im-core/internal/tracing"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	tr "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func AddClientTrace(clientName, agentHostPort string, plugins client.PluginContainer) (*tr.TracerProvider, context.Context) {
	tracer, ctx, err := tracing.InitJaeger(clientName, agentHostPort)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Jaeger: %v", err))
	}
	ts := otel.Tracer(clientName)
	plugins.Add(NewClientTracingPlugin(ts))
	return tracer, ctx
}

type clientTracingPlugin struct {
	tracer      trace.Tracer
	propagators propagation.TextMapPropagator
}

func NewClientTracingPlugin(tracer trace.Tracer) *clientTracingPlugin {
	return &clientTracingPlugin{
		tracer:      tracer,
		propagators: otel.GetTextMapPropagator(),
	}
}

func (p *clientTracingPlugin) PreCall(ctx context.Context, servicePath, serviceMethod string, args interface{}) error {
	spanCtx := tracing.Extract(ctx, p.propagators)
	ctx0 := trace.ContextWithSpanContext(ctx, spanCtx)

	ctx1, span := p.tracer.Start(ctx0, "rpcx.client."+servicePath+"."+serviceMethod)
	tracing.Inject(ctx1, p.propagators)

	ctx.(*share.Context).SetValue(tracing.OpenTelemetryKey, span)

	span.AddEvent("PreCall")
	span.SetAttributes(attribute.String("rpcx.ServicePath", servicePath))
	span.SetAttributes(attribute.String("rpcx.ServiceMethod", serviceMethod))
	span.SetAttributes(attribute.String("rpcx.MessageType", "request"))
	return nil
}

func (p *clientTracingPlugin) PostCall(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}, err error) error {
	span := ctx.Value(tracing.OpenTelemetryKey).(trace.Span)
	defer span.End()

	span.AddEvent("PostCall")
	span.SetAttributes(attribute.String("rpcx.MessageType", "response"))
	if err != nil {
		span.SetAttributes(attribute.String("rpcx.Error", err.Error()))
	}
	return nil
}
