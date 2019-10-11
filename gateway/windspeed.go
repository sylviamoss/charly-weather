package gateway

import (
	"net/http"
	"os"
	"time"
)

type WindspeedGateway struct {
	baseURL    string
	httpClient *http.Client
}

func NewWindspeedGateway() *WindspeedGateway {
	return &WindspeedGateway{
		baseURL: os.Getenv("WINDSPEED_BASE_URL"),
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}
