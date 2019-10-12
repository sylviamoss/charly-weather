package weather

import (
	"encoding/json"
	"net/http"
	"os"
)

type Gateway struct {
	baseURL    string
	httpClient *HttpClient
}

func NewTemperatureGateway(client *HttpClient) *Gateway {
	return &Gateway{
		baseURL:    os.Getenv("TEMPERATURE_BASE_URL"),
		httpClient: client,
	}
}

func NewWindspeedGateway(client *HttpClient) *Gateway {
	return &Gateway{
		baseURL:    os.Getenv("WINDSPEED_BASE_URL"),
		httpClient: client,
	}
}

func (g *Gateway) GetResourceAt(date string, resource interface{}) *HttpError {
	body, httpError := g.httpClient.MakeRequest(http.MethodGet, g.baseURL+"?at="+date, nil)
	if httpError != nil {
		return httpError
	}

	if err := json.Unmarshal(body, resource); err != nil {
		return HttpErrorBuilder().From(http.StatusText(http.StatusInternalServerError), "Failed to unmarshal resource response.")
	}

	return nil
}
