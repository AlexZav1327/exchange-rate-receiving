package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AlexZav1327/exchange-rate-receiving/internal/models"
	"github.com/sirupsen/logrus"
)

var ErrCurrentNotFound = errors.New("no such current")

type Service struct {
	log *logrus.Entry
}

type CurrenciesData struct {
	currenciesList  []models.Currency
	currency        models.Currency
	lastRequestTime time.Time
}

var currenciesData = &CurrenciesData{} //nolint:gochecknoglobals

func New(log *logrus.Logger) *Service {
	return &Service{
		log: log.WithField("module", "service"),
	}
}

func (s *Service) GetCurrenciesList(ctx context.Context) ([]models.Currency, error) {
	var err error

	if currenciesData.currenciesList == nil || time.Since(currenciesData.lastRequestTime).Minutes() >= 10 {
		currenciesData.currenciesList, currenciesData.lastRequestTime, err = s.sendRequest(ctx)
		if err != nil {
			return nil, fmt.Errorf("sendRequest: %w", err)
		}

		return currenciesData.currenciesList, nil
	}

	return currenciesData.currenciesList, nil
}

func (s *Service) GetCurrency(ctx context.Context, id string) (models.Currency, error) {
	var err error

	if currenciesData.currenciesList == nil || time.Since(currenciesData.lastRequestTime).Minutes() >= 10 {
		currenciesData.currenciesList, currenciesData.lastRequestTime, err = s.sendRequest(ctx)
		if err != nil {
			return models.Currency{}, fmt.Errorf("sendRequest: %w", err)
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

func (s *Service) sendRequest(ctx context.Context) ([]models.Currency, time.Time, error) {
	endpoint := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1"

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("http.DefaultClient.Do: %w", err)
	}

	defer func() {
		err = response.Body.Close()
		if err != nil {
			s.log.Warningf("resp.Body.Close: %s", err)
		}
	}()

	var currenciesList []models.Currency

	err = json.NewDecoder(response.Body).Decode(&currenciesList)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("json.NewDecoder.Decode: %w", err)
	}

	lastRequestTime := time.Now()

	return currenciesList, lastRequestTime, nil
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
