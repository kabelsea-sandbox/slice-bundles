package asynq

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"

	"github.com/kabelsea-sandbox/slice/pkg/run"
)

type ServerWorker interface {
	run.Worker
}

type serverWorker struct {
	server *asynq.Server
	mux    *asynq.ServeMux
	stop   chan struct{}
	done   chan struct{}
}

func NewServerWorker(server *asynq.Server, mux *asynq.ServeMux) ServerWorker {
	return &serverWorker{
		server: server,
		mux:    mux,
		stop:   make(chan struct{}, 1),
		done:   make(chan struct{}, 1),
	}
}

func (s *serverWorker) Run(ctx context.Context) error {
	if err := s.server.Run(s.mux); err != nil {
		return errors.Wrap(err, "server run failed")
	}

	<-s.stop

	s.server.Shutdown()

	close(s.done)

	return nil
}

func (s *serverWorker) Stop(err error) {
	close(s.stop)

	// wait done
	<-s.done
}
