# Receipt Processing Web Service
Backend API service built with GO and Gin framework. It takes a receipt in JSON format and calculates reward points based on given rules.

## About the webservice
### Features
* Submit receipts via JSON POST request
* Validate receipts with custom validation rules
* Calculate and return reward points for the given receipt
* In-memory map to map IDs to Receipt objects and IDs to Points

### Tech Stack
* GO with Gin framework.
* GO UUID package for generating unique IDs for receipts.
  
### API Endpoints:
`POST /receipts/process`
* POST request with receipt data in JSON format.
* Returns a JSON object with a unique ID for the given receipt.
  
`GET /receipts/{id}/points`
* GET request with ID generated from the POST request.
* Returns a JSON object with points calculated for the corresponding receipt.
  
### API Design
User submits receipt JSON objects via POST request. The receipt must follow these validation rules:
* All fields are required. No field must be missing otherwise a **400 BadRequest** is returned
* **retailer** name should not be an empty field
* **purchaseDate** should be of valid yyyy-mm-dd format
* **purchaseTime** should be of valid hh:mm in 24 hour clock format
* **items** array should not be empty
* item **shortDescription** should not be empty
* item **Price** should be valid
* **total** price should not be an empty field
* price values should not start with a 0, unless immediately followed by a decimal point. There should only be 2 digits after a decimal point from 00 to 99

Once all of the above conditions are satisfied, a unique ID is generated, and the receipt object is mapped to the ID in an in-memory map.

If everything is successful, a status 200 is returned along with the ID as a JSON object. This ID is used to make a GET request.

If points are already cached, return the points. Otherwise, first retrieve the corresponding receipt object and calculate reward points based on the given rules.
Then, we map the ID to calculated points to retrieve later.

If everything is successful, a status 200 is returned along with the points as a JSON object. Otherwise, a status 404 is returned.
  
## Instructions to run the project
### With GO installed
#### Clone the repository
```
git clone https://github.com/shreyas0511/receipt-processor.git
```
#### Navigate to the project folder
```
cd receipt-processor
```
#### Run the main.go file
```
go run main.go
```
or
```
go run .
```
### With Docker
#### Build the docker image
```
docker build -t receipt-processor .
```
#### Run the docker container
```
docker run -p 8080:8080 receipt-processor
```
The service will be exposed at http://localhost:8080

### To test the API
#### Run tests in the test folder
```
go test ./test
```
#### With Postman, curl or any API client
I used a VScode REST client to send GET and POST requests for simplicity. It will work with any API client or command line tools like curl

Send a POST request with an appropriate receipt
```http
POST http://localhost:8080/receipts/process
Content-Type: application/json

{
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
}
```
This will return status 200 and an ID if the receipt is valid
```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 12 Apr 2025 21:55:47 GMT
Content-Length: 45
Connection: close

{
  "id": "3e009cab-db3d-4237-83f3-d9eaaaadfd4e"
}
```
Use the returned id to make a get request
```http
GET http://localhost:8080/receipts/3e009cab-db3d-4237-83f3-d9eaaaadfd4e/points
```
This will return a status 200 and calculated points if ID is valid
```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sat, 12 Apr 2025 21:56:46 GMT
Content-Length: 13
Connection: close

{
  "points": 28
}
```

## About me
My name is Shreyas Shridhar. I graduated with a Master's in Computer Science from Rochester Institute of Technology in December 2024. I have a strong interest in backend development, cloud infrastructure and building scalable systems. You can reach me at:
* Email: ss9531@rit.edu
* Linkedin: https://linkedin.com/in/shreyasshridhar0511

This project was completed as a part of Fetch Backend Engineer Apprenticeship position.
