package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	// Initialize Gin router
	router = gin.Default()

	// Root GET endpoint
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Event Weather Guard API!")
	})

	// Configuring server
	port := 8080
	address := fmt.Sprintf(":%d", port)
	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	log.Printf("Server running on port %d", port)
}
