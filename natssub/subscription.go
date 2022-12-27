package natssub

import (
	"context"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Pointer type helper
type ptr[T any] interface {
	*T
}

// MessageHandler interface
type MessageHandler func(ctx context.Context, msg *nats.Msg) error

// Subscription interface
type Subscription interface {
	// Run
	Run(ctx context.Context) error

	// Stop
	Stop(err error)
}

// Subscription interface implementation
type subscription struct {
	logger   *zap.Logger
	nc       *nats.Conn
	jsm      nats.JetStreamContext
	sub      *nats.Subscription
	options  *Options
	messages chan *nats.Msg
	stop     chan struct{}
	done     chan struct{}
}

// NewSubscription constructor
func NewSubscription(
	logger *zap.Logger,
	nc *nats.Conn,
	jsm nats.JetStreamContext,
	options *Options,
) (Subscription, error) {
	logger = logger.With(
		zap.Namespace("nats/subscription"),
		zap.String("subject", options.subject),
	)

	if options.queue != "" {
		logger = logger.With(zap.String("queue", options.queue))
	}

	sub := &subscription{
		logger:  logger,
		nc:      nc,
		jsm:     jsm,
		options: options,
		stop:    make(chan struct{}, 1),
		done:    make(chan struct{}, 1),
	}

	if options.buffer != 0 {
		sub.messages = make(chan *nats.Msg, options.buffer)
	} else {
		sub.messages = make(chan *nats.Msg)
	}

	return sub, nil
}

// Start implement Subscription interface
func (s *subscription) Run(ctx context.Context) error {
	var (
		wg  sync.WaitGroup
		sub *nats.Subscription
		err error

		queue       = s.options.queue
		subject     = s.options.subject
		concurrency = s.options.concurrency
	)

	logger := s.logger.With(
		zap.String("queue", queue),
		zap.String("subject", subject),
		zap.Int("concurrency", concurrency),
	)

	if s.options.jetstream {
		opts := []nats.SubOpt{
			nats.ManualAck(),
		}

		if queue != "" {
			sub, err = s.jsm.ChanQueueSubscribe(subject, queue, s.messages, opts...)
		} else {
			sub, err = s.jsm.ChanSubscribe(subject, s.messages, opts...)
		}
	} else {
		if queue != "" {
			sub, err = s.nc.ChanQueueSubscribe(subject, queue, s.messages)
		} else {
			sub, err = s.nc.ChanSubscribe(subject, s.messages)
		}
	}

	if err != nil {
		return errors.Wrap(err, "subscription failed")
	}

	s.sub = sub

	// create some workers
	for i := 1; i <= concurrency; i++ {
		wg.Add(1)

		logger.Info("[SUBSCRIPTIONS] Created")

		go func() {
			defer wg.Done()

			for {
				select {
				case m := <-s.messages:
					s.handle(ctx, m)
				case <-s.stop:
					logger.Info("graceful shutdown")

					if err := s.sub.Unsubscribe(); err != nil {
						logger.Warn("unsubscribe failed", zap.Error(err))
					}

					// empty messages channel
					for m := range s.messages {
						s.handle(ctx, m)
					}
					return
				}
			}
		}()
	}

	// wait all workers
	wg.Wait()

	// send done signal
	s.done <- struct{}{}

	return nil
}

// Stop implement Subscription interface
func (s *subscription) Stop(err error) {
	close(s.stop)

	// wait done
	<-s.done

	s.logger.Info("stopped")
}

// Handle implement Subscription interface
func (s *subscription) handle(ctx context.Context, m *nats.Msg) {
	logger := s.logger.With(
		zap.Any("message", m),
	)

	if err := s.options.handler(ctx, m); err != nil {
		logger.Error("message handle failed", []zap.Field{
			zap.Error(err),
		}...)
	}

	if s.options.jetstream {
		if err := m.Ack(); err != nil {
			logger.Error("ack message failed", zap.Error(err))
		}
	}
}
