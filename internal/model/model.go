// Package model contains models of using entities
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

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

// Deal is a struct for creating new deals
type Deal struct {
	DealID        uuid.UUID       `json:"-"`            // id of deal - each deal have new id
	SharesCount  decimal.Decimal `json:"sharescount"` // amount of shares of the selected transaction company
	ProfileID     uuid.UUID       `json:"profileid" validate:"required,uuid"`
	Company       string          `json:"company" validate:"required"`
	PurchasePrice decimal.Decimal `json:"-"`
	StopLoss      decimal.Decimal `json:"stoploss" validate:"required"`
	TakeProfit    decimal.Decimal `json:"takeprofit" validate:"required"`
	DealTime      time.Time       `json:"-"`
	EndDealTime   time.Time       `json:"-"`
	Profit        decimal.Decimal `json:"-"`
}
