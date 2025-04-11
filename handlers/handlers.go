package handlers

import (
	"net/http"
	"receipt-processor/models"
	"receipt-processor/services"
	"receipt-processor/utils"
	"receipt-processor/validators"

	"github.com/gin-gonic/gin"
)

func ProcessReceipt(context *gin.Context) {
	var reciept models.Receipt
	err := context.ShouldBindJSON(&reciept)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid receipt"})
		return
	}

	// Validate json data
	err = validators.ValidateReceipt(reciept)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate unique id for the receipt
	var id models.ID
	id.ID = utils.GenerateUniqueId()
	id.Receipt = reciept

	models.Receipts[id.ID] = id.Receipt
	context.JSON(http.StatusOK, gin.H{"id": id.ID})
}

func GetPointsFromID(context *gin.Context) {
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
	calculatedPoints, err := services.CalculatePoints(receipt)
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
