package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type ForecastResponse struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	CurrentUnits         struct {
		Time          string `json:"time"`
		Interval      string `json:"interval"`
		Temperature2M string `json:"temperature_2m"`
	} `json:"current_units"`
	Current struct {
		Time          string  `json:"time"`
		Interval      int     `json:"interval"`
		Temperature2M float64 `json:"temperature_2m"`
	} `json:"current"`
}

type GeocodingResponse struct {
	Results []CityInformation `json:"results"`
}

type CityInformation struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (c *ApiClient) GetCityInformation(cityName string) (geocodingResponse *GeocodingResponse, err error) {
	params := url.Values{}
	params.Add("name", cityName)
	params.Add("count", "1")

	var result GeocodingResponse

	err = c.fetch(c.LocationEndpoint, "/v1/search", params, &result)

	if err != nil {
		return nil, err
	}

	return &result, err
}

func (c *ApiClient) GetWeatherInformationForCity(longitude float64, latitude float64) (forecastResponse *ForecastResponse, err error) {
	params := url.Values{}
	params.Add("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	params.Add("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	params.Add("current", "temperature_2m")
	params.Add("timezone", "auto")

	var result ForecastResponse

	err = c.fetch(c.ForecastEndpoint, "/v1/forecast", params, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *ApiClient) fetch(baseUri string, endpoint string, params url.Values, target any) error {
	response, err := c.get(baseUri, endpoint, params)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error: StatusCode: %d", response.StatusCode)
	}

	return decodeResult(response.Body, target)
}

func decodeResult(responseBody io.Reader, target any) error {
	err := json.NewDecoder(responseBody).Decode(target)

	if err != nil {
		return fmt.Errorf("JSON could not be decoded %w", err)
	}

	return nil
}
