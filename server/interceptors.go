package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const authTokenKey string = "auth_token"
const authTokenValue string = "authd"
const grpcService = 5
const grpcMethod = 7

func logCalls(l *log.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			msg = fmt.Sprintf("DEBUG :%v", msg)
		case logging.LevelInfo:
			msg = fmt.Sprintf("INFO :%v", msg)
		case logging.LevelWarn:
			msg = fmt.Sprintf("WARN :%v", msg)
		case logging.LevelError:
			msg = fmt.Sprintf("ERROR :%v", msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}

		l.Println(msg, fields[grpcService], fields[grpcMethod])
	})
}

func validateAuthToken(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	if t, ok := md[authTokenKey]; ok {
		switch {
		case len(t) != 1:
			return nil, status.Errorf(
				codes.InvalidArgument,
				"%s should contain only 1 value", authTokenKey,
			)
		case t[0] != "authd":
			return nil, status.Errorf(
				codes.Unauthenticated,
				"incorect %s", authTokenKey,
			)
		}
	} else {
		return nil, status.Errorf(
			codes.Unauthenticated,
			"failed to get %s", authTokenKey,
		)
	}
	return ctx, nil
}

// replaced by using go-grpc-middleware/v2/interceptors/auth which
// now shares the same validateAuthToken call above for unary and streaming
//
//func unaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//	if err := validateAuthToken(ctx); err != nil {
//		return nil, err
//	}
//	return handler(ctx, req)
//}
//
//func streamAuthInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
//	if err := validateAuthToken(ss.Context()); err != nil {
//		return err
//	}
//	return handler(srv, ss)
//}

// replaced by go-grpc-middleware logging interceptor and logCalls function
//
//func unaryLogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//	log.Println(info.FullMethod, "called")
//	return handler(ctx, req)
//}
//
//func streamLogInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
//	log.Println(info.FullMethod, "called")
//	return handler(srv, ss)
//}
