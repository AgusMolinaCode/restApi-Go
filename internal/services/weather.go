package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AgusMolinaCode/restApi-Go.git/internal/models"
)

func GetWeather(lat, lon float64, date string) (*models.WeatherResponse, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&units=metric&appid=%s", lat, lon, apiKey)

	// Convertir la fecha del formato "DD/MM/YYYY" a "YYYY-MM-DD"
	parsedDate, err := time.Parse("02/01/2006", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}
	formattedDate := parsedDate.Format("2006-01-02")

	resp, err := http.Get(fmt.Sprintf("%s&date=%s", url, formattedDate))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var forecast struct {
		City struct {
			Name string `json:"name"`
		} `json:"city"`
		List []struct {
			Dt   int64 `json:"dt"`
			Main struct {
				Temp float64 `json:"temp"`
			} `json:"main"`
			Weather []struct {
				Main        string `json:"main"`
				Description string `json:"description"`
			} `json:"weather"`
		} `json:"list"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return nil, err
	}

	var closestForecast *models.WeatherResponse
	minDiff := int64(1<<63 - 1)
	for _, entry := range forecast.List {
		diff := abs(entry.Dt - parsedDate.Unix())
		if diff < minDiff {
			minDiff = diff
			closestForecast = &models.WeatherResponse{
				City: struct {
					Name string `json:"name"`
				}{Name: forecast.City.Name},
				Main: struct {
					Temp float64 `json:"temp"`
				}{Temp: entry.Main.Temp},
				Weather: []struct {
					Main        string `json:"main"`
					Description string `json:"description"`
				}{{Main: entry.Weather[0].Main, Description: entry.Weather[0].Description}},
				Link: fmt.Sprintf("https://openweathermap.org/weathermap?basemap=map&cities=true&layer=temperature&lat=%f&lon=%f&zoom=10", lat, lon),
			}
		}
	}

	return closestForecast, nil
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
