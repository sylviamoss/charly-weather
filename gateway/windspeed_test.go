package gateway

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
)

type WindspeedGatewayTestSuite struct {
	suite.Suite
	gateway *WindspeedGateway
}

func TestWindspeedGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(WindspeedGatewayTestSuite))
}

func (suite *WindspeedGatewayTestSuite) SetupTest() {
	os.Setenv("WINDSPEED_BASE_URL", "http://baseurl.com")
}

func (suite *WindspeedGatewayTestSuite) TestGatewayShouldReturnWindspeedForAGivenDate() {
	// Given
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(suite.T(), r.Method, http.MethodGet)
		assert.Equal(suite.T(), "application/json", r.Header.Get("Content-Type"))
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
	speed, err := suite.gateway.GetWindspeedAt("2018-08-12T12:00:00Z")

	// Then
	assert.Assert(suite.T(), err == nil)
	assert.DeepEqual(suite.T(), *speed, Windspeed{
		North: -17.46941232124016,
		West:  16.46941232124016,
		Date:  "2018-08-12T00:00:00Z",
	})
}

func (suite *WindspeedGatewayTestSuite) TestGatewayShouldReturnErrorWhenStatusCodeIsNotOK() {
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
	suite.gateway = NewWindspeedGateway(httpClient)

	// When
	speed, err := suite.gateway.GetWindspeedAt("")

	// Then
	assert.DeepEqual(suite.T(), *err, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Date is not a valid RFC3339 DateTime",
	})
	assert.Assert(suite.T(), speed == nil)
}

func (suite *WindspeedGatewayTestSuite) TestGatewayShouldReturnErrorWhenStatusBodyIsEmpty() {
	// Given
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	httpClient := &HttpClient{
		client: NewHttpClientForTesting(handler),
	}
	suite.gateway = NewWindspeedGateway(httpClient)

	// When
	speed, err := suite.gateway.GetWindspeedAt("")

	// Then
	assert.DeepEqual(suite.T(), *err, HttpError{
		Type:    http.StatusText(http.StatusInternalServerError),
		Message: "Failed to unmarshal windspeed response.",
	})
	assert.Assert(suite.T(), speed == nil)
}
