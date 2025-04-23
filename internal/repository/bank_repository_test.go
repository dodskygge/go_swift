package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dodskygge/go_swift/internal/model"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

// Helper function to set up a test database
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)

	// Create the `banks` table
	_, err = db.Exec(`
        CREATE TABLE banks (
            swift_code TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            address TEXT NOT NULL,
            country_iso2_code TEXT NOT NULL,
            country_name TEXT NOT NULL,
            is_headquarter BOOLEAN NOT NULL
        )
    `)
	assert.NoError(t, err)

	return db
}

// Unit test for GetBySwiftCode
func TestGetBySwiftCode(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &MySQLSwiftRepository{DB: db}

	// Insert test data
	_, err := db.Exec(`
        INSERT INTO banks (swift_code, name, address, country_iso2_code, country_name, is_headquarter)
        VALUES ('TESTUS33XXX', 'Test Bank', '123 Main St', 'US', 'United States', TRUE)
    `)
	assert.NoError(t, err)

	// Test GetBySwiftCode
	entity, err := repo.GetBySwiftCode(context.Background(), "TESTUS33XXX")
	assert.NoError(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, "Test Bank", entity.BankName)
	assert.Equal(t, true, entity.IsHeadquarter)

	// Test non-existent SWIFT code
	entity, err = repo.GetBySwiftCode(context.Background(), "NONEXISTENT")
	assert.NoError(t, err)
	assert.Nil(t, entity)
}

// Unit test for GetBranchesByHqSwiftCode
func TestGetBranchesByHqSwiftCode(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &MySQLSwiftRepository{DB: db}

	// Insert test data
	_, err := db.Exec(`
        INSERT INTO banks (swift_code, name, address, country_iso2_code, country_name, is_headquarter)
        VALUES 
            ('TESTUS33XXX', 'Test Bank HQ', '123 Main St', 'US', 'United States', TRUE),
            ('TESTUS33ABC', 'Test Bank Branch', '456 Branch St', 'US', 'United States', FALSE)
    `)
	assert.NoError(t, err)

	// Test GetBranchesByHqSwiftCode
	branches, err := repo.GetBranchesByHqSwiftCode(context.Background(), "TESTUS33")
	assert.NoError(t, err)
	assert.Len(t, branches, 1)
	assert.Equal(t, "Test Bank Branch", branches[0].BankName)
	assert.Equal(t, false, branches[0].IsHeadquarter)
}

// Unit test for GetByCountry
func TestGetByCountry(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &MySQLSwiftRepository{DB: db}

	// Insert test data
	_, err := db.Exec(`
        INSERT INTO banks (swift_code, name, address, country_iso2_code, country_name, is_headquarter)
        VALUES 
            ('TESTUS33XXX', 'Test Bank HQ', '123 Main St', 'US', 'United States', TRUE),
            ('TESTUS33ABC', 'Test Bank Branch', '456 Branch St', 'US', 'United States', FALSE)
    `)
	assert.NoError(t, err)

	// Test GetByCountry
	entities, err := repo.GetByCountry(context.Background(), "US")
	assert.NoError(t, err)
	assert.Len(t, entities, 2)
	assert.Equal(t, "Test Bank HQ", entities[0].BankName)
	assert.Equal(t, "Test Bank Branch", entities[1].BankName)
}

// Unit test for Create
func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &MySQLSwiftRepository{DB: db}

	// Test Create
	entity := &model.SwiftEntity{
		SwiftCode:     "TESTUS33XXX",
		BankName:      "Test Bank",
		Address:       "123 Main St",
		CountryISO2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
	}
	err := repo.Create(context.Background(), entity)
	assert.NoError(t, err)

	// Verify data in the database
	row := db.QueryRow(`
        SELECT swift_code, name, address, country_iso2_code, country_name, is_headquarter
        FROM banks
        WHERE swift_code = 'TESTUS33XXX'
    `)
	var swiftCode, name, address, countryISO2, countryName string
	var isHeadquarter bool
	err = row.Scan(&swiftCode, &name, &address, &countryISO2, &countryName, &isHeadquarter)
	assert.NoError(t, err)
	assert.Equal(t, "TESTUS33XXX", swiftCode)
	assert.Equal(t, "Test Bank", name)
	assert.Equal(t, true, isHeadquarter)
}

// Unit test for Delete
func TestDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &MySQLSwiftRepository{DB: db}

	// Insert test data
	_, err := db.Exec(`
        INSERT INTO banks (swift_code, name, address, country_iso2_code, country_name, is_headquarter)
        VALUES ('TESTUS33XXX', 'Test Bank', '123 Main St', 'US', 'United States', TRUE)
    `)
	assert.NoError(t, err)

	// Test Delete
	err = repo.Delete(context.Background(), "TESTUS33XXX")
	assert.NoError(t, err)

	// Verify data is deleted
	row := db.QueryRow(`
        SELECT swift_code FROM banks WHERE swift_code = 'TESTUS33XXX'
    `)
	err = row.Scan(new(string))
	assert.Equal(t, sql.ErrNoRows, err)
}
