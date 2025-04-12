package models

type Receipt struct {
	Retailer     string  `json:"retailer" binding:"required"`
	PurchaseDate string  `json:"purchaseDate" binding:"required"`
	PurchaseTime string  `json:"purchaseTime" binding:"required"`
	Items        []Items `json:"items" binding:"required"`
	Total        string  `json:"total" binding:"required"`
}

// Map of {id: Receipt} to retrieve Receipt given an id
var Receipts = map[string]Receipt{}
