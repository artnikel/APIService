// Package model contains models of using entities
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// User contains an info about the user and will be written in a users table
type User struct {
	ID           uuid.UUID `json:"-"`                                      // unique id of user
	Login        string    `json:"login" validate:"required,min=5,max=20"` // username of user account
	Password     []byte    `json:"password" validate:"required,min=8"`     // password of user account
	RefreshToken string    `json:"-"`                                      // token that storing the session in the database
}

// TokenPair contains two tokens for authorization of user
type TokenPair struct {
	AccessToken  string `json:"accesstoken" form:"accesstoken"`   // token that gives access to secure methods
	RefreshToken string `json:"refreshtoken" form:"refreshtoken"` // same token as in struct User
}

// Balance contains an info about the balance and will be written in a balance table
type Balance struct {
	BalanceID uuid.UUID `json:"-" validate:"required,uuid"`                          // id of balance operation - each operation have new id
	ProfileID uuid.UUID `json:"-" validate:"required,uuid"`                          // same value as ID in struct User
	Operation float64   `json:"operation" validate:"required,gt=0" form:"operation"` // sum of money to be deposit or withdraw
}

// Deal is a struct for creating new deals
type Deal struct {
	DealID        uuid.UUID       `json:"dealid"`                         // id of deal - each deal have new id
	SharesCount   decimal.Decimal `json:"sharescount"`                    // amount of shares of the selected transaction company
	ProfileID     uuid.UUID       `json:"-" validate:"required,uuid"`     // id of user/profile
	Company       string          `json:"company" validate:"required"`    // name of company in share
	PurchasePrice decimal.Decimal `json:"purchaseprice"`                  // entry price in position
	StopLoss      decimal.Decimal `json:"stoploss" validate:"required"`   // lower limit where the price can go
	TakeProfit    decimal.Decimal `json:"takeprofit" validate:"required"` // upper limit where the price can go
	DealTime      time.Time       `json:"dealtime"`                       // entry time in position
	EndDealTime   time.Time       `json:"-"`                              // time of closing position
	Profit        decimal.Decimal `json:"-"`                              // revenue of position
}
