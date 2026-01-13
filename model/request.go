package model

import (
	"github.com/relvacode/iso8601"
)

type EventForecastRequest struct {
	Name      string        `json:"name" binding:"required"`
	Location  Location      `json:"location" binding:"required"`
	StartTime *iso8601.Time `json:"start_time" binding:"required"`
	EndTime   *iso8601.Time `json:"end_time" binding:"required"`
}

type Location struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}
