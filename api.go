package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// handleCrawlRequest verarbeitet eingehende API-Anfragen zum Crawlen eines Linktree-Profils
func handleCrawlRequest(w http.ResponseWriter, r *http.Request) {
	// Sicherstellen, dass die Anfrage eine GET-Anfrage ist
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Holen des Profilnamens aus den Query-Parametern (z.B. ?profile=myprofile)
	profile := r.URL.Query().Get("profile")
	if profile == "" {
		http.Error(w, "Profile name is required", http.StatusBadRequest)
		return
	}

	// Erstellen des Links zum Crawlen
	link := fmt.Sprintf("https://linktr.ee/%s", profile)

	// Crawlen des Profils (Verwendung der Funktion aus crawler.go)
	result, err := crawlLinktreeProfile(link)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error during crawling: %v", err), http.StatusInternalServerError)
		return
	}

	// Das Ergebnis als JSON zurückgeben
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
