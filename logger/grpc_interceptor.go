package logger

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

func GetGRPCInterceptor(mainLogger *Logger) func(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		mainLogger.Tracef("call=%v time=%v err=%v",
			method, time.Since(start), err)
		return err
	}
}
