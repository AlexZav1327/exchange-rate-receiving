package httpserver

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	log *logrus.Entry
}

func NewHandler(log *logrus.Logger) *Handler {
	return &Handler{
		log: log.WithField("module", "handler"),
	}
}

func (h *Handler) getList(_ http.ResponseWriter, _ *http.Request) {
}

func (h *Handler) get(_ http.ResponseWriter, _ *http.Request) {
}
