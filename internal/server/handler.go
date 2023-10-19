package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AlexZav1327/exchange-rate-receiving/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service CurrenciesService
	log     *logrus.Entry
}

type CurrenciesService interface {
	GetCurrenciesList(ctx context.Context) ([]models.Currency, error)
	GetCurrency(ctx context.Context, id string) (models.Currency, error)
}

func NewHandler(service CurrenciesService, log *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log.WithField("module", "handler"),
	}
}

func (h *Handler) getList(w http.ResponseWriter, r *http.Request) {
	currenciesList, err := h.service.GetCurrenciesList(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		h.log.Infof("err: %s", err)

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(currenciesList)
	if err != nil {
		h.log.Warningf("json.NewEncoder.Encode: %s", err)
	}
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	currency, err := h.service.GetCurrency(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(currency)
	if err != nil {
		h.log.Warningf("json.NewEncoder.Encode: %s", err)
	}
}
