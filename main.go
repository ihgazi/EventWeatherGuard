// @title           EventWeatherGuard API
// @version         1.0
// @description     API for weather risk assessment for outdoor events.
// @host      localhost:8080
// @BasePath  /
package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "github.com/ihgazi/EventWeatherGuard/docs"
	"github.com/ihgazi/EventWeatherGuard/handler"
	"github.com/ihgazi/EventWeatherGuard/logger"
)

// Gin engine used to define HTTP routes and middleware
var router *gin.Engine

func main() {
	// Initialize Gin router
	router = gin.Default()

	// Initialize logger
	logger.Init()
	defer logger.Log.Sync()

	// Setup API routes
	api := router.Group("/")
	{
		api.POST("/event-forecast", handler.EventForecastHandler)
	}
	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Configuring and start HTTP server
	port := 8080
	address := fmt.Sprintf(":%d", port)
	if err := router.Run(address); err != nil {
		logger.Log.Error("Failed to run server: %v", zap.Error(err))
	}
}
