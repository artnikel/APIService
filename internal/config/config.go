// Package config with environment variables
package config

import "github.com/caarlos0/env"

// Variables is a struct with environment variables
type Variables struct {
	TokenSignature    string `env:"TOKEN_SIGNATURE"`
	APIPort           int    `env:"API_PORT"`
	RedisPriceAddress string `env:"REDIS_PRICE_ADDRESS"`
	TradingAddress    string `env:"TRADING_ADDRESS"`
	ProfileAddress    string `env:"PROFILE_ADDRESS"`
	BalanceAddress    string `env:"BALANCE_ADDRESS"`
}

// New returns parsed object of config
func New() (*Variables, error) {
	cfg := &Variables{}
	err := env.Parse(cfg)
	return cfg, err
}
