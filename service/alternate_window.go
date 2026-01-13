package service

import (
	"sort"
	"time"

	"github.com/ihgazi/EventWeatherGuard/model"
)

// FindTopKWindows returns the top K time windows with the most suitable weather conditions.
func FindTopKWindows(
	hourly []model.HourlyForecast,
	eventDuration int,
	k int,
) []model.EventWindow {
	if len(hourly) < eventDuration || k <= 0 {
		return nil
	}

	limit := min(len(hourly), 24)
	candidates := []model.EventWindow{}

	for i := 0; i+eventDuration <= limit; i++ {
		window := hourly[i : i+eventDuration]

		result := ClassifyEvent(window)

		// Ignore Unsafe / Risky time windows
		if result.Severity >= 50 {
			continue
		}

		candidates = append(candidates, model.EventWindow{
			StartTime: window[0].Time,
			EndTime:   window[len(window)-1].Time.Add(time.Hour),
			Score:     result.Severity,
		})
	}

	// Sort by severity, then by earliest time
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Score < candidates[j].Score {
			return candidates[i].Score < candidates[j].Score
		}
		return candidates[i].StartTime.Before(candidates[j].StartTime)
	})

	if k > len(candidates) {
		k = len(candidates)
	}

	return candidates[:k]
}
