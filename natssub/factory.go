package natssub

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type FactoryBuilder func(opts ...Option) *Options

func Factory(opts ...Option) FactoryBuilder {
	options := &Options{}

	for _, o := range opts {
		o(options)
	}

	return func(opts ...Option) *Options {
		return options
	}
}

func wrappedJetStreamFactoryBuilder(
	logger *zap.Logger,
	jsm nats.JetStreamContext,
	builder FactoryBuilder,
) (Subscription, error) {
	return NewSubscription(logger, nil, jsm, builder())
}

func wrappedNatsFactoryBuilder(
	logger *zap.Logger,
	nc *nats.Conn,
	builder FactoryBuilder,
) (Subscription, error) {
	return NewSubscription(logger, nc, nil, builder())
}
