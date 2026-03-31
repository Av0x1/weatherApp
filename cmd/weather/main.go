package main

import (
	"fmt"
	"os"
	"time"
	"weatherApp/internal/weather"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Please input city name.")
		return
	}

	apiClient := weather.NewApiClient()

	geocodingResponse, err := apiClient.GetCityInformation(args[1])

	if err != nil {
		fmt.Printf("an error occurred: %v \n", err)
		return
	}

	if len(geocodingResponse.Results) == 0 {
		fmt.Printf("No city was found.")
		return
	}

	cityInformation := geocodingResponse.Results[0]

	fmt.Printf("Stadt: %s (ID: %d). Längengrad: %f. Breitengrad: %f \n", cityInformation.Name, cityInformation.Id, cityInformation.Longitude, cityInformation.Latitude)

	forecastResponse, err := apiClient.GetWeatherInformationForCity(cityInformation.Longitude, cityInformation.Latitude)

	if err != nil {
		fmt.Printf("an error occurred: %s", err.Error())
		return
	}

	t, err := time.Parse("2006-01-02T15:04", forecastResponse.Current.Time)

	if err != nil {
		fmt.Println("an error occurred while trying to format datetime", err)
	}

	fmt.Printf("Stadt: %s. Temperatur: %.1f C°. Zeit: %s",
		cityInformation.Name, forecastResponse.Current.Temperature2M, t.Format("02.01.2006 15:04"))
}
