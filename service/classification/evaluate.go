// Package classification provides rule-based event weather classification logic.
//
// This package evaluates weather conditions for events using a configurable set of rules.
// Each rule defines criteria for classifying weather as suitable or unsuitable for an event.
// The evaluation process selects the most appropriate rule based on event context and weather data.
package classification

import (
	"github.com/ihgazi/EventWeatherGuard/model"
)

// EvaluateHourlyRisk assesses the weather risk for a single hourly forecast.
// It returns an HourlyEvaluation containing the risk level, reason, and severity score.
func EvaluateHourlyRisk(
	h model.HourlyForecast,
	t SeverityThresholds,
	w SeverityWeights,
) HourlyEvaluation {

	var selectedRule *RiskRule

	// Find most severe matching rule
	for _, rule := range ClassificationRules {
		if rule.Matches(h, t) {
			if selectedRule == nil ||
				riskPriority(rule.Level) > riskPriority(selectedRule.Level) {
				selectedRule = &rule
			}
		}
	}

	if selectedRule != nil {
		return HourlyEvaluation{
			Level:    selectedRule.Level,
			Reason:   selectedRule.Description(h),
			Severity: computeSeverity(h, w, t),
		}
	}

	return HourlyEvaluation{
		Level:    Safe,
		Reason:   "Favorable weather conditions.",
		Severity: computeSeverity(h, w, t),
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

// Evaluate priority of different risk levels
func riskPriority(level RiskLevel) int {
	switch level {
	case Unsafe:
		return 3
	case Risky:
		return 2
	case Safe:
		return 1
	default:
		return 0
	}
}
