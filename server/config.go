package server

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	ADDR string `required:"true"`

	POSTGRES_URL string `required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("SVR", &cfg); err != nil {
		return nil, errors.Wrap(err, "cannot process environment")
	}
	return &cfg, nil
}
