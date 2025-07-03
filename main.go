package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

var numbers = []int{1, 2, 3, 4, 5}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/numbers", numbersHandler)

	fmt.Println("SecureBox API is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "OK"}`))
}

func numbersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getNumbers(w, r)
	case http.MethodPost:
		postNumber(w, r)
	case http.MethodPut:
		updateNumber(w, r)
	case http.MethodDelete:
		deleteNumber(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /numbers
func getNumbers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(numbers)
}

// POST /numbers
// Body: { "number": 42 }
func postNumber(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Number int `json:"number"`
	}
	json.NewDecoder(r.Body).Decode(&input)
	numbers = append(numbers, input.Number)
	w.Write([]byte("Number added"))
}

// PUT /numbers?index=1
// Body: { "number": 99 }
func updateNumber(w http.ResponseWriter, r *http.Request) {
	indexStr := r.URL.Query().Get("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 || index >= len(numbers) {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	var input struct {
		Number int `json:"number"`
	}
	json.NewDecoder(r.Body).Decode(&input)
	numbers[index] = input.Number
	w.Write([]byte("Number updated"))
}

// DELETE /numbers?index=2
func deleteNumber(w http.ResponseWriter, r *http.Request) {
	indexStr := r.URL.Query().Get("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 || index >= len(numbers) {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	numbers = append(numbers[:index], numbers[index+1:]...)
	w.Write([]byte("Number deleted"))
}
