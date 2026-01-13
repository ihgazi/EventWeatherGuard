package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/ihgazi/EventWeatherGuard/handler"
)

var router *gin.Engine

func main() {
	// Initialize Gin router
	router = gin.Default()

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
