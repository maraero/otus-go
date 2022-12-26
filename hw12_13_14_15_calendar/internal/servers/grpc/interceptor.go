package server_grpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func loggerInterceptor(logger logger.Log) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		md, ok := metadata.FromIncomingContext(ctx)
		result, err := handler(ctx, req)
		logger.Info(fmt.Sprintf("%s %s %s", method(info.FullMethod), latency(start), userAgent(md, ok)))
		return result, err
	}
}

func method(full string) string {
	return full[strings.LastIndex(full, "/"):]
}

func latency(start time.Time) string {
	return fmt.Sprintf("%dms", time.Since(start).Milliseconds())
}

func userAgent(md metadata.MD, ok bool) (result string) {
	if !ok {
		return
	}
	agents, ok := md["user-agent"]
	if !ok {
		return
	}
	if len(agents) == 0 {
		return
	}
	return `"` + agents[0] + `"`
}
