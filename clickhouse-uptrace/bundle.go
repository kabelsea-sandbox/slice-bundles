package clickhouse

import (
	"time"

	"github.com/uptrace/go-clickhouse/ch"

	"github.com/kabelsea-sandbox/slice"
)

type Bundle struct {
	Addr string `envconfig:"CLICKHOUSE_ADDR" default:"clickhouse:9000" required:"true"`
	Auth struct {
		Database string `envconfig:"DATABASE" required:"true"`
		Username string `envconfig:"USERNAME" default:"default"`
		Password string `envconfig:"PASSWORD"`
	} `envconfig:"CLICKHOUSE_AUTH"`
	DialTimeout     time.Duration `envconfig:"CLICKHOUSE_DIAL_TIMEOUT" default:"5s" required:"True"`
	PoolSize        int           `envconfig:"CLICKHOUSE_POOL_SIZE" default:"2" required:"True"`
	ConnMaxLifetime time.Duration `envconfig:"CLICKHOUSE_CONN_MAX_LIFETIME" default:"10m" required:"True"`
}

// Build implements Bundle.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewClient)
}

// NewClient creates clickhouse client.
func (b *Bundle) NewClient(logger slice.Logger) (*ch.DB, error) {
	logger.Infof("clickhouse", "Create uptrace clickhouse connection")

	return ch.Connect(
		ch.WithAddr(b.Addr),

		// auth
		ch.WithDatabase(b.Auth.Database),
		ch.WithUser(b.Auth.Username),
		ch.WithPassword(b.Auth.Password),

		ch.WithDialTimeout(b.DialTimeout),
		ch.WithPoolSize(b.PoolSize),
		ch.WithConnMaxIdleTime(b.ConnMaxLifetime),
	), nil
}
