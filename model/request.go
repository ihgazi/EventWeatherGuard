package model

import "time"

type EventForecastRequest struct {
	Name      string    `json:"name" binding:"required"`
	Location  Location  `json:"location" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

type Location struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}
