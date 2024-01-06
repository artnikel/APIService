package handler

// var (
// 	testUser = model.User{
// 		ID:       uuid.New(),
// 		Login:    "testLogin",
// 		Password: "testPassword",
// 	}
// 	testBalance = model.Balance{
// 		BalanceID: uuid.New(),
// 		ProfileID: uuid.New(),
// 		Operation: 637.81,
// 	}
// 	testDeal = model.Deal{
// 		SharesCount: decimal.NewFromFloat(1.5),
// 		Company:     "Apple",
// 		StopLoss:    decimal.NewFromFloat(180.5),
// 		TakeProfit:  decimal.NewFromFloat(500.5),
// 		Profit:      decimal.NewFromFloat(150),
// 	}
// 	testShare = model.Share{
// 		Company: "Apple",
// 		Price:   195.5,
// 	}
// 	v = validator.New()
// )

// func TestSignUp(t *testing.T) {
// 	srv := new(mocks.UserService)
// 	hndl := NewHandler(srv, nil, nil, v)

// 	jsonData, err := json.Marshal(testUser)
// 	require.NoError(t, err)

// 	srv.On("SignUp", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
// 	e := echo.New()

// 	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	err = hndl.SignUp(c)
// 	require.NoError(t, err)
// 	srv.AssertExpectations(t)
// }

// func TestLogin(t *testing.T) {
// 	srv := new(mocks.UserService)
// 	hndl := NewHandler(srv, nil, nil, v)

// 	jsonData, err := json.Marshal(testUser)
// 	require.NoError(t, err)

// 	srv.On("Login", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
// 	e := echo.New()

// 	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	err = hndl.Login(c)
// 	require.NoError(t, err)

// 	require.JSONEq(t, "", rec.Body.String())
// 	srv.AssertExpectations(t)
// }

// func TestDeleteAccount(t *testing.T) {
// 	usrv := new(mocks.UserService)
// 	bsrv := new(mocks.BalanceService)
// 	hndl := NewHandler(usrv, bsrv, nil, v)
// 	jsonData, err := json.Marshal(testBalance.ProfileID)
// 	require.NoError(t, err)
// 	bsrv.On("GetIDByToken", mock.Anything, mock.AnythingOfType("string")).Return(testBalance.ProfileID, nil).Once()
// 	usrv.On("DeleteAccount", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testUser.ID.String(), nil).Once()
// 	e := echo.New()

// 	req := httptest.NewRequest(http.MethodDelete, "/delete", bytes.NewReader(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	err = hndl.DeleteAccount(c)
// 	require.NoError(t, err)
// 	require.Contains(t, rec.Body.String(), testUser.ID.String())
// 	usrv.AssertExpectations(t)
// 	bsrv.AssertExpectations(t)
// }

// func TestDeposit(t *testing.T) {
// 	srv := new(mocks.BalanceService)
// 	hndl := NewHandler(nil, srv, nil, v)

// 	jsonData, err := json.Marshal(testBalance)
// 	require.NoError(t, err)
// 	srv.On("GetIDByToken", mock.Anything, mock.AnythingOfType("string")).Return(testBalance.ProfileID, nil).Once()
// 	srv.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation, nil).Once()

// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/deposit", bytes.NewReader(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	strOperation := strconv.FormatFloat(testBalance.Operation, 'f', -1, 64)
// 	err = hndl.Deposit(c)
// 	require.NoError(t, err)
// 	require.Contains(t, rec.Body.String(), strOperation)
// 	srv.AssertExpectations(t)
// }

// func TestWithdraw(t *testing.T) {
// 	srv := new(mocks.BalanceService)
// 	hndl := NewHandler(nil, srv, nil, v)

// 	jsonData, err := json.Marshal(testBalance)
// 	require.NoError(t, err)
// 	srv.On("GetIDByToken", mock.Anything, mock.AnythingOfType("string")).Return(testBalance.ProfileID, nil).Once()
// 	srv.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation, nil).Once()

// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewReader(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	strOperation := strconv.FormatFloat(testBalance.Operation, 'f', -1, 64)
// 	err = hndl.Withdraw(c)
// 	require.NoError(t, err)
// 	require.Contains(t, rec.Body.String(), strOperation)
// 	srv.AssertExpectations(t)
// }

// func TestGetBalance(t *testing.T) {
// 	srv := new(mocks.BalanceService)
// 	hndl := NewHandler(nil, srv, nil, v)

