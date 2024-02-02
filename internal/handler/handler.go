// Package handler is the top level of the application and it contains request handlers
package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/artnikel/APIService/internal/config"
	berrors "github.com/artnikel/APIService/internal/errors"
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
}

// TradingService is an interface that defines the method on Trading entity.
type TradingService interface {
	CreatePosition(ctx context.Context, deal *model.Deal) error
	ClosePositionManually(ctx context.Context, dealid, profileid uuid.UUID) (float64, error)
	GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error)
	GetClosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error)
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
		cfg:            cfg,
	}
}

// NewRedisStore creates a new Redis storage instance for sessions
func NewRedisStore(cfg config.Variables) *redistore.RediStore {
	// nolint gonmd
	store, err := redistore.NewRediStore(10, "tcp", cfg.RedisPriceAddress, "", []byte(cfg.TokenSignature))
	if err != nil {
		log.Fatalf("failed to create redis store: %v", err)
	}
	return store
}

// getProfileID is method for getting id of profile from session
func (h *Handler) getProfileID(c echo.Context) (uuid.UUID, error) {
	cookie, err := c.Cookie("SESSION_ID")
	if err != nil {
		logrus.Errorf("getProfileID: %v", err)
		return uuid.Nil, c.Redirect(http.StatusSeeOther, "/")
	}
	store := NewRedisStore(h.cfg)
	session, err := store.Get(c.Request(), cookie.Name)
	if err != nil {
		logrus.Errorf("getProfileID: %v", err)
		return uuid.Nil, echo.ErrNotFound
	}
	if len(session.Values) == 0 {
		return uuid.Nil, c.Redirect(http.StatusSeeOther, "/")
	}
	profileid := session.Values["id"].(string)
	profileUUID, err := uuid.Parse(profileid)
	if err != nil {
		logrus.Errorf("getProfileID: %v", err)
		return uuid.Nil, echo.ErrInternalServerError
	}
	return profileUUID, nil
}

// Auth is endpoint for auth page
func (h *Handler) Auth(c echo.Context) error {
	tmpl, err := template.ParseFiles("templates/auth/auth.html")
	if err != nil {
		return echo.ErrNotFound
	}
	return tmpl.ExecuteTemplate(c.Response().Writer, "auth", nil)
}

// Index is endpoint for main page
func (h *Handler) Index(c echo.Context) error {
	type PageData struct {
		Orders []*model.Deal
	}
	tmpl, err := template.ParseFiles("templates/index/index.html")
	if err != nil {
		return echo.ErrNotFound
	}
	profileID, err := h.getProfileID(c)
	if err != nil {
		return echo.ErrUnauthorized
	}
	balance, err := h.balanceService.GetBalance(c.Request().Context(), profileID)
	if err != nil {
		logrus.Errorf("index: %v", err)
		return echo.ErrInternalServerError
	}
	orders, err := h.tradingService.GetUnclosedPositions(c.Request().Context(), profileID)
	if err != nil {
		logrus.Errorf("index: %v", err)
		return echo.ErrInternalServerError
	}
	return tmpl.ExecuteTemplate(c.Response().Writer, "index", struct {
		Balance  float64
		PageData PageData
	}{
		Balance:  balance,
		PageData: PageData{Orders: orders},
	})
}

// SignUp calls method of Service by handler
func (h *Handler) SignUp(c echo.Context) error {
	tmpl, err := template.ParseFiles("templates/auth/auth.html")
	if err != nil {
		return echo.ErrNotFound
	}
	var user model.User
	if errBind := c.Bind(&user); errBind != nil {
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "Failed to read fields",
		})
	}
	tempPassword := user.Password
	err = h.validate.StructCtx(c.Request().Context(), user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Login":    user.Login,
			"Password": user.Password,
		}).Errorf("signUp: %v", err)
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "Invalid fields! The fields have not been validated",
		})
	}
	err = h.userService.SignUp(c.Request().Context(), &user)
	if err != nil {
		var e *berrors.BusinessError
		if errors.As(err, &e) {
			return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
				"errorMsg": e.Message,
			})
		}
		logrus.Errorf("signUp: %v", err)
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "Failed to sign up",
		})
	}
	user.Password = tempPassword
	userID, err := h.userService.GetByLogin(c.Request().Context(), &user)
	if err != nil {
		logrus.Errorf("signUp: %v", err)
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "Failed to log in",
		})
	}
	store := NewRedisStore(h.cfg)
	session, err := store.Get(c.Request(), "SESSION_ID")
	if err != nil {
		logrus.Errorf("signUp: %v", err)
		return echo.ErrNotFound
	}
	session.Values["id"] = userID.String()
	session.Values["login"] = user.Login
	session.Values["password"] = user.Password
	if err = session.Save(c.Request(), c.Response()); err != nil {
		logrus.Errorf("signUp: %v", err)
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "Error saving session",
		})
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
	if errBind := c.Bind(&user); errBind != nil {
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "Failed to read fields",
		})
	}
	err = h.validate.StructCtx(c.Request().Context(), user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Login":    user.Login,
			"Password": user.Password,
		}).Errorf("login: %v", err)
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "Invalid fields! The fields have not been validated",
		})
	}
	userID, err := h.userService.GetByLogin(c.Request().Context(), &user)
	if err != nil {
		logrus.Errorf("login: %v", err)
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "Wrong login or password",
		})
	}
	store := NewRedisStore(h.cfg)
	session, err := store.Get(c.Request(), "SESSION_ID")
	if err != nil {
		logrus.Errorf("login: %v", err)
		return echo.ErrNotFound
	}
	session.Values["id"] = userID.String()
	session.Values["login"] = user.Login
	session.Values["password"] = user.Password
	if err = session.Save(c.Request(), c.Response().Writer); err != nil {
		logrus.Errorf("login: %v", err)
		return tmpl.ExecuteTemplate(c.Response().Writer, "auth", map[string]string{
			"errorMsg": "Error saving session",
		})
	}
	return c.Redirect(http.StatusSeeOther, "/index")
}

// DeleteAccount calls method of Service by handler
func (h *Handler) DeleteAccount(c echo.Context) error {
	profileID, err := h.getProfileID(c)
	if err != nil {
		return echo.ErrUnauthorized
	}
	_, err = h.userService.DeleteAccount(c.Request().Context(), profileID)
	if err != nil {
		var e *berrors.BusinessError
		if errors.As(err, &e) {
			return c.HTML(http.StatusBadRequest, `<script>alert('`+e.Message+`');
			window.location.href = '/';</script>`)
		}
		logrus.WithFields(logrus.Fields{
			"ID": profileID,
		}).Errorf("deleteAccount: %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Failed to delete your account');
		window.location.href = '/index';</script>`)
	}
	return c.HTML(http.StatusOK, `<script>alert('Your account has been successfully deleted!');
	 window.location.href = '/';</script>`)
}

// Deposit calls method of Service by handler
func (h *Handler) Deposit(c echo.Context) error {
	profileID, err := h.getProfileID(c)
	if err != nil {
		return echo.ErrUnauthorized
	}
	sumOfMoney, err := strconv.ParseFloat(c.FormValue("operation"), 64)
	if err != nil {
		logrus.Errorf("deposit: %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Invalid sum of money');
		window.location.href = '/index';</script>`)
	}
	balance := model.Balance{
		ProfileID: profileID,
		Operation: sumOfMoney,
	}
	_, err = h.balanceService.BalanceOperation(c.Request().Context(), &balance)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"BalanceId": balance.BalanceID,
			"ProfileId": balance.ProfileID,
			"Operation": balance.Operation,
		}).Errorf("deposit: %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Failed to made balance operation');
		 window.location.href = '/index';</script>`)
	}
	return c.HTML(http.StatusOK, `<script>alert('Deposit of `+fmt.Sprintf("%.2f", sumOfMoney)+
		`$ approved!'); window.location.href = '/index';</script>`)
}

// Withdraw calls method of Service by handler
func (h *Handler) Withdraw(c echo.Context) error {
	profileID, err := h.getProfileID(c)
	if err != nil {
		return echo.ErrUnauthorized
	}
	sumOfMoney, err := strconv.ParseFloat(c.FormValue("operation"), 64)
	if err != nil {
		logrus.Errorf("withdraw: %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Invalid sum of money');
		window.location.href = '/index';</script>`)
	}
	balance := model.Balance{
		ProfileID: profileID,
		Operation: sumOfMoney,
	}
	balance.Operation = -balance.Operation
	_, err = h.balanceService.BalanceOperation(c.Request().Context(), &balance)
	if err != nil {
		var e *berrors.BusinessError
		if errors.As(err, &e) {
			return c.HTML(http.StatusBadRequest, `<script>alert('`+e.Message+`');
			window.location.href = '/index';</script>`)
		}
		logrus.WithFields(logrus.Fields{
			"BalanceId": balance.BalanceID,
			"ProfileId": balance.ProfileID,
			"Operation": balance.Operation,
		}).Errorf("withdraw: %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Failed to made balance operation');
		 window.location.href = '/index';</script>`)
	}
	return c.HTML(http.StatusOK, `<script>alert('Withdraw of `+fmt.Sprintf("%.2f", sumOfMoney)+
		`$ approved!'); window.location.href = '/index';</script>`)
}

