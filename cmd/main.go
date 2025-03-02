package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger := createLogger()
	zap.ReplaceGlobals(logger)
	defer logger.Sync() // flushes buffer, if any

	zap.L().Info("app is starting..")

	fiberApp := fiber.New()
	fiberApp.Get("/", func(c *fiber.Ctx) error {
		zap.L().Info("request received")
		return c.SendString("Hello, World!")
	})

	fiberApp.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	if err := fiberApp.Listen(":3000"); err != nil {
		zap.L().Fatal("failed to start server", zap.Error(err))
	}
	zap.L().Info("server started")
}

func createLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	return zap.Must(config.Build())
}
