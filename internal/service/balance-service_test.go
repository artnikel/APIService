package service

import (
	"context"
	"testing"

	"github.com/artnikel/APIService/internal/model"
	"github.com/artnikel/APIService/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testBalance = &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: uuid.New(),
		Operation: 100.9,
	}
)

func TestBalanceOperation(t *testing.T) {
	rep := new(mocks.BalanceRepository)
	srv := NewBalanceService(rep)
	rep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation, nil).Once()
	_, err := srv.BalanceOperation(context.Background(), testBalance)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

func TestGetBalanceAndOperation(t *testing.T) {
	rep := new(mocks.BalanceRepository)
	srv := NewBalanceService(rep)
	rep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation, nil).Once()
	money, err := srv.GetBalance(context.Background(), testBalance.ProfileID)
	require.NoError(t, err)
	require.Equal(t, money, testBalance.Operation)
	rep.AssertExpectations(t)
}

func TestGetIDByToken(t *testing.T) {
	rep := new(mocks.BalanceRepository)
	srv := NewBalanceService(rep)
	urep := new(mocks.UserRepository)
	usrv := NewUserService(urep)
	rep.On("GetIDByToken", mock.AnythingOfType("string")).Return(testBalance.ProfileID, nil).Once()
	tokens, err := usrv.GenerateTokenPair(testBalance.ProfileID)
	require.NoError(t, err)
	profileid, err := srv.GetIDByToken("Bearer " + tokens.AccessToken)
	require.NoError(t, err)
	require.Equal(t, profileid, testBalance.ProfileID)
}
