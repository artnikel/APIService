package service

import (
	"context"
	"fmt"

	"github.com/artnikel/APIService/internal/model"
)

// TradingRepository is an interface that contains methods for long or short strategies
type TradingRepository interface {
	Strategies(ctx context.Context, strategy string, deal *model.Deal) (float64, error)
}

// BalanceService contains BalanceRepository interface
type TradingService struct {
	tRep TradingRepository
}

// NewTradingService accepts TradingRepository object and returnes an object of type *TradingService
func NewTradingService(tRep TradingRepository) *TradingService {
	return &TradingService{tRep: tRep}
}

// Strategies is a method of TradingService calls method of Repository
func (ts *TradingService) Strategies(ctx context.Context, strategy string, deal *model.Deal) (float64, error) {
	profit, err := ts.tRep.Strategies(ctx, strategy, deal)
	if err != nil {
		return 0, fmt.Errorf("TradingService-Strategies: error:%w", err)
	}
	return profit, nil
}
