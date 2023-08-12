// Package repository is a lower level of project
package repository

import (
	"context"
	"fmt"

	"github.com/artnikel/APIService/internal/model"
	uproto "github.com/artnikel/ProfileService/proto"
	"github.com/google/uuid"
)

// ProfileRepository represents the client of UserService repository implementation.
type ProfileRepository struct {
	client uproto.UserServiceClient
}

// NewProfileRepository creates and returns a new instance of ProfileRepository, using the provided proto.UserServiceClient.
func NewProfileRepository(client uproto.UserServiceClient) *ProfileRepository {
	return &ProfileRepository{
		client: client,
	}
}

// SignUp call a method of ProfileService.
func (p *ProfileRepository) SignUp(ctx context.Context, user *model.User) error {
	_, err := p.client.SignUp(ctx, &uproto.SignUpRequest{User: &uproto.User{
		Login:    user.Login,
		Password: string(user.Password),
	}})
	if err != nil {
		return fmt.Errorf("ProfileRepository-SignUp: error:%w", err)
	}
	return nil
}

// GetByLogin call a method of ProfileService.
func (p *ProfileRepository) GetByLogin(ctx context.Context, login string) ([]byte, uuid.UUID, error) {
	resp, err := p.client.GetByLogin(ctx, &uproto.GetByLoginRequest{Login: login})
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("ProfileRepository-SignUp: error:%w", err)
	}
	idUUID, err := uuid.Parse(resp.Id)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("ProfileRepository-SignUp: failed to parse:%w", err)
	}
	return []byte(resp.Password), idUUID, nil
}

// AddRefreshToken call a method of ProfileService.
func (p *ProfileRepository) AddRefreshToken(ctx context.Context, id uuid.UUID, refreshToken string) error {
	_, err := p.client.AddRefreshToken(ctx, &uproto.AddRefreshTokenRequest{
		Id:           id.String(),
		RefreshToken: refreshToken,
	})
	if err != nil {
		return fmt.Errorf("ProfileRepository-AddRefreshToken: error:%w", err)
	}
	return nil
}

// GetRefreshTokenByID call a method of ProfileService.
func (p *ProfileRepository) GetRefreshTokenByID(ctx context.Context, id uuid.UUID) (string, error) {
	resp, err := p.client.GetRefreshTokenByID(ctx, &uproto.GetRefreshTokenByIDRequest{Id: id.String()})
	if err != nil {
		return "", fmt.Errorf("ProfileRepository-GetRefreshTokenByID: error:%w", err)
	}
	return resp.RefreshToken, nil
}

// DeleteAccount call a method of ProfileService.
func (p *ProfileRepository) DeleteAccount(ctx context.Context, id uuid.UUID) (string, error) {
	resp, err := p.client.DeleteAccount(ctx, &uproto.DeleteAccountRequest{Id: id.String()})
	if err != nil {
		return "", fmt.Errorf("ProfileRepository-DeleteAccount: error:%w", err)
	}
	return resp.Id, nil
}
