// Package handler is the top level of the application and it contains request handlers
package handler

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"text/template"

	"github.com/artnikel/APIService/internal/config"
	"github.com/artnikel/APIService/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gopkg.in/boj/redistore.v1"
)

// UserService is an interface that defines the methods on User entity.
type UserService interface {
	SignUp(ctx context.Context, user *model.User) error
	GetByLogin(ctx context.Context, user *model.User) (uuid.UUID, error)
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
	CreatePosition(ctx context.Context, deal *model.Deal) error
	ClosePositionManually(ctx context.Context, dealid, profileid uuid.UUID) (float64, error)
	GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error)
	GetPrices(ctx context.Context) ([]model.Share, error)
}

// Handler is responsible for handling HTTP requests related to entities.
type Handler struct {
	userService    UserService
	balanceService BalanceService
	tradingService TradingService
	validate       *validator.Validate
	cfg            config.Variables
}

// NewHandler creates a new instance of the Handler struct.
func NewHandler(userService UserService, balanceService BalanceService, tradingService TradingService, v *validator.Validate, cfg config.Variables) *Handler {
	return &Handler{
		userService:    userService,
		balanceService: balanceService,
		tradingService: tradingService,
		validate:       v,
		cfg: cfg,
	}
}

func NewRedisStore(cfg config.Variables) *redistore.RediStore {
	store, err := redistore.NewRediStore(10, "tcp", cfg.RedisPriceAddress, "", []byte(cfg.TokenSignature))
	if err != nil {
		log.Fatalf("Failed to create redis store: %v", err)
	}
	return store
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

func (h *Handler) Auth(c echo.Context) error {
	tmpl, err := template.ParseFiles("templates/auth/auth.html")
	if err != nil {
		return echo.ErrNotFound
	}
	return tmpl.ExecuteTemplate(c.Response().Writer, "auth", nil)
}

// SignUp calls method of Service by handler
func (h *Handler) SignUp(c echo.Context) error {
	tmpl, err := template.ParseFiles("templates/auth/auth.html")
	if err != nil {
		return echo.ErrNotFound
	}
	var user model.User
	if err := c.Bind(&user); err != nil {
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "Failed to bind fields",
		})
	}
	tempPassword := user.Password
	err = h.validate.StructCtx(c.Request().Context(), user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Login":    user.Login,
			"Password": user.Password,
		}).Errorf("signUp %v", err)
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "The fields have not been validated",
		})
	}
	err = h.userService.SignUp(c.Request().Context(), &user)
	if err != nil {
		logrus.Errorf("signUp %v", err)
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "Failed to sign up",
		})
	}
	user.Password = tempPassword
	userID, err := h.userService.GetByLogin(c.Request().Context(), &user)
	if err != nil {
		logrus.Errorf("signUp %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to log in")
	}
	store := NewRedisStore(h.cfg)
	session, err := store.Get(c.Request(), "SESSION_ID")
	if err != nil {
		logrus.Errorf("signUp %v", err)
		return echo.ErrNotFound
	}
	session.Values["id"] = userID.String()
	session.Values["login"] = user.Login
	session.Values["password"] = user.Password
	if err = session.Save(c.Request(), c.Response()); err != nil {
		logrus.Errorf("signUp %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "error saving session")
	}
	return c.Redirect(http.StatusSeeOther, "/index")
}

// Login calls method of Service by handler
func (h *Handler) Login(c echo.Context) error {
	tmpl, err := template.ParseFiles("templates/auth/auth.html")
	if err != nil {
		return echo.ErrNotFound
	}
	var user model.User
	if err := c.Bind(&user); err != nil {
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "Failed to bind fields",
		})
	}
	err = h.validate.StructCtx(c.Request().Context(), user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Login":    user.Login,
			"Password": user.Password,
		}).Errorf("login %v", err)
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "The fields have not been validated",
		})
	}
	userID, err := h.userService.GetByLogin(c.Request().Context(), &user)
	if err != nil {
		logrus.Errorf("login %v", err)
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "Wrong login or password",
		})
	}
	store := NewRedisStore(h.cfg)
	session, err := store.Get(c.Request(), "SESSION_ID")
	if err != nil {
		logrus.Errorf("login %v", err)
		return echo.ErrNotFound
	}
	session.Values["id"] = userID.String()
	session.Values["login"] = user.Login
	session.Values["password"] = user.Password
	if err = session.Save(c.Request(), c.Response().Writer); err != nil {
		logrus.Errorf("login %v", err)
		return c.String(http.StatusBadRequest, "error saving session")
	}
	return c.Redirect(http.StatusSeeOther, "/index")
}

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

// Deposit calls method of Service by handler
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

// nolint dupl // in swagger can't bind two routers to one method
// Long calls method of Service by handler
func (h *Handler) Long(c echo.Context) error {
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
	err = h.tradingService.CreatePosition(c.Request().Context(), deal)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Long: failed to create position long")
	}
	return c.JSON(http.StatusOK, "Position long created.")
}

// nolint dupl // in swagger can't bind two routers to one method
// Short calls method of Service by handler
func (h *Handler) Short(c echo.Context) error {
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
	err = h.tradingService.CreatePosition(c.Request().Context(), deal)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-Short: failed to create position short")
	}
	return c.JSON(http.StatusOK, "Position short created.")
}

// ClosePositionManually calls method of Service by handler
func (h *Handler) ClosePositionManually(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	profileid, err := h.balanceService.GetIDByToken(authHeader)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-ClosePositionManually-GetIDByToken: failed to get ID by token")
	}
	var requestData closeData
	err = c.Bind(&requestData)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-ClosePositionManually: invalid request payload")
	}
	dealUUID, err := uuid.Parse(requestData.DealID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-ClosePositionManually: failed to parse id")
	}
	profit, err := h.tradingService.ClosePositionManually(c.Request().Context(), dealUUID, profileid)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-ClosePositionManually: failed to close position")
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("Position closed with profit %f", profit))
}

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

// GetPrices calls method of Service by handler
func (h *Handler) GetPrices(c echo.Context) error {
	shares, err := h.tradingService.GetPrices(c.Request().Context())
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-GetPrices: failed to get shares")
	}
	return c.JSON(http.StatusOK, shares)
}
