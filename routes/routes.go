package routes

import (
	"receipt-processor/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// Register the POST route
	server.POST("/receipts/process", handlers.ProcessReceipt)

	// Register the GET route
	server.GET("/receipts/:id/points", handlers.GetPointsFromID)
}
