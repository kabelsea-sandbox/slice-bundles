package clickhouse

import (
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"

	"github.com/kabelsea-sandbox/slice"
)

type Bundle struct {
	Addrs []string `envconfig:"CLICKHOUSE_ADDRS" default:"clickhouse:9000" required:"true"`
	Auth  struct {
		Database string `envconfig:"DATABASE" required:"true"`
		Username string `envconfig:"USERNAME" default:"default"`
		Password string `envconfig:"PASSWORD"`
	} `envconfig:"CLICKHOUSE_AUTH"`
	Debug           bool          `envconfig:"CLICKHOUSE_DEBUG" default:"false"`
	DialTimeout     time.Duration `envconfig:"CLICKHOUSE_DIAL_TIMEOUT" default:"5s" required:"True"`
	MaxOpenConns    int           `envconfig:"CLICKHOUSE_MAX_OPEN_CONNS" default:"2" required:"True"`
	MaxIdleConns    int           `envconfig:"CLICKHOUSE_MAX_IDLE_CONNS" default:"2" required:"True"`
	ConnMaxLifetime time.Duration `envconfig:"CLICKHOUSE_CONN_MAX_LIFETIME" default:"10m" required:"True"`
	BlockBufferSize uint8         `envconfig:"CLICKHOUSE_BLOCK_BUFFER_SIZE" default:"10" required:"True"`
	Settings        struct {
		MaxExecutionTime int `envconfig:"MAX_EXECUTION_TIME" default:"15" required:"True"`
	} `envconfig:"CLICKHOUSE_SETTINGS"`
}

// Build implements Bundle.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewClient)
}

// NewClient creates clickhouse client.
func (b *Bundle) NewClient(logger slice.Logger) (clickhouse.Conn, error) {
	logger.Infof("clickhouse", "Create clickhouse connection")

	return clickhouse.Open(&clickhouse.Options{
		Addr: b.Addrs,
		Auth: clickhouse.Auth{
			Database: b.Auth.Database,
			Username: b.Auth.Username,
			Password: b.Auth.Password,
		},
		Debug: b.Debug,
		Settings: clickhouse.Settings{
			"max_execution_time": b.Settings.MaxExecutionTime,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:      b.DialTimeout,
		MaxOpenConns:     b.MaxOpenConns,
		MaxIdleConns:     b.MaxIdleConns,
		ConnMaxLifetime:  b.ConnMaxLifetime,
		BlockBufferSize:  b.BlockBufferSize,
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
	})
}
