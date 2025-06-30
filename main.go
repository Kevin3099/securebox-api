package main

import (
	"fmt"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(("Content-Type"), "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "SecureBox API is running")
	w.Write([]byte(`{"status": "OK"}`))
}

func main() {
	http.HandleFunc("/health", healthHandler)

	fmt.Println("SecureBox API is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
