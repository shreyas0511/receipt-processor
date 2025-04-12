package services

import (
	"fmt"
	"math"
	"receipt-processor/models"
	"receipt-processor/validators"
	"strconv"
	"strings"
	"unicode"
)

// Calculates the total points for a receipt based on the rules
func CalculatePoints(r models.Receipt) (int64, error) {
	totalPoints := 0

	fmt.Println("Breakdown:")

	// One point for every alphanumeric character in the retailer name.
	count := 0
	for _, ch := range r.Retailer {
		if unicode.IsDigit(ch) || unicode.IsLetter(ch) {
			count++
		}
	}
	totalPoints += count

	// 50 points if the total is a round dollar amount with no cents.
	total := r.Total
	if total[len(total)-1] == '0' && total[len(total)-2] == '0' {
		totalPoints += 50
	}

	// 25 points if the total is a multiple of 0.25.
	parsedTotal, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return 0, err
	}
	if math.Mod(parsedTotal, 0.25) == 0 {
		totalPoints += 25
	}

	// 5 points for every two items on the receipt.
	totalItems := len(r.Items)
	totalPoints += (totalItems / 2) * 5

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	count = 0
	for _, item := range r.Items {
		trimmedItem := strings.TrimSpace(item.ShortDescription)
		trimmedLength := len(trimmedItem)
		price, err := strconv.ParseFloat(item.Price, 64)

		if err != nil {
			return 0, err
		}

		if trimmedLength%3 == 0 {
			count += int(math.Ceil(price * 0.2))
		}
	}
	totalPoints += count

	// 6 points if the day in the purchase date is odd.
	purchaseDate, err := validators.ValidateDate(r.PurchaseDate)
	if err != nil {
		return 0, err
	}
	if purchaseDate.Day()%2 == 1 {
		totalPoints += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime, err := validators.ValidateTime(r.PurchaseTime)
	if err != nil {
		return 0, err
	}
	if (purchaseTime.Hour() >= 14 && purchaseTime.Hour() <= 15) && (purchaseTime.Minute() >= 1 && purchaseTime.Minute() <= 59) {
		totalPoints += 10
	}

	return int64(totalPoints), nil
}
