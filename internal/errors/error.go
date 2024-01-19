package errors

import "fmt"

const (
	LoginAlreadyExist = "LOGIN_ALREADY_EXIST"
	UserDoesntExists  = "USER_DOESNT_EXISTS"
	PurchasePriceOut = "PURCHASE_PRICE_OUT"
	NotEnoughMoney   = "NOT_ENOUGH_MONEY"
)

type BusinessError struct {
	Code    string
	Message string
}

func New(code, message string) *BusinessError {
	return &BusinessError{Code: code, Message: message}
}

func (bs *BusinessError) Error() string {
	return fmt.Sprintf("code: %s, message: %s", bs.Code, bs.Message)
}
