package main

import (
	. "jvh_local/TEST/api"
	"log"
	"net/http"
	"os"
)


func main() {

	// Get port for Heroku
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set up handlers
	http.HandleFunc("/igcinfo/api/", InfoHandler)
	http.HandleFunc("/igcinfo/api/igc", ApiHandler)
	http.HandleFunc("/igcinfo/api/igc/", IdHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}