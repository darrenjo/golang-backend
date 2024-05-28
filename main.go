package main

import (
	"backend-api/database"
	"log"
	"net/http"
	"os"

	"backend-api/router"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize and migrate the database
	database.Init()
	database.Migrate()
}

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Set up routes using router from the router package
	r = router.NewRouter()

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port if not specified
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
