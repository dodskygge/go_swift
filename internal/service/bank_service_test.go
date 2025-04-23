package service

import (
	"context"
	"testing"

	"github.com/dodskygge/go_swift/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository for testing
type MockSwiftCodeRepository struct {
	mock.Mock
}

func (m *MockSwiftCodeRepository) GetBySwiftCode(ctx context.Context, code string) (*model.SwiftEntity, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SwiftEntity), args.Error(1)
}

func (m *MockSwiftCodeRepository) GetBranchesByHqSwiftCode(ctx context.Context, hqCode string) ([]*model.SwiftEntity, error) {
	args := m.Called(ctx, hqCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.SwiftEntity), args.Error(1)
}

func (m *MockSwiftCodeRepository) GetByCountry(ctx context.Context, countryISO2 string) ([]*model.SwiftEntity, error) {
	args := m.Called(ctx, countryISO2)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.SwiftEntity), args.Error(1)
}

func (m *MockSwiftCodeRepository) Create(ctx context.Context, swift *model.SwiftEntity) error {
	args := m.Called(ctx, swift)
	return args.Error(0)
}

func (m *MockSwiftCodeRepository) Delete(ctx context.Context, swiftCode string) error {
	args := m.Called(ctx, swiftCode)
	return args.Error(0)
}

// Unit test for GetSwiftCodeDetails
func TestGetSwiftCodeDetails(t *testing.T) {
	mockRepo := new(MockSwiftCodeRepository)
	service := NewSwiftCodeService(mockRepo)

	// Mock data
	mockEntity := &model.SwiftEntity{
		SwiftCode:     "TESTUS33XXX",
		BankName:      "Test Bank",
		Address:       "123 Main St",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
	}
	mockBranches := []*model.SwiftEntity{
		{
			SwiftCode:     "TESTUS33ABC",
			BankName:      "Test Branch",
			Address:       "456 Branch St",
			CountryISO2:   "US",
			CountryName:   "United States",
			IsHeadquarter: false,
		},
	}

	// Define mock behavior
	mockRepo.On("GetBySwiftCode", mock.Anything, "TESTUS33XXX").Return(mockEntity, nil)
	mockRepo.On("GetBranchesByHqSwiftCode", mock.Anything, "TESTUS33").Return(mockBranches, nil)

	// Call service
	result, err := service.GetSwiftCodeDetails(context.Background(), "TESTUS33XXX")

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Bank", result.BankName)
	assert.Len(t, result.Branches, 1)
	assert.Equal(t, "Test Branch", result.Branches[0].BankName)

	mockRepo.AssertExpectations(t)
}

// Unit test for GetSwiftCodesByCountry
func TestGetSwiftCodesByCountry(t *testing.T) {
	mockRepo := new(MockSwiftCodeRepository)
	service := NewSwiftCodeService(mockRepo)

	// Mock data
	mockEntities := []*model.SwiftEntity{
		{
			SwiftCode:     "TESTUS33XXX",
			BankName:      "Test Bank",
			Address:       "123 Main St",
			CountryISO2:   "US",
			CountryName:   "United States",
			IsHeadquarter: true,
		},
	}

	// Define mock behavior
	mockRepo.On("GetByCountry", mock.Anything, "US").Return(mockEntities, nil)

	// Call service
	result, err := service.GetSwiftCodesByCountry(context.Background(), "US")

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "US", result.CountryISO2)
	assert.Len(t, result.SwiftCodes, 1)
	assert.Equal(t, "Test Bank", result.SwiftCodes[0].BankName)

	mockRepo.AssertExpectations(t)
}

// Unit test for CreateSwiftCode
func TestCreateSwiftCode(t *testing.T) {
	mockRepo := new(MockSwiftCodeRepository)
	service := NewSwiftCodeService(mockRepo)

	// Mock request
	req := model.CreateSwiftCodeRequest{
		SwiftCode:     "TESTUS33XXX",
		BankName:      "Test Bank",
		Address:       "123 Main St",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
	}

	// Define mock behavior
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	// Call service
	err := service.CreateSwiftCode(context.Background(), req)

	// Assert results
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// Unit test for DeleteSwiftCode
func TestDeleteSwiftCode(t *testing.T) {
	mockRepo := new(MockSwiftCodeRepository)
	service := NewSwiftCodeService(mockRepo)

	// Define mock behavior
	mockRepo.On("Delete", mock.Anything, "TESTUS33XXX").Return(nil)

	// Call service
	err := service.DeleteSwiftCode(context.Background(), "TESTUS33XXX")

	// Assert results
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// Unit test for CreateSwiftCode with validation error
func TestCreateSwiftCodeValidationError(t *testing.T) {
	mockRepo := new(MockSwiftCodeRepository)
	service := NewSwiftCodeService(mockRepo)

	// Mock request with invalid SWIFT code
	req := model.CreateSwiftCodeRequest{
		SwiftCode:     "INVALID",
		BankName:      "Test Bank",
		Address:       "123 Main St",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
	}

	// Call service
	err := service.CreateSwiftCode(context.Background(), req)

	// Assert validation error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid SWIFT code")
}
