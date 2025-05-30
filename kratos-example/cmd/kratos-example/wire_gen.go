// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-example/internal/biz"
	"kratos-example/internal/conf"
	"kratos-example/internal/data"
	"kratos-example/internal/server"
	"kratos-example/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewDB(confData)
	dataData, cleanup, err := data.NewData(confData, logger, db)
	if err != nil {
		return nil, nil, err
	}
	exampleRepo := data.NewExampleRepo(dataData, logger)
	exampleUsecase := biz.NewExampleUsecase(exampleRepo, logger)
	exampleService := service.NewExampleService(exampleUsecase)
	grpcServer := server.NewGRPCServer(confServer, exampleService, logger)
	httpServer := server.NewHTTPServer(confServer, exampleService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
