package model

// EventForecastResponse represents the API response for event weather forecast.
//
// swagger:model EventForecastResponse
type EventForecastResponse struct {
	Classification   string           `json:"classification"`
	Severity         int              `json:"severity"`
	Summary          string           `json:"summary"`
	Reasons          []string         `json:"reasons"`
	ForecastWindow   []HourlyForecast `json:"forecast_window"`
	AlternateWindows []EventWindow    `json:"alternate_timings,omitempty"`
}
