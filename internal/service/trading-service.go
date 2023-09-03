package service

import (
	"context"
	"fmt"

	"github.com/artnikel/APIService/internal/model"
	"github.com/google/uuid"
)

// TradingRepository is an interface that contains methods for long or short strategies
type TradingRepository interface {
	CreatePosition(ctx context.Context, deal *model.Deal) error
	ClosePositionManually(ctx context.Context, dealid, profileid uuid.UUID) (float64, error)
	GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error)
	GetPrices(ctx context.Context) ([]model.Share, error)
}

// TradingService contains BalanceRepository interface
type TradingService struct {
	tRep TradingRepository
}

// NewTradingService accepts TradingRepository object and returnes an object of type *TradingService
func NewTradingService(tRep TradingRepository) *TradingService {
	return &TradingService{tRep: tRep}
}

// CreatePosition is a method of TradingService calls method of Repository
func (ts *TradingService) CreatePosition(ctx context.Context, deal *model.Deal) error {
	err := ts.tRep.CreatePosition(ctx, deal)
	if err != nil {
		return fmt.Errorf("TradingService-CreatePosition: error:%w", err)
	}
	return nil
}

// ClosePositionManually is a method of TradingService calls method of Repository
func (ts *TradingService) ClosePositionManually(ctx context.Context, dealid, profileid uuid.UUID) (float64, error) {
	profit, err := ts.tRep.ClosePositionManually(ctx, dealid, profileid)
	if err != nil {
		return 0, fmt.Errorf("TradingService-ClosePositionManually: error:%w", err)
	}
	return profit, nil
}

// GetUnclosedPositions is a method of TradingService calls method of Repository
func (ts *TradingService) GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error) {
	unclosedDeals, err := ts.tRep.GetUnclosedPositions(ctx, profileid)
	if err != nil {
		return nil, fmt.Errorf("TradingService-GetUnclosedPositions: error:%w", err)
	}
	return unclosedDeals, nil
}

// GetPrices is a method of TradingService calls method of Repository
func (ts *TradingService) GetPrices(ctx context.Context) ([]model.Share, error) {
	shares, err := ts.tRep.GetPrices(ctx)
	if err != nil {
		return nil, fmt.Errorf("TradingService-GetPrices: error:%w", err)
	}
	return shares, nil
}
