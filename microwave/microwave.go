package microwave

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sqrt-7/microwave/log"
)

const (
	MsgBootError = "BOOT_ERROR"
)

var (
	ErrNamespaceMissing = errors.New("namespace missing")
)

type Microwave struct {
	namespace string
	logger    log.Logger
	envs      map[string]string
	servers   []ServerWrapper

	wg         *sync.WaitGroup
	shutdownCh chan struct{}
	errCh      chan error
}

func New(namespace string, options ...Option) (*Microwave, error) {
	if namespace == "" {
		return nil, ErrNamespaceMissing
	}

	s := &Microwave{
		namespace:  namespace,
		envs:       make(map[string]string),
		servers:    make([]ServerWrapper, 0),
		wg:         &sync.WaitGroup{},
		shutdownCh: make(chan struct{}),
		errCh:      make(chan error),
	}

	for _, opt := range options {
		if err := opt.apply(s); err != nil {
			return nil, errors.Wrap(err, MsgBootError)
		}
	}

	if s.logger == nil {
		logger, err := log.NewDefault(s.namespace)
		if err != nil {
			return nil, errors.Wrap(err, MsgBootError)
		}

		s.logger = logger
	}

	return s, nil
}

func (s Microwave) Namespace() string {
	return s.namespace
}

func (s Microwave) Logger() log.Logger {
	return s.logger
}

func (s Microwave) Env(env string) string {
	v, ok := s.envs[env]
	if !ok {
		return ""
	}

	return v
}

func (s Microwave) WaitGroup() *sync.WaitGroup {
	return s.wg
}

func (s Microwave) ShutdownChan() chan struct{} {
	return s.shutdownCh
}

func (s Microwave) ErrorChan() chan error {
	return s.errCh
}

func (s *Microwave) Start() {
	defer s.Stop()

	// Listen to system signals
	osSig := make(chan os.Signal, 2)
	signal.Notify(osSig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		time.Sleep(time.Second * 3)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	for _, srv := range s.servers {
		s.wg.Add(1)
		go srv.Run(s)
	}

mainLoop:
	for {
		select {
		case e := <-s.errCh:
			{
				s.logger.Error("SERVICE_ERROR").WithField("error", e.Error()).Send()
				break mainLoop
			}
		case sig := <-osSig:
			{
				s.logger.Info("SERVICE_TERM").WithField("sig", sig.String()).Send()
				break mainLoop
			}
		}
	}
}

func (s *Microwave) Stop() {
	close(s.shutdownCh)
	s.wg.Wait()
}

func (s *Microwave) AddServer(srv ServerWrapper) {
	s.servers = append(s.servers, srv)
}
