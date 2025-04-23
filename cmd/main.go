package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dodskygge/go_swift/internal/db"
	"github.com/dodskygge/go_swift/internal/handler"
	"github.com/dodskygge/go_swift/internal/repository"
	"github.com/dodskygge/go_swift/internal/service"
)

func main() {
	fmt.Println("SWIFT REST API")
	fmt.Println("Server is starting...")

	// Load environment variables from .env file for local development
	//if err := godotenv.Load(); err != nil {
	//	fmt.Println("Error loading .env file:", err)
	//	os.Exit(1)
	//}

	// Connect to the database
	database, err := db.ConnectDB()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		os.Exit(1)
	}
	defer database.Close()

	// Initialize repository, service, and set the global service variable
	repo := &repository.MySQLSwiftRepository{DB: database}
	swiftService := service.NewSwiftCodeService(repo)
	handler.SwiftService = swiftService

	// Setup HTTP handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", handler.HealthCheckHandler)                          // Health check endpoint
	mux.HandleFunc("/api/v1/swift-codes/country/", handler.GetSwiftCodesByCountryHandler) // Get SWIFT codes by country
	mux.HandleFunc("/api/v1/swift-codes", handler.CreateSwiftCodeHandler)                 // Create SWIFT code
	mux.HandleFunc("/api/v1/swift-codes/", func(w http.ResponseWriter, r *http.Request) { // Handle GET and DELETE for SWIFT codes
		if r.Method == http.MethodGet {
			handler.GetSwiftCodeHandler(w, r)
		} else if r.Method == http.MethodDelete {
			handler.DeleteSwiftCodeHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Started successfully. Listening on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
