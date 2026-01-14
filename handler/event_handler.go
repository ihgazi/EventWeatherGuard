// Packge handler provides HTTP handlers for the EventWeatherGuard API.
// This file defines the handler for the event weather forecast requests.
package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/ihgazi/EventWeatherGuard/client"
	"github.com/ihgazi/EventWeatherGuard/logger"
	"github.com/ihgazi/EventWeatherGuard/model"
	"github.com/ihgazi/EventWeatherGuard/service"
)

// EventForecastHandler handles POST requests for event weather forecasts.
//
// @Summary      Get event weather forecast and risk classification
// @Description  Returns weather risk assessment for a given event location and time window. Optionally fetches alternate time windows, in case current window is Unsafe or Risky.
// @Tags         event
// @Accept       json
// @Produce      json
// @Param        request  body      model.EventForecastRequest  true  "Event forecast request"
// @Success      200      {object}  model.EventForecastResponse
// @Failure      400      {object}  map[string]string
// @Failure 	 404 	  {object}  map[string]string
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
			"error": "Invalid event timings: Event duration must be positive and lie within next 6 days.",
		})
		return
	}

	// Validate location of event
	if err := validateLocation(req.Location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Handle response in case of no forecast received
	if len(forecast) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data unavailable for the give forecast duration"})
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

	if req.ListAlters && response.Classification != "Safe" {
		response.AlternateWindows = alternateWindows(ctx, req, weatherSvc)
	}

	c.JSON(http.StatusOK, response)
}

func validateEventTime(start, end time.Time) bool {
	now := time.Now().UTC()

	if !start.After(now) || !end.After(start) {
		return false
	}

	// Default limit of 6 days from current time
	maxAllowed := now.Add(6 * 24 * time.Hour)
	if end.After(maxAllowed) {
		return false
	}

	return true
}

func validateLocation(l model.Location) error {
	if l.Latitude < -90 || l.Latitude > 90 || l.Longitude < -180 || l.Longitude > 180 {
		return errors.New("Invalid Coordinates: latitude must be in [-90, 90] and longitude must be in [-180, 180].")
	}

	return nil
}

// alternateWindows suggests alternate time windows for an event with optimal weather conditions.
//
// Given an event request, this function analyzes the weather forecast for the next 24 hours
// and returns up to three alternate time slots that best match the event's duration and weather suitability.
func alternateWindows(
	ctx context.Context,
	req model.EventForecastRequest,
	weatherSvc *service.WeatherService,
) []model.EventWindow {
	eventHours := int(req.EndTime.Sub(req.StartTime.Time).Hours())

	winStart := time.Now().UTC()
	winEnd := winStart.Add(24 * time.Hour)
	oneDayForecast, err := weatherSvc.GetEventForecast(ctx, req.Location.Latitude, req.Location.Longitude, winStart, winEnd)
	if err != nil {
		logger.Log.Error("Failed to fetch alternate times: ", zap.Error(err))
	}

	alternates := service.FindTopKWindows(
		oneDayForecast,
		eventHours,
		3, // Fetch best 3 possible alternate timings
	)

	return alternates
}
