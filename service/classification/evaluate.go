// Package classification provides the core logic for evaluating weather risk levels
// for individual hours based on configurable thresholds and weights.
package classification

import (
	"fmt"

	"github.com/ihgazi/EventWeatherGuard/model"
)

// EvaluateHourlyRisk assesses the weather risk for a single hourly forecast.
// It returns an HourlyEvaluation containing the risk level, reason, and severity score.
func EvaluateHourlyRisk(
	h model.HourlyForecast,
	t SeverityThresholds,
	w SeverityWeights,
) HourlyEvaluation {
	if h.Precipitation >= t.UnsafeRainMM ||
		h.WindKmh >= t.UnsafeWindKmh ||
		h.Weather == "Thunderstorm" {

		return HourlyEvaluation{
			Level:    Unsafe,
			Reason:   buildReason(h, t),
			Severity: computeSeverity(h, w, t),
		}
	}

	if h.Precipitation >= t.RiskyRainMM ||
		h.WindKmh >= t.RiskyWindKmh ||
		h.RainProb >= t.RiskyRainProb ||
		h.Weather == "Heavy Rain" {

		return HourlyEvaluation{
			Level:    Risky,
			Reason:   buildReason(h, t),
			Severity: computeSeverity(h, w, t),
		}
	}

	return HourlyEvaluation{
		Level:    Safe,
		Reason:   "Favorable weather conditions, Have a great day!",
		Severity: computeSeverity(h, w, t),
	}
}

// buildReason generates a human-readable explanation for the assigned risk level
// based on the hourly weather data and thresholds.
func buildReason(h model.HourlyForecast, t SeverityThresholds) string {
	switch {
	case h.Weather == "Thunderstorm":
		return fmt.Sprintf("Thunderstorm predicted at %s", h.Time.Format("12:00"))
	case h.Precipitation >= t.RiskyRainMM:
		return fmt.Sprintf(
			"Expected %.1f mm of rain with probability %d%% at %s",
			h.Precipitation,
			h.RainProb,
			h.Time.Format("15:04"),
		)
	case h.RainProb >= t.RiskyRainProb:
		return fmt.Sprintf(
			"High chance of rain (%d%%) at %s",
			h.RainProb,
			h.Time.Format("15:04"),
		)
	case h.WindKmh >= t.RiskyWindKmh:
		return fmt.Sprintf(
			"Expected wind speed of %.1f km/h at %s",
			h.WindKmh,
			h.Time.Format("15:04"),
		)
	default:
		return "No significant wind or rain expected."
	}
}

// computeSeverity calculates a normalized severity score (0.0–1.0) for the hour
// based on precipitation, wind, and rain probability, weighted by the provided configuration.
func computeSeverity(
	h model.HourlyForecast,
	w SeverityWeights,
	t SeverityThresholds) float64 {
	sr := min(1.0, h.Precipitation/t.UnsafeRainMM)
	sw := min(1.0, h.WindKmh/t.UnsafeWindKmh)
	sp := min(1.0, float64(h.RainProb)/100.0)

	score := w.RainMM*sr + w.Wind*sw + w.RainProb*sp
	score = max(score, wmoCap(h, w))

	return min(1.0, score)
}

// wmoCap returns a minimum severity score based on the WMO weather code label.
func wmoCap(h model.HourlyForecast, w SeverityWeights) float64 {
	switch h.Weather {
	case "Thunderstorm":
		return w.Storm[3]
	case "Heavy Rain":
		return w.Storm[2]
	case "Rain Showers":
		return w.Storm[1]
	default:
		return w.Storm[0]
	}
}
