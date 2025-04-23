package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dodskygge/go_swift/internal/model"
	"github.com/dodskygge/go_swift/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock service for testing
type MockSwiftCodeService struct {
	mock.Mock
}

func (m *MockSwiftCodeService) GetBySwiftCode(ctx context.Context, code string) (*model.SwiftEntity, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SwiftEntity), args.Error(1)
}

func (m *MockSwiftCodeService) GetBranchesByHqSwiftCode(ctx context.Context, hqCode string) ([]*model.SwiftEntity, error) {
	args := m.Called(ctx, hqCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.SwiftEntity), args.Error(1)
}

func (m *MockSwiftCodeService) GetByCountry(ctx context.Context, countryISO2 string) ([]*model.SwiftEntity, error) {
	args := m.Called(ctx, countryISO2)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.SwiftEntity), args.Error(1)
}

func (m *MockSwiftCodeService) Create(ctx context.Context, swift *model.SwiftEntity) error {
	args := m.Called(ctx, swift)
	return args.Error(0)
}

func (m *MockSwiftCodeService) Delete(ctx context.Context, swiftCode string) error {
	args := m.Called(ctx, swiftCode)
	return args.Error(0)
}

func TestGetSwiftCodeHandler(t *testing.T) {
	mockService := new(MockSwiftCodeService)
	SwiftService = service.NewSwiftCodeService(mockService)

	// Mock response for GetBySwiftCode
	mockResponse := &model.SwiftEntity{
		Address:       "123 Main St",
		BankName:      "Test Bank",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
		SwiftCode:     "TESTUS33XXX",
	}

	// Mock response for GetBranchesByHqSwiftCode
	mockBranches := []*model.SwiftEntity{
		{
			Address:       "456 Branch St",
			BankName:      "Test Bank Branch",
			CountryISO2:   "US",
			IsHeadquarter: false,
			SwiftCode:     "TESTUS33ABC",
		},
	}

	// Define mock behavior
	mockService.On("GetBySwiftCode", mock.Anything, "TESTUS33XXX").Return(mockResponse, nil)
	mockService.On("GetBranchesByHqSwiftCode", mock.Anything, "TESTUS33").Return(mockBranches, nil)

	// Create request and recorder
	req := httptest.NewRequest(http.MethodGet, "/api/v1/swift-codes/TESTUS33XXX", nil)
	rec := httptest.NewRecorder()

	// Call handler
	GetSwiftCodeHandler(rec, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)
	var response model.SwiftCodeResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	expectedResponse := model.SwiftCodeResponse{
		Address:       "123 Main St",
		BankName:      "Test Bank",
		CountryISO2:   "US",
		CountryName:   "UNITED STATES",
		IsHeadquarter: true,
		SwiftCode:     "TESTUS33XXX",
		Branches: []model.SwiftCodeBranch{
			{
				Address:       "456 Branch St",
				BankName:      "Test Bank Branch",
				CountryISO2:   "US",
				IsHeadquarter: false,
				SwiftCode:     "TESTUS33ABC",
			},
		},
	}
	assert.Equal(t, expectedResponse, response)

	mockService.AssertExpectations(t)
}

func TestGetSwiftCodesByCountryHandler(t *testing.T) {
	mockService := new(MockSwiftCodeService)
	SwiftService = service.NewSwiftCodeService(mockService)

	// Mock response for GetByCountry
	mockResponse := []*model.SwiftEntity{
		{
			Address:       "123 Main St",
			BankName:      "Test Bank",
			CountryISO2:   "US",
			CountryName:   "United States", //Not using uppercase for testing
			IsHeadquarter: true,
			SwiftCode:     "TESTUS33XXX",
		},
	}

	// Define mock behavior
	mockService.On("GetByCountry", mock.Anything, "US").Return(mockResponse, nil)

	// Create request and recorder
	req := httptest.NewRequest(http.MethodGet, "/api/v1/swift-codes/country/US", nil)
	rec := httptest.NewRecorder()

	// Call handler
	GetSwiftCodesByCountryHandler(rec, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)
	var response model.SwiftCodesByCountryResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	expectedResponse := model.SwiftCodesByCountryResponse{
		CountryISO2: "US",
		CountryName: "UNITED STATES",
		SwiftCodes: []model.SwiftCodeMinimalResponse{
			{
				Address:       "123 Main St",
				BankName:      "Test Bank",
				CountryISO2:   "US",
				IsHeadquarter: true,
				SwiftCode:     "TESTUS33XXX",
			},
		},
	}
	assert.Equal(t, expectedResponse, response)

	mockService.AssertExpectations(t)
}

func TestCreateSwiftCodeHandler(t *testing.T) {
	mockService := new(MockSwiftCodeService)
	SwiftService = service.NewSwiftCodeService(mockService)

	// Mock request body
	requestBody := `{
        "address": "123 Main St",
        "bankName": "Test Bank",
        "countryISO2": "US",
        "countryName": "United States",
        "isHeadquarter": true,
        "swiftCode": "TESTUS33XXX"
    }`

	// Define expected entity
	expectedEntity := &model.SwiftEntity{
		Address:       "123 Main St",
		BankName:      "Test Bank",
		CountryISO2:   "US",
		CountryName:   "UNITED STATES",
		IsHeadquarter: true,
		SwiftCode:     "TESTUS33XXX",
	}

	// Define mock behavior
	mockService.On("Create", mock.Anything, expectedEntity).Return(nil)

	// Create request and recorder
	req := httptest.NewRequest(http.MethodPost, "/api/v1/swift-codes", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Call handler
	CreateSwiftCodeHandler(rec, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, rec.Code)
	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "SWIFT code created successfully", response["message"])

	mockService.AssertExpectations(t)
}

func TestDeleteSwiftCodeHandler(t *testing.T) {
	mockService := new(MockSwiftCodeService)
	SwiftService = service.NewSwiftCodeService(mockService)

	mockService.On("Delete", mock.Anything, "TESTUS33XXX").Return(nil)

	// Create request and recorder
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/swift-codes/TESTUS33XXX", nil)
	rec := httptest.NewRecorder()

	// Call handler
	DeleteSwiftCodeHandler(rec, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rec.Code)
	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "SWIFT code deleted successfully", response["message"])

	mockService.AssertExpectations(t)
}
