package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://afonso.dev", http.StatusTemporaryRedirect)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", redirectHandler)

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
