package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"flamingsteve/pkg/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	*grpc.Server
	log       logger.Logger
	registers []RegisterFunc
	gateway   http.Handler
	webWrap   http.Handler

	opts Options
}

type Options struct {
	RegisterFcts []RegisterFunc
	Port         uint16
}

type RegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

func NewServer(opts Options) *Server {
	s := &Server{
		log:  logFactory("grpc"),
		opts: opts,
	}

	s.Server = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			//UnaryLogInterceptor(log.Desugar()),
			grpc_validator.StreamServerInterceptor(),
			//grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			UnaryLogInterceptor(s.log),
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	reflection.Register(s.Server)

	return s
}

func (s *Server) Listen() {

	go s.serve()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit)
	<-quit
}

func (s *Server) serve() {
	port := s.opts.Port
	endpoint := fmt.Sprintf(":%d", port)

	mux := http.NewServeMux()

	ctx := context.Background()
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	gwmux := runtime.NewServeMux()
	for _, register := range s.opts.RegisterFcts {
		err := register(ctx, gwmux, endpoint, opts)
		if err != nil {
			s.log.Panicf("failed to register gateway handler: %v", err)
		}
	}
	mux.Handle("/", gwmux)
	s.gateway = mux

	s.webWrap = grpcweb.WrapServer(s.Server,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return strings.Contains(origin, "localhost") ||
				strings.Contains(origin, "127.0.0.1") ||
				strings.Contains(origin, "pibrain")
		}),
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
	)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    endpoint,
		Handler: h2c.NewHandler(http.HandlerFunc(s.grpcHandlerFunc), &http2.Server{}),
	}

	s.log.Infof("starting grpc on port %v", port)
	err = srv.Serve(conn)
	if err != nil {
		s.log.Panicf("cannot serve grpc: %v", err)
	}
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func (s *Server) grpcHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
	//log.Debug("header", r.Header.Get("Content-Type"))
	//log.Debug("proto: ", r.Proto)
	//log.Debug("body: ", r.Body)
	if r.ProtoMajor == 2 && reqHasHeaderValue(r, "Content-Type", "application/grpc") {
		s.Server.ServeHTTP(w, r)
	} else if reqHasHeaderValue(r, "Access-Control-Request-Headers", "x-grpc-web") || reqHasHeaderValue(r, "X-Grpc-Web", "") {
		s.webWrap.ServeHTTP(w, r)
	} else {
		s.gateway.ServeHTTP(w, r)
	}
}

func reqHasHeaderValue(r *http.Request, header string, value string) bool {
	h := r.Header.Get(header)
	if h != "" {
		return strings.Contains(h, value)
	}
	return false
}
