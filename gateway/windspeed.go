package gateway

import (
	"encoding/json"
	"net/http"
	"os"
)

type Windspeed struct {
	North float32 `json:"north,omitempty"`
	West  float32 `json:"west,omitempty"`
	Date  string  `json:"date,omitempty"`
}

type WindspeedGateway struct {
	baseURL    string
	httpClient *HttpClient
}

func NewWindspeedGateway(client *HttpClient) *WindspeedGateway {
	return &WindspeedGateway{
		baseURL:    os.Getenv("WINDSPEED_BASE_URL"),
		httpClient: client,
	}
}

func (g *WindspeedGateway) GetWindspeedAt(date string) (*Windspeed, *HttpError) {
	body, httpError := g.httpClient.MakeRequest(http.MethodGet, g.baseURL+"?at="+date, nil)
	if httpError != nil {
		return nil, httpError
	}

	var speed Windspeed
	if err := json.Unmarshal(body, &speed); err != nil {
		return nil, HttpErrorBuilder().From(http.StatusText(http.StatusInternalServerError), "Failed to unmarshal windspeed response.")
	}

	return &speed, nil
}
