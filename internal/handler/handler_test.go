package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

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
		Password: []byte("testPassword"),
		RefreshToken: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
		eyJleHAiOjE2OTE1MzE2NzAsImlkIjoiMjE5NDkxNjctNTRhOC00NjAwLTk1NzMtM2EwYzAyZTE4NzFjIn0.
		RI9lxDrDlj0RS3FAtNSdwFGz14v9NX1tOxmLjSpZ2dU`,
	}
	tokens = model.TokenPair{
		AccessToken: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
		eyJleHAiOjE2OTEyNzI5MjMsImlkIjoiMjE5NDkxNjctNTRhOC00NjAwLTk1NzMtM2EwYzAyZTE4NzFjIn0.
		X8EOWD4iisVSilCDqxR0kHyaEbplhS5ZitmP9RbUtKk`,
		RefreshToken: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
		eyJleHAiOjE2OTE1MzEyMjMsImlkIjoiMjE5NDkxNjctNTRhOC00NjAwLTk1NzMtM2EwYzAyZTE4NzFjIn0.
		3UGwETfRPcsctV_smpsaq5CQV0MgYACJNHJ91sz9ISk`,
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
	}
	v = validator.New()
)

func TestSignUp(t *testing.T) {
	srv := new(mocks.UserService)
	hndl := NewHandler(srv, nil, nil, v)

	jsonData, err := json.Marshal(testUser)
	require.NoError(t, err)

	srv.On("SignUp", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = hndl.SignUp(c)
	require.NoError(t, err)
	srv.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	srv := new(mocks.UserService)
	hndl := NewHandler(srv, nil, nil, v)

	jsonData, err := json.Marshal(testUser)
	require.NoError(t, err)

	srv.On("Login", mock.Anything, mock.AnythingOfType("*model.User")).Return(&tokens, nil).Once()
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = hndl.Login(c)
	require.NoError(t, err)

	expectedResp := map[string]interface{}{
		"Access Token : ":  tokens.AccessToken,
		"Refresh Token : ": tokens.RefreshToken,
	}
	expectedJSON, err := json.Marshal(expectedResp)
	require.NoError(t, err)

	require.JSONEq(t, string(expectedJSON), rec.Body.String())
	srv.AssertExpectations(t)
}

func TestRefresh(t *testing.T) {
	srv := new(mocks.UserService)
	hndl := NewHandler(srv, nil, nil, v)

	jsonData, err := json.Marshal(tokens)
	require.NoError(t, err)

	srv.On("Refresh", mock.Anything, mock.AnythingOfType("*model.TokenPair")).Return(&tokens, nil).Once()
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = hndl.Refresh(c)
	require.NoError(t, err)

	expectedResp := map[string]interface{}{
		"Access Token : ":  tokens.AccessToken,
		"Refresh Token : ": tokens.RefreshToken,
	}
	expectedJSON, err := json.Marshal(expectedResp)
	require.NoError(t, err)

	require.JSONEq(t, string(expectedJSON), rec.Body.String())
	srv.AssertExpectations(t)
}

func TestDeleteAccount(t *testing.T) {
	srv := new(mocks.UserService)
	hndl := NewHandler(srv, nil, nil, v)

	srv.On("DeleteAccount", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testUser.ID.String(), nil).Once()
	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/delete/:id", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(testUser.ID.String())

	err := hndl.DeleteAccount(c)
	require.NoError(t, err)
	require.Contains(t, rec.Body.String(), testUser.ID.String())
	srv.AssertExpectations(t)
}

func TestDeposit(t *testing.T) {
	srv := new(mocks.BalanceService)
	hndl := NewHandler(nil, srv, nil, v)

	jsonData, err := json.Marshal(testBalance)
	require.NoError(t, err)
	srv.On("GetIDByToken", mock.Anything, mock.AnythingOfType("string")).Return(testBalance.ProfileID, nil).Once()
	srv.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/deposit", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
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
	hndl := NewHandler(nil, srv, nil, v)

	jsonData, err := json.Marshal(testBalance)
	require.NoError(t, err)
	srv.On("GetIDByToken", mock.Anything, mock.AnythingOfType("string")).Return(testBalance.ProfileID, nil).Once()
	srv.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	strOperation := strconv.FormatFloat(testBalance.Operation, 'f', -1, 64)
	err = hndl.Withdraw(c)
	require.NoError(t, err)
	require.Contains(t, rec.Body.String(), strOperation)
	srv.AssertExpectations(t)
}

func TestGetBalance(t *testing.T) {
	srv := new(mocks.BalanceService)
	hndl := NewHandler(nil, srv, nil, v)

	jsonData, err := json.Marshal(testBalance.ProfileID)
	require.NoError(t, err)
	srv.On("GetIDByToken", mock.Anything, mock.AnythingOfType("string")).Return(testBalance.ProfileID, nil).Once()
	srv.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/getbalance", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	strOperation := strconv.FormatFloat(testBalance.Operation, 'f', -1, 64)
	err = hndl.GetBalance(c)
	require.NoError(t, err)
	require.Contains(t, rec.Body.String(), strOperation)
	srv.AssertExpectations(t)
}

func TestLong(t *testing.T) {
	srv := new(mocks.BalanceService)
	hndl := NewHandler(nil, srv, nil, v)
	jsonData, err := json.Marshal(testDeal)
	require.NoError(t, err)
	srv.On("GetIDByToken", mock.Anything, mock.AnythingOfType("string")).Return(testBalance.ProfileID, nil).Once()
	srv.On("Strategies", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*model.Deal")).Return(testDeal.Profit.InexactFloat64(), nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/long", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = hndl.Long(c)
	require.NoError(t, err)
}

func TestShort(t *testing.T) {
	srv := new(mocks.BalanceService)
	hndl := NewHandler(nil, srv, nil, v)
	jsonData, err := json.Marshal(testDeal)
	require.NoError(t, err)
	srv.On("GetIDByToken", mock.Anything, mock.AnythingOfType("string")).Return(testBalance.ProfileID, nil).Once()
	srv.On("Strategies", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*model.Deal")).Return(testDeal.Profit.InexactFloat64(), nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/short", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = hndl.Short(c)
	require.NoError(t, err)
}