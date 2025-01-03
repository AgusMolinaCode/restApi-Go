package models

type WeatherResponse struct {
	City struct {
		Name string `json:"name"`
	} `json:"city"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Link string `json:"link"`
}
