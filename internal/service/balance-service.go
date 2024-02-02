package service

import (
	"context"
	"fmt"

	"github.com/artnikel/APIService/internal/config"
	berrors "github.com/artnikel/APIService/internal/errors"
	"github.com/artnikel/APIService/internal/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// BalanceRepository is an interface that contains methods for user manipulation
type BalanceRepository interface {
	BalanceOperation(ctx context.Context, balance *model.Balance) (float64, error)
	GetBalance(ctx context.Context, profileid uuid.UUID) (float64, error)
}

// BalanceService contains BalanceRepository interface
type BalanceService struct {
	bRep BalanceRepository
	cfg  config.Variables
}

// NewBalanceService accepts BalanceRepository object and returnes an object of type *BalanceService
func NewBalanceService(bRep BalanceRepository, cfg config.Variables) *BalanceService {
	return &BalanceService{bRep: bRep, cfg: cfg}
}

// BalanceOperation is a method of BalanceService calls method of Repository
func (bs *BalanceService) BalanceOperation(ctx context.Context, balance *model.Balance) (float64, error) {
	if decimal.NewFromFloat(balance.Operation).IsNegative() {
		money, err := bs.GetBalance(ctx, balance.ProfileID)
		if err != nil {
			return 0, fmt.Errorf("balanceOperation %w", err)
		}
		if decimal.NewFromFloat(money).Cmp(decimal.NewFromFloat(balance.Operation).Abs()) == 1 {
			operation, err := bs.bRep.BalanceOperation(ctx, balance)
			if err != nil {
				return 0, fmt.Errorf("balanceOperation %w", err)
			}
			return operation, nil
		}
		return 0, berrors.New(berrors.NotEnoughMoney, "Not enough money")
	}
	operation, err := bs.bRep.BalanceOperation(ctx, balance)
	if err != nil {
		return 0, fmt.Errorf("balanceOperation %w", err)
	}
	return operation, nil
}

// GetBalance is a method of BalanceService calls method of Repository
func (bs *BalanceService) GetBalance(ctx context.Context, profileid uuid.UUID) (float64, error) {
	money, err := bs.bRep.GetBalance(ctx, profileid)
	if err != nil {
		return 0, fmt.Errorf("getBalance %w", err)
	}
	return money, nil
}
