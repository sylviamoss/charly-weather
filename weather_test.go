package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
)

type WeatherTestSuite struct {
	suite.Suite
	echo   *echo.Echo
	module *Module
}

func TestWeatherTestSuite(t *testing.T) {
	suite.Run(t, new(WeatherTestSuite))
}

func (suite *WeatherTestSuite) SetupTest() {
	suite.module = NewModule()
	suite.echo = echo.New()
	suite.module.RegisterRoutes(suite.echo)
	suite.populateModuleWithFakeData()
}

func (suite *WeatherTestSuite) TestGetTemperaturesOrderedByDate() {
	// Given
	req := httptest.NewRequest("GET", "/temperatures?start=2018-08-01T12:00:00Z&end=2018-08-02T11:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetTemperature(context)

	// Then
	var temps []Temperature
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusOK)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &temps))
	assert.Equal(suite.T(), len(temps), 2)
	assert.DeepEqual(suite.T(), temps[0], Temperature{
		Temp: 10.5353456000000,
		Date: "2018-08-01T00:00:00Z",
	})
	assert.DeepEqual(suite.T(), temps[1], Temperature{
		Temp: 13.5353456555445,
		Date: "2018-08-02T00:00:00Z",
	})
}

func (suite *WeatherTestSuite) TestGetTemperatureReturnBadRequestWhenMissingStartDate() {
	// Given
	req := httptest.NewRequest("GET", "/temperatures?end=2018-08-02T11:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetTemperature(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusBadRequest)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Please provide both start and end dates",
	})
}

func (suite *WeatherTestSuite) TestGetTemperatureReturnBadRequestWhenMissingEndDate() {
	// Given
	req := httptest.NewRequest("GET", "/temperatures?start=2018-08-01T12:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetTemperature(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusBadRequest)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Please provide both start and end dates",
	})
}

func (suite *WeatherTestSuite) TestGetTemperatureReturnBadRequestWhenStartDateIsMalformed() {
	// Given
	req := httptest.NewRequest("GET", "/temperatures?start=2018-08-01&end=2018-08-01T12:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetTemperature(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusBadRequest)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Please provide dates with format ISO8601 DateTime (eg. 2018-08-12T12:00:00Z)",
	})
}

func (suite *WeatherTestSuite) TestGetTemperatureReturnBadRequestWhenEndDateIsMalformed() {
	// Given
	req := httptest.NewRequest("GET", "/temperatures?start=2018-08-01T12:00:00Z&end=2018-08-20", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetTemperature(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusBadRequest)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Please provide dates with format ISO8601 DateTime (eg. 2018-08-12T12:00:00Z)",
	})
}

func (suite *WeatherTestSuite) TestGetTemperatureReturnInternalServerErrorWhenDataIsNotFound() {
	// Given
	req := httptest.NewRequest("GET", "/temperatures?start=2018-08-01T12:00:00Z&end=2018-08-03T11:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetTemperature(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusInternalServerError)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusNotFound),
		Message: "Resource not found for 2018-08-03T00:00:00Z",
	})
}

func (suite *WeatherTestSuite) TestGetSpeedsOrderedByDate() {
	// Given
	req := httptest.NewRequest("GET", "/speeds?start=2018-08-01T12:00:00Z&end=2018-08-02T11:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetSpeed(context)

	// Then
	var speeds []Windspeed
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusOK)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &speeds))
	assert.Equal(suite.T(), len(speeds), 2)
	assert.DeepEqual(suite.T(), speeds[0], Windspeed{
		North: 9.5353456087290,
		West:  -13.5353456037382,
		Date:  "2018-08-01T00:00:00Z",
	})
	assert.DeepEqual(suite.T(), speeds[1], Windspeed{
		North: 10.5353456026384,
		West:  -15.5353456074028,
		Date:  "2018-08-02T00:00:00Z",
	})
}

func (suite *WeatherTestSuite) TestGetSpeedReturnBadRequestWhenMissingStartDate() {
	// Given
	req := httptest.NewRequest("GET", "/speeds?end=2018-08-02T11:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetSpeed(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusBadRequest)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Please provide both start and end dates",
	})
}

func (suite *WeatherTestSuite) TestGetSpeedReturnBadRequestWhenMissingEndDate() {
	// Given
	req := httptest.NewRequest("GET", "/speeds?start=2018-08-01T12:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetSpeed(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusBadRequest)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Please provide both start and end dates",
	})
}

func (suite *WeatherTestSuite) TestGetSpeedReturnBadRequestWhenStartDateIsMalformed() {
	// Given
	req := httptest.NewRequest("GET", "/speeds?start=2018-08-01&end=2018-08-01T12:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetSpeed(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusBadRequest)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Please provide dates with format ISO8601 DateTime (eg. 2018-08-12T12:00:00Z)",
	})
}

func (suite *WeatherTestSuite) TestGetSpeedReturnBadRequestWhenEndDateIsMalformed() {
	// Given
	req := httptest.NewRequest("GET", "/speeds?start=2018-08-01T12:00:00Z&end=2018-08-20", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetSpeed(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusBadRequest)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Please provide dates with format ISO8601 DateTime (eg. 2018-08-12T12:00:00Z)",
	})
}

func (suite *WeatherTestSuite) TestGetSpeedReturnInternalServerErrorWhenDataIsNotFound() {
	// Given
	req := httptest.NewRequest("GET", "/speeds?start=2018-08-01T12:00:00Z&end=2018-08-03T11:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetSpeed(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusInternalServerError)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusNotFound),
		Message: "Resource not found for 2018-08-03T00:00:00Z",
	})
}

