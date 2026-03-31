package weather

import (
	"fmt"
	"net/http"
	"net/url"
)

type ApiClient struct {
	HttpClient       *http.Client
	LocationEndpoint string
	ForecastEndpoint string
}

func NewApiClient() *ApiClient {
	return &ApiClient{
		HttpClient:       &http.Client{},
		LocationEndpoint: "https://geocoding-api.open-meteo.com",
		ForecastEndpoint: "https://api.open-meteo.com",
	}
}

func (c *ApiClient) get(baseUri string, endpoint string, params url.Values) (*http.Response, error) {
	requestUri, err := getUriWithParams(baseUri, endpoint, params)

	if err != nil {
		return nil, err
	}

	response, err := c.HttpClient.Get(requestUri)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func getUriWithParams(baseUri string, endpoint string, params url.Values) (string, error) {
	if baseUri == "" || baseUri == " " {
		return "", fmt.Errorf("baseUri must not be empty")
	}

	if endpoint == "" || endpoint == " " {
		return baseUri, nil
	}

	requestUri, err := url.JoinPath(baseUri, endpoint)

	if err != nil {
		return "", err
	}

	if params != nil {
		queryString := params.Encode()
		requestUri = requestUri + "?" + queryString
	}

	return requestUri, nil
}
