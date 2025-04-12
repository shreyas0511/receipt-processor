package handlers

import (
	"net/http"
	"receipt-processor/models"
	"receipt-processor/services"
	"receipt-processor/utils"
	"receipt-processor/validators"

	"github.com/gin-gonic/gin"
)

// Process the POST JSON request
func ProcessReceipt(context *gin.Context) {
	var reciept models.Receipt
	err := context.ShouldBindJSON(&reciept)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "The receipt is invalid."})
		// context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Validate json data
	err = validators.ValidateReceipt(reciept)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "The receipt is invalid."})
		// context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate unique id for the receipt
	var id models.ID
	id.ID = utils.GenerateUniqueId()

	// Store id and receipt in a map
	id.Receipt = reciept
	models.Receipts[id.ID] = id.Receipt

	// 200 response, return the id
	context.JSON(http.StatusOK, gin.H{"id": id.ID})
}

// Process the GET request and calculate points
func GetPointsFromID(context *gin.Context) {
	// extract id from the get request header
	id := context.Param("id")

	// check if points have alredy been calculated for the given id
	totalPoints, ok := models.PointsForId[id]

	if ok {
		context.JSON(http.StatusOK, gin.H{"points": totalPoints.Points})
		return
	}

	// get the receipt for the given id if id is valid
	receipt, exists := models.Receipts[id]

	if !exists {
		context.JSON(http.StatusNotFound, gin.H{"message": "No receipt found for that ID."})
		return
	}

	// calculate points for the corresponding receipt
	calculatedPoints, err := services.CalculatePoints(receipt)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "No receipt found for that ID."})
		return
	}

	// store the calculated points in the Points struct
	var points models.Points
	points.Points = calculatedPoints

	// Store id and points in a map
	models.PointsForId[id] = points

	// 200 response, return the calculated points
	context.JSON(http.StatusOK, gin.H{"points": points.Points})
}
