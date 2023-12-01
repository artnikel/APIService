// Package service contains business logic of a project
package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/artnikel/APIService/internal/config"
	"github.com/artnikel/APIService/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository is an interface that contains methods for user manipulation
type UserRepository interface {
	SignUp(ctx context.Context, user *model.User) error
	GetByLogin(ctx context.Context, username string) ([]byte, uuid.UUID, error)
	AddRefreshToken(ctx context.Context, id uuid.UUID, refreshToken string) error
	GetRefreshTokenByID(ctx context.Context, id uuid.UUID) (string, error)
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

// Expiration time for an access and a refresh tokens
const (
	accessTokenExpiration  = 15 * time.Minute
	refreshTokenExpiration = 72 * time.Hour
	bcryptCost             = 14
)

// SignUp is a method of UserService that hashed password and calls method of Repository
func (us *UserService) SignUp(ctx context.Context, user *model.User) error {
	var err error
	user.Password, err = us.GenerateHash(user.Password)
	if err != nil {
		return fmt.Errorf("generateHash %w", err)
	}
	err = us.uRep.SignUp(ctx, user)
	if err != nil {
		return fmt.Errorf("signUp %w", err)
	}
	return nil
}

// Login is a method of UserService that getting password and id, then checked password hash, generating tokens and added refresh token to database.
func (us *UserService) Login(ctx context.Context, user *model.User) (*model.TokenPair, error) {
	hash, id, err := us.uRep.GetByLogin(ctx, user.Login)
	user.ID = id
	if err != nil {
		return &model.TokenPair{}, fmt.Errorf("getByLogin %w", err)
	}
	verified, err := us.CheckPasswordHash(hash, user.Password)
	if err != nil || !verified {
		return &model.TokenPair{}, fmt.Errorf("checkPasswordHash %w", err)
	}
	tokenPair, err := us.GenerateTokenPair(user.ID)
	if err != nil {
		return &model.TokenPair{}, fmt.Errorf("generateTokenPair %w", err)
	}
	sum := sha256.Sum256([]byte(tokenPair.RefreshToken))
	hashedRefreshToken, err := us.GenerateHash(sum[:])
	if err != nil {
		return &model.TokenPair{}, fmt.Errorf("generateHash %w", err)
	}
	user.RefreshToken = string(hashedRefreshToken)
	err = us.uRep.AddRefreshToken(context.Background(), user.ID, user.RefreshToken)
	if err != nil {
		return &model.TokenPair{}, fmt.Errorf("addRefreshToken %w", err)
	}
	return tokenPair, nil
}

// Refresh is a method of UserService that refeshes access token and refresh token
func (us *UserService) Refresh(ctx context.Context, tokenPair *model.TokenPair) (*model.TokenPair, error) {
	id, err := us.TokensIDCompare(tokenPair)
	if err != nil {
		return &model.TokenPair{}, fmt.Errorf("tokensIDCompare %w", err)
	}
	hash, err := us.uRep.GetRefreshTokenByID(ctx, id)
	if err != nil {
		return &model.TokenPair{}, fmt.Errorf("getPasswordByUsernsame %w", err)
	}
	sum := sha256.Sum256([]byte(tokenPair.RefreshToken))
	verified, err := us.CheckPasswordHash([]byte(hash), sum[:])
	if err != nil || !verified {
		return &model.TokenPair{}, fmt.Errorf("refreshToken invalid")
	}
	tokenPair, err = us.GenerateTokenPair(id)
	if err != nil {
		return &model.TokenPair{}, fmt.Errorf("generateTokenPair %w", err)
	}
	sum = sha256.Sum256([]byte(tokenPair.RefreshToken))
	hashedRefreshToken, err := us.GenerateHash(sum[:])
	if err != nil {
		return &model.TokenPair{}, fmt.Errorf("generateHash %w", err)
	}
	var user model.User
	user.RefreshToken = string(hashedRefreshToken)
	user.ID = id
	err = us.uRep.AddRefreshToken(context.Background(), user.ID, user.RefreshToken)
	if err != nil {
		return &model.TokenPair{}, fmt.Errorf("addRefreshToken %w", err)
	}
	return tokenPair, nil
}

// DeleteAccount is a method from UserService that deleted account by id
func (us *UserService) DeleteAccount(ctx context.Context, id uuid.UUID) (string, error) {
	idString, err := us.uRep.DeleteAccount(ctx, id)
	if err != nil {
		return "", fmt.Errorf("deleteAccount %w", err)
	}
	return idString, nil
}

// TokensIDCompare compares IDs from refresh and access token for being equal
func (us *UserService) TokensIDCompare(tokenPair *model.TokenPair) (uuid.UUID, error) {
	accessToken, err := us.ValidateToken(tokenPair.AccessToken, us.cfg.TokenSignature)
	if err != nil {
		return uuid.Nil, fmt.Errorf("validateToken %w", err)
	}
	var accessID uuid.UUID
	var uuidID uuid.UUID
	if claims, ok := accessToken.Claims.(jwt.MapClaims); ok && accessToken.Valid {
		uuidID, err = uuid.Parse(claims["id"].(string))
		if err != nil {
			return uuid.Nil, fmt.Errorf("parse %w", err)
		}
		accessID = uuidID
	}
	refreshToken, err := us.ValidateToken(tokenPair.RefreshToken, us.cfg.TokenSignature)
	if err != nil {
		return uuid.Nil, fmt.Errorf("validateToken %w", err)
	}
	var refreshID uuid.UUID
	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		exp := claims["exp"].(float64)
		uuidID, err = uuid.Parse(claims["id"].(string))
		if err != nil {
			return uuid.Nil, fmt.Errorf("parse %w", err)
		}
		refreshID = uuidID
		if exp < float64(time.Now().Unix()) {
			return uuid.Nil, fmt.Errorf("token is expired")
		}
	}
	if accessID != refreshID {
		return uuid.Nil, fmt.Errorf("user ID in access token doesn't equal user ID in refresh token")
	}
	return accessID, nil
}

// GenerateHash is a method that makes from bytes hashed value
func (us *UserService) GenerateHash(password []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword(password, bcryptCost)
	if err != nil {
		return bytes, fmt.Errorf("generateFromPassword %w", err)
	}
	return bytes, nil
}

// CheckPasswordHash is a method  that checks if hash is equal hash from given password
func (us *UserService) CheckPasswordHash(hash, password []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		return false, fmt.Errorf("compareHashAndPassword %w", err)
	}
	return true, nil
}

// GenerateTokenPair generates pair of access and refresh tokens
func (us *UserService) GenerateTokenPair(id uuid.UUID) (*model.TokenPair, error) {
	accessToken, err := us.GenerateJWTToken(accessTokenExpiration, id)
	if err != nil {
		return &model.TokenPair{}, fmt.Errorf("generateJWTToken %w", err)
	}
	refreshToken, err := us.GenerateJWTToken(refreshTokenExpiration, id)
	if err != nil {
		return &model.TokenPair{}, fmt.Errorf("generateJWTToken %w", err)
	}
	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GenerateJWTToken is a method that generate JWT token with given expiration with user id
func (us *UserService) GenerateJWTToken(expiration time.Duration, id uuid.UUID) (string, error) {
	claims := &jwt.MapClaims{
		"exp": time.Now().Add(expiration).Unix(),
		"id":  id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(us.cfg.TokenSignature))
	if err != nil {
		return "", fmt.Errorf("token.SignedString %w", err)
	}
	return tokenString, nil
}

// ValidateToken validates the JWT token using the secret key, it checks the signing method and returns the parsed token if it is valid.
func (us *UserService) ValidateToken(tokenString, secretKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
