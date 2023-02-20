package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"stream/config"
	"stream/util/path"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	rootPath = path.RootPath()

	configPath string
	conf *config.Configuration
	logger *zap.Logger
	loggerWriter *lumberjack.Logger
)

func init() {
	pflag.StringVarP(&configPath, "conf", "", filepath.Join(rootPath, "config.yaml"), "config path, eg: --conf config.yaml")

	cobra.OnInitialize(func() {
		initConfig()
		initLogger()
	})	
}

func main() {
	rootCmd := &cobra.Command{
		Use: "app",
		Run: func(cmd *cobra.Command, args []string) {
			app, cleanup, err := wireApp(conf, loggerWriter, logger)
			if err != nil {
                panic(err)
            }
            defer cleanup()

			log.Printf("start app %s ...", conf.App.Port)
			if err := app.Run(); err != nil {
                panic(err)
            }

			quit := make(chan os.Signal)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<- quit

			log.Println("shutdown app!")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := app.Stop(ctx); err != nil {
				panic(err)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
        panic(err)
    }
}

func initConfig()  {
	if !filepath.IsAbs(configPath) {
        configPath = filepath.Join(rootPath, "conf", configPath)
    }

	fmt.Println("load config: ", configPath)

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %s", err))
	}

	if err := v.Unmarshal(&conf); err != nil {
		panic(err)
	}
}

func initLogger() {
    var level zapcore.Level // zap 日志等级
    var options []zap.Option // zap 配置项

    logFileDir := conf.Log.RootDir
    if !filepath.IsAbs(logFileDir) {
        logFileDir = filepath.Join(rootPath, logFileDir)
    }

    if ok, _ := path.Exists(logFileDir); !ok {
        _ = os.Mkdir(conf.Log.RootDir, os.ModePerm)
    }

    switch conf.Log.Level {
    case "debug":
        level = zap.DebugLevel
        options = append(options, zap.AddStacktrace(level))
    case "info":
        level = zap.InfoLevel
    case "warn":
        level = zap.WarnLevel
    case "error":
        level = zap.ErrorLevel
        options = append(options, zap.AddStacktrace(level))
    case "dpanic":
        level = zap.DPanicLevel
    case "panic":
        level = zap.PanicLevel
    case "fatal":
        level = zap.FatalLevel
    default:
        level = zap.InfoLevel
    }

    // 调整编码器默认配置
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
        encoder.AppendString(time.Format("2006-01-02 15:04:05.000"))
    }
    encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
        encoder.AppendString(conf.App.Env + "." + l.String())
    }

    loggerWriter = &lumberjack.Logger{
        Filename:   filepath.Join(logFileDir, conf.Log.Filename),
        MaxSize:    conf.Log.MaxSize,
        MaxBackups: conf.Log.MaxBackups,
        MaxAge:     conf.Log.MaxAge,
        Compress:   conf.Log.Compress,
    }

    logger = zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(loggerWriter), level), options...)
}
