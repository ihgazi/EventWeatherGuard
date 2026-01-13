package service

import (
	"context"
	"time"

	"github.com/relvacode/iso8601"
	"go.uber.org/zap"

	"github.com/ihgazi/EventWeatherGuard/client"
	"github.com/ihgazi/EventWeatherGuard/logger"
	"github.com/ihgazi/EventWeatherGuard/model"
)

// WeatherService provides methods to fetch and process weather forecasts
// using an external API client.
type WeatherService struct {
	client *client.OpenMeteoClient
}

func NewWeatherService(client *client.OpenMeteoClient) *WeatherService {
	return &WeatherService{
		client: client,
	}
}

// GetEventForecast retrieves and processes hourly weather forecasts for a given event location and time window.
// It filters the forecast data to only include hours within the specified start and end times.
func (s *WeatherService) GetEventForecast(
	ctx context.Context,
	lat, long float64,
	start, end time.Time,
) ([]model.HourlyForecast, error) {

	raw, err := s.client.FetchWeatherData(ctx, lat, long)
	if err != nil {
		return nil, err
	}

	var result []model.HourlyForecast

	for i, t := range raw.Hourly.Time {
		parsed, err := iso8601.ParseString(t)
		if err != nil {
			logger.Log.Error("Failed to parse time: %v", zap.Error(err))
			continue
		}

		// Time interval falls outside window
		if parsed.Before(start) || parsed.After(end) || parsed.Equal(end) {
			continue
		}

		result = append(result, model.HourlyForecast{
			Time:          parsed,
			RainProb:      raw.Hourly.PrecipitationProbability[i],
			Precipitation: raw.Hourly.Rain[i],
			WindKmh:       raw.Hourly.WindSpeed10m[i],
			Weather:       weatherCodeToLabel(raw.Hourly.WeatherCode[i]),
		})
	}

	return result, nil
}

// Map WMO weather codes to human-readable labels
func weatherCodeToLabel(code int) string {
	switch {
	case code >= 95:
		return "Thunderstorm"
	case code >= 80:
		return "Heavy Rain"
	case code >= 60:
		return "Rain Showers"
	default:
		return "Clear"
	}
}
