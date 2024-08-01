package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// main function program starts here
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", GetPoints).Methods("GET")

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Creating a structure to store receipt data
type Receipt struct {
	ID           string `json:"id"` // Id is added to use in GetPoints API
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
	Points       int    `json:"points"` //points is added, to store calculatePoints value.
}

// creating a structure to store item data
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var receipts map[string]Receipt // Global variable to store receipts

func init() {
	receipts = make(map[string]Receipt)
}

func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt

	err := json.NewDecoder(r.Body).Decode(&receipt) // to decode json and find any errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	receipt.ID = uuid.New().String()
	receipt.Points = calculatePoints(receipt)
	receipts[receipt.ID] = receipt

	response := map[string]string{"id": receipt.ID}
	json.NewEncoder(w).Encode(response)
}

func GetPoints(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	receipt, found := receipts[id]
	if !found {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	response := map[string]int{"points": receipt.Points}
	json.NewEncoder(w).Encode(response)
}

// Calculate points based on rules mentioned in the exercise

func calculatePoints(receipt Receipt) int {
	// One point for every alphanumeric character in the retailer name
	points := countAlphanumeric(receipt.Retailer)

	total, _ := strconv.ParseFloat(receipt.Total, 64)

	// 50 points if the total is a round dollar amount with no cents.
	if isRoundDollar(total) {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25
	if isMultipleOfQuarter(total) {
		points += 25
	}

	// 5 points for every two items on the receipt.
	points += (len(receipt.Items) / 2) * 5

	// if the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range receipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		price, _ := strconv.ParseFloat(item.Price, 64)
		if trimmedLen%3 == 0 {
			points += int(math.Ceil(price * 0.2))
		}
	}

	// 6 points if the day in the purchase date is odd.
	purchaseDay := receipt.PurchaseDate[len(receipt.PurchaseDate)-2:]
	if dayToInt(purchaseDay)%2 != 0 {
		points += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm
	purchaseTime := strings.Split(receipt.PurchaseTime, ":")
	hour := purchaseTime[0]
	if hourInt := hourToInt(hour); hourInt >= 14 && hourInt < 16 {
		points += 10
	}

	return points
}

// to count the alphanumeric characters in the retailer's name on receipt
func countAlphanumeric(str string) int {
	count := 0
	for _, c := range str {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			count++
		}
	}
	return count
}

// to get a round figure
func isRoundDollar(amount float64) bool {
	return amount == float64(int(amount))
}

// multiple of 0.25
func isMultipleOfQuarter(amount float64) bool {
	return math.Mod(amount, 0.25) == 0
}

// hours to int
func hourToInt(hour string) int {
	var h int
	fmt.Sscanf(hour, "%d", &h)
	return h
}

// day to int
func dayToInt(day string) int {
	var d int
	fmt.Sscanf(day, "%d", &d)
	return d
}
