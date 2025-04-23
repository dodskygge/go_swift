package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/dodskygge/go_swift/internal/model"
)

// SwiftCodeRepository defines the interface for repository operations related to SWIFT codes.
type SwiftCodeRepository interface {
	GetBySwiftCode(ctx context.Context, code string) (*model.SwiftEntity, error)
	GetBranchesByHqSwiftCode(ctx context.Context, hqCode string) ([]*model.SwiftEntity, error)
	GetByCountry(ctx context.Context, countryISO2 string) ([]*model.SwiftEntity, error)
	Create(ctx context.Context, swift *model.SwiftEntity) error
	Delete(ctx context.Context, swiftCode string) error
}

// SwiftCodeService provides business logic for SWIFT code operations.
type SwiftCodeService struct {
	repo SwiftCodeRepository
}

// NewSwiftCodeService initializes a new SwiftCodeService with the given repository.
func NewSwiftCodeService(repo SwiftCodeRepository) *SwiftCodeService {
	return &SwiftCodeService{repo: repo}
}

// GetSwiftCodeDetails retrieves details for a specific SWIFT code, including branches if it's a headquarters.
func (s *SwiftCodeService) GetSwiftCodeDetails(ctx context.Context, swiftCode string) (*model.SwiftCodeResponse, error) {
	entity, err := s.repo.GetBySwiftCode(ctx, swiftCode)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}

	// Normalize country codes and names to uppercase.
	entity.CountryISO2 = strings.ToUpper(entity.CountryISO2)
	entity.CountryName = strings.ToUpper(entity.CountryName)

	response := &model.SwiftCodeResponse{
		Address:       entity.Address,
		BankName:      entity.BankName,
		CountryISO2:   entity.CountryISO2,
		CountryName:   entity.CountryName,
		IsHeadquarter: entity.IsHeadquarter,
		SwiftCode:     entity.SwiftCode,
		Branches:      []model.SwiftCodeBranch{},
	}

	// If HQ, retrieve branches
	if entity.IsHeadquarter {
		branches, err := s.repo.GetBranchesByHqSwiftCode(ctx, entity.SwiftCode[:8])
		if err != nil {
			return nil, err
		}
		for _, b := range branches {
			// Pomijamy HQ wśród oddziałów
			if b.IsHeadquarter {
				continue
			}

			b.CountryISO2 = strings.ToUpper(b.CountryISO2)

			branch := model.SwiftCodeBranch{
				Address:       b.Address,
				BankName:      b.BankName,
				CountryISO2:   b.CountryISO2,
				IsHeadquarter: b.IsHeadquarter,
				SwiftCode:     b.SwiftCode,
			}
			response.Branches = append(response.Branches, branch)
		}
	}

	return response, nil
}

// GetSwiftCodesByCountry retrieves all SWIFT codes for a specific country.
func (s *SwiftCodeService) GetSwiftCodesByCountry(ctx context.Context, countryISO2 string) (*model.SwiftCodesByCountryResponse, error) {
	entities, err := s.repo.GetByCountry(ctx, countryISO2)
	if err != nil {
		return nil, err
	}
	if len(entities) == 0 {
		return nil, nil
	}

	// Normalize country codes and names to uppercase
	countryISO2 = strings.ToUpper(countryISO2)
	countryName := strings.ToUpper(entities[0].CountryName)

	response := &model.SwiftCodesByCountryResponse{
		CountryISO2: countryISO2,
		CountryName: countryName,
		SwiftCodes:  []model.SwiftCodeMinimalResponse{},
	}

	for _, entity := range entities {
		swiftCode := model.SwiftCodeMinimalResponse{
			Address:       entity.Address,
			BankName:      entity.BankName,
			CountryISO2:   entity.CountryISO2,
			IsHeadquarter: entity.IsHeadquarter,
			SwiftCode:     entity.SwiftCode,
		}
		response.SwiftCodes = append(response.SwiftCodes, swiftCode)
	}

	return response, nil
}

// CreateSwiftCode validates and creates a new SWIFT code entry in the database.
func (s *SwiftCodeService) CreateSwiftCode(ctx context.Context, req model.CreateSwiftCodeRequest) error {
	// Validate input data.
	if len(req.SwiftCode) < 8 {
		return fmt.Errorf("invalid SWIFT code: must be at least 8 characters")
	}
	if req.CountryISO2 == "" || req.CountryName == "" {
		return fmt.Errorf("countryISO2 and countryName cannot be empty")
	}
	if req.BankName == "" || req.Address == "" {
		return fmt.Errorf("bankName and address cannot be empty")
	}

	// Check if the SWIFT code is valid
	isHQ := len(req.SwiftCode) == 11 && req.SwiftCode[8:] == "XXX"
	if isHQ != req.IsHeadquarter {
		return fmt.Errorf("SWIFT code does not match the provided isHeadquarter value")
	}

	// Create a new SWIFT entity.
	entity := &model.SwiftEntity{
		SwiftCode:     req.SwiftCode,
		BankName:      req.BankName,
		Address:       req.Address,
		CountryISO2:   strings.ToUpper(req.CountryISO2),
		CountryName:   strings.ToUpper(req.CountryName),
		IsHeadquarter: req.IsHeadquarter,
	}

	// Save the entity in the database.
	err := s.repo.Create(ctx, entity)
	if err != nil {
		return fmt.Errorf("failed to create SWIFT code: %w", err)
	}

	return nil
}

// DeleteSwiftCode deletes a SWIFT code entry from the database.
func (s *SwiftCodeService) DeleteSwiftCode(ctx context.Context, swiftCode string) error {
	// Validate the SWIFT code.
	if len(swiftCode) < 8 {
		return fmt.Errorf("invalid SWIFT code: must be at least 8 characters")
	}

	// Delete the SWIFT code from the database.
	err := s.repo.Delete(ctx, swiftCode)
	if err != nil {
		return fmt.Errorf("failed to delete SWIFT code: %w", err)
	}

	return nil
}
