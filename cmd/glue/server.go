package main

import (
	"context"

	"flamingsteve/cmd"
	"flamingsteve/cmd/glue/api"
	"flamingsteve/pkg/grpc"
	"github.com/draeron/gopkgs/logger"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/pflag"
	grpc2 "google.golang.org/grpc"
)

// Generate Golang grpc server & client
//go:generate protoc fpixel.proto -I . -I $GOPATH/src -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:./api
//go:generate protoc fpixel.proto -I . -I $GOPATH/src -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=./api
//go:generate protoc fpixel.proto -I . -I $GOPATH/src -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --govalidators_out=./api
//go:generate protoc fpixel.proto -I . -I $GOPATH/src -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=./api

// Generate python client code
//go:generate python -m grpc_tools.protoc -I . -I $GOPATH/src -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --python_out=./py --grpc_python_out=./py fpixel.proto
//go:generate python -m grpc_tools.protoc -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway --python_out=./py google/api/annotations.proto google/api/http.proto protoc-gen-swagger/options/annotations.proto protoc-gen-swagger/options/openapiv2.proto

var args struct {
	port uint16
}

func init() {
	pflag.Uint16Var(&args.port, "port", 8080, "port to listen for requests")
}

func main() {
	pflag.Parse()
	cmd.SetupLoggers()
	log := logger.New("main")

	log.Info("glue started")
	defer log.Info("glue stopped")

	svr := grpc.NewServer(grpc.Options{
		Port: args.port,
		RegisterFcts: []grpc.RegisterFunc{
			func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc2.DialOption) error {
				return fpixels.RegisterFlamePixelsHandlerFromEndpoint(ctx, mux, endpoint, opts)
			},
		},
	})

	fpixels.RegisterFlamePixelsServer(svr.Server, NewService())

	defer svr.GracefulStop()
	go svr.Listen()

	cmd.WaitForCtrlC()
}
