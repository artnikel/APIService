package repository

import (
	"context"
	"strconv"
	"testing"

	"github.com/artnikel/APIService/internal/model"
	bproto "github.com/artnikel/BalanceService/proto"
	"github.com/artnikel/BalanceService/proto/mocks"
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
	client := new(mocks.BalanceServiceClient)
	strOperation := strconv.FormatFloat(testBalance.Operation, 'f', -1, 64)
	client.On("BalanceOperation", mock.Anything, mock.Anything).
		Return(&bproto.BalanceOperationResponse{Operation: strOperation}, nil)
	rep := NewBalanceRepository(client)
	_, err := rep.BalanceOperation(context.Background(), testBalance)
	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestGetBalance(t *testing.T) {
	client := new(mocks.BalanceServiceClient)
	client.On("GetBalance", mock.Anything, mock.Anything).
		Return(&bproto.GetBalanceResponse{Money: testBalance.Operation}, nil)
	rep := NewBalanceRepository(client)
	_, err := rep.GetBalance(context.Background(), testBalance.ProfileID)
	require.NoError(t, err)
	client.AssertExpectations(t)
}
