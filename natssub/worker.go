package natssub

import (
	"sync"

	"go.uber.org/zap"
	"golang.org/x/net/context"

	"github.com/kabelsea-sandbox/slice/pkg/run"
)

// Worker interface
type Worker interface {
	run.Worker
	Append(s Subscription)
}

// Worker interface implementation
type worker struct {
	logger        *zap.Logger
	subscriptions []Subscription
	stop          chan struct{}
	done          chan struct{}
}

func NewWorker(logger *zap.Logger) Worker {
	return &worker{
		logger:        logger,
		subscriptions: make([]Subscription, 0, 100),
		stop:          make(chan struct{}, 1),
		done:          make(chan struct{}, 1),
	}
}

// Append build Subscription and add to subscriptions pool
func (w *worker) Append(s Subscription) {
	w.subscriptions = append(w.subscriptions, s)
}

// Start implements Worker interface
func (w *worker) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	for _, s := range w.subscriptions {
		wg.Add(1)

		sub := s

		go func() {
			defer wg.Done()

			go func() {
				if err := sub.Run(ctx); err != nil {
					w.logger.Fatal("subscription run failed", zap.Error(err))
				}
			}()
		}()
	}

	<-w.stop

	// wait all workers
	wg.Wait()

	// close done channel
	close(w.done)

	return nil
}

// Stop implements Worker interface
func (w *worker) Stop(err error) {
	close(w.stop)

	for _, sub := range w.subscriptions {
		sub.Stop(err)
	}

	// wait done
	<-w.done
}
