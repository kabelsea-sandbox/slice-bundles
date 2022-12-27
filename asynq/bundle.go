package asynq

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/kabelsea-sandbox/slice"
	"github.com/kabelsea-sandbox/slice/pkg/di"
	"github.com/kabelsea-sandbox/slice/pkg/run"
)

// Bundle struct
type Bundle struct {
	client *asynq.Client
	server *asynq.Server

	Client bool // if set True bootstrap asynq client
	Server bool // if set True bootstrap asynq server

	Redis struct {
		Addr         string        `envconfig:"ADDR" default:"redis:6379" required:"true"`
		Username     string        `envconfig:"USERNAME" default:""`
		Password     string        `envconfig:"PASSWORD" default:""`
		Database     int           `envconfig:"DATABASE"`
		DialTimeout  time.Duration `envconfig:"DIAL_TIMEOUT" default:"5s" required:"True"`
		ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"2s" required:"True"`
		WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"2s" required:"True"`
		PoolSize     int           `envconfig:"POOL_SIZE" default:"10"`
	} `envconfig:"ASYNQ_REDIS" required:"True"`

	ServerConfig struct {
		Concurrency     int            `envconfig:"CONCURRENCY" default:"1" required:"True"`
		Queues          map[string]int `envconfig:"QUEUES"`
		StrictPriority  bool           `envconfig:"STRICT_PRIORITY"`
		ShutdownTimeout time.Duration  `envconfig:"SHUTDOWN_TIMEOUT" default:"5s" required:"True"`
	} `envconfig:"ASYNQ_CONFIG"`
}

// Build implements Bundle.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	if b.Server {
		builder.Provide(NewServerWorker, di.As(new(run.Worker)))

		builder.Provide(b.NewServer)
		builder.Provide(b.NewServeMux)
	}

	if b.Client {
		builder.Provide(b.NewClient)
	}
}

// Shutdown implements Bundle interface.
func (b *Bundle) Shutdown(_ context.Context, _ slice.Container) (err error) {
	if b.Client {
		_ = b.client.Close()
	}
	return nil
}

func (b *Bundle) redisOpts() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:         b.Redis.Addr,
		Username:     b.Redis.Username,
		Password:     b.Redis.Password,
		DB:           b.Redis.Database,
		DialTimeout:  b.Redis.DialTimeout,
		ReadTimeout:  b.Redis.ReadTimeout,
		WriteTimeout: b.Redis.WriteTimeout,
	}
}

func (b *Bundle) NewClient() *asynq.Client {
	return asynq.NewClient(b.redisOpts())
}

func (b *Bundle) NewServer() *asynq.Server {
	return asynq.NewServer(
		b.redisOpts(),
		asynq.Config{
			Concurrency:     b.ServerConfig.Concurrency,
			Queues:          b.ServerConfig.Queues,
			StrictPriority:  b.ServerConfig.StrictPriority,
			ShutdownTimeout: b.ServerConfig.ShutdownTimeout,
		},
	)
}

func (b *Bundle) NewServeMux() *asynq.ServeMux {
	return asynq.NewServeMux()
}
