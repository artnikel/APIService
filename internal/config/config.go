// Package config with environment variables
package config

// Variables is a struct with environment variables
type Variables struct {
	TokenSignature string `env:"TOKEN_SIGNATURE"`
}
