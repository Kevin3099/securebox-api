package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

var numbers = []int{1, 2, 3, 4, 5}
var secrets = map[string]string{} // New map to store secrets

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/numbers", numbersHandler)
	http.HandleFunc("secrets", secretsHandler)

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

// Secrets Handler

func secretsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSecrets(w, r)
	case http.MethodPost:
		postSecret(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getSecrets(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(secrets)
}

func postSecret(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.Key == "" || input.Value == "" {
		http.Error(w, "Key and value cannot be empty", http.StatusBadRequest)
		return
	}

	secrets[input.Key] = input.Value
	w.Write([]byte("Secret added"))
}
