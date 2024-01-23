// Package errors contains business errors
package errors

import "fmt"

const (
	// LoginAlreadyExist is error code if login already exist
	LoginAlreadyExist = "LOGIN_ALREADY_EXIST"
	// UserDoesntExists is error code if user doesn`t exist
	UserDoesntExists = "USER_DOESNT_EXISTS"
	// PurchasePriceOut is error code if takeprofit or stoploss out of limit
	PurchasePriceOut = "PURCHASE_PRICE_OUT"
	// NotEnoughMoney is error code if user don`t have enough money
	NotEnoughMoney = "NOT_ENOUGH_MONEY"
)

// BusinessError is struct for business errors
type BusinessError struct {
	Code    string
	Message string
}

// New is constructor for manage business errors
func New(code, message string) *BusinessError {
	return &BusinessError{Code: code, Message: message}
}

// Error is method for creating business errors
func (bs *BusinessError) Error() string {
	return fmt.Sprintf("code: %s, message: %s", bs.Code, bs.Message)
}
