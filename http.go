package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"time"
)

type HttpClient struct {
	client *http.Client
}

type HttpError struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

func HttpErrorBuilder() *HttpError {
	return &HttpError{}
}

func (e *HttpError) From(errorType string, message string) *HttpError {
	e.Type = errorType
	e.Message = message
	return e
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout:   time.Second * 10,
			Transport: &http.Transport{MaxConnsPerHost: 50},
		},
	}
}

func NewHttpClientForTesting(handler http.Handler) *http.Client {
	s := httptest.NewServer(handler)
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}
}

func (c *HttpClient) MakeRequest(method string, url string) ([]byte, *HttpError) {
	response, err := c.client.Get(url)
	if err != nil {
		return nil, HttpErrorBuilder().From(http.StatusText(http.StatusInternalServerError), err.Error())
	}
	defer response.Body.Close()

	httpError := validateResponseStatus(response)
	if httpError != nil {
		return nil, httpError
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, HttpErrorBuilder().From(http.StatusText(http.StatusInternalServerError), err.Error())
	}

	return responseBody, nil
}

func validateResponseStatus(response *http.Response) *HttpError {
	if response.StatusCode >= 400 {
		var httpError HttpError
		body, err := ioutil.ReadAll(response.Body)
		_ = json.Unmarshal(body, &httpError)
		if err != nil {
			return HttpErrorBuilder().From(http.StatusText(response.StatusCode), "Something went wrong...")
		}
		return HttpErrorBuilder().From(http.StatusText(response.StatusCode), httpError.Message)
	}
	return nil
}
