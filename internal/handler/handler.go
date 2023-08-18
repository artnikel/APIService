// Package handler is the top level of the application and it contains request handlers
package handler

import (
	"context"
	"fmt"
	"math"
	"net/http"

	"github.com/artnikel/APIService/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// UserService is an interface that defines the methods on User entity.
type UserService interface {
	SignUp(ctx context.Context, user *model.User) error
	Login(ctx context.Context, user *model.User) (*model.TokenPair, error)
	Refresh(ctx context.Context, tokenPair *model.TokenPair) (*model.TokenPair, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) (string, error)
}

// BalanceService is an interface that defines the methods on Balance entity.
type BalanceService interface {
	BalanceOperation(ctx context.Context, balance *model.Balance) (float64, error)
	GetBalance(ctx context.Context, profileid uuid.UUID) (float64, error)
	GetIDByToken(authHeader string) (uuid.UUID, error)
}

// TradingService is an interface that defines the method on Trading entity.
type TradingService interface {
	GetProfit(ctx context.Context, strategy string, deal *model.Deal) (float64, error)
	ClosePosition(ctx context.Context, dealid, profileid uuid.UUID) error
	GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error)
}

// Handler is responsible for handling HTTP requests related to entities.
type Handler struct {
	userService    UserService
	balanceService BalanceService
	tradingService TradingService
	validate       *validator.Validate
}

// NewHandler creates a new instance of the Handler struct.
func NewHandler(userService UserService, balanceService BalanceService, tradingService TradingService, v *validator.Validate) *Handler {
	return &Handler{
		userService:    userService,
		balanceService: balanceService,
		tradingService: tradingService,
		validate:       v,
	}
}

