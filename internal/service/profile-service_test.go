package service

import (
	"context"
	"crypto/sha256"
	"log"
	"os"
	"testing"

	"github.com/artnikel/APIService/internal/config"
	"github.com/artnikel/APIService/internal/model"
	"github.com/artnikel/APIService/internal/service/mocks"
	"github.com/caarlos0/env"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var (
	testUser = model.User{
		ID:           uuid.New(),
		Login:        "testLogin",
		Password:     []byte("testPassword"),
		RefreshToken: "",
	}
	cfg config.Variables
)

func TestMain(m *testing.M) {
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("could not parse config: ", err)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestSignUp(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)
	rep.On("SignUp", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
	err := srv.SignUp(context.Background(), &testUser)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	rep := new(mocks.UserRepository)

	hashedbytes, err := bcrypt.GenerateFromPassword(testUser.Password, bcryptCost)
	require.NoError(t, err)

	rep.On("GetByLogin", mock.Anything, mock.AnythingOfType("string")).
		Return(hashedbytes, testUser.ID, nil)
	rep.On("AddRefreshToken", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("string")).
		Return(nil)

	srv := NewUserService(rep, &cfg)

	_, err = srv.Login(context.Background(), &testUser)
	require.NoError(t, err)
}

func TestRefresh(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)

	tokenPair, err := srv.GenerateTokenPair(testUser.ID)
	require.NoError(t, err)
	sum := sha256.Sum256([]byte(tokenPair.RefreshToken))

	hashedbytes, err := bcrypt.GenerateFromPassword(sum[:], bcryptCost)
	require.NoError(t, err)

	rep.On("GetRefreshTokenByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(string(hashedbytes), nil)
	rep.On("AddRefreshToken", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("string")).
		Return(nil)

	_, err = srv.Refresh(context.Background(), tokenPair)
	require.NoError(t, err)
}

func TestDeleteAccount(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)

	rep.On("DeleteAccount", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(testUser.ID.String(), nil)
	_, err := srv.DeleteAccount(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestTokensIDCompare(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)
	tokenPair, err := srv.GenerateTokenPair(testUser.ID)
	require.NoError(t, err)
	id, err := srv.TokensIDCompare(tokenPair)
	require.NoError(t, err)
	require.Equal(t, testUser.ID, id)
}

func TestGenerateHash(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)
	testBytes := []byte("test")
	_, err := srv.GenerateHash(testBytes)
	require.NoError(t, err)
}

func TestCheckPasswordHash(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)
	testBytes := []byte("test")
	hashedBytes, err := srv.GenerateHash(testBytes)
	require.NoError(t, err)
	isEqual, err := srv.CheckPasswordHash(hashedBytes, testBytes)
	require.NoError(t, err)
	require.True(t, isEqual)
}

func TestGenerateTokenPair(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)
	_, err := srv.GenerateTokenPair(testUser.ID)
	require.NoError(t, err)
}

func TestGenerateJWTToken(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)
	_, err := srv.GenerateJWTToken(accessTokenExpiration, testUser.ID)
	require.NoError(t, err)
}