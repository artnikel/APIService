package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/artnikel/APIService/internal/config"
	"github.com/artnikel/APIService/internal/handler/mocks"
	"github.com/artnikel/APIService/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testUser = model.User{
		ID:       uuid.New(),
		Login:    "testLogin",
		Password: "testPassword",
	}
	testBalance = model.Balance{
		BalanceID: uuid.New(),
		ProfileID: uuid.New(),
		Operation: 637.81,
	}
	testDeal = model.Deal{
		SharesCount: decimal.NewFromFloat(1.5),
		Company:     "Apple",
		StopLoss:    decimal.NewFromFloat(180.5),
		TakeProfit:  decimal.NewFromFloat(500.5),
		Profit:      decimal.NewFromFloat(150),
	}
	testShare = model.Share{
		Company: "Apple",
		Price:   195.5,
	}
	v   = validator.New()
	cfg *config.Variables
)

func TestMain(m *testing.M) {
	var err error
	cfg, err = config.New()
	if err != nil {
		log.Fatalf("could not parse config: %v", err)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestSignUp(t *testing.T) {
	srv := new(mocks.UserService)
	hndl := NewHandler(srv, nil, nil, v, cfg)

	formData := url.Values{}
	formData.Set("login", testUser.Login)
	formData.Set("password", testUser.Password)

	originalDir, err := os.Getwd()
	require.NoError(t, err)

	defer func() {
		errCh := os.Chdir(originalDir)
		require.NoError(t, errCh)
	}()

	formDataReader := strings.NewReader(formData.Encode())
	err = os.Chdir("../../../APIService")
	require.NoError(t, err)

	srv.On("SignUp", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
	srv.On("GetByLogin", mock.Anything, mock.AnythingOfType("*model.User")).Return(testUser.ID, nil).Once()
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signup", formDataReader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = hndl.SignUp(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusSeeOther, rec.Code)
	srv.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	srv := new(mocks.UserService)
	hndl := NewHandler(srv, nil, nil, v, cfg)

	formData := url.Values{}
	formData.Set("login", testUser.Login)
	formData.Set("password", testUser.Password)

	originalDir, err := os.Getwd()
	require.NoError(t, err)

	defer func() {
		errCh := os.Chdir(originalDir)
		require.NoError(t, errCh)
	}()

	formDataReader := strings.NewReader(formData.Encode())
	err = os.Chdir("../../../APIService")
	require.NoError(t, err)

	srv.On("GetByLogin", mock.Anything, mock.AnythingOfType("*model.User")).Return(testUser.ID, nil).Once()
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/login", formDataReader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = hndl.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusSeeOther, rec.Code)
	srv.AssertExpectations(t)
}

func TestDeleteAccount(t *testing.T) {
	usrv := new(mocks.UserService)
	bsrv := new(mocks.BalanceService)
	hndl := NewHandler(usrv, bsrv, nil, v, cfg)
	jsonData, err := json.Marshal(testBalance.ProfileID)
	require.NoError(t, err)
	usrv.On("DeleteAccount", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testUser.ID.String(), nil).Once()
	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/delete", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = hndl.DeleteAccount(c)
	require.NoError(t, err)
	usrv.AssertExpectations(t)
	bsrv.AssertExpectations(t)
}

func TestDeposit(t *testing.T) {
	srv := new(mocks.BalanceService)
	hndl := NewHandler(nil, srv, nil, v, cfg)
	store := NewRedisStore(cfg)

	srv.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/deposit", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	req.Form = url.Values{}
	req.Form.Add("operation", strconv.FormatFloat(testBalance.Operation, 'f', -1, 64))
	session, err := store.Get(req, "SESSION_ID")
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "SESSION_ID", Value: session.ID})

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	strOperation := strconv.FormatFloat(testBalance.Operation, 'f', -1, 64)
	err = hndl.Deposit(c)
	require.NoError(t, err)
	require.Contains(t, rec.Body.String(), strOperation)
	srv.AssertExpectations(t)
}

func TestWithdraw(t *testing.T) {
	srv := new(mocks.BalanceService)
	hndl := NewHandler(nil, srv, nil, v, cfg)
	store := NewRedisStore(cfg)

	srv.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/withdraw", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	req.Form = url.Values{}
	req.Form.Add("operation", strconv.FormatFloat(testBalance.Operation, 'f', -1, 64))
	session, err := store.Get(req, "SESSION_ID")
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "SESSION_ID", Value: session.ID})

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	strOperation := strconv.FormatFloat(testBalance.Operation, 'f', -1, 64)
	err = hndl.Withdraw(c)
	require.NoError(t, err)
	require.Contains(t, rec.Body.String(), strOperation)
	srv.AssertExpectations(t)
}

func TestCreatePosition(t *testing.T) {
	srv := new(mocks.TradingService)
	hndl := NewHandler(nil, nil, srv, v, cfg)
	store := NewRedisStore(cfg)

	srv.On("CreatePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/long", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	req.Form = url.Values{}
	req.Form.Add("company", testDeal.Company)
	req.Form.Add("sharescount", testDeal.SharesCount.String())
	req.Form.Add("stoploss", testDeal.StopLoss.String())
	req.Form.Add("takeprofit", testDeal.TakeProfit.String())
	session, err := store.Get(req, "SESSION_ID")
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "SESSION_ID", Value: session.ID})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = hndl.CreatePosition(c)
	require.NoError(t, err)
	srv.AssertExpectations(t)
}

func TestClosePositionManually(t *testing.T) {
	tsrv := new(mocks.TradingService)
	bsrv := new(mocks.BalanceService)
	hndl := NewHandler(nil, bsrv, tsrv, v, cfg)
	store := NewRedisStore(cfg)

	tsrv.On("ClosePositionManually", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(testDeal.Profit.InexactFloat64(), nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/closeposition", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	req.Form = url.Values{}
	req.Form.Add("dealid", testDeal.DealID.String())
	session, err := store.Get(req, "SESSION_ID")
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "SESSION_ID", Value: session.ID})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = hndl.ClosePositionManually(c)
	require.NoError(t, err)
	bsrv.AssertExpectations(t)
	tsrv.AssertExpectations(t)
}

func TestGetUnclosedPositions(t *testing.T) {
	tsrv := new(mocks.TradingService)
	bsrv := new(mocks.BalanceService)
	hndl := NewHandler(nil, bsrv, tsrv, v, cfg)

	var testDeals []*model.Deal
	testDeals = append(testDeals, &testDeal)
	tsrv.On("GetUnclosedPositions", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testDeals, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/getunclosed", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := hndl.GetUnclosedPositions(c)
	require.NoError(t, err)
	tsrv.AssertExpectations(t)
	bsrv.AssertExpectations(t)
}

func TestGetClosedPositions(t *testing.T) {
	tsrv := new(mocks.TradingService)
	bsrv := new(mocks.BalanceService)
	hndl := NewHandler(nil, bsrv, tsrv, v, cfg)

	var testDeals []*model.Deal
	testDeals = append(testDeals, &testDeal)
	tsrv.On("GetClosedPositions", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testDeals, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/getclosed", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := hndl.GetClosedPositions(c)
	require.NoError(t, err)
	tsrv.AssertExpectations(t)
	bsrv.AssertExpectations(t)
}

func TestGetPrices(t *testing.T) {
	srv := new(mocks.TradingService)
	hndl := NewHandler(nil, nil, srv, v, cfg)
	var testShares []model.Share
	testShares = append(testShares, testShare)
	srv.On("GetPrices", mock.Anything).Return(testShares, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/getprices", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := hndl.GetPrices(c)
	require.NoError(t, err)
	srv.AssertExpectations(t)
}
