package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const baseURL = "https://api.etherscan.io/api"

func GetBalance(address string) (string, error) {
	key := os.Getenv("ETHERSCAN_API_KEY")
	url := fmt.Sprintf("%s?module=account&action=balance&address=%s&tag=latest&apikey=%s", baseURL, address, key)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var r struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}
	if r.Status != "1" {
		return "", fmt.Errorf("etherscan error: %s", r.Message)
	}
	return r.Result, nil
}

type Tx struct {
	Hash      string `json:"hash"`
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
	TimeStamp string `json:"timeStamp"`
}

func GetTransactions(address string) ([]Tx, error) {
	key := os.Getenv("ETHERSCAN_API_KEY")
	url := fmt.Sprintf("%s?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s",
		baseURL, address, key)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  []Tx   `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if r.Status != "1" {
		return nil, fmt.Errorf("etherscan error: %s", r.Message)
	}
	return r.Result, nil
}
