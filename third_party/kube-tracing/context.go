package kubetracing

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

type contextKey struct{}

var activeSpanKey = contextKey{}

func SpanContextStringFromContext(ctx context.Context, opName string) string {
	traceCtx := ctx.Value(activeSpanKey)
	if traceCtx == nil {
		return opName
	}

	if str, ok := traceCtx.(string); !ok || str == "" {
		return opName
	}

	return traceCtx.(string)
}

func WrapContextWithSpan(ctx context.Context, value interface{}) context.Context {
	if value == nil {
		return ctx
	}

	switch value := value.(type) {
	case opentracing.Span:
		return context.WithValue(ctx, activeSpanKey, ContextStringFromSpan(value))
	case string:
		return context.WithValue(ctx, activeSpanKey, value)
	}

	return ctx
}
