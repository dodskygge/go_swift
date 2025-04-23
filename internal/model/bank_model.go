package model

// Response for a single SWIFT code
type SwiftCodeResponse struct {
	Address       string            `json:"address"`
	BankName      string            `json:"bankName"`
	CountryISO2   string            `json:"countryISO2"`
	CountryName   string            `json:"countryName"`
	IsHeadquarter bool              `json:"isHeadquarter"`
	SwiftCode     string            `json:"swiftCode"`
	Branches      []SwiftCodeBranch `json:"branches,omitempty"`
}

// Response for a branch of a SWIFT code
type SwiftCodeBranch struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

// Response for SWIFT codes by country
type SwiftCodesByCountryResponse struct {
	CountryISO2 string                     `json:"countryISO2"`
	CountryName string                     `json:"countryName"`
	SwiftCodes  []SwiftCodeMinimalResponse `json:"swiftCodes"`
}

// Minimal response for a SWIFT code
type SwiftCodeMinimalResponse struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

// Request to create a new SWIFT code
type CreateSwiftCodeRequest struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

// Entity representing a SWIFT code in the database
type SwiftEntity struct {
	Address       string
	BankName      string
	CountryISO2   string
	CountryName   string
	IsHeadquarter bool
	SwiftCode     string
}
