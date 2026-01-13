package service

import (
	"fmt"

	"github.com/ihgazi/EventWeatherGuard/model"
	cls "github.com/ihgazi/EventWeatherGuard/service/classification"
)

type ClassificationResult struct {
	Classification cls.RiskLevel
	Reason         []string
	Summary        string
	Severity       int
}

func ClassifyEvent(hours []model.HourlyForecast) ClassificationResult {
	finalLevel := cls.Safe
	var reasons []string
	maxSeverity := 0.0
	var peakReport cls.HourlyEvaluation

	// Generate hourly evaluations and find the worst case
	// Reasons are aggregated over all hourly windows
	for _, h := range hours {
		eval := cls.EvaluateHourlyRisk(h, cls.DefaultThresholds, cls.DefaultWeights)

		if eval.Level == cls.Unsafe {
			finalLevel = cls.Unsafe
		} else if eval.Level == cls.Risky && finalLevel != cls.Unsafe {
			finalLevel = cls.Risky
		}

		if eval.Level != cls.Safe {
			reasons = append(reasons, eval.Reason)
		}

		if eval.Severity > maxSeverity {
			maxSeverity = eval.Severity
			peakReport = eval
		}
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "No significant wind or rain expected.")
	}

	fmt.Printf("Peak Report: %+v\n", peakReport)
	return ClassificationResult{
		Classification: finalLevel,
		Reason:         reasons,
		Summary:        buildSummary(peakReport),
		Severity:       int(maxSeverity * 100),
	}
}

func buildSummary(peakReport cls.HourlyEvaluation) string {
	switch peakReport.Level {
	case cls.Safe:
		return "Weather conditions are safe throughout the event."
	case cls.Unsafe:
		return "Severe thunderstorms are expected during the event."
	default:
		return "Moderate rainfall and winds are expected during the event."
	}
}
