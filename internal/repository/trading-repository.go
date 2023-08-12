package repository

import (
	"context"
	"fmt"

	"github.com/artnikel/APIService/internal/model"
	tproto "github.com/artnikel/TradingService/proto"
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

// Strategies call a method of TradingService.
func (r *TradingRepository) Strategies(ctx context.Context, strategy string, deal *model.Deal) (float64, error) {
	resp, err := r.client.Strategies(ctx, &tproto.StrategiesRequest{
		Strategy: strategy,
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
		return 0, fmt.Errorf("TradingRepository-Strategies: error:%w", err)
	}
	return resp.Profit, nil
}