// 	jsonData, err := json.Marshal(testBalance.ProfileID)
// 	require.NoError(t, err)
// 	srv.On("GetIDByToken", mock.Anything, mock.AnythingOfType("string")).Return(testBalance.ProfileID, nil).Once()
// 	srv.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation, nil).Once()

// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodGet, "/getbalance", bytes.NewReader(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	strOperation := strconv.FormatFloat(testBalance.Operation, 'f', -1, 64)
// 	err = hndl.GetBalance(c)
// 	require.NoError(t, err)
// 	require.Contains(t, rec.Body.String(), strOperation)
// 	srv.AssertExpectations(t)
// }

// func TestLong(t *testing.T) {
// 	srv := new(mocks.BalanceService)
// 	hndl := NewHandler(nil, srv, nil, v)
// 	jsonData, err := json.Marshal(testDeal)
// 	require.NoError(t, err)
// 	srv.On("GetIDByToken", mock.Anything, mock.AnythingOfType("string")).Return(testBalance.ProfileID, nil).Once()
// 	srv.On("CreatePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(testDeal.Profit.InexactFloat64(), nil).Once()

// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodGet, "/long", bytes.NewReader(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	err = hndl.Long(c)
// 	require.NoError(t, err)
// }

// func TestShort(t *testing.T) {
// 	srv := new(mocks.BalanceService)
// 	hndl := NewHandler(nil, srv, nil, v)
// 	jsonData, err := json.Marshal(testDeal)
// 	require.NoError(t, err)
// 	srv.On("GetIDByToken", mock.Anything, mock.AnythingOfType("string")).Return(testBalance.ProfileID, nil).Once()
// 	srv.On("CreatePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(testDeal.Profit.InexactFloat64(), nil).Once()

// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodGet, "/short", bytes.NewReader(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	err = hndl.Short(c)
// 	require.NoError(t, err)
// }

// func TestClosePositionManually(t *testing.T) {
// 	tsrv := new(mocks.TradingService)
// 	bsrv := new(mocks.BalanceService)
// 	hndl := NewHandler(nil, bsrv, tsrv, v)
// 	idParams := &model.Deal{
// 		DealID:    testDeal.DealID,
// 		ProfileID: testDeal.ProfileID,
// 	}
// 	jsonData, err := json.Marshal(idParams)
// 	require.NoError(t, err)
// 	bsrv.On("GetIDByToken", mock.Anything, mock.AnythingOfType("string")).Return(testBalance.ProfileID, nil).Once()
// 	tsrv.On("ClosePositionManually", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
// 		Return(testDeal.Profit.InexactFloat64(), nil).Once()

// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/closeposition", bytes.NewReader(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	err = hndl.ClosePositionManually(c)
// 	require.NoError(t, err)
// 	bsrv.AssertExpectations(t)
// 	tsrv.AssertExpectations(t)
// }

// func TestGetUnclosedPositions(t *testing.T) {
// 	tsrv := new(mocks.TradingService)
// 	bsrv := new(mocks.BalanceService)
// 	hndl := NewHandler(nil, bsrv, tsrv, v)

// 	jsonData, err := json.Marshal(testBalance)
// 	require.NoError(t, err)
// 	var testDeals []*model.Deal
// 	testDeals = append(testDeals, &testDeal)
// 	bsrv.On("GetIDByToken", mock.Anything, mock.AnythingOfType("string")).Return(testBalance.ProfileID, nil).Once()
// 	tsrv.On("GetUnclosedPositions", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testDeals, nil).Once()

// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodGet, "/getunclosed", bytes.NewReader(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	err = hndl.GetUnclosedPositions(c)
// 	require.NoError(t, err)
// 	tsrv.AssertExpectations(t)
// 	bsrv.AssertExpectations(t)
// }

// func TestGetPrices(t *testing.T) {
// 	srv := new(mocks.TradingService)
// 	hndl := NewHandler(nil, nil, srv, v)

// 	jsonData, err := json.Marshal(testBalance)
// 	require.NoError(t, err)
// 	var testShares []model.Share
// 	testShares = append(testShares, testShare)
// 	srv.On("GetPrices", mock.Anything).Return(testShares, nil).Once()

// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodGet, "/getprices", bytes.NewReader(jsonData))
// 	req.Header.Set("Content-Type", "application/json")
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	err = hndl.GetPrices(c)
// 	require.NoError(t, err)
// 	srv.AssertExpectations(t)
// }
