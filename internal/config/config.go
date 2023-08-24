// Package config with environment variables
package config

import "github.com/caarlos0/env"

// Variables is a struct with environment variables
type Variables struct {
	TokenSignature string `env:"TOKEN_SIGNATURE"`
	TradingApiPort int    `env:"TRADING_API_PORT" `
}

// New returns parsed object of config
func New() (*Variables, error) {
	cfg := &Variables{}
	err := env.Parse(cfg)
	return cfg, err
}
