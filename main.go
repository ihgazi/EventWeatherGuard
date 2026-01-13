package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/ihgazi/EventWeatherGuard/handler"
	"github.com/ihgazi/EventWeatherGuard/logger"
)

var router *gin.Engine

func main() {
	// Initialize Gin router
	router = gin.Default()

	// Initialize logger
	logger.Init()
	defer logger.Log.Sync()

	api := router.Group("/")
	{
		api.POST("/event-forecast", handler.EventForecastHandler)
	}

	// Configuring server
	port := 8080
	address := fmt.Sprintf(":%d", port)
	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
