// Package main of a project
package main

import (
	"fmt"
	"log"

	"github.com/artnikel/APIService/internal/config"
	"github.com/artnikel/APIService/internal/handler"
	"github.com/artnikel/APIService/internal/repository"
	"github.com/artnikel/APIService/internal/service"
	bproto "github.com/artnikel/BalanceService/proto"
	uproto "github.com/artnikel/ProfileService/proto"
	tproto "github.com/artnikel/TradingService/proto"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// nolint funlen
func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("could not parse config: %v", err)
	}
	v := validator.New()
	uconn, err := grpc.Dial("localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	bconn, err := grpc.Dial("localhost:8095", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	tconn, err := grpc.Dial("localhost:8088", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer func() {
		errConnClose := uconn.Close()
		if err != nil {
			log.Fatalf("could not close connection: %v", errConnClose)
		}
		errConnClose = bconn.Close()
		if err != nil {
			log.Fatalf("could not close connection: %v", errConnClose)
		}
		errConnClose = tconn.Close()
		if err != nil {
			log.Fatalf("could not close connection: %v", errConnClose)
		}
	}()
	uclient := uproto.NewUserServiceClient(uconn)
	bclient := bproto.NewBalanceServiceClient(bconn)
	tclient := tproto.NewTradingServiceClient(tconn)
	urep := repository.NewProfileRepository(uclient)
	brep := repository.NewBalanceRepository(bclient)
	trep := repository.NewTradingRepository(tclient)
	usrv := service.NewUserService(urep, *cfg)
	bsrv := service.NewBalanceService(brep, *cfg)
	tsrv := service.NewTradingService(trep)
	hndl := handler.NewHandler(usrv, bsrv, tsrv, v, *cfg)
	fmt.Println("API Service started")
	e := echo.New()
	e.Static("/templates", "templates")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	store := handler.NewRedisStore(*cfg)
	store.SetMaxAge(10 * 24 * 3600)
	e.Use(session.Middleware(store))
	e.GET("/", hndl.Auth)
	e.GET("/index", hndl.Index)
	e.POST("/signup", hndl.SignUp)
	e.POST("/login", hndl.Login)
	e.POST("/delete", hndl.DeleteAccount)
	e.POST("/deposit", hndl.Deposit)
	e.POST("/withdraw", hndl.Withdraw)
	e.GET("/getbalance", hndl.GetBalance)
	e.POST("/long", hndl.CreatePosition)
	e.POST("/short", hndl.CreatePosition)
	e.POST("/closeposition", hndl.ClosePositionManually)
	e.GET("/getunclosed", hndl.GetUnclosedPositions)
	e.GET("/getprices", hndl.GetPrices)
	e.POST("/logout", hndl.Logout)
	address := fmt.Sprintf(":%d", cfg.TradingAPIPort)
	e.Logger.Fatal(e.Start(address))
}
