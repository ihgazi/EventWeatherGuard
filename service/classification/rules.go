package classification

import (
	"fmt"

	"github.com/ihgazi/EventWeatherGuard/model"
)

// Rule defines the criteria for classifying weather severity levels
//
// Each rule specifies thresholds and conditions for weather parameters (e.g., precipitation,
// wind speed) and is associated with specific event types.
type RiskRule struct {
	ID          string
	Level       RiskLevel
	Matches     func(h model.HourlyForecast, t SeverityThresholds) bool
	Description func(h model.HourlyForecast) string
}

var ClassificationRules = []RiskRule{
	{
		ID:    "UNSAFE_THUNDERSTORM",
		Level: Unsafe,
		Matches: func(h model.HourlyForecast, t SeverityThresholds) bool {
			return h.Weather == "Thunderstorm"
		},
		Description: func(h model.HourlyForecast) string {
			return fmt.Sprintf("Thunderstorm predicted at %s", h.Time.Format("15:04"))
		},
	},
	{
		ID:    "UNSAFE_EXTREME_RAIN_WIND",
		Level: Unsafe,
		Matches: func(h model.HourlyForecast, t SeverityThresholds) bool {
			return h.Precipitation >= t.UnsafeRainMM || h.WindKmh >= t.UnsafeWindKmh
		},
		Description: func(h model.HourlyForecast) string {
			return fmt.Sprintf(
				"Extreme weather: %.1f mm rain and %.1f km/h wind at %s",
				h.Precipitation,
				h.WindKmh,
				h.Time.Format("15:04"),
			)
		},
	},
	{
		ID:    "RISKY_MODERATE_RAIN_WIND",
		Level: Risky,
		Matches: func(h model.HourlyForecast, t SeverityThresholds) bool {
			return h.Precipitation >= t.RiskyRainMM || h.WindKmh >= t.RiskyWindKmh || h.RainProb >= t.RiskyRainProb
		},
		Description: func(h model.HourlyForecast) string {
			return fmt.Sprintf(
				"Moderate risk: %.1f mm rain, %.1f km/h wind, %d%% rain probability at %s",
				h.Precipitation,
				h.WindKmh,
				h.RainProb,
				h.Time.Format("15:04"),
			)
		},
	},
	{
		ID:    "RISKY_HEAVY_RAIN",
		Level: Risky,
		Matches: func(h model.HourlyForecast, t SeverityThresholds) bool {
			return h.Weather == "Heavy Rain"
		},
		Description: func(h model.HourlyForecast) string {
			return fmt.Sprintf("Heavy rain predicted at %s", h.Time.Format("15:04"))
		},
	},
}
