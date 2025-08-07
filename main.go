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
	http.HandleFunc("/store", storeHandler)   // POST
	http.HandleFunc("/get", getHandler)       // GET
	http.HandleFunc("/update", updateHandler) // PUT
	http.HandleFunc("/delete", deleteHandler) // DELETE

	log.Println("SecureBoxAPI running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// storeHandler handles POST /store
func storeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID     string `json:"id"`
		Secret string `json:"secret"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == "" || req.Secret == "" {
		http.Error(w, "Missing id or secret", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if _, exists := secrets[req.ID]; exists {
		http.Error(w, "Secret already exists. Use update instead.", http.StatusConflict)
		return
	}
	secrets[req.ID] = req.Secret

	json.NewEncoder(w).Encode(map[string]string{"status": "stored"})
}

// getHandler handles GET /get?id=...
func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

// updateHandler handles PUT /update
func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID     string `json:"id"`
		Secret string `json:"secret"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == "" || req.Secret == "" {
		http.Error(w, "Missing id or secret", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if _, exists := secrets[req.ID]; !exists {
		http.Error(w, "Secret not found", http.StatusNotFound)
		return
	}
	secrets[req.ID] = req.Secret

	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// deleteHandler handles DELETE /delete?id=...
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if _, exists := secrets[id]; !exists {
		http.Error(w, "Secret not found", http.StatusNotFound)
		return
	}
	delete(secrets, id)

	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
