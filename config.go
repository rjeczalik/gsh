package gsh

import (
	"context"
	"time"

	"golang.org/x/crypto/ssh"
)

type Option interface {
	Apply(*Config) *Config
}

type Heartbeat struct {
	Interval time.Duration
	MaxCount int
}

type Retry struct {
	Timeout  time.Duration
	MaxCount int
}

type Config struct {
	ssh.ClientConfig `json:"-" yaml:"-"`

	Network     string
	Address     string
	KeepAlive   bool
	Connection  Retry
	ServerAlive Heartbeat
}

func (cfg *Config) With(opts ...Option) *Config {
	for _, opt := range opts {
		cfg = opt.Apply(cfg)
	}
	return cfg
}

func (cfg *Config) WithAuth(methods ...ssh.AuthMethod) *Config {
	cfg.ClientConfig.Auth = append(cfg.ClientConfig.Auth, methods...)
	return cfg
}

var _ = &Client{ConfigCallback: new(Config).Callback()}

func (cfg *Config) Callback() ConfigCallback {
	return func(context.Context, string, string) (*Config, error) { return cfg, nil }
}
