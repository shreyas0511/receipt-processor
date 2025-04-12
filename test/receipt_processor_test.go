package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-processor/handlers"
	"testing"

	"github.com/gin-gonic/gin"
)

func setUpRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/receipts/process", handlers.ProcessReceipt)
	router.GET("/receipts/:id/points", handlers.GetPointsFromID)

	return router
}

func TestValidReceipt(t *testing.T) {
	router := setUpRouter()
	receipt1 := `{
		"retailer": "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": [
			{
			"shortDescription": "Mountain Dew 12PK",
			"price": "6.49"
			}
		],
		"total": "35.35"
		}`
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(receipt1)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 Ok, received %d instead", w.Code)
	}
}

func TestInvalidReceiptFields(t *testing.T) {
	router := setUpRouter()
	tests := []struct {
		testName    string
		receiptJson string
		statusCode  int
	}{
		{testName: "Invalid Retailer",
			receiptJson: `{
			"retailer": "",
			"purchaseDate": "2022-01-01",
			"purchaseTime": "13:01",
			"items": [
				{
				"shortDescription": "Mountain Dew 12PK",
				"price": "6.49"
				}
			],
			"total": "6.49"
			}`,
			statusCode: http.StatusBadRequest,
		}, {
			testName: "Invalid Time",
			receiptJson: `{
			"retailer": "Target",
			"purchaseDate": "2022-01-01",
			"purchaseTime": "25:01",
			"items": [
				{
				"shortDescription": "Mountain Dew 12PK",
				"price": "6.49"
				}
			],
			"total": "6.49"
			}`,
			statusCode: http.StatusBadRequest,
		}, {
			testName: "Invalid Date",
			receiptJson: `{
			"retailer": "Target",
			"purchaseDate": "2021-02-29",
			"purchaseTime": "13:01",
			"items": [
				{
				"shortDescription": "Mountain Dew 12PK",
				"price": "6.49"
				}
			],
			"total": "6.49"
			}`,
			statusCode: http.StatusBadRequest,
		}, {
			testName: "Invalid Price",
			receiptJson: `{
			"retailer": "Target",
			"purchaseDate": "2022-01-01",
			"purchaseTime": "13:01",
			"items": [
				{
				"shortDescription": "Mountain Dew 12PK",
				"price": "6.100"
				}
			],
			"total": "7.00"
			}`,
			statusCode: http.StatusBadRequest,
		}, {
			testName: "Invalid Total",
			receiptJson: `{
			"retailer": "Target",
			"purchaseDate": "2022-01-01",
			"purchaseTime": "13:01",
			"items": [
				{
				"shortDescription": "Mountain Dew 12PK",
				"price": "6.49"
				}
			],
			"total": "asdf"
			}`,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(test.receiptJson)))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.statusCode {
				t.Errorf("Expected status %d, received %d instead", test.statusCode, w.Code)
			}
		})
	}
}

func TestMissingReceiptFields(t *testing.T) {
	router := setUpRouter()
	tests := []struct {
		testName    string
		receiptJson string
		statusCode  int
	}{
		{testName: "Missing Retailer",
			receiptJson: `{
			"purchaseDate": "2022-01-01",
			"purchaseTime": "13:01",
			"items": [
				{
				"shortDescription": "Mountain Dew 12PK",
				"price": "6.49"
				}
			],
			"total": "6.49"
			}`,
			statusCode: http.StatusBadRequest,
		}, {
			testName: "Missing Time",
			receiptJson: `{
			"retailer": "Target",
			"purchaseDate": "2022-01-01",
			"items": [
				{
				"shortDescription": "Mountain Dew 12PK",
				"price": "6.49"
				}
			],
			"total": "6.49"
			}`,
			statusCode: http.StatusBadRequest,
		}, {
			testName: "Missing Date",
			receiptJson: `{
			"retailer": "Target",
			"purchaseTime": "13:01",
			"items": [
				{
				"shortDescription": "Mountain Dew 12PK",
				"price": "6.49"
				}
			],
			"total": "6.49"
			}`,
			statusCode: http.StatusBadRequest,
		}, {
			testName: "Missing Price",
			receiptJson: `{
			"retailer": "Target",
			"purchaseDate": "2022-01-01",
			"purchaseTime": "13:01",
			"items": [
				{
				"shortDescription": "Mountain Dew 12PK"
				}
			],
			"total": "7.00"
			}`,
			statusCode: http.StatusBadRequest,
		}, {
			testName: "Missing Total",
			receiptJson: `{
			"retailer": "Target",
			"purchaseDate": "2022-01-01",
			"purchaseTime": "13:01",
			"items": [
				{
				"shortDescription": "Mountain Dew 12PK",
				"price": "6.49"
				}
			]
			}`,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(test.receiptJson)))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.statusCode {
				t.Errorf("Expected status %d bad request, received %d instead", test.statusCode, w.Code)
			}
		})
	}
}

func TestInvalidID(t *testing.T) {
	router := setUpRouter()
	testId := "3bbfe121-ef25-4412-ba56-a1be43e06b4a"
	req, _ := http.NewRequest("GET", "/receipts/"+testId+"/points", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, not found, received %d instead", w.Code)
	}
}

func TestCalculatePoints(t *testing.T) {
	router := setUpRouter()
	receipt1 := `{
	"retailer": "Target",
	"purchaseDate": "2022-01-01",
	"purchaseTime": "13:01",
	"items": [
		{
		"shortDescription": "Mountain Dew 12PK",
		"price": "6.49"
		},{
		"shortDescription": "Emils Cheese Pizza",
		"price": "12.25"
		},{
		"shortDescription": "Knorr Creamy Chicken",
		"price": "1.26"
		},{
		"shortDescription": "Doritos Nacho Cheese",
		"price": "3.35"
		},{
		"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
		"price": "12.00"
		}
	],
	"total": "35.35"
	}`
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(receipt1)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 Ok, received %d instead", w.Code)
	}

	// get the id
	var res map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("Failed to parse POST response: %v", err)
	}

	testId := res["id"]
	req, _ = http.NewRequest("GET", "/receipts/"+testId+"/points", nil)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 ok, received %d instead", w.Code)
	}

	var pointsRes map[string]int
	_ = json.Unmarshal(w.Body.Bytes(), &pointsRes)

	expectedPoints := 28
	if pointsRes["points"] != expectedPoints {
		t.Errorf("Expected %d points, got %d instead", expectedPoints, pointsRes["points"])
	}
}
