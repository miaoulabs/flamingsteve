package grpc

import (
	"context"
	"flamingsteve/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

var logFactory logger.LoggerFactory = func(name string) logger.Logger {
	return logger.Dummy()
}

func SetLoggerFactory(newLogger func(name string) logger.Logger) {
	logFactory = newLogger
}

// UnaryServerInterceptor returns a new unary server interceptors that adds zap.Logger to the context.
func UnaryLogInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		resp, err := next(ctx, req)
		dur := time.Now().Sub(startTime)
		code := status.Convert(err)

		if code.Err() != nil {
			logger.Errorf("GRPC: %s, Code: %v, Duration: %v ms, Msg: %v", info.FullMethod, code.Code(), dur.Milliseconds(), code.Message())
		} else {
			logger.Infof("GRPC: %s, Code: %v, Duration: %v", info.FullMethod, code.Code(), dur.Milliseconds())
		}

		return resp, err
	}
}

//// StreamServerInterceptor returns a new streaming server interceptor that adds zap.Logger to the context.
//func StreamServerInterceptor(logger *zap.Logger, opts ...Option) grpc.StreamServerInterceptor {
//	o := evaluateServerOpt(opts)
//	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
//		startTime := time.Now()
//		newCtx := newLoggerForCall(stream.Context(), logger, info.FullMethod, startTime)
//		wrapped := grpc_middleware.WrapServerStream(stream)
//		wrapped.WrappedContext = newCtx
//
//		err := handler(srv, wrapped)
//		if !o.shouldLog(info.FullMethod, err) {
//			return err
//		}
//		code := o.codeFunc(err)
//		level := o.levelFunc(code)
//
//		// re-extract logger from newCtx, as it may have extra fields that changed in the holder.
//		ctxzap.Extract(newCtx).Check(level, "finished streaming call with code "+code.String()).Write(
//			zap.Error(err),
//			zap.String("grpc.code", code.String()),
//			o.durationFunc(time.Since(startTime)),
//		)
//
//		return err
//	}
//}