// inputData is a struct for binding login and password.
type inputData struct {
	Login    string `json:"login" form:"login" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

// dealData is a struct for binding new deal.
type dealData struct {
	SharesCount float64 `json:"sharescount" form:"sharescount" validate:"required"`
	Company     string  `json:"company" form:"company" validate:"required"`
	StopLoss    float64 `json:"stoploss" form:"stoploss" validate:"required"`
	TakeProfit  float64 `json:"takeprofit" form:"takeprofit" validate:"required"`
}

type closeData struct {
	DealID string `json:"dealid" form:"dealid" validate:"required,uuid"`
}

// @Summary Registration
// @ID create-account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body inputData true "Please fill the login (5-20 symbols) and password (minimum 8 symbols) fields."
// @Success 201 {string} string "token"
// @Failure 400 {object} error
// @Router /signup [post]
// SignUp calls method of Service by handler
func (h *Handler) SignUp(c echo.Context) error {
	var newUser model.User
	requestData := &inputData{}
	err := c.Bind(requestData)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-SignUpUser: Invalid request payload")
	}
	newUser.Login = requestData.Login
	newUser.Password = []byte(requestData.Password)
	err = h.validate.StructCtx(c.Request().Context(), newUser)
	if err != nil {
		logrus.Errorf("Handler-SignUp: error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to validate")
	}
	err = h.userService.SignUp(c.Request().Context(), &newUser)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ID":           newUser.ID,
			"Login":        newUser.Login,
			"Password":     newUser.Password,
			"RefreshToken": newUser.RefreshToken,
		}).Errorf("Handler-SignUp: error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to signUp")
	}
	return c.JSON(http.StatusCreated, "Account created.")
}

// @Summary Account login
// @ID login
// @Tags auth
// @Accept json
// @Produce json
// @Param input body inputData true "Please fill the login and password fields."
// @Success 200 {string} string "tokens"
// @Failure 400 {string} error
// @Router /login [post]
// Login calls method of Service by handler
func (h *Handler) Login(c echo.Context) error {
	var requestData inputData
	err := c.Bind(&requestData)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-GetByLogin: Invalid request payload")
	}
	var user model.User
	user.Login = requestData.Login
	user.Password = []byte(requestData.Password)
	err = h.validate.VarCtx(c.Request().Context(), requestData.Login, "required")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Not valid data: login field is empty")
	}
	err = h.validate.VarCtx(c.Request().Context(), requestData.Password, "required")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Not valid data: password field is empty")
	}
	tokenPair, err := h.userService.Login(c.Request().Context(), &user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ID":           user.ID,
			"Login":        user.Login,
			"Password":     user.Password,
			"RefreshToken": user.RefreshToken,
		}).Errorf("Handler-SignUp: error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to log in")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"Access Token : ":  tokenPair.AccessToken,
		"Refresh Token : ": tokenPair.RefreshToken,
	})
}

// @Summary Refresh tokens
// @ID refresh-token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.TokenPair true "Please fill the access and refresh tokens fields."
// @Success 200 {string} string "tokens"
// @Failure 400 {string} error
// @Router /refresh [post]
// Refresh calls method of Service by handler
func (h *Handler) Refresh(c echo.Context) error {
	var requestData model.TokenPair
	err := c.Bind(&requestData)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Refresh: Invalid request payload")
	}
	var tokens model.TokenPair
	tokens.AccessToken = requestData.AccessToken
	tokens.RefreshToken = requestData.RefreshToken
	err = h.validate.VarCtx(c.Request().Context(), requestData.AccessToken, "required")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Not valid data: access token field is empty")
	}
	err = h.validate.VarCtx(c.Request().Context(), requestData.RefreshToken, "required")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Not valid data: refresh token field is empty")
	}
	newTokens, err := h.userService.Refresh(c.Request().Context(), &tokens)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Access Token":  tokens.AccessToken,
			"Refresh Token": tokens.RefreshToken,
		}).Errorf("Handler-Refresh: error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to refresh")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"Access Token : ":  newTokens.AccessToken,
		"Refresh Token : ": newTokens.RefreshToken,
	})
}

// @Summary Delete your account
// @Security ApiKeyAuth
// @ID delete-account
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Failure 400 {string} error
// @Router /delete [delete]
// DeleteAccount calls method of Service by handler
func (h *Handler) DeleteAccount(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	id, err := h.balanceService.GetIDByToken(authHeader)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-GetBalance-GetIDByToken: failed to get ID by token")
	}
	str, err := h.userService.DeleteAccount(c.Request().Context(), id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Id": id,
		}).Errorf("Handler-Refresh: error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete")
	}
	return c.JSON(http.StatusOK, str)
}

// @Summary Balance operation deposit
// @Security ApiKeyAuth
// @ID deposit
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.Balance true "Please input the amount of money to deposit"
// @Success 200 {string} string
// @Failure 400 {string} error
// @Router /deposit [post]
// BalanceOperation calls method of Service by handler
func (h *Handler) Deposit(c echo.Context) error {
	var (
		operation   float64
		requestData = model.Balance{
			Operation: operation,
		}
		operationType = "Deposit"
		output        = func(money float64) string {
			return fmt.Sprintf("%s of %.2f successfully made.", operationType, math.Abs(money))
		}
	)
	err := c.Bind(&requestData)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Deposit: invalid request payload")
	}
	err = h.validate.VarCtx(c.Request().Context(), requestData.Operation, "required,gt=0")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Deposit: failed to validate operation")
	}
	authHeader := c.Request().Header.Get("Authorization")
	profileid, err := h.balanceService.GetIDByToken(authHeader)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Deposit-GetIDByToken: failed to get ID by token")
	}
	balance := model.Balance{
		ProfileID: profileid,
		Operation: requestData.Operation,
	}
	money, err := h.balanceService.BalanceOperation(c.Request().Context(), &balance)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"BalanceId": balance.BalanceID,
			"ProfileId": balance.ProfileID,
			"Operation": balance.Operation,
		}).Errorf("Handler-Deposit: error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Handler-Deposit: failed to made balance operation")
	}
	return c.JSON(http.StatusOK, output(money))
}

// @Summary Balance operation withdraw
// @Security ApiKeyAuth
// @ID withdraw
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.Balance true "Please input the amount of money to withdraw"
// @Success 200 {string} string
// @Failure 400 {string} error
// @Router /withdraw [post]
// Withdraw calls method of Service by handler
func (h *Handler) Withdraw(c echo.Context) error {
	var (
		operation   float64
		requestData = model.Balance{
			Operation: operation,
		}
		operationType = "Withdraw"
		output        = func(money float64) string {
			return fmt.Sprintf("%s of %.2f successfully made.", operationType, math.Abs(money))
		}
	)
	err := c.Bind(&requestData)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Withdraw: invalid request payload")
	}
	err = h.validate.VarCtx(c.Request().Context(), requestData.Operation, "required,gt=0")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Withdraw: failed to validate operation")
	}
	authHeader := c.Request().Header.Get("Authorization")
	profileid, err := h.balanceService.GetIDByToken(authHeader)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Withdraw-GetIDByToken: failed to get ID by token")
	}
	balance := model.Balance{
		ProfileID: profileid,
		Operation: requestData.Operation,
	}
	balance.Operation = -balance.Operation
	money, err := h.balanceService.BalanceOperation(c.Request().Context(), &balance)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"BalanceId": balance.BalanceID,
			"ProfileId": balance.ProfileID,
			"Operation": balance.Operation,
		}).Errorf("Handler-Withdraw: error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Handler-Withdraw: failed to made balance operation")
	}
	return c.JSON(http.StatusOK, output(money))
}

// @Summary Check sum of money on your balance
// @Security ApiKeyAuth
// @ID get-operation
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Failure 400 {string} error
// @Router /getbalance [get]
// GetBalance calls method of Service by handler
func (h *Handler) GetBalance(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	profileid, err := h.balanceService.GetIDByToken(authHeader)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-GetBalance-GetIDByToken: failed to get ID by token")
	}
	money, err := h.balanceService.GetBalance(c.Request().Context(), profileid)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ProfileId": profileid,
		}).Errorf("Handler-GetBalance: error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get balance")
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("Balance: %f", money))
}

// @Summary Invest money with strategy "long"
// @Security ApiKeyAuth
// @ID long
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dealData true "Please fill the SraresCount,Company,StopLoss and TakeProfit fields."
// @Success 200 {string} string
// @Failure 400 {string} error
// @Router /long [post]
// Long calls method of Service by handler
func (h *Handler) Long(c echo.Context) error {
	strategy := "long"
	authHeader := c.Request().Header.Get("Authorization")
	profileid, err := h.balanceService.GetIDByToken(authHeader)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Long-GetIDByToken: failed to get ID by token")
	}
	var requestData dealData
	err = c.Bind(&requestData)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Long: invalid request payload")
	}
	deal := &model.Deal{
		ProfileID:   profileid,
		SharesCount: decimal.NewFromFloat(requestData.SharesCount),
		Company:     requestData.Company,
		StopLoss:    decimal.NewFromFloat(requestData.StopLoss),
		TakeProfit:  decimal.NewFromFloat(requestData.TakeProfit),
	}
	err = h.validate.StructCtx(c.Request().Context(), requestData)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Long: failed to validate operation")
	}
	profit, err := h.tradingService.GetProfit(c.Request().Context(), strategy, deal)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Long: failed to run strategies")
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("Profit: %f", profit))
}

// @Summary Invest money with strategy "short"
// @Security ApiKeyAuth
// @ID short
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dealData true "Please fill the fields about Share."
// @Success 200 {string} string
// @Failure 400 {string} error
// @Router /short [post]
// Short calls method of Service by handler
func (h *Handler) Short(c echo.Context) error {
	strategy := "short"
	authHeader := c.Request().Header.Get("Authorization")
	profileid, err := h.balanceService.GetIDByToken(authHeader)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Short-GetIDByToken: failed to get ID by token")
	}
	var requestData dealData
	err = c.Bind(&requestData)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Short: invalid request payload")
	}
	deal := &model.Deal{
		ProfileID:   profileid,
		SharesCount: decimal.NewFromFloat(requestData.SharesCount),
		Company:     requestData.Company,
		StopLoss:    decimal.NewFromFloat(requestData.StopLoss),
		TakeProfit:  decimal.NewFromFloat(requestData.TakeProfit),
	}
	err = h.validate.StructCtx(c.Request().Context(), requestData)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Short: failed to validate operation")
	}
	profit, err := h.tradingService.GetProfit(c.Request().Context(), strategy, deal)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Short: failed to run strategies")
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("Profit: %f", profit))
}

// @Summary Close the position
// @Security ApiKeyAuth
// @ID close-position
// @Tags auth
// @Accept json
// @Produce json
// @Param input body closeData true "Please fill the id of deal"
// @Success 200 {string} string
// @Failure 400 {string} error
// @Router /closeposition [post]
// ClosePosition calls method of Service by handler
func (h *Handler) ClosePosition(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	profileid, err := h.balanceService.GetIDByToken(authHeader)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-ClosePosition-GetIDByToken: failed to get ID by token")
	}
	var requestData closeData
	err = c.Bind(&requestData)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-ClosePosition: invalid request payload")
	}
	dealUUID, err := uuid.Parse(requestData.DealID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-ClosePosition: failed to parse id")
	}
	err = h.tradingService.ClosePosition(c.Request().Context(), dealUUID, profileid)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-ClosePosition: failed to close position")
	}
	return c.JSON(http.StatusOK, "Position closed.")
}

// @Summary Show all unclosed positions
// @Security ApiKeyAuth
// @ID get-unclosed
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} model.Deal
// @Failure 400 {string} error
// @Router /getunclosed [get]
// GetUnclosedPositions calls method of Service by handler
func (h *Handler) GetUnclosedPositions(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	profileid, err := h.balanceService.GetIDByToken(authHeader)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-GetUnclosedPositions-GetIDByToken: failed to get ID by token")
	}
	unclosedDeals, err := h.tradingService.GetUnclosedPositions(c.Request().Context(), profileid)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-GetUnclosedPositions: failed to get unclosed positions")
	}
	return c.JSON(http.StatusOK, unclosedDeals)
}
