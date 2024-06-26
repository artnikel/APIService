package repository

import (
	"context"
	"fmt"
	"strconv"

	berrors "github.com/artnikel/APIService/internal/errors"
	"github.com/artnikel/APIService/internal/model"
	bproto "github.com/artnikel/BalanceService/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

// BalanceRepository represents the client of Balance Service repository implementation.
type BalanceRepository struct {
	client bproto.BalanceServiceClient
}

// NewBalanceRepository creates and returns a new instance of BalanceRepository, using the provided proto.BalanceServiceClient.
func NewBalanceRepository(client bproto.BalanceServiceClient) *BalanceRepository {
	return &BalanceRepository{
		client: client,
	}
}

// BalanceOperation call a method of BalanceService.
func (b *BalanceRepository) BalanceOperation(ctx context.Context, balance *model.Balance) (float64, error) {
	resp, err := b.client.BalanceOperation(ctx, &bproto.BalanceOperationRequest{Balance: &bproto.Balance{
		Balanceid: balance.BalanceID.String(),
		Profileid: balance.ProfileID.String(),
		Operation: balance.Operation,
	}})
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if ok && grpcStatus.Message() == berrors.NotEnoughMoney {
			return 0, berrors.New(berrors.NotEnoughMoney, "Not enough money")
		}
		return 0, fmt.Errorf("balanceOperation %w", err)
	}
	operation, err := strconv.ParseFloat(resp.Operation, 64)
	if err != nil {
		return 0, fmt.Errorf("parseFloat %w", err)
	}
	return operation, nil
}

// GetBalance call a method of BalanceService.
func (b *BalanceRepository) GetBalance(ctx context.Context, profileid uuid.UUID) (float64, error) {
	resp, err := b.client.GetBalance(ctx, &bproto.GetBalanceRequest{Profileid: profileid.String()})
	if err != nil {
		return 0, fmt.Errorf("getBalance %w", err)
	}
	return resp.Money, nil
}
