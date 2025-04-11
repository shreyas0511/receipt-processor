package main

import (
	"net/http"
	"receipt-processor/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	server := gin.Default()
	server.POST("/receipts/process", processReceipt)
	server.GET("/receipts/:id/points", getPointsFromID)
	server.Run(":8080")
}

func processReceipt(context *gin.Context) {
	var reciept models.Receipt
	err := context.ShouldBindJSON(&reciept)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid receipt"})
		return
	}

	// Generate UIUD later
	var id models.ID
	id.ID = uuid.New().String()
	id.Receipt = reciept

	models.Receipts[id.ID] = id.Receipt
	context.JSON(http.StatusOK, gin.H{"id": id.ID})
}

func getPointsFromID(context *gin.Context) {
	// extract id from the get request header
	id := context.Param("id")

	// check if points have alredy been calculated for the particular id
	totalPoints, ok := models.PointsForId[id]

	if ok {
		context.JSON(http.StatusOK, gin.H{"points": totalPoints.Points})
		return
	}

	// get the receipt for the respective id if id is valid
	receipt, exists := models.Receipts[id]

	if !exists {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}

	// calculate points for the respective receipt
	calculatedPoints, err := receipt.CalculatePoints()
	if err != nil {
		context.AbortWithStatus(http.StatusNotFound)
	}

	// store the calculated points in the Points struct
	var points models.Points
	points.Points = calculatedPoints

	// map the id to calculated points to retrieve later
	models.PointsForId[id] = points

	if err != nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}

	context.JSON(http.StatusOK, gin.H{"points": points.Points})
}