func (suite *WeatherTestSuite) TestGetWeathersOrderedByDate() {
	// Given
	req := httptest.NewRequest("GET", "/weather?start=2018-08-01T12:00:00Z&end=2018-08-02T11:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetWeather(context)

	// Then
	var weathers []Weather
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusOK)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &weathers))
	assert.Equal(suite.T(), len(weathers), 2)
	assert.DeepEqual(suite.T(), weathers[0], Weather{
		North: 9.5353456087290,
		West:  -13.5353456037382,
		Temp:  10.5353456000000,
		Date:  "2018-08-01T00:00:00Z",
	})
	assert.DeepEqual(suite.T(), weathers[1], Weather{
		North: 10.5353456026384,
		West:  -15.5353456074028,
		Temp:  13.5353456555445,
		Date:  "2018-08-02T00:00:00Z",
	})
}

func (suite *WeatherTestSuite) TestGetWeatherReturnBadRequestWhenMissingStartDate() {
	// Given
	req := httptest.NewRequest("GET", "/weather?end=2018-08-02T11:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetWeather(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusBadRequest)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Please provide both start and end dates",
	})
}

func (suite *WeatherTestSuite) TestGetWeatherReturnBadRequestWhenMissingEndDate() {
	// Given
	req := httptest.NewRequest("GET", "/weather?start=2018-08-01T12:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetWeather(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusBadRequest)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Please provide both start and end dates",
	})
}

func (suite *WeatherTestSuite) TestGetWeatherReturnBadRequestWhenStartDateIsMalformed() {
	// Given
	req := httptest.NewRequest("GET", "/weather?start=2018-08-01&end=2018-08-01T12:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetWeather(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusBadRequest)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Please provide dates with format ISO8601 DateTime (eg. 2018-08-12T12:00:00Z)",
	})
}

func (suite *WeatherTestSuite) TestGetWeatherReturnBadRequestWhenEndDateIsMalformed() {
	// Given
	req := httptest.NewRequest("GET", "/weather?start=2018-08-01T12:00:00Z&end=2018-08-20", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetWeather(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusBadRequest)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusBadRequest),
		Message: "Please provide dates with format ISO8601 DateTime (eg. 2018-08-12T12:00:00Z)",
	})
}

func (suite *WeatherTestSuite) TestGetWeatherReturnInternalServerErrorWhenDataIsNotFound() {
	// Given
	req := httptest.NewRequest("GET", "/weather?start=2018-08-01T12:00:00Z&end=2018-08-03T11:00:00Z", nil)
	rec := httptest.NewRecorder()
	context := suite.echo.NewContext(req, rec)

	// When
	err := suite.module.GetWeather(context)

	// Then
	var httpError HttpError
	assert.NilError(suite.T(), err)
	assert.Equal(suite.T(), rec.Code, http.StatusInternalServerError)
	assert.NilError(suite.T(), json.Unmarshal(rec.Body.Bytes(), &httpError))
	assert.DeepEqual(suite.T(), httpError, HttpError{
		Type:    http.StatusText(http.StatusNotFound),
		Message: "Resource not found for 2018-08-03T00:00:00Z",
	})
}

func (suite *WeatherTestSuite) populateModuleWithFakeData() {
	windspeeds := make(map[string]Windspeed)
	windspeeds["2018-08-02T00:00:00Z"] = Windspeed{
		North: 10.5353456026384,
		West:  -15.5353456074028,
		Date:  "2018-08-02T00:00:00Z",
	}
	windspeeds["2018-08-01T00:00:00Z"] = Windspeed{
		North: 9.5353456087290,
		West:  -13.5353456037382,
		Date:  "2018-08-01T00:00:00Z",
	}
	suite.module.speeds = &WindspeedGatewayMock{
		speeds: windspeeds,
	}

	temperatures := make(map[string]Temperature)
	temperatures["2018-08-02T00:00:00Z"] = Temperature{
		Temp: 13.5353456555445,
		Date: "2018-08-02T00:00:00Z",
	}
	temperatures["2018-08-01T00:00:00Z"] = Temperature{
		Temp: 10.5353456000000,
		Date: "2018-08-01T00:00:00Z",
	}
	suite.module.temperatures = &TemperatureGatewayMock{
		temperatures: temperatures,
	}
}

type TemperatureGatewayMock struct {
	temperatures map[string]Temperature
}

func (g *TemperatureGatewayMock) GetResourceAt(date string, resource interface{}) *HttpError {
	if reflect.DeepEqual(g.temperatures[date], Temperature{}) {
		return &HttpError{http.StatusText(http.StatusNotFound), "Resource not found for " + date}
	}
	jsonTemp, _ := json.Marshal(g.temperatures[date])
	_ = json.Unmarshal(jsonTemp, resource)
	return nil
}

type WindspeedGatewayMock struct {
	speeds map[string]Windspeed
}

func (g *WindspeedGatewayMock) GetResourceAt(date string, resource interface{}) *HttpError {
	if reflect.DeepEqual(g.speeds[date], Windspeed{}) {
		return &HttpError{http.StatusText(http.StatusNotFound), "Resource not found for " + date}
	}
	jsonTemp, _ := json.Marshal(g.speeds[date])
	_ = json.Unmarshal(jsonTemp, resource)
	return nil
}
