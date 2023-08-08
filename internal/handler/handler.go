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
	"github.com/sirupsen/logrus"
)

// UserService is an interface that defines the methods on User entity.
type UserService interface {
	SignUp(ctx context.Context, user *model.User) error
	Login(ctx context.Context, user *model.User) (*model.TokenPair, error)
	Refresh(ctx context.Context, tokenPair *model.TokenPair) (*model.TokenPair, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) (string, error)
}

type BalanceService interface {
	BalanceOperation(ctx context.Context, balance *model.Balance) (float64, error)
	GetBalance(ctx context.Context, profileid uuid.UUID) (float64, error)
	GetIDByToken(authHeader string) (uuid.UUID, error)
}

// Handler is responsible for handling HTTP requests related to entities.
type Handler struct {
	userService    UserService
	balanceService BalanceService
	validate       *validator.Validate
}

// NewHandler creates a new instance of the Handler struct.
func NewHandler(userService UserService, balanceService BalanceService, v *validator.Validate) *Handler {
	return &Handler{
		userService:    userService,
		balanceService: balanceService,
		validate:       v,
	}
}

// InputData is a struct for binding login and password.
type InputData struct {
	Login    string `json:"login" form:"login"`
	Password string `json:"password" form:"password"`
}

// SignUp calls method of Service by handler
func (h *Handler) SignUp(c echo.Context) error {
	var newUser model.User
	requestData := &InputData{}
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

// Login calls method of Service by handler
func (h *Handler) Login(c echo.Context) error {
	var requestData InputData
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

// DeleteAccount calls method of Service by handler
func (h *Handler) DeleteAccount(c echo.Context) error {
	id := c.Param("id")
	err := h.validate.VarCtx(c.Request().Context(), id, "required,uuid")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Not valid id or field id is empty")
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Handler-DeleteAccount: failed to parse id")
	}
	str, err := h.userService.DeleteAccount(c.Request().Context(), idUUID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Id": id,
		}).Errorf("Handler-Refresh: error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete")
	}
	return c.JSON(http.StatusOK, str)
}

// BalanceOperation calls method of Service by handler
func (h *Handler) BalanceOperation(c echo.Context) error {
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
		return c.JSON(http.StatusBadRequest, "Handler-BalanceOperation: invalid request payload")
	}
	err = h.validate.VarCtx(c.Request().Context(), requestData.Operation, "required,gt=0")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-BalanceOperation: failed to validate operation")
	}
	authHeader := c.Request().Header.Get("Authorization")
	profileid, err := h.balanceService.GetIDByToken(authHeader)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return c.JSON(http.StatusBadRequest, "Handler-BalanceOperation-GetIDByToken: failed to get ID by token")
	}
	balance := model.Balance{
		ProfileID: profileid,
		Operation: requestData.Operation,
	}
	if c.Request().RequestURI == "/withdraw" {
		balance.Operation = -balance.Operation
		operationType = "Withdraw"
	}
	money, err := h.balanceService.BalanceOperation(c.Request().Context(), &balance)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"BalanceId": balance.BalanceID,
			"ProfileId": balance.ProfileID,
			"Operation": balance.Operation,
		}).Errorf("Handler-BalanceOperation: error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Handler-BalanceOperation: failed to made balance operation")
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
