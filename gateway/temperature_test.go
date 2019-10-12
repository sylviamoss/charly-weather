package gateway

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
)

type TemperatureGatewayTestSuite struct {
	suite.Suite
	gateway *TemperatureGateway
}

func TestTemperatureGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(TemperatureGatewayTestSuite))
}

func (suite *TemperatureGatewayTestSuite) SetupTest() {
	os.Setenv("TEMPERATURE_BASE_URL", "http://baseurl.com")
}

func (suite *TemperatureGatewayTestSuite) TestGatewayShouldReturnTemperatureForAGivenDate() {
	// Given
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), r.Method, http.MethodGet)
		assert.Equal(suite.T(), "application/json", r.Header.Get("Content-Type"))
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
	temperature, err := suite.gateway.GetTemperatureAt("2018-08-12T12:00:00Z")

	// Then
	assert.Assert(suite.T(), err == nil)
	assert.DeepEqual(suite.T(), *temperature, Temperature{
		Temp: 10.46941232124016,
		Date: "2018-08-12T00:00:00Z",
	})
}

func (suite *TemperatureGatewayTestSuite) TestGatewayShouldReturnErrorWhenStatusCodeIsNotOK() {
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
	temperature, err := suite.gateway.GetTemperatureAt("")

	// Then
	assert.DeepEqual(suite.T(), *err, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Date is not a valid RFC3339 DateTime",
	})
	assert.Assert(suite.T(), temperature == nil)
}

func (suite *TemperatureGatewayTestSuite) TestGatewayShouldReturnErrorWhenStatusBodyIsEmpty() {
	// Given
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	httpClient := &HttpClient{
		client: NewHttpClientForTesting(handler),
	}
	suite.gateway = NewTemperatureGateway(httpClient)

	// When
	temperature, err := suite.gateway.GetTemperatureAt("")

	// Then
	assert.DeepEqual(suite.T(), *err, HttpError{
		Type:    http.StatusText(http.StatusInternalServerError),
		Message: "Failed to unmarshal temperature response.",
	})
	assert.Assert(suite.T(), temperature == nil)
}
