package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

var (
	secrets = make(map[string]string)
	mu      sync.RWMutex
)

func main() {
	http.HandleFunc("/store", storeHandler)
	http.HandleFunc("/get", getHandler)

	log.Println("Simple SecretBox running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     string `json:"id"`
		Secret string `json:"secret"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == "" || req.Secret == "" {
		http.Error(w, "Missing id or secret", http.StatusBadRequest)
		return
	}

	mu.Lock()
	secrets[req.ID] = req.Secret
	mu.Unlock()

	json.NewEncoder(w).Encode(map[string]string{"status": "stored"})
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	mu.RLock()
	secret, ok := secrets[id]
	mu.RUnlock()

	if !ok {
		http.Error(w, "Secret not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"secret": secret})
}
