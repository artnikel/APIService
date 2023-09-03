package service

import (
	"context"
	"testing"
	"time"

	"github.com/artnikel/APIService/internal/model"
	"github.com/artnikel/APIService/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testDeal = &model.Deal{
		DealID:        uuid.New(),
		SharesCount:   decimal.NewFromFloat(1.5),
		ProfileID:     testBalance.ProfileID,
		Company:       "Apple",
		PurchasePrice: decimal.NewFromFloat(1350),
		StopLoss:      decimal.NewFromFloat(1500),
		TakeProfit:    decimal.NewFromFloat(1000),
		DealTime:      time.Now().UTC(),
	}
	testShare = &model.Share{
		Company: "Microsoft",
		Price:   999,
	}
	testProfit = 150.0
)

func TestCreatePosition(t *testing.T) {
	rep := new(mocks.TradingRepository)
	srv := NewTradingService(rep)
	rep.On("CreatePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).
		Return(nil).Once()
	err := srv.CreatePosition(context.Background(), testDeal)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

func TestClosePositionManually(t *testing.T) {
	rep := new(mocks.TradingRepository)
	srv := NewTradingService(rep)
	rep.On("ClosePositionManually", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(testProfit, nil).Once()
	_, err := srv.ClosePositionManually(context.Background(), testDeal.DealID, testDeal.ProfileID)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

func TestGetUnclosedPositions(t *testing.T) {
	rep := new(mocks.TradingRepository)
	srv := NewTradingService(rep)
	var testDeals []*model.Deal
	testDeals = append(testDeals, testDeal)
	rep.On("GetUnclosedPositions", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(testDeals, nil)
	_, err := srv.GetUnclosedPositions(context.Background(), testDeal.ProfileID)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

func TestGetPrices(t *testing.T) {
	rep := new(mocks.TradingRepository)
	srv := NewTradingService(rep)
	var testPrices []model.Share
	testPrices = append(testPrices, *testShare)
	rep.On("GetPrices", mock.Anything).Return(testPrices, nil)
	_, err := srv.GetPrices(context.Background())
	require.NoError(t, err)
	rep.AssertExpectations(t)
}
