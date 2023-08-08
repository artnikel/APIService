package repository

import (
	"context"
	"testing"

	"github.com/artnikel/APIService/internal/model"
	"github.com/artnikel/ProfileService/uproto"
	"github.com/artnikel/ProfileService/uproto/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testUser = model.User{
		ID:       uuid.New(),
		Login:    "testLogin",
		Password: []byte("testPassword"),
		RefreshToken: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
		eyJleHAiOjE2OTE1MzE2NzAsImlkIjoiMjE5NDkxNjctNTRhOC00NjAwLTk1NzMtM2EwYzAyZTE4NzFjIn0.
		RI9lxDrDlj0RS3FAtNSdwFGz14v9NX1tOxmLjSpZ2dU`,
	}
)

func TestSignUp(t *testing.T) {
	client := new(mocks.UserServiceClient)
	client.On("SignUp", mock.Anything, mock.Anything).
		Return(&uproto.SignUpResponse{Id: testUser.ID.String()}, nil)
	rep := NewProfileRepository(client)
	err := rep.SignUp(context.Background(), &testUser)
	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestGetByLogin(t *testing.T) {
	client := new(mocks.UserServiceClient)
	client.On("GetByLogin", mock.Anything, mock.Anything).
		Return(&uproto.GetByLoginResponse{Password: string(testUser.Password), Id: testUser.ID.String()}, nil)
	rep := NewProfileRepository(client)
	password, id, err := rep.GetByLogin(context.Background(), testUser.Login)
	require.NoError(t, err)
	require.Equal(t, password, testUser.Password)
	require.Equal(t, id, testUser.ID)
	client.AssertExpectations(t)
}

func TestAddRefreshToken(t *testing.T) {
	client := new(mocks.UserServiceClient)
	client.On("AddRefreshToken", mock.Anything, mock.Anything).
		Return(nil, nil)
	rep := NewProfileRepository(client)
	err := rep.AddRefreshToken(context.Background(), testUser.ID, testUser.RefreshToken)
	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestGetRefreshTokenByID(t *testing.T) {
	client := new(mocks.UserServiceClient)
	client.On("GetRefreshTokenByID", mock.Anything, mock.Anything).
		Return(&uproto.GetRefreshTokenByIDResponse{RefreshToken: testUser.RefreshToken}, nil)
	rep := NewProfileRepository(client)
	refreshToken, err := rep.GetRefreshTokenByID(context.Background(), testUser.ID)
	require.NoError(t, err)
	require.Equal(t, refreshToken, testUser.RefreshToken)
	client.AssertExpectations(t)
}

func TestDeleteAccount(t *testing.T) {
	client := new(mocks.UserServiceClient)
	client.On("DeleteAccount", mock.Anything, mock.Anything).
		Return(&uproto.DeleteAccountResponse{Id: testUser.ID.String()}, nil)
	rep := NewProfileRepository(client)
	id, err := rep.DeleteAccount(context.Background(), testUser.ID)
	require.Equal(t, id, testUser.ID.String())
	require.NoError(t, err)
	client.AssertExpectations(t)
}
