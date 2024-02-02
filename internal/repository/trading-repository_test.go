package repository

import (
	"context"
	"testing"
	"time"

	"github.com/artnikel/APIService/internal/model"
	tproto "github.com/artnikel/TradingService/proto"
	"github.com/artnikel/TradingService/proto/mocks"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		Profit:        decimal.NewFromFloat(300),
		EndDealTime:   time.Now().UTC(),
	}
	testShare = &model.Share{
		Company: "Microsoft",
		Price:   999,
	}
	testProfit = 150.0
)

func TestCreatePosition(t *testing.T) {
	client := new(mocks.TradingServiceClient)
	client.On("CreatePosition", mock.Anything, mock.Anything).
		Return(&tproto.CreatePositionResponse{}, nil)
	rep := NewTradingRepository(client)
	err := rep.CreatePosition(context.Background(), testDeal)
	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestClosePositionManually(t *testing.T) {
	client := new(mocks.TradingServiceClient)
	client.On("ClosePositionManually", mock.Anything, mock.Anything).
		Return(&tproto.ClosePositionManuallyResponse{Profit: testProfit}, nil)
	rep := NewTradingRepository(client)
	_, err := rep.ClosePositionManually(context.Background(), testDeal.DealID, testDeal.ProfileID)
	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestGetUnclosedPositions(t *testing.T) {
	client := new(mocks.TradingServiceClient)
	var testDeals []*tproto.Deal
	protoDeal := tproto.Deal{
		DealID:        testDeal.DealID.String(),
		SharesCount:   testDeal.SharesCount.InexactFloat64(),
		ProfileID:     testDeal.ProfileID.String(),
		Company:       testDeal.Company,
		PurchasePrice: testDeal.PurchasePrice.InexactFloat64(),
		StopLoss:      testDeal.StopLoss.InexactFloat64(),
		TakeProfit:    testDeal.TakeProfit.InexactFloat64(),
		DealTime:      timestamppb.New(testDeal.DealTime),
	}
	testDeals = append(testDeals, &protoDeal)
	client.On("GetUnclosedPositions", mock.Anything, mock.Anything).
		Return(&tproto.GetUnclosedPositionsResponse{Deal: testDeals}, nil)
	rep := NewTradingRepository(client)
	_, err := rep.GetUnclosedPositions(context.Background(), testDeal.ProfileID)
	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestGetClosedPositions(t *testing.T) {
	client := new(mocks.TradingServiceClient)
	var testDeals []*tproto.Deal
	protoDeal := tproto.Deal{
		DealID:        testDeal.DealID.String(),
		SharesCount:   testDeal.SharesCount.InexactFloat64(),
		ProfileID:     testDeal.ProfileID.String(),
		Company:       testDeal.Company,
		PurchasePrice: testDeal.PurchasePrice.InexactFloat64(),
		StopLoss:      testDeal.StopLoss.InexactFloat64(),
		TakeProfit:    testDeal.TakeProfit.InexactFloat64(),
		DealTime:      timestamppb.New(testDeal.DealTime),
		Profit:        testDeal.Profit.InexactFloat64(),
		EndDealTime:   timestamppb.New(testDeal.EndDealTime),
	}
	testDeals = append(testDeals, &protoDeal)
	client.On("GetClosedPositions", mock.Anything, mock.Anything).
		Return(&tproto.GetClosedPositionsResponse{Deal: testDeals}, nil)
	rep := NewTradingRepository(client)
	_, err := rep.GetClosedPositions(context.Background(), testDeal.ProfileID)
	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestGetPrices(t *testing.T) {
	client := new(mocks.TradingServiceClient)
	var testShares []*tproto.TradingShare
	protoShare := tproto.TradingShare{
		Company: testShare.Company,
		Price:   testShare.Price,
	}
	testShares = append(testShares, &protoShare)
	client.On("GetPrices", mock.Anything, mock.Anything).
		Return(&tproto.GetPricesResponse{Share: testShares}, nil)
	rep := NewTradingRepository(client)
	_, err := rep.GetPrices(context.Background())
	require.NoError(t, err)
	client.AssertExpectations(t)
}
