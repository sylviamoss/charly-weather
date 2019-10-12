package main

import (
	"encoding/json"
	"net/http"
	"os"
)

type Gateway interface {
	GetResourceAt(date string, resource interface{}) *HttpError
}

type GatewayModule struct {
	baseURL    string
	httpClient *HttpClient
}

func NewTemperatureGateway(client *HttpClient) *GatewayModule {
	return &GatewayModule{
		baseURL:    os.Getenv("TEMPERATURE_BASE_URL"),
		httpClient: client,
	}
}

func NewWindspeedGateway(client *HttpClient) *GatewayModule {
	return &GatewayModule{
		baseURL:    os.Getenv("WINDSPEED_BASE_URL"),
		httpClient: client,
	}
}

func (g *GatewayModule) GetResourceAt(date string, resource interface{}) *HttpError {
	body, httpError := g.httpClient.MakeRequest(http.MethodGet, g.baseURL+"?at="+date)
	if httpError != nil {
		return httpError
	}

	if err := json.Unmarshal(body, resource); err != nil {
		return HttpErrorBuilder().From(http.StatusText(http.StatusInternalServerError), "Failed to unmarshal resource response.")
	}

	return nil
}
