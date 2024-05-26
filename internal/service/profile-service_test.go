package service

import (
	"testing"

	"github.com/artnikel/APIService/internal/config"
	"github.com/artnikel/APIService/internal/service/mocks"
	"github.com/stretchr/testify/require"
)

var (
	cfg config.Variables
)

func TestGenerateHash(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)
	testBytes := []byte("test")
	_, err := srv.GenerateHash(string(testBytes))
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

func TestCheckPasswordHash(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)
	testBytes := []byte("test")
	hashedBytes, err := srv.GenerateHash(string(testBytes))
	require.NoError(t, err)
	isEqual, err := srv.CheckPasswordHash(hashedBytes, string(testBytes))
	require.NoError(t, err)
	require.True(t, isEqual)
	rep.AssertExpectations(t)
}
