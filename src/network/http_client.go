package network

import (
	"net/http"
	"time"
)

const requestTimeout = 1 * time.Second

func NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: requestTimeout,
	}
}
