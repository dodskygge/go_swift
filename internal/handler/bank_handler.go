package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dodskygge/go_swift/internal/model"
	"github.com/dodskygge/go_swift/internal/service"
)

var SwiftService *service.SwiftCodeService

// Handles GET /api/v1/swift-codes/{swift-code}
func GetSwiftCodeHandler(w http.ResponseWriter, r *http.Request) {
	prefix := "/api/v1/swift-codes/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.NotFound(w, r)
		return
	}
	swiftCode := strings.TrimPrefix(r.URL.Path, prefix)
	if swiftCode == "" {
		http.NotFound(w, r)
		return
	}

	result, err := SwiftService.GetSwiftCodeDetails(r.Context(), swiftCode)
	if err != nil || result == nil {
		http.Error(w, "SWIFT code not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Handles GET /api/v1/swift-codes/country/{countryISO2code}
func GetSwiftCodesByCountryHandler(w http.ResponseWriter, r *http.Request) {
	prefix := "/api/v1/swift-codes/country/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.NotFound(w, r)
		return
	}

	countryCode := strings.TrimPrefix(r.URL.Path, prefix)
	if countryCode == "" {
		http.NotFound(w, r)
		return
	}

	results, err := SwiftService.GetSwiftCodesByCountry(r.Context(), countryCode)
	if err != nil {
		http.Error(w, "Error fetching SWIFT codes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// Handles POST /api/v1/swift-codes
func CreateSwiftCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req model.CreateSwiftCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	if err := SwiftService.CreateSwiftCode(r.Context(), req); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create SWIFT code: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "SWIFT code created successfully"})
}

// Handles DELETE /api/v1/swift-codes/{swift-code}
func DeleteSwiftCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	prefix := "/api/v1/swift-codes/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.NotFound(w, r)
		return
	}

	swiftCode := strings.TrimPrefix(r.URL.Path, prefix)
	if swiftCode == "" {
		http.NotFound(w, r)
		return
	}

	err := SwiftService.DeleteSwiftCode(r.Context(), swiftCode)
	if err != nil {
		if strings.Contains(err.Error(), "no SWIFT code found") {
			http.Error(w, "SWIFT code not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to delete SWIFT code: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "SWIFT code deleted successfully"})
}
