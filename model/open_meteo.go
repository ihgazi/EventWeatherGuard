package model

type OpenMeteoResponse struct {
	Hourly struct {
		Time                     []string  `json:"time"`
		PrecipitationProbability []int     `json:"precipitation_probability"`
		Rain                     []float64 `json:"rain"`
		WindSpeed10m             []float64 `json:"wind_speed_10m"`
		WeatherCode              []int     `json:"weather_code"`
	} `json:"hourly"`
}
