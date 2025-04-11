package main

import (
	"net/http"
	"receipt-processor/models"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.POST("/receipts/process", processReceipt)
	server.Run(":8080")
}

func processReceipt(context *gin.Context) {
	var reciept models.Receipt
	err := context.ShouldBindJSON(&reciept)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid receipt"})
	}

	// Generate UIUD later
	var id models.ID
	id.ID = "1"
	id.Receipt = reciept

	models.Receipts[id.ID] = id.Receipt
	context.JSON(http.StatusOK, gin.H{"message": "receipt processed", "id": id.ID})
}
