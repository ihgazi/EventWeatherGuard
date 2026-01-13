package model

import "time"

// HourlyForecast represents weather data for a single hour.
//
// swagger:model HourlyForecast
type HourlyForecast struct {
	Time          time.Time `json:"time"`
	RainProb      int       `json:"rain_prob"`
	Precipitation float64   `json:"precip_mm"`
	WindKmh       float64   `json:"wind_kmh"`
	Weather       string    `json:"weather"`
}
