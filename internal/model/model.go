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

// Balance contains an info about the balance and will be written in a balance table
type Balance struct {
	BalanceID uuid.UUID `json:"balanceid" validate:"required,uuid"`                  // id of balance operation - each operation have new id
	ProfileID uuid.UUID `json:"profileid" validate:"required,uuid"`                  // same value as ID in struct User
	Operation float64   `json:"operation" validate:"required,gt=0" form:"operation"` // sum of money to be deposit or withdraw
}
