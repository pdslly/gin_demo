package main

import (
	"context"
	"net/http"
	"stream/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	conf *config.Configuration
	logger *zap.Logger
	httpSrv *http.Server
}

func newHttpServer(
    conf *config.Configuration,
    router *gin.Engine,
    ) *http.Server {
    return &http.Server{
        Addr:    ":" + conf.App.Port,
        Handler: router,
    }
}

func newApp(
	conf *config.Configuration,
	logger *zap.Logger,
	httpSrv *http.Server,
) *App  {
	return &App{conf, logger, httpSrv}
}

func (a *App) Run() error {
	a.logger.Info("http server started")
	err := a.httpSrv.ListenAndServe()
	return err
}

func (a *App) Stop(ctx context.Context) error  {
	return nil
}