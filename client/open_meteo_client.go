// package client provides API clients for external weather services.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ihgazi/EventWeatherGuard/model"
)

type OpenMeteoClient struct {
	httpClient *http.Client
}

func NewOpenMeteoClient() *OpenMeteoClient {
	return &OpenMeteoClient{
		httpClient: &http.Client{Timeout: 1 * time.Minute},
	}
}

// FetchWeatherData retrieves weather forecast data from the Open-Meteo API
// for the specified latitude and longitude. It returns a parsed OpenMeteoResponse
// or an error if the request fails.
func (c *OpenMeteoClient) FetchWeatherData(ctx context.Context, lat, long float64) (*model.OpenMeteoResponse, error) {
	// Open-Meteo API URL with required parameters
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&hourly=precipitation_probability,rain,wind_speed_10m,weather_code&timezone=UTC",
		lat, long,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("open-meteo returned status %d", resp.StatusCode)
	}

	var data model.OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
