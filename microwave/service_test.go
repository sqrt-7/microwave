package microwave_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/sqrt-7/microwave/microwave"
	"github.com/sqrt-7/microwave/server"
)

func TestService1(t *testing.T) {
	s1 := func(s *microwave.Service) server.Server {
		srv := &server.GRPCServer{
			Logger:     s.Logger(),
			Port:       s.Env("GRPC_PORT"),
			WaitGroup:  s.WaitGroup(),
			ShutdownCh: s.ShutdownChan(),
			ErrCh:      s.ErrorChan(),
		}

		return srv.Default(nil, nil)
	}

	os.Setenv("GRPC_PORT", "1234")

	mw, err := microwave.New("my-service").
		CheckEnv("GRPC_PORT").
		AddServer(s1).
		Init()

	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(-1)
	}

	mw.Start()
}
