package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func GetTokenPrices(ids []string, vsCurrency string) (map[string]struct {
	Price       float64 `json:"current_price"`
	LastUpdated string  `json:"last_updated"`
}, error) {
	url := fmt.Sprintf("%s/coins/markets?vs_currency=%s&ids=%s",
		os.Getenv("COINGECKO_API_URL"), vsCurrency, strings.Join(ids, ","))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []struct {
		ID            string  `json:"id"`
		CurrentPrice  float64 `json:"current_price"`
		LastUpdatedAt string  `json:"last_updated"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	m := make(map[string]struct {
		Price       float64
		LastUpdated string
	})
	for _, coin := range data {
		m[coin.ID] = struct {
			Price       float64
			LastUpdated string
		}{coin.CurrentPrice, coin.LastUpdatedAt}
	}
	return m, nil
}
