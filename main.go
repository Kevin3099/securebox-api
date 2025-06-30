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
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/delete", deleteHandler)

	fmt.Println("SecureBox API is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(("Content-Type"), "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "OK"}`))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	// Allow only GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET Allowed", http.StatusMethodNotAllowed)
		return
	}

	json.NewEncoder(w).Encode(numbers) // Encode the numbers slice as JSON and write it to the response
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST Allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Number int `json:"number"`
	}
	json.NewDecoder(r.Body).Decode(&input)  // Decode the JSON request body into the input struct
	numbers = append(numbers, input.Number) // Append the new number to the slice
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get index from query param
	indexStr := r.URL.Query().Get("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 || index >= len(numbers) {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	// Get number from body
	var input struct {
		Number int `json:"number"`
	}
	// Convert back from JSON to struct
	json.NewDecoder(r.Body).Decode(&input)

	// Replace the number
	numbers[index] = input.Number
	w.Write([]byte("Number updated"))
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}

	indexStr := r.URL.Query().Get("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 || index >= len(numbers) {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	// Remove the item at the index
	numbers = append(numbers[:index], numbers[index+1:]...)
	w.Write([]byte("Number deleted"))
}
