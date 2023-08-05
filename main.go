// Package main of a project
package main

import (
	"log"

	"github.com/artnikel/APIService/internal/config"
	"github.com/artnikel/APIService/internal/handler"
	"github.com/artnikel/APIService/internal/repository"
	"github.com/artnikel/APIService/internal/service"
	"github.com/artnikel/ProfileService/proto"
	"github.com/caarlos0/env"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var cfg config.Variables
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("could not parse config: ", err)
	}
	v := validator.New()
	conn, err := grpc.Dial("localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("could not connect: %v", err)
	}
	defer func() {
		errConnClose := conn.Close()
		if err != nil {
			logrus.Fatalf("could not close connection: %v", errConnClose)
		}
	}()
	client := proto.NewUserServiceClient(conn)
	rep := repository.NewProfileRepository(client)
	srv := service.NewUserService(rep, &cfg)
	hndl := handler.NewHandler(srv,v)
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/signup",hndl.SignUp)
	e.POST("/login",hndl.Login)
	e.POST("/refresh",hndl.Refresh)
	e.DELETE("/delete/:id",hndl.DeleteAccount)
	e.Logger.Fatal(e.Start(":7777"))

}
