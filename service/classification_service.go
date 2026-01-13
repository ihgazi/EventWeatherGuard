package service

import (
	"fmt"

	"github.com/ihgazi/EventWeatherGuard/model"
)

type ClassificationResult struct {
	Classification RiskLevel
	Reason         []string
	Severity       int
}

func ClassifyEvent(hours []model.HourlyForecast) ClassificationResult {
	finalLevel := Safe
	var reasons []string
	maxSeverity := 0.0

	for _, h := range hours {
		eval := evaluateHourlyRisk(h, defaultThresholds, defaultWeights)

		if eval.Level == Unsafe {
			finalLevel = Unsafe
		} else if eval.Level == Risky && finalLevel != Unsafe {
			finalLevel = Risky
		}

		if eval.Level != Safe {
			reasons = append(reasons, eval.Reason)
		}

		if eval.Severity > maxSeverity {
			maxSeverity = eval.Severity
		}
	}

	return ClassificationResult{
		Classification: finalLevel,
		Reason:         reasons,
		Severity:       int(maxSeverity * 100),
	}
}

func evaluateHourlyRisk(
	h model.HourlyForecast,
	t thresholds,
	w severityWeights,
) HourlyEvaluation {
	if h.Precipitation >= t.UnsafeRainMM ||
		h.WindKmh >= t.UnsafeWindKmh ||
		h.Weather == "Thunderstorm" {

		return HourlyEvaluation{
			Level:    Unsafe,
			Reason:   buildReason(h),
			Severity: computeSeverity(h, w),
		}
	}

	if h.Precipitation >= t.RiskyRainMM ||
		h.WindKmh >= t.RiskyWindKmh ||
		h.RainProb >= t.RiskyRainProb {

		return HourlyEvaluation{
			Level:    Risky,
			Reason:   buildReason(h),
			Severity: computeSeverity(h, w),
		}
	}

	return HourlyEvaluation{
		Level:    Safe,
		Reason:   "Favorable weather conditions, Have a great day!",
		Severity: computeSeverity(h, w),
	}
}

func buildReason(h model.HourlyForecast) string {
	switch {
	case h.Weather == "Thunderstorm":
		return fmt.Sprintf("Thunderstorm predicted at %s", h.Time.Format("12:00"))
	case h.Precipitation > 0:
		return fmt.Sprintf(
			"Expected %.1f mm of rain with probability %d%% at %s",
			h.Precipitation,
			h.RainProb,
			h.Time.Format("15:04"),
		)
	case h.WindKmh > 0:
		return fmt.Sprintf(
			"Expected wind speed of %.1f km/h at %s",
			h.WindKmh,
			h.Time.Format("15:04"),
		)
	default:
		return "Favorable weather conditions, Have a great day!"
	}
}

func computeSeverity(h model.HourlyForecast, w severityWeights) float64 {
	sr := min(1.0, h.Precipitation/10.0)
	sw := min(1.0, h.WindKmh/50.0)
	ss := 0.0

	if h.Weather == "Thunderstorm" {
		ss = 1.0
	}

	score := w.Rain*sr + w.Wind*sw + w.Storm*ss
	return min(1.0, score)
}
