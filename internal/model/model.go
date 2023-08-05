// Package model contains models of using entities
package model

import "github.com/google/uuid"

// User contains an info about the user and will be written in a users table
type User struct {
	ID           uuid.UUID `json:"-"`
	Login        string    `json:"login" validate:"required,min=5,max=20"`
	Password     []byte    `json:"password" validate:"required,min=8"`
	RefreshToken string    `json:"-" `
}

// TokenPair contains two tokens for authorization of user
type TokenPair struct {
	AccessToken  string `json:"accesstoken" form:"accesstoken"`
	RefreshToken string `json:"refreshtoken" form:"refreshtoken"`
}
