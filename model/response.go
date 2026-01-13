package model

type EventForecastResponse struct {
	Classification string           `json:"classification"`
	Summary        string           `json:"summary"`
	Reasons        []string         `json:"reasons"`
	ForecastWindow []HourlyForecast `json:"forecast_window"`
}
