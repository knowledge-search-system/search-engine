package grpc

import (
	"context"
	"errors"
	"runtime/debug"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	LangMetadataKey = "x-lang"
	defaultLang     = "ru"
)

type Translator interface {
	Translate(messageKey, lang string) string
}

type SentinelError interface {
	error
	Code() int
	MessageKey() string
}

func LoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		fields := []zap.Field{zap.String("method", info.FullMethod)}
		if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
			fields = append(fields, zap.String("trace_id", span.TraceID().String()))
		}

		resp, err := handler(ctx, req)
		if err != nil {
			logger.Error("grpc request failed", append(fields, zap.Error(err))...)
		} else {
			logger.Info("grpc request handled", fields...)
		}

		return resp, err
	}
}

func RecoveryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("grpc handler panicked",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
					zap.String("stack", string(debug.Stack())),
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

func ErrorTranslationInterceptor(translator Translator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		var sentinel SentinelError
		if !errors.As(err, &sentinel) {
			return resp, err
		}

		return resp, status.Error(grpcCodeFor(sentinel.Code()), translator.Translate(sentinel.MessageKey(), langFromContext(ctx)))
	}
}

func langFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return defaultLang
	}

	values := md.Get(LangMetadataKey)
	if len(values) == 0 || values[0] == "" {
		return defaultLang
	}

	return values[0]
}

func grpcCodeFor(code int) codes.Code {
	switch code {
	case 1:
		return codes.InvalidArgument
	case 2:
		return codes.NotFound
	case 3:
		return codes.AlreadyExists
	case 4:
		return codes.Internal
	case 5:
		return codes.Unavailable
	default:
		return codes.Unknown
	}
}
