package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Starte den API-Server und registriere die Endpunkte
	http.HandleFunc("/crawl", handleCrawlRequest)

	fmt.Println("Server is running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
