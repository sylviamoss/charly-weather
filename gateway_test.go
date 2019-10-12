package main

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
)

type GatewayTestSuite struct {
	suite.Suite
	gateway *GatewayModule
}

func TestGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(GatewayTestSuite))
}

func (suite *GatewayTestSuite) SetupTest() {
	os.Setenv("TEMPERATURE_BASE_URL", "http://baseurl.com")
	os.Setenv("WINDSPEED_BASE_URL", "http://baseurl.com")
}

func (suite *GatewayTestSuite) TestTemperatureGatewayShouldReturnTemperatureForAGivenDate() {
	// Given
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), r.Method, http.MethodGet)
		assert.Equal(suite.T(), "at=2018-08-12T12:00:00Z", r.URL.RawQuery)

		w.WriteHeader(http.StatusOK)
		json, _ := json.Marshal(Temperature{
			Temp: 10.46941232124016,
			Date: "2018-08-12T00:00:00Z",
		})
		w.Write([]byte(json))
	})
	httpClient := &HttpClient{
		client: NewHttpClientForTesting(handler),
	}
	suite.gateway = NewTemperatureGateway(httpClient)

	// When
	var temperature Temperature
	err := suite.gateway.GetResourceAt("2018-08-12T12:00:00Z", &temperature)

	// Then
	assert.Assert(suite.T(), err == nil)
	assert.DeepEqual(suite.T(), temperature, Temperature{
		Temp: 10.46941232124016,
		Date: "2018-08-12T00:00:00Z",
	})
}

func (suite *GatewayTestSuite) TestWindspeedGatewayShouldReturnWindspeedForAGivenDate() {
	// Given
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), r.Method, http.MethodGet)
		assert.Equal(suite.T(), "at=2018-08-12T12:00:00Z", r.URL.RawQuery)

		w.WriteHeader(http.StatusOK)
		json, _ := json.Marshal(Windspeed{
			North: -17.46941232124016,
			West:  16.46941232124016,
			Date:  "2018-08-12T00:00:00Z",
		})
		w.Write([]byte(json))
	})
	httpClient := &HttpClient{
		client: NewHttpClientForTesting(handler),
	}
	suite.gateway = NewWindspeedGateway(httpClient)

	// When
	var speed Windspeed
	err := suite.gateway.GetResourceAt("2018-08-12T12:00:00Z", &speed)

	// Then
	assert.Assert(suite.T(), err == nil)
	assert.DeepEqual(suite.T(), speed, Windspeed{
		North: -17.46941232124016,
		West:  16.46941232124016,
		Date:  "2018-08-12T00:00:00Z",
	})
}

func (suite *GatewayTestSuite) TestGatewayShouldReturnErrorWhenStatusCodeIsNotOK() {
	// Given
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(HttpError{
			Message: "Date is not a valid RFC3339 DateTime",
		})
		w.Write([]byte(json))
	})
	httpClient := &HttpClient{
		client: NewHttpClientForTesting(handler),
	}
	suite.gateway = NewTemperatureGateway(httpClient)

	// When
	var temperature Temperature
	err := suite.gateway.GetResourceAt("", &temperature)

	// Then
	assert.DeepEqual(suite.T(), *err, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Date is not a valid RFC3339 DateTime",
	})
	assert.DeepEqual(suite.T(), temperature, Temperature{})
}

func (suite *GatewayTestSuite) TestGatewayShouldReturnErrorWhenStatusBodyIsEmpty() {
	// Given
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	httpClient := &HttpClient{
		client: NewHttpClientForTesting(handler),
	}
	suite.gateway = NewTemperatureGateway(httpClient)

	// When
	var temperature Temperature
	err := suite.gateway.GetResourceAt("", &temperature)

	// Then
	assert.DeepEqual(suite.T(), *err, HttpError{
		Type:    http.StatusText(http.StatusInternalServerError),
		Message: "Failed to unmarshal resource response.",
	})
	assert.DeepEqual(suite.T(), temperature, Temperature{})
}
