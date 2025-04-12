package validators

import (
	"errors"
	"receipt-processor/models"
	"regexp"
	"time"
)

// regex for valid price format
var PriceRegEx = regexp.MustCompile(`^(0|[1-9][0-9]*)\.\d{2}$`)

func ValidateReceipt(receipt models.Receipt) error {
	// Retailer name should not be an empty string
	if len(receipt.Retailer) == 0 {
		return errors.New("invalid retailer")
	}

	// Date should be of a valid format
	_, err := ValidateDate(receipt.PurchaseDate)
	if err != nil {
		return errors.New("invalid date")
	}

	// Time should be of a valid fomat
	_, err = ValidateTime(receipt.PurchaseTime)
	if err != nil {
		return errors.New("invalid time")
	}

	// Validate Items
	// Items list should not be empty
	if len(receipt.Items) == 0 {
		return errors.New("empty list of items")
	}

	// Every Item should have a valid description and price
	for _, item := range receipt.Items {
		description := item.ShortDescription
		price := item.Price

		if len(description) == 0 || !ValidatePrice(price) {
			return errors.New("invalid item")
		}
	}

	// Validate price/total amount
	if !ValidatePrice(receipt.Total) {
		return errors.New("invalid price")
	}

	return nil
}

func ValidateDate(purchaseDate string) (time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", purchaseDate)
	return parsedDate, err
}

func ValidateTime(purchaseTime string) (time.Time, error) {
	parsedTime, err := time.Parse("15:04", purchaseTime)
	return parsedTime, err
}

func ValidatePrice(amount string) bool {
	return PriceRegEx.MatchString(amount)
}
