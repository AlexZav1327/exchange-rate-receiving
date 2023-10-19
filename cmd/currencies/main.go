package main

import (
	"context"
	"os/signal"
	"syscall"

	httpserver "github.com/AlexZav1327/exchange-rate-receiving/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	logger := logrus.StandardLogger()

	server := httpserver.New("", 8080, logger)

	err := server.Run(ctx)
	if err != nil {
		logger.Panicf("server.Run: %s", err)
	}
}
