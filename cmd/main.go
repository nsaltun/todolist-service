package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/nsaltun/todolist-service/config"
	_ "github.com/nsaltun/todolist-service/pkg/logging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.NewAppConfig()
	defer zap.L().Sync()

	zap.L().Info("app is starting..")

	fiberApp := fiber.New()
	fiberApp.Get("/", func(c *fiber.Ctx) error {
		zap.L().Info("request received")
		return c.SendString("Hello, World!")
	})

	fiberApp.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	fiberApp.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	go func() {
		if err := fiberApp.Listen(fmt.Sprintf(":%s", appConfig.HTTPPort)); err != nil {
			zap.L().Fatal("failed to start server", zap.Error(err))
		}
	}()
	zap.L().Info("server started on port", zap.String("port", appConfig.HTTPPort))

	gracefulShutdown(fiberApp)
}

func gracefulShutdown(fiberApp *fiber.App) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	zap.L().Info("shutting down server..")
	if err := fiberApp.ShutdownWithTimeout(5 * time.Second); err != nil {
		zap.L().Fatal("failed to shutdown server", zap.Error(err))
	}
	zap.L().Info("server gracefully stopped")
}
