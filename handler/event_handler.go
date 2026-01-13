package handler

import (
	"net/http"

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

	forecast, err := weatherSvc.GetEventForecast(
		req.Location.Latitude,
		req.Location.Longitude,
		req.StartTime.UTC(),
		req.EndTime.UTC(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to get forecast from weather service": err.Error()})
		return
	}

	// TODO: Apply classification rules to get risk level

	c.JSON(http.StatusOK, gin.H{"forecast": forecast})
}
