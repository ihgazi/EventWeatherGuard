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

func EventForecastHandler(c *gin.Context) {
	var req model.EventForecastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Add checks to validate whether event time is in 16 day window

	weatherSvc := service.NewWeatherService(
		client.NewOpenMeteoClient(),
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Minute)
	defer cancel()

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
		Summary:        string(result.Classification) + " weather expected.",
		Reasons:        result.Reason,
		Severity:       result.Severity,
		ForecastWindow: forecast,
	}

	c.JSON(http.StatusOK, response)
}
