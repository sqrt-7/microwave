package microwave

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pkg/errors"
	"github.com/sqrt-7/microwave/log"
	"github.com/sqrt-7/microwave/server"
	"github.com/sqrt-7/microwave/tools"
)

const (
	MsgBootError = "BOOT_ERROR"
)

var (
	ErrNamespaceMissing = errors.New("namespace missing")
)

type Service struct {
	namespace  string
	logger     log.Logger
	envs       map[string]string
	servers    []server.Server
	wg         *sync.WaitGroup
	shutdownCh chan struct{}
	errCh      chan error

	envsToCheck []string
	serverFns   []ServerFn
}

type ServerFn func(s *Service) server.Server

func New(namespace string) *Service {
	s := &Service{
		namespace:  namespace,
		envs:       make(map[string]string),
		wg:         &sync.WaitGroup{},
		shutdownCh: make(chan struct{}),
		errCh:      make(chan error),
	}

	return s
}

func (s *Service) Logger() log.Logger {
	return s.logger
}

func (s *Service) Env(env string) string {
	v, ok := s.envs[env]
	if !ok {
		return ""
	}

	return v
}

func (s *Service) WaitGroup() *sync.WaitGroup {
	return s.wg
}

func (s *Service) ShutdownChan() chan struct{} {
	return s.shutdownCh
}

func (s *Service) ErrorChan() chan error {
	return s.errCh
}

func (s *Service) WithLogger(logger log.Logger) *Service {
	s.logger = logger
	return s
}

func (s *Service) CheckEnv(envs ...string) *Service {
	s.envsToCheck = envs
	return s
}

func (s *Service) AddServer(srvFn ServerFn) *Service {
	if s.serverFns == nil {
		s.serverFns = []ServerFn{srvFn}
	} else {
		s.serverFns = append(s.serverFns, srvFn)
	}

	return s
}

func (s *Service) Init() (*Service, error) {
	if s.namespace == "" {
		return nil, ErrNamespaceMissing
	}

	if s.logger == nil {
		logger, err := log.NewDefault(s.namespace)
		if err != nil {
			return nil, errors.Wrap(err, MsgBootError)
		}

		s.logger = logger
	}

	if len(s.envsToCheck) > 0 {
		envs, err := tools.EnvLookup(s.envsToCheck...)
		if err != nil {
			return nil, errors.Wrap(err, MsgBootError)
		}
		s.envs = envs
	}

	for _, srvFn := range s.serverFns {
		s.wg.Add(1)
		srv := srvFn(s)
		if s.servers == nil {
			s.servers = []server.Server{srv}
		} else {
			s.servers = append(s.servers, srv)
		}
	}

	return s, nil
}

func (s *Service) Start() {
	defer s.Stop()

	// Listen to system signals
	osSig := make(chan os.Signal)
	signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM)

	for _, srv := range s.servers {
		go srv.Run()
	}

mainLoop:
	for {
		select {
		case e := <-s.errCh:
			{
				s.logger.Error("SERVICE_ERROR").WithField("error", e.Error())
				break mainLoop
			}
		case sig := <-osSig:
			{
				s.logger.Info("SERVICE_TERM").WithField("sig", sig.String())
				break mainLoop
			}
		}
	}
}

func (s *Service) Stop() {
	close(s.shutdownCh)
	s.wg.Wait()
}