// CreatePosition calls method of Service by handler
func (h *Handler) CreatePosition(c echo.Context) error {
	profileID, err := h.getProfileID(c)
	if err != nil {
		return echo.ErrUnauthorized
	}
	sharesCount, err := decimal.NewFromString(c.FormValue("sharescount"))
	if err != nil {
		logrus.Errorf("createPosition: %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Invalid shares count value');
		 window.location.href = '/index';</script>`)
	}
	stopLoss, err := decimal.NewFromString(c.FormValue("stoploss"))
	if err != nil {
		logrus.Errorf("createPosition: %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Invalid stop loss value');
		 window.location.href = '/index';</script>`)
	}
	takeProfit, err := decimal.NewFromString(c.FormValue("takeprofit"))
	if err != nil {
		logrus.Errorf("createPosition: %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Invalid take profit value');
		 window.location.href = '/index';</script>`)
	}
	strategy := "long"
	if stopLoss.Cmp(takeProfit) == 1 {
		strategy = "short"
	}
	deal := &model.Deal{
		ProfileID:   profileID,
		SharesCount: sharesCount,
		Company:     c.FormValue("company"),
		StopLoss:    stopLoss,
		TakeProfit:  takeProfit,
	}
	err = h.tradingService.CreatePosition(c.Request().Context(), deal)
	if err != nil {
		var e *berrors.BusinessError
		if errors.As(err, &e) {
			return c.HTML(http.StatusBadRequest, `<script>alert('`+e.Message+`');
			window.location.href = '/index';</script>`)
		}
		logrus.Errorf("createPosition: %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Failed to create position');
		 window.location.href = '/index';</script>`)
	}
	return c.HTML(http.StatusOK, `<script>alert('Position `+strategy+` created!');
	 window.location.href = '/index';</script>`)
}

// ClosePositionManually calls method of Service by handler
func (h *Handler) ClosePositionManually(c echo.Context) error {
	profileID, err := h.getProfileID(c)
	if err != nil {
		return echo.ErrUnauthorized
	}
	dealUUID, err := uuid.Parse(c.FormValue("dealid"))
	if err != nil {
		logrus.Errorf("closePositionManually: %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Invalid deal ID');
		 window.location.href = '/index';</script>`)
	}
	profit, err := h.tradingService.ClosePositionManually(c.Request().Context(), dealUUID, profileID)
	if err != nil {
		logrus.Errorf("closePositionManually: %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Failed to close position');
		 window.location.href = '/index';</script>`)
	}
	return c.HTML(http.StatusOK, `<script>alert('Position closed with profit `+fmt.Sprintf("%.2f", profit)+`');
	 window.location.href = '/index';</script>`)
}

// GetUnclosedPositions calls method of Service by handler
func (h *Handler) GetUnclosedPositions(c echo.Context) error {
	profileID, err := h.getProfileID(c)
	if err != nil {
		return echo.ErrUnauthorized
	}
	unclosedPositions, err := h.tradingService.GetUnclosedPositions(c.Request().Context(), profileID)
	if err != nil {
		logrus.Errorf("getUnclosedPositions: %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Failed to get positions');
		 window.location.href = '/index';</script>`)
	}
	return c.JSON(http.StatusOK, unclosedPositions)
}

// GetClosedPositions calls method of Service by handler
func (h *Handler) GetClosedPositions(c echo.Context) error {
	profileID, err := h.getProfileID(c)
	if err != nil {
		return echo.ErrUnauthorized
	}
	closedPositions, err := h.tradingService.GetClosedPositions(c.Request().Context(), profileID)
	if err != nil {
		logrus.Errorf("getClosedPositions: %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Failed to get positions');
		 window.location.href = '/index';</script>`)
	}
	return c.JSON(http.StatusOK, closedPositions)
}

// GetPrices calls method of Service by handler
func (h *Handler) GetPrices(c echo.Context) error {
	shares, err := h.tradingService.GetPrices(c.Request().Context())
	if err != nil {
		logrus.Infof("getPrices: %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Failed to get shares');
		 window.location.href = '/index';</script>`)
	}
	return c.JSON(http.StatusOK, shares)
}

// Logout delete session of user
func (h *Handler) Logout(c echo.Context) error {
	store := NewRedisStore(h.cfg)
	cookie, err := c.Cookie("SESSION_ID")
	if err != nil {
		logrus.Errorf("logout %v", err)
		return c.HTML(http.StatusInternalServerError, `<script>alert('Failed to get cookie');
		window.location.href = '/';</script>`)
	}
	session, err := store.Get(c.Request(), cookie.Name)
	if err != nil {
		logrus.Errorf("logout %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Failed to get your session');
		window.location.href = '/';</script>`)
	}
	session.Options.MaxAge = -1
	if err = session.Save(c.Request(), c.Response().Writer); err != nil {
		logrus.Errorf("logout %v", err)
		return c.HTML(http.StatusBadRequest, `<script>alert('Failed to log out');
		 window.location.href = '/index';</script>`)
	}
	return c.Redirect(http.StatusSeeOther, "/")
}
