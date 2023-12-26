package service

import (
	"context"
	"testing"

	"github.com/artnikel/APIService/internal/config"
	"github.com/artnikel/APIService/internal/model"
	"github.com/artnikel/APIService/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testUser = model.User{
		ID:       uuid.New(),
		Login:    "testLogin",
		Password: "testPassword",
	}
	cfg config.Variables
)

func TestSignUp(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, cfg)
	rep.On("SignUp", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
	err := srv.SignUp(context.Background(), &testUser)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

// func TestLogin(t *testing.T) {
// 	rep := new(mocks.UserRepository)

// 	hashedbytes, err := bcrypt.GenerateFromPassword(testUser.Password, bcryptCost)
// 	require.NoError(t, err)

// 	rep.On("GetByLogin", mock.Anything, mock.AnythingOfType("string")).
// 		Return(hashedbytes, testUser.ID, nil)

// 	srv := NewUserService(rep, cfg)

// 	err = srv.Login(context.Background(), &testUser)
// 	require.NoError(t, err)
// 	rep.AssertExpectations(t)
// }

func TestDeleteAccount(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, cfg)

	rep.On("DeleteAccount", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(testUser.ID.String(), nil)
	_, err := srv.DeleteAccount(context.Background(), testUser.ID)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

// func TestGenerateHash(t *testing.T) {
// 	rep := new(mocks.UserRepository)
// 	srv := NewUserService(rep, cfg)
// 	testBytes := []byte("test")
// 	_, err := srv.GenerateHash(testBytes)
// 	require.NoError(t, err)
// 	rep.AssertExpectations(t)
// }

// func TestCheckPasswordHash(t *testing.T) {
// 	rep := new(mocks.UserRepository)
// 	srv := NewUserService(rep, cfg)
// 	testBytes := []byte("test")
// 	hashedBytes, err := srv.GenerateHash(testBytes)
// 	require.NoError(t, err)
// 	isEqual, err := srv.CheckPasswordHash(hashedBytes, testBytes)
// 	require.NoError(t, err)
// 	require.True(t, isEqual)
// 	rep.AssertExpectations(t)
// }
