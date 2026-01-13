package model

import "time"

type HourlyForecast struct {
	Time     time.Time `json:"time"`
	RainProb int       `json:"rain_probability"`
	WindKmh  float64   `json:"wind_kmh"`
	Weather  string    `json:"weather"`
}
