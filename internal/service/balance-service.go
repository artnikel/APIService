package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/artnikel/APIService/internal/config"
	"github.com/artnikel/APIService/internal/model"
	"github.com/golang-jwt/jwt"
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
		return 0, fmt.Errorf("not enough money")
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

// GetIDByToken is a method that get id by access token
func (bs *BalanceService) GetIDByToken(authHeader string) (uuid.UUID, error) {
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return uuid.Nil, fmt.Errorf("authorization header is invalid")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(bs.cfg.TokenSignature), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("error jwt parse %w", err)
	}
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("access token is invalid")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if id, ok := claims["id"].(string); ok {
			profileid, err := uuid.Parse(id)
			if err != nil {
				return uuid.Nil, fmt.Errorf("parse %w", err)
			}
			return profileid, nil
		}
	}
	return uuid.Nil, nil
}
