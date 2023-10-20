package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/AlexZav1327/exchange-rate-receiving/internal/repo"
	httpserver "github.com/AlexZav1327/exchange-rate-receiving/internal/server"
	"github.com/AlexZav1327/exchange-rate-receiving/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	logger := logrus.StandardLogger()
	currenciesRepo := repo.New(logger)
	currencyService := service.New(currenciesRepo, logger)
	server := httpserver.New("", 8080, currencyService, logger)

	err := server.Run(ctx)
	if err != nil {
		logger.Panicf("server.Run: %s", err)
	}
}
