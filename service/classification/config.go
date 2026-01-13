package classification

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

// severityThresholds for classifying weather conditions
// Build new severityThresholds to adjust classification sensitivity
type SeverityThresholds struct {
	UnsafeRainMM  float64
	UnsafeWindKmh float64
	RiskyRainMM   float64
	RiskyWindKmh  float64
	RiskyRainProb int
}

var DefaultThresholds = SeverityThresholds{
	UnsafeRainMM:  10.0,
	UnsafeWindKmh: 40.0,
	RiskyRainMM:   2.5,
	RiskyWindKmh:  30.0,
	RiskyRainProb: 40,
}

// Weights assigned to different weather factors for severity calculation
type SeverityWeights struct {
	Storm    []float64 // Different weights for varying WMO weather codes
	RainMM   float64
	RainProb float64
	Wind     float64
}

var DefaultWeights = SeverityWeights{
	Storm:    []float64{0, 0.25, 0.5, 1.0},
	RainMM:   0.2,
	RainProb: 0.3,
	Wind:     0.5,
}
