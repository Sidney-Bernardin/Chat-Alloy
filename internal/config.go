package internal

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	ADDR string `required:"true"`

	SESSION_DURATION      time.Duration `required:"true"`
	SESSION_COOKIE_DOMAIN string        `required:"true"`

	POSTGRES_URL string `required:"true"`
	REDIS_ADDR   string `required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("APP", &cfg); err != nil {
		return nil, errors.Wrap(err, "cannot process environment")
	}
	return &cfg, nil
}
