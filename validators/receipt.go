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

	// date time should not be after current date time
	// combinedDateTime := time.Date(purchaseDate.Year(), purchaseDate.Month(), purchaseDate.Day(), purchaseTime.Hour(), purchaseTime.Minute(), purchaseTime.Second(), 0, time.UTC)
	// if combinedDateTime.After(time.Now()) {
	// 	return errors.New("invalid date/time")
	// }

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
	// Valid prices 0.99, 1.05, 10.50, 0.00
	// Amount should not start with 0, unless immediately followed by a decimal point
	// Amount should have exactly 2 digits after the decimal point, from 00 to 99, as 100 cents is 1 dollar
	return PriceRegEx.MatchString(amount)
}
