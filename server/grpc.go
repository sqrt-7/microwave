package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"runtime/debug"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sqrt-7/microwave/log"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server

	Logger     log.Logger
	Port       string
	WaitGroup  *sync.WaitGroup
	ShutdownCh chan struct{}
	ErrCh      chan error
}

func (s *GRPCServer) Default(optUnaryInterceptors []grpc.UnaryServerInterceptor, optStreamInterceptors []grpc.StreamServerInterceptor) *GRPCServer {
	unaryLogger := initUnaryLogger(s.Logger)
	streamLogger := initStreamLogger(s.Logger)

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		unaryLogger,
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(panicHandler(s.Logger))),
		grpc_tags.UnaryServerInterceptor(),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		streamLogger,
		grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(panicHandler(s.Logger))),
		grpc_tags.StreamServerInterceptor(),
	}

	if optUnaryInterceptors != nil {
		unaryInterceptors = append(unaryInterceptors, optUnaryInterceptors...)
	}

	if optStreamInterceptors != nil {
		streamInterceptors = append(streamInterceptors, optStreamInterceptors...)
	}

	return s.NewGRPCServer(optUnaryInterceptors, optStreamInterceptors)
}

func (s *GRPCServer) NewGRPCServer(optUnaryInterceptors []grpc.UnaryServerInterceptor, optStreamInterceptors []grpc.StreamServerInterceptor) *GRPCServer {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(optUnaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(optStreamInterceptors...)),
		grpc.StatsHandler(&ocgrpc.ServerHandler{}), // this enables TraceID propagation
	)

	s.server = server

	return s
}

func (s *GRPCServer) Run() {
	defer s.WaitGroup.Done()

	if string(s.Port[0]) != ":" {
		s.Port = ":" + s.Port
	}

	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		s.Logger.Error("GRPC_SERVER_ERROR").WithField("error", err.Error()).Send()
		s.ErrCh <- err
		return
	}

	go func() {
		if err := s.server.Serve(lis); err != nil {
			s.Logger.Error("GRPC_SERVER_ERROR").WithField("error", err.Error()).Send()
			s.ErrCh <- err
		}
	}()

	s.Logger.Info("GRPC_SERVER_STARTED").WithField("port", s.Port).Send()

	<-s.ShutdownCh

	s.server.GracefulStop()
	s.Logger.Info("GRPC_SERVER_STOPPED").WithField("port", s.Port).Send()
}

// initUnaryLogger creates a logger for unary requests
func initUnaryLogger(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Info("GRPC_IN").For(ctx).WithField("method", info.FullMethod).Send()

		resp, err := handler(ctx, req)

		// todo
		//err = errorHandler(ctx, logger, err, info.FullMethod)

		// Log response
		logger.Info("GRPC_OUT").For(ctx).WithField("method", info.FullMethod).Send()

		return resp, err
	}
}

// initStreamLogger returns the logger for stream requests
func initStreamLogger(logger log.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		logger.Info("GRPC_STREAM_IN").For(ss.Context()).WithField("method", info.FullMethod).Send()

		err := handler(srv, ss)

		// todo
		//err = errorHandler(ss.Context(), logger, err, info.FullMethod)

		logger.Info("GRPC_STREAM_OUT").For(ss.Context()).WithField("method", info.FullMethod).Send()

		return err
	}
}

// panicHandler logs the panic & stack trace
func panicHandler(logger log.Logger) grpc_recovery.RecoveryHandlerFunc {
	return func(p interface{}) error {
		logger.Error("GRPC_PANIC").
			WithField("panic", fmt.Sprint(p)).
			WithField("stack", base64.RawStdEncoding.EncodeToString(debug.Stack())).
			Send()

		return fmt.Errorf("grpc_panic: %v", p)
	}
}
