package service

type RiskLevel string

// Risk levels for weather classification
const (
	Safe   RiskLevel = "Safe"
	Risky  RiskLevel = "Risky"
	Unsafe RiskLevel = "Unsafe"
)

// Report for hourly weather risk evaluation
type HourlyEvaluation struct {
	Level    RiskLevel
	Reason   string
	Severity float64
}

// thresholds for classifying weather conditions
// Build new thresholds to adjust classification sensitivity
type thresholds struct {
	UnsafeRainMM  float64
	UnsafeWindKmh float64
	RiskyRainMM   float64
	RiskyWindKmh  float64
	RiskyRainProb int
}

var defaultThresholds = thresholds{
	UnsafeRainMM:  5.0,
	UnsafeWindKmh: 40.0,
	RiskyRainMM:   1.0,
	RiskyWindKmh:  30.0,
	RiskyRainProb: 60,
}

// Weights assigned to different weather factors for severity calculation
type severityWeights struct {
	Storm float64
	Rain  float64
	Wind  float64
}

var defaultWeights = severityWeights{
	Storm: 0.5,
	Rain:  0.3,
	Wind:  0.2,
}
