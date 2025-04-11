package routes

import (
	"receipt-processor/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/receipts/process", handlers.ProcessReceipt)
	server.GET("/receipts/:id/points", handlers.GetPointsFromID)
}
