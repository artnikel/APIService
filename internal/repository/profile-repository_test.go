package repository

import (
	"context"
	"testing"

	"github.com/artnikel/APIService/internal/model"
	uproto "github.com/artnikel/ProfileService/proto"
	"github.com/artnikel/ProfileService/proto/mocks"
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
		Return(&uproto.GetByLoginResponse{Password: testUser.Password, Id: testUser.ID.String()}, nil)
	rep := NewProfileRepository(client)
	_, id, err := rep.GetByLogin(context.Background(), testUser.Login)
	require.NoError(t, err)
	require.Equal(t, id, testUser.ID)
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
