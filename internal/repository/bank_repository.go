package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dodskygge/go_swift/internal/model"
)

type MySQLSwiftRepository struct {
	DB *sql.DB
}

// Retrieves a SWIFT code by its value
func (repo *MySQLSwiftRepository) GetBySwiftCode(ctx context.Context, swiftCode string) (*model.SwiftEntity, error) {
	query := `
        SELECT swift_code, name, address, country_iso2_code, country_name, is_headquarter
        FROM banks
        WHERE swift_code = ?
    `
	row := repo.DB.QueryRowContext(ctx, query, swiftCode)

	entity := new(model.SwiftEntity)
	err := row.Scan(
		&entity.SwiftCode,
		&entity.BankName,
		&entity.Address,
		&entity.CountryISO2,
		&entity.CountryName,
		&entity.IsHeadquarter,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No result found
	}
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// Retrieves all branches for a given headquarters SWIFT code
func (repo *MySQLSwiftRepository) GetBranchesByHqSwiftCode(ctx context.Context, hqCode string) ([]*model.SwiftEntity, error) {
	query := `
        SELECT swift_code, name, address, country_iso2_code, country_name, is_headquarter
        FROM banks
        WHERE swift_code LIKE ? AND is_headquarter = FALSE
    `
	rows, err := repo.DB.QueryContext(ctx, query, hqCode+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var entities []*model.SwiftEntity
	for rows.Next() {
		entity := new(model.SwiftEntity)
		err := rows.Scan(
			&entity.SwiftCode,
			&entity.BankName,
			&entity.Address,
			&entity.CountryISO2,
			&entity.CountryName,
			&entity.IsHeadquarter,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		entities = append(entities, entity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return entities, nil
}

// Retrieves all SWIFT codes for a given country
func (repo *MySQLSwiftRepository) GetByCountry(ctx context.Context, countryISO2 string) ([]*model.SwiftEntity, error) {
	query := `
        SELECT swift_code, name, address, country_iso2_code, country_name, is_headquarter
        FROM banks
        WHERE country_iso2_code = ?
    `
	rows, err := repo.DB.QueryContext(ctx, query, countryISO2)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var entities []*model.SwiftEntity
	for rows.Next() {
		entity := new(model.SwiftEntity)
		err := rows.Scan(
			&entity.SwiftCode,
			&entity.BankName,
			&entity.Address,
			&entity.CountryISO2,
			&entity.CountryName,
			&entity.IsHeadquarter,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		entities = append(entities, entity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return entities, nil
}

// Creates a new SWIFT code entry
func (repo *MySQLSwiftRepository) Create(ctx context.Context, swift *model.SwiftEntity) error {
	query := `
        INSERT INTO banks (swift_code, name, address, country_iso2_code, country_name, is_headquarter)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err := repo.DB.ExecContext(ctx, query,
		swift.SwiftCode,
		swift.BankName,
		swift.Address,
		swift.CountryISO2,
		swift.CountryName,
		swift.IsHeadquarter,
	)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}

// Deletes a SWIFT code entry
func (repo *MySQLSwiftRepository) Delete(ctx context.Context, swiftCode string) error {
	query := `
        DELETE FROM banks
        WHERE swift_code = ?
    `
	result, err := repo.DB.ExecContext(ctx, query, swiftCode)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no SWIFT code found with the given value")
	}

	return nil
}
