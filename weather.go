package main

import (
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Temperature struct {
	Temp float64 `json:"temp,omitempty"`
	Date string  `json:"date,omitempty"`
}

type Windspeed struct {
	North float64 `json:"north,omitempty"`
	West  float64 `json:"west,omitempty"`
	Date  string  `json:"date,omitempty"`
}

type Weather struct {
	North float64 `json:"north,omitempty"`
	West  float64 `json:"west,omitempty"`
	Temp  float64 `json:"temp,omitempty"`
	Date  string  `json:"date,omitempty"`
}

type Module struct {
	logger       zerolog.Logger
	temperatures Gateway
	speeds       Gateway
}

func NewModule() *Module {
	return &Module{
		logger:       log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
		temperatures: NewTemperatureGateway(NewHttpClient()),
		speeds:       NewWindspeedGateway(NewHttpClient()),
	}
}

func (m *Module) RegisterRoutes(e *echo.Echo) {
	e.GET("/temperatures", m.GetTemperature)
	e.GET("/speeds", m.GetSpeed)
	e.GET("/weather", m.GetWeather)
}

func (m *Module) GetTemperature(c echo.Context) error {
	startDate, endDate, err := getStartdAndEndDateFromRequest(c)
	if err != nil {
		m.logger.Error().Msg(err.Type + " " + err.Message)
		return c.JSON(http.StatusBadRequest, err)
	}

	var wg sync.WaitGroup
	ch := make(chan Temperature)
	errc := make(chan HttpError)
	defer close(ch)
	defer close(errc)

	for !startDate.After(endDate) {
		wg.Add(1)
		go func(startDate time.Time) {
			var temp Temperature
			err := m.temperatures.GetResourceAt(startDate.Format("2006-01-02T15:04:05Z"), &temp)
			if err != nil {
				m.logger.Error().Msg(err.Type + " " + err.Message)
				errc <- *err
				return
			}
			ch <- temp
		}(startDate)
		startDate = startDate.Add(time.Hour * 24)
	}

	var temperatures []Temperature
	go func() {
		for temp := range ch {
			temperatures = append(temperatures, temp)
			wg.Done()
		}
	}()

	var httpErros []HttpError
	go func() {
		for httpError := range errc {
			httpErros = append(httpErros, httpError)
			wg.Done()
		}
	}()
	wg.Wait()

	if len(httpErros) > 0 {
		return c.JSON(http.StatusInternalServerError, httpErros[0])
	}

	sort.Slice(temperatures, func(i, j int) bool {
		return temperatures[i].Date < temperatures[j].Date
	})
	return c.JSON(http.StatusOK, temperatures)
}

func (m *Module) GetSpeed(c echo.Context) error {
	startDate, endDate, err := getStartdAndEndDateFromRequest(c)
	if err != nil {
		m.logger.Error().Msg(err.Type + " " + err.Message)
		return c.JSON(http.StatusBadRequest, err)
	}

	var wg sync.WaitGroup
	ch := make(chan Windspeed)
	errc := make(chan HttpError)
	defer close(ch)
	defer close(errc)

	for !startDate.After(endDate) {
		wg.Add(1)
		go func(startDate time.Time) {
			var speed Windspeed
			err := m.speeds.GetResourceAt(startDate.Format("2006-01-02T15:04:05Z"), &speed)
			if err != nil {
				m.logger.Error().Msg(err.Type + " " + err.Message)
				errc <- *err
				return
			}
			ch <- speed
		}(startDate)
		startDate = startDate.Add(time.Hour * 24)
	}

	var speeds []Windspeed
	go func() {
		for speed := range ch {
			speeds = append(speeds, speed)
			wg.Done()
		}
	}()

	var httpErros []HttpError
	go func() {
		for httpError := range errc {
			httpErros = append(httpErros, httpError)
			wg.Done()
		}
	}()
	wg.Wait()

	if len(httpErros) > 0 {
		return c.JSON(http.StatusInternalServerError, httpErros[0])
	}

	sort.Slice(speeds, func(i, j int) bool {
		return speeds[i].Date < speeds[j].Date
	})
	return c.JSON(http.StatusOK, speeds)
}

func (m *Module) GetWeather(c echo.Context) error {
	return nil
}

func getStartdAndEndDateFromRequest(c echo.Context) (time.Time, time.Time, *HttpError) {
	start := c.QueryParam("start")
	end := c.QueryParam("end")
	if start == "" || end == "" {
		err := HttpErrorBuilder().From(http.StatusText(http.StatusBadRequest), "Please provide both start and end dates")
		return time.Now(), time.Now(), err
	}

	startDate, startErr := time.Parse("2006-01-02T15:04:05Z", start)
	endDate, err := time.Parse("2006-01-02T15:04:05Z", end)
	if startErr != nil || err != nil {
		err := HttpErrorBuilder().From(http.StatusText(http.StatusBadRequest), "Please provide dates with format ISO8601 DateTime (eg. 2018-08-12T12:00:00Z)")
		return time.Now(), time.Now(), err
	}
	startDate = startDate.Truncate(24 * time.Hour)
	endDate = endDate.Truncate(24 * time.Hour)
	return startDate, endDate, nil
}
