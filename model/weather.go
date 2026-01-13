package model

import "time"

type HourlyForecast struct {
	Time          time.Time `json:"time"`
	RainProb      int       `json:"rain_probability"`
	Precipitation float64   `json:"precip_mm"`
	WindKmh       float64   `json:"wind_kmh"`
	Weather       string    `json:"weather"`
}
