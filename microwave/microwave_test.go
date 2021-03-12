package microwave_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/sqrt-7/microwave/log"
	"github.com/sqrt-7/microwave/microwave"
)

type TestWrapper struct {
}

func (s TestWrapper) Run(m *microwave.Microwave) {

}

func TestService1(t *testing.T) {
	namespace := "my-service"
	envs := []string{
		"GRPC_PORT",
		"HTTP_PORT",
		"ANOTHER_HTTP_PORT",
	}

	os.Setenv("GRPC_PORT", "1234")
	os.Setenv("HTTP_PORT", "1235")
	os.Setenv("ANOTHER_HTTP_PORT", "1236")

	cfg, _ := log.ObservedConfig(namespace)
	l, err := log.New(namespace, cfg)

	mw, err := microwave.New(
		namespace,
		microwave.CustomLogger(l),
		microwave.RequireEnvs(envs...),
	)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(-1)
	}

	grpcSrv := microwave.NewGRPCServer(microwave.DefaultGRPCInterceptors(mw.Logger()))
	grpcWrapper := &microwave.GRPCWrapper{
		Port:   mw.Env("GRPC_PORT"),
		Server: grpcSrv,
	}

	httpSrv := microwave.NewHTTPServer()
	httpWrapper := &microwave.HTTPWrapper{
		Port:   mw.Env("HTTP_PORT"),
		Server: httpSrv,
	}

	anotherHttpSrv := microwave.NewHTTPServer()
	anotherHttpWrapper := &microwave.HTTPWrapper{
		Port:   mw.Env("ANOTHER_HTTP_PORT"),
		Server: anotherHttpSrv,
	}

	mw.AddServer(grpcWrapper)
	mw.AddServer(httpWrapper)
	mw.AddServer(anotherHttpWrapper)

	mw.Start()
}
