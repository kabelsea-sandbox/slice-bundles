package geoip

import (
	"github.com/kabelsea-sandbox/slice"

	maxminddb "github.com/oschwald/maxminddb-golang"
)

type Bundle struct {
	MaxmindDatabase string `envconfig:"MAXMIND_DATABASE" required:"true"`
}

// Build implements Bundle.
func (b *Bundle) Build(builder slice.ContainerBuilder) {
	builder.Provide(b.NewConnection)
}

// NewConnection creates maxmind reader.
func (b *Bundle) NewConnection(logger slice.Logger) (*maxminddb.Reader, error) {
	logger.Infof("geoip", "Create geoip connection")

	return maxminddb.Open(b.MaxmindDatabase)
}
