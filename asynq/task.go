package asynq

import (
	"github.com/hibiken/asynq"
)

type Handler interface {
	asynq.Handler

	// Task return task name
	Task() string
}
