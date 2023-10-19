package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	host   string
	port   int
	Server *http.Server
	log    *logrus.Entry
}

func New(host string, port int, log *logrus.Logger) *Server {
	server := Server{
		host: host,
		port: port,
		log:  log.WithField("module", "http"),
	}

	server.Server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
		ReadHeaderTimeout: 30 * time.Second,
	}

	return &server
}

func (s *Server) Run(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	go func() {
		<-ctx.Done()

		err := s.Server.Shutdown(shutdownCtx)
		if err != nil {
			s.log.Warningf("Server.Shutdown: %s", err)
		}
	}()

	err := s.Server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("Server.ListenAndServe: %w", err)
	}

	return nil
}
