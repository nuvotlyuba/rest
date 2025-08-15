package server

import (
	"net"
	"strconv"
	"time"
)

type (
	TimeoutConfig struct {
		Handler time.Duration `env:"handler"`
		Read    time.Duration `env:"read"`
	}

	Config struct {
		Host    string        `env:"host"`
		Port    int           `env:"port"`
		Timeout TimeoutConfig `env:"timeout"`
	}
)

const (
	timeoutHandler = 5 * time.Second
	timeoutRead    = 5 * time.Second

	DefaultPort = 8080
)

func NewConfig() Config {
	return Config{
		Host: "0.0.0.0",
		Port: DefaultPort,
		Timeout: TimeoutConfig{
			Handler: timeoutHandler,
			Read:    timeoutRead,
		},
	}
}

func (c Config) getTimeoutHandler() time.Duration {
	if c.Timeout.Handler <= 0 {
		return timeoutHandler
	}

	return c.Timeout.Handler
}

func (c Config) getTimeoutRead() time.Duration {
	if c.Timeout.Read <= 0 {
		return timeoutHandler
	}

	return c.Timeout.Read
}

func (c Config) getPort() int {
	if c.Port <= 0 {
		return DefaultPort
	}
	return c.Port
}

func (c Config) getAddress() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.getPort()))
}
