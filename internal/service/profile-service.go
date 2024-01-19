// Package service contains business logic of a project
package service

import (
	"context"
	"fmt"

	"github.com/artnikel/APIService/internal/config"
	"github.com/artnikel/APIService/internal/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository is an interface that contains methods for user manipulation
type UserRepository interface {
	SignUp(ctx context.Context, user *model.User) error
	GetByLogin(ctx context.Context, username string) ([]byte, uuid.UUID, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) (string, error)
}

// UserService contains UserRepository interface
type UserService struct {
	uRep UserRepository
	cfg  config.Variables
}

// NewUserService accepts UserRepository object and returnes an object of type *UserService
func NewUserService(uRep UserRepository, cfg config.Variables) *UserService {
	return &UserService{uRep: uRep, cfg: cfg}
}

const (
	bcryptCost = 14
)

// SignUp is a method of UserService that hashed password and calls method of Repository
func (us *UserService) SignUp(ctx context.Context, user *model.User) error {
	var errHash error
	user.Password, errHash = us.GenerateHash(user.Password)
	if errHash != nil {
		return fmt.Errorf("generateHash %w", errHash)
	}
	err := us.uRep.SignUp(ctx, user)
	if err != nil {

		return fmt.Errorf("signUp %w", err)
	}
	return nil
}

// Login is a method of UserService that getting password and id, then checked password hash, generating tokens and added refresh token to database.
func (us *UserService) GetByLogin(ctx context.Context, user *model.User) (uuid.UUID, error) {
	hash, id, err := us.uRep.GetByLogin(ctx, user.Login)
	user.ID = id
	if err != nil {
		return uuid.Nil, fmt.Errorf("getByLogin %w", err)
	}
	verified, err := us.CheckPasswordHash(string(hash), user.Password)
	if err != nil || !verified {
		return uuid.Nil, fmt.Errorf("checkPasswordHash %w", err)
	}
	return id, nil
}

// DeleteAccount is a method from UserService that deleted account by id
func (us *UserService) DeleteAccount(ctx context.Context, id uuid.UUID) (string, error) {
	idString, err := us.uRep.DeleteAccount(ctx, id)
	if err != nil {
		return "", fmt.Errorf("deleteAccount %w", err)
	}
	return idString, nil
}

// GenerateHash is a method that makes from bytes hashed value
func (us *UserService) GenerateHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return string(bytes), fmt.Errorf("generateFromPassword %w", err)
	}
	return string(bytes), nil
}

// CheckPasswordHash is a method  that checks if hash is equal hash from given password
func (us *UserService) CheckPasswordHash(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, fmt.Errorf("compareHashAndPassword %w", err)
	}
	return true, nil
}
