package gateway

import (
	"encoding/json"
	"net/http"
	"os"
)

type Temperature struct {
	Temp float32 `json:"temp,omitempty"`
	Date string  `json:"date,omitempty"`
}

type TemperatureGateway struct {
	baseURL    string
	httpClient *HttpClient
}

func NewTemperatureGateway(client *HttpClient) *TemperatureGateway {
	return &TemperatureGateway{
		baseURL:    os.Getenv("TEMPERATURE_BASE_URL"),
		httpClient: client,
	}
}

func (g *TemperatureGateway) GetTemperatureAt(date string) (*Temperature, *HttpError) {
	body, httpError := g.httpClient.MakeRequest(http.MethodGet, g.baseURL+"?at="+date, nil)
	if httpError != nil {
		return nil, httpError
	}

	var temperature Temperature
	if err := json.Unmarshal(body, &temperature); err != nil {
		return nil, HttpErrorBuilder().From(http.StatusText(http.StatusInternalServerError), "Failed to unmarshal temperature response.")
	}

	return &temperature, nil
}
