package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/nsaltun/todolist-service/app/healthcheck"
	"github.com/nsaltun/todolist-service/app/todoitem"
	"github.com/nsaltun/todolist-service/config"
	"github.com/nsaltun/todolist-service/infra/postgres"
	"github.com/nsaltun/todolist-service/pkg/httphandler"
	_ "github.com/nsaltun/todolist-service/pkg/logging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.NewAppConfig()
	defer zap.L().Sync()

	zap.L().Info("app is starting..")

	postgresConn := postgres.NewPostgresConnection(appConfig.PostgresConfig)
	defer postgresConn.Close()

	todoRepo := postgres.NewTodoRepository(postgresConn)

	healthCheckHandler := healthcheck.NewHealthCheckHandler()
	todoItemsGetHandler := todoitem.NewGetTodoItemsHandler(todoRepo)
	todoItemsCreateHandler := todoitem.NewCreateTodoItemHandler(todoRepo)

	fiberApp := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Concurrency:  256 * 1024,
	})

	fiberApp.Get("/", func(c *fiber.Ctx) error {
		zap.L().Info("request received")
		return c.SendString("Hello, World!")
	})

	fiberApp.Get("/healthcheck", httphandler.Handle(healthCheckHandler))
	fiberApp.Get("/todoitems", httphandler.Handle(todoItemsGetHandler))
	fiberApp.Post("/todoitems", httphandler.Handle(todoItemsCreateHandler))

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

func httpClient() *http.Client {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	return httpClient
}
