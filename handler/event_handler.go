// Packge handler provides HTTP handlers for the EventWeatherGuard API.
// This file defines the handler for the event weather forecast requests.
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ihgazi/EventWeatherGuard/client"
	"github.com/ihgazi/EventWeatherGuard/model"
	"github.com/ihgazi/EventWeatherGuard/service"
)

// EventForecastHandler handles POST requests for event weather forecasts.
//
// @Summary      Get event weather forecast and risk classification
// @Description  Returns weather risk assessment for a given event location and time window.
// @Tags         event
// @Accept       json
// @Produce      json
// @Param        request  body      model.EventForecastRequest  true  "Event forecast request"
// @Success      200      {object}  model.EventForecastResponse
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /event-forecast [post]
func EventForecastHandler(c *gin.Context) {
	var req model.EventForecastRequest

	// Bind and validate JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate time window of event
	if !validateEventTime(req.StartTime.Time, req.EndTime.Time) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "event window must be between the next 6 days",
		})
		return
	}

	// Initialize weather service with Open-Meteo client
	weatherSvc := service.NewWeatherService(
		client.NewOpenMeteoClient(),
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Minute)
	defer cancel()

	// Fetch the weather forecast for the event location and time window
	forecast, err := weatherSvc.GetEventForecast(
		ctx,
		req.Location.Latitude,
		req.Location.Longitude,
		req.StartTime.UTC(),
		req.EndTime.UTC(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to get forecast from weather service": err.Error()})
		return
	}

	result := service.ClassifyEvent(forecast)

	response := model.EventForecastResponse{
		Classification: string(result.Classification),
		Summary:        result.Summary,
		Reasons:        result.Reason,
		Severity:       result.Severity,
		ForecastWindow: forecast,
	}

	c.JSON(http.StatusOK, response)
}

func validateEventTime(start, end time.Time) bool {
	now := time.Now().UTC()

	if !start.After(now) || !end.After(start) {
		return false
	}

	maxAllowed := now.Add(6 * 24 * time.Hour)
	if end.After(maxAllowed) {
		return false
	}

	return true
}
