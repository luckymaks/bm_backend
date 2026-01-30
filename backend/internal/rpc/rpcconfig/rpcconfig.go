package rpcconfig

import (
	"context"

	"github.com/caarlos0/env/v11"
	"github.com/cockroachdb/errors"
)

type EnvConfig struct {
	DeploymentIdent string `env:"DEPLOYMENT_IDENT" envDefault:"Dev"`
	MainTableName   string `env:"MAIN_TABLE_NAME"`
	AWSRegion       string `env:"AWS_REGION"`
}

type Config struct {
	Env EnvConfig
}

type Loader struct {
	env EnvConfig
}

func NewLoader() (*Loader, error) {
	envCfg, err := env.ParseAs[EnvConfig]()
	if err != nil {
		return nil, errors.Wrap(err, "parsing environment config")
	}

	return &Loader{
		env: envCfg,
	}, nil
}

func (l *Loader) Load(_ context.Context) (*Config, error) {
	return &Config{
		Env: l.env,
	}, nil
}
