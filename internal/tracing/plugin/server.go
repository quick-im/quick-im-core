package plugin

import (
	"context"

	"github.com/quick-im/quick-im-core/internal/tracing"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/share"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type serverTracingPlugin struct {
	tracer      trace.Tracer
	propagators propagation.TextMapPropagator
}

func NewServerTracingPlugin(tracer trace.Tracer) *serverTracingPlugin {
	return &serverTracingPlugin{
		tracer:      tracer,
		propagators: otel.GetTextMapPropagator(),
	}
}

// func (p *serverTracingPlugin) PreCall(ctx context.Context, serviceName, methodName string, args interface{}) (interface{}, error) {
// 	println(1)
// 	return args, nil
// }

// func (p *serverTracingPlugin) PostCall(ctx context.Context, serviceName, methodName string, args interface{}, reply interface{}) (interface{}, error) {
// 	println(4)
// 	// span := opentracing.SpanFromContext(ctx)
// 	// span.Finish()
// 	return reply, nil
// }

func (p *serverTracingPlugin) PreHandleRequest(ctx context.Context, r *protocol.Message) error {
	spanCtx := tracing.Extract(ctx, p.propagators)
	ctx0 := trace.ContextWithSpanContext(ctx, spanCtx)

	ctx1, span := p.tracer.Start(ctx0, "rpcx.service."+r.ServicePath+"."+r.ServiceMethod)
	tracing.Inject(ctx1, p.propagators)

	ctx.(*share.Context).SetValue(tracing.OpenTelemetryKey, span)

	span.AddEvent("PreHandleRequest")
	span.SetAttributes(attribute.String("rpcx.ServicePath", r.ServicePath))
	span.SetAttributes(attribute.String("rpcx.ServiceMethod", r.ServiceMethod))
	span.SetAttributes(attribute.String("rpcx.MessageType", "request"))

	return nil
}

func (p *serverTracingPlugin) PostWriteResponse(ctx context.Context, req *protocol.Message, res *protocol.Message, e error) error {
	span := ctx.Value(tracing.OpenTelemetryKey).(trace.Span)
	span.AddEvent("PostWriteResponse")
	span.SetAttributes(attribute.String("rpcx.Error", e.Error()))
	defer span.End()
	return nil
}
