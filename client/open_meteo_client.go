package client

import (
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

func (c *OpenMeteoClient) FetchWeatherData(lat, long float64) (*model.OpenMeteoResponse, error) {
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&hourly=precipitation_probability,rain,wind_speed_10m,weather_code&timezone=UTC",
		lat, long,
	)

	resp, err := c.httpClient.Get(url)
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
