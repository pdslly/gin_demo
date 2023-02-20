package router

import (
	"path/filepath"
	"stream/app/middleware"
	"stream/config"
	"stream/util/path"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var ProviderSet = wire.NewSet(NewRouter)

func NewRouter(conf *config.Configuration, logger *zap.Logger, ginM *middleware.GinLogger) *gin.Engine {
	if conf.App.Env == "production" {
        gin.SetMode(gin.ReleaseMode)
    }

	router := gin.New()
	router.Use(ginM.Handler(logger))

	rootDir := path.RootPath()

	// 前端网页模板
	router.LoadHTMLGlob(filepath.Join(rootDir, "template/*"))

    // 前端项目静态资源
    router.StaticFile("/", filepath.Join(rootDir, "static/index.html"))
    router.Static("/assets", filepath.Join(rootDir, "static/assets"))
    router.StaticFile("/favicon.ico", filepath.Join(rootDir, "static/favicon.ico"))

	return router
}