package repository

import (
	"context"
	"fmt"

	berrors "github.com/artnikel/APIService/internal/errors"
	"github.com/artnikel/APIService/internal/model"
	tproto "github.com/artnikel/TradingService/proto"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TradingRepository represents the client of Trading Service repository implementation.
type TradingRepository struct {
	client tproto.TradingServiceClient
}

// NewTradingRepository creates and returns a new instance of TradingRepository, using the provided proto.TradingServiceClient.
func NewTradingRepository(client tproto.TradingServiceClient) *TradingRepository {
	return &TradingRepository{
		client: client,
	}
}

// CreatePosition call a method of TradingService.
func (r *TradingRepository) CreatePosition(ctx context.Context, deal *model.Deal) error {
	_, err := r.client.CreatePosition(ctx, &tproto.CreatePositionRequest{
		Deal: &tproto.Deal{
			DealID:        deal.DealID.String(),
			SharesCount:   deal.SharesCount.InexactFloat64(),
			ProfileID:     deal.ProfileID.String(),
			Company:       deal.Company,
			PurchasePrice: deal.PurchasePrice.InexactFloat64(),
			StopLoss:      deal.StopLoss.InexactFloat64(),
			TakeProfit:    deal.TakeProfit.InexactFloat64(),
			DealTime:      timestamppb.New(deal.DealTime),
			EndDealTime:   timestamppb.New(deal.EndDealTime),
			Profit:        deal.Profit.InexactFloat64(),
		},
	})
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if ok && grpcStatus.Message() == berrors.NotEnoughMoney {
			return berrors.New(berrors.NotEnoughMoney, "Not enough money")
		}
		if ok && grpcStatus.Message() == berrors.PurchasePriceOut {
			return berrors.New(berrors.PurchasePriceOut, "Purchase price out of stoploss/takeprofit")
		}
		return fmt.Errorf("createPosition %w", err)
	}
	return nil
}

// ClosePositionManually call a method of TradingService.
func (r *TradingRepository) ClosePositionManually(ctx context.Context, dealid, profileid uuid.UUID) (float64, error) {
	resp, err := r.client.ClosePositionManually(ctx, &tproto.ClosePositionManuallyRequest{
		Dealid:    dealid.String(),
		Profileid: profileid.String(),
	})
	if err != nil {
		return 0, fmt.Errorf("closePositionManually %w", err)
	}
	return resp.Profit, nil
}

// GetUnclosedPositions call a method of TradingService.
func (r *TradingRepository) GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error) {
	resp, err := r.client.GetUnclosedPositions(ctx, &tproto.GetUnclosedPositionsRequest{
		Profileid: profileid.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("getUncosedPositions %w", err)
	}
	unclosedDeals := make([]*model.Deal, len(resp.Deal))
	for i, deal := range resp.Deal {
		dealUUID, err := uuid.Parse(deal.DealID)
		if err != nil {
			return nil, fmt.Errorf("parse %w", err)
		}
		unclosedDeal := &model.Deal{
			DealID:        dealUUID,
			SharesCount:   decimal.NewFromFloat(deal.SharesCount),
			Company:       deal.Company,
			PurchasePrice: decimal.NewFromFloat(deal.PurchasePrice),
			StopLoss:      decimal.NewFromFloat(deal.StopLoss),
			TakeProfit:    decimal.NewFromFloat(deal.TakeProfit),
			DealTime:      deal.DealTime.AsTime(),
		}
		unclosedDeals[i] = unclosedDeal
	}
	return unclosedDeals, nil
}

// GetClosedPositions call a method of TradingService.
func (r *TradingRepository) GetClosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error) {
	resp, err := r.client.GetClosedPositions(ctx, &tproto.GetClosedPositionsRequest{
		Profileid: profileid.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("getClosedPositions %w", err)
	}
	closedDeals := make([]*model.Deal, len(resp.Deal))
	for i, deal := range resp.Deal {
		dealUUID, err := uuid.Parse(deal.DealID)
		if err != nil {
			return nil, fmt.Errorf("parse %w", err)
		}
		closedDeal := &model.Deal{
			DealID:        dealUUID,
			SharesCount:   decimal.NewFromFloat(deal.SharesCount),
			Company:       deal.Company,
			PurchasePrice: decimal.NewFromFloat(deal.PurchasePrice),
			StopLoss:      decimal.NewFromFloat(deal.StopLoss),
			TakeProfit:    decimal.NewFromFloat(deal.TakeProfit),
			DealTime:      deal.DealTime.AsTime(),
			Profit:        decimal.NewFromFloat(deal.Profit),
			EndDealTime:   deal.EndDealTime.AsTime(),
		}
		closedDeals[i] = closedDeal
	}
	return closedDeals, nil
}

// GetPrices call a method of TradingService.
func (r *TradingRepository) GetPrices(ctx context.Context) ([]model.Share, error) {
	resp, err := r.client.GetPrices(ctx, &tproto.GetPricesRequest{})
	if err != nil {
		return nil, fmt.Errorf("getPrices %w", err)
	}
	allShares := make([]model.Share, len(resp.Share))
	for i, share := range resp.Share {
		allShare := model.Share{
			Company: share.Company,
			Price:   share.Price,
		}
		allShares[i] = allShare
	}
	return allShares, nil
}
