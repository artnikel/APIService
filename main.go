// Package main of a project
package main

import (
	"log"

	"github.com/artnikel/APIService/internal/config"
	"github.com/artnikel/APIService/internal/handler"
	custommiddleware "github.com/artnikel/APIService/internal/middleware"
	"github.com/artnikel/APIService/internal/repository"
	"github.com/artnikel/APIService/internal/service"
	bproto "github.com/artnikel/BalanceService/proto"
	uproto "github.com/artnikel/ProfileService/proto"
	tproto "github.com/artnikel/TradingService/proto"
	"github.com/caarlos0/env"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var cfg config.Variables
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("could not parse config: ", err)
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
	tconn, err := grpc.Dial("localhost:8099", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	usrv := service.NewUserService(urep, &cfg)
	bsrv := service.NewBalanceService(brep)
	tsrv := service.NewTradingService(trep)
	hndl := handler.NewHandler(usrv, bsrv, tsrv,v)
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/signup", hndl.SignUp)
	e.POST("/login", hndl.Login)
	e.POST("/refresh", hndl.Refresh)
	e.DELETE("/delete/:id", hndl.DeleteAccount)
	e.POST("/deposit", hndl.BalanceOperation, custommiddleware.JWTMiddleware)
	e.POST("/withdraw", hndl.BalanceOperation, custommiddleware.JWTMiddleware)
	e.GET("/getbalance", hndl.GetBalance, custommiddleware.JWTMiddleware)
	e.POST("/long",hndl.Strategies,custommiddleware.JWTMiddleware)
	e.POST("/short",hndl.Strategies,custommiddleware.JWTMiddleware)
	e.Logger.Fatal(e.Start(":7777"))
}
