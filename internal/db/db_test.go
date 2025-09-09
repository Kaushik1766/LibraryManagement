package db

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetDB(t *testing.T) {
	originalDBURL := os.Getenv("DATABASE_URL")
	defer func() {
		if originalDBURL != "" {
			os.Setenv("DATABASE_URL", originalDBURL)
		} else {
			os.Unsetenv("DATABASE_URL")
		}
	}()

	tests := []struct {
		name        string
		dbURL       string
		expectPanic bool
		panicMsg    string
	}{
		{
			name:        "valid database URL",
			dbURL:       "postgres://user:password@localhost:5432/testdb?sslmode=disable",
			expectPanic: false,
		},
		{
			name:        "missing database URL",
			dbURL:       "",
			expectPanic: true,
			panicMsg:    "database url not specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dbURL != "" {
				os.Setenv("DATABASE_URL", tt.dbURL)
			} else {
				os.Unsetenv("DATABASE_URL")
			}

			if tt.expectPanic {
				defer func() {
					if r := recover(); r != nil {
						if panicMsg, ok := r.(string); ok {
							if panicMsg != tt.panicMsg {
								t.Errorf("GetDB() panic message = %v, want %v", panicMsg, tt.panicMsg)
							}
						} else {
							t.Errorf("GetDB() panic type = %T, want string", r)
						}
					} else {
						t.Errorf("GetDB() expected panic but didn't get one")
					}
				}()
			}

			if !tt.expectPanic {
				t.Skip("skipping as db con req")
			} else {
				GetDB()
			}
		})
	}
}

func TestCreateTables(t *testing.T) {
	// Create a temporary SQL file for testing
	tempSQLContent := `
		CREATE TABLE IF NOT EXISTS test_table (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50)
		);
	`

	// Create temporary file
	tempFile, err := os.CreateTemp("", "test_init_*.sql")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up

	// Write test SQL content
	if _, err := tempFile.WriteString(tempSQLContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Save original working directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd) // Restore working directory

	// Change to temp directory for relative path testing
	tempDir := os.TempDir()
	os.Chdir(tempDir)

	// Copy our temp file to the expected location
	expectedPath := "./sql/InitDB.sql"
	os.MkdirAll("./sql", 0755)
	tempSQLFile, err := os.Create(expectedPath)
	if err != nil {
		t.Fatalf("Failed to create test SQL file: %v", err)
	}
	if _, err := tempSQLFile.WriteString(tempSQLContent); err != nil {
		t.Fatalf("Failed to write test SQL file: %v", err)
	}
	tempSQLFile.Close()
	defer os.Remove(expectedPath)
	defer os.Remove("./sql")

	tests := []struct {
		name        string
		setupDB     func() *sql.DB
		expectError bool
		expectPanic bool
		errorMsg    string
	}{
		{
			name: "successful table creation",
			setupDB: func() *sql.DB {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("Failed to create mock database: %v", err)
				}
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS test_table").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			expectError: false,
		},
		{
			name: "database execution error",
			setupDB: func() *sql.DB {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("Failed to create mock database: %v", err)
				}
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS test_table").WillReturnError(errors.New("database error"))
				return db
			},
			expectError: true,
			errorMsg:    "error creating tables: database error",
		},
		{
			name: "nil database connection",
			setupDB: func() *sql.DB {
				return nil
			},
			expectError: true,
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.setupDB()
			if db != nil {
				defer db.Close()
			}

			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("CreateTables() expected panic but didn't get one")
					}
				}()
			}

			err := CreateTables(db)

			if !tt.expectPanic {
				if tt.expectError {
					if err == nil {
						t.Errorf("CreateTables() expected error but got none")
					} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
						t.Errorf("CreateTables() error = %v, want %v", err.Error(), tt.errorMsg)
					}
				} else {
					if err != nil {
						t.Errorf("CreateTables() unexpected error = %v", err)
					}
				}
			}
		})
	}
}

func TestCreateTables_FileNotFound(t *testing.T) {
	// Save original working directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	// Change to a directory that doesn't have the sql file
	tempDir := os.TempDir()
	os.Chdir(tempDir)

	// Ensure the file doesn't exist
	os.Remove("./sql/InitDB.sql")

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// This should panic due to missing file
	defer func() {
		if r := recover(); r != nil {
			if panicMsg, ok := r.(string); ok {
				if !contains(panicMsg, "error reading sql file") {
					t.Errorf("CreateTables() panic message = %v, expected to contain 'error reading sql file'", panicMsg)
				}
			} else {
				t.Errorf("CreateTables() panic type = %T, want string", r)
			}
		} else {
			t.Errorf("CreateTables() expected panic due to missing file but didn't get one")
		}
	}()

	CreateTables(db)
}

func TestGetDB_ConnectionError(t *testing.T) {
	// Note: This test is challenging because sql.Open doesn't immediately validate
	// the connection. The actual connection validation happens when the DB is used.
	// For testing purposes, we'll skip this as it requires actual database setup.

	t.Skip("Skipping connection error test - requires actual database setup for proper testing")
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsAt(s, substr)))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
