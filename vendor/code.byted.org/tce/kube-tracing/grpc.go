package kubetracing

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.Pairs()
	}
	if spanStrs, ok := md[strings.ToLower(KubetracingSpanContextKey)]; ok && len(spanStrs) >= 1 {
		ctx = WrapContextWithSpan(ctx, spanStrs[0])
	}

	return handler(ctx, req)
}

func UnaryClientInterceptor(ctx context.Context,
	method string, req, resp interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) (err error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.Pairs()
	}

	value := ctx.Value(activeSpanKey)
	if spanStr, ok := value.(string); ok && spanStr != "" {
		md[strings.ToLower(KubetracingSpanContextKey)] = []string{spanStr}
	}

	return invoker(metadata.NewOutgoingContext(ctx, md), method, req, resp, cc, opts...)
}
