package service

import (
	"context"
	"fmt"

	"github.com/artnikel/APIService/internal/model"
	"github.com/google/uuid"
)

// TradingRepository is an interface that contains methods for long or short strategies
type TradingRepository interface {
	GetProfit(ctx context.Context, strategy string, deal *model.Deal) (float64, error)
	ClosePosition(ctx context.Context, dealid, profileid uuid.UUID) error
	GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error)
}

// TradingService contains BalanceRepository interface
type TradingService struct {
	tRep TradingRepository
}

// NewTradingService accepts TradingRepository object and returnes an object of type *TradingService
func NewTradingService(tRep TradingRepository) *TradingService {
	return &TradingService{tRep: tRep}
}

// GetProfit is a method of TradingService calls method of Repository
func (ts *TradingService) GetProfit(ctx context.Context, strategy string, deal *model.Deal) (float64, error) {
	profit, err := ts.tRep.GetProfit(ctx, strategy, deal)
	if err != nil {
		return 0, fmt.Errorf("TradingService-GetProfit: error:%w", err)
	}
	return profit, nil
}

func (ts *TradingService) ClosePosition(ctx context.Context, dealid, profileid uuid.UUID) error {
	err := ts.tRep.ClosePosition(ctx, dealid,profileid)
	if err != nil {
		return fmt.Errorf("TradingService-ClosePosition: error:%w", err)
	}
	return nil
}

func (ts *TradingService) GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error) {
	unclosedDeals, err := ts.tRep.GetUnclosedPositions(ctx, profileid)
	if err != nil {
		return nil, fmt.Errorf("TradingService-GetUnclosedPositions: error:%w", err)
	}
	return unclosedDeals, nil
}
