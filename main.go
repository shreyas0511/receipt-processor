package main

import (
	"receipt-processor/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the server
	server := gin.Default()

	// Register the required routes
	routes.RegisterRoutes(server)

	// Run the server
	server.Run(":8080")
}
