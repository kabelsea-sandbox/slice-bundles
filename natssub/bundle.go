package natssub

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/kabelsea-sandbox/slice"
	"github.com/kabelsea-sandbox/slice/pkg/di"
	"github.com/kabelsea-sandbox/slice/pkg/run"
)

type Bundle struct{}

// Build implements Bundle interface
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(NewWorker, di.As(new(run.Worker)))
}

// Boot implements Bundle interface
func (b *Bundle) Boot(_ context.Context, container slice.Container) (err error) {
	var factories []FactoryBuilder

	if container.Has(&factories) {
		if err := container.Resolve(&factories); err != nil {
			return err
		}

		for _, factory := range factories {
			options := factory()

			if options.jetstream {
				if err := container.Invoke(b.RegisterJetStreamSubscriptions(factory)); err != nil {
					return err
				}
			} else {
				if err := container.Invoke(b.RegisterNatsSubscriptions(factory)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Shutdown implements Bundle interface.
func (b *Bundle) Shutdown(_ context.Context, _ slice.Container) (err error) {
	return nil
}

func (b *Bundle) RegisterNatsSubscriptions(factory FactoryBuilder) func(
	logger *zap.Logger,
	ns *nats.Conn,
	worker Worker,
) error {
	return func(logger *zap.Logger, nc *nats.Conn, worker Worker) error {
		sub, err := wrappedNatsFactoryBuilder(logger, nc, factory)
		if err != nil {
			return err
		}
		worker.Append(sub)
		return nil
	}
}

func (b *Bundle) RegisterJetStreamSubscriptions(factory FactoryBuilder) func(
	logger *zap.Logger,
	jsm nats.JetStreamContext,
	worker Worker,
) error {
	return func(logger *zap.Logger, jsm nats.JetStreamContext, worker Worker) error {
		sub, err := wrappedJetStreamFactoryBuilder(logger, jsm, factory)
		if err != nil {
			return err
		}
		worker.Append(sub)
		return nil
	}
}
