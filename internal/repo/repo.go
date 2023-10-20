package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AlexZav1327/exchange-rate-receiving/internal/models"
	"github.com/sirupsen/logrus"
)

type Repo struct {
	log *logrus.Entry
}

func New(log *logrus.Logger) *Repo {
	return &Repo{
		log: log.WithField("module", "repo"),
	}
}

func (r *Repo) SendRequest(ctx context.Context) ([]models.Currency, time.Time, error) {
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
			r.log.Warningf("resp.Body.Close: %s", err)
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
