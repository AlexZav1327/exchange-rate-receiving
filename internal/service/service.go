package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AlexZav1327/exchange-rate-receiving/internal/models"
	"github.com/sirupsen/logrus"
)

var ErrCurrentNotFound = errors.New("no such current")

type Service struct {
	repo CurrenciesRepo
	log  *logrus.Entry
}

type CurrenciesRepo interface {
	SendRequest(ctx context.Context) ([]models.Currency, time.Time, error)
}

type CurrenciesData struct {
	currenciesList  []models.Currency
	currency        models.Currency
	lastRequestTime time.Time
}

var currenciesData = &CurrenciesData{} //nolint:gochecknoglobals

func New(repo CurrenciesRepo, log *logrus.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log.WithField("module", "service"),
	}
}

func (s *Service) GetCurrenciesList(ctx context.Context) ([]models.Currency, error) {
	var err error

	if currenciesData.currenciesList == nil || time.Since(currenciesData.lastRequestTime).Minutes() >= 10 {
		currenciesData.currenciesList, currenciesData.lastRequestTime, err = s.repo.SendRequest(ctx)
		if err != nil {
			return nil, fmt.Errorf("repo.SendRequest: %w", err)
		}

		return currenciesData.currenciesList, nil
	}

	return currenciesData.currenciesList, nil
}

func (s *Service) GetCurrency(ctx context.Context, id string) (models.Currency, error) {
	var err error

	if currenciesData.currenciesList == nil || time.Since(currenciesData.lastRequestTime).Minutes() >= 10 {
		currenciesData.currenciesList, currenciesData.lastRequestTime, err = s.repo.SendRequest(ctx)
		if err != nil {
			return models.Currency{}, fmt.Errorf("repo.SendRequest: %w", err)
		}

		currency, err := s.extractCurrencyFromList(id)
		if err != nil {
			return models.Currency{}, fmt.Errorf("extractCurrencyFromList: %w", err)
		}

		return currency, nil
	}

	currency, err := s.extractCurrencyFromList(id)
	if err != nil {
		return models.Currency{}, fmt.Errorf("extractCurrencyFromList: %w", err)
	}

	return currency, nil
}

func (*Service) extractCurrencyFromList(id string) (models.Currency, error) {
	for i, v := range currenciesData.currenciesList {
		if v.ID == id {
			currenciesData.currency = currenciesData.currenciesList[i]

			return currenciesData.currency, nil
		}
	}

	return models.Currency{}, ErrCurrentNotFound
}
