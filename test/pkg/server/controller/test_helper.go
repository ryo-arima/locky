package controller

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestHelper provides utilities for testing
type TestHelper struct {
	DB         *gorm.DB
	BaseConfig config.BaseConfig
	MockDB     sqlmock.Sqlmock
	SqlDB      *sql.DB
}

// NewTestHelper creates a new test helper with mock database
func NewTestHelper() *TestHelper {
	// Create a mock database connection
	db, mockDB, sqlDB := setupTestDB()

	baseConfig := config.BaseConfig{
		DBConnection: db,
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Common: config.Common{},
				Server: config.Server{
					Admin: config.Admin{
						Emails: []string{"admin@test.com"},
					},
				},
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
					UserEmail:      "test@test.com",
					UserPassword:   "testpass",
				},
			},
			MySQL: config.MySQL{
				Host: "localhost",
				User: "test",
				Pass: "test",
				Port: "3306",
				Db:   "test_db",
			},
		},
	}

	return &TestHelper{
		DB:         db,
		BaseConfig: baseConfig,
		MockDB:     mockDB,
		SqlDB:      sqlDB,
	}
}

// setupTestDB creates a test database connection with mock
func setupTestDB() (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	// Create a mock database for testing
	return createMockDB()
}

// createMockDB creates a mock database connection for testing using sqlmock
func createMockDB() (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	// Create a mock SQL database
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic(fmt.Sprintf("Failed to create mock database: %v", err))
	}

	// Set up basic expectations for GORM initialization
	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("8.0.0"))

	// Create GORM DB with the mock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create GORM database: %v", err))
	}

	return gormDB, mock, sqlDB
}

// CleanupDB cleans up the test database
func (th *TestHelper) CleanupDB() {
	if th.SqlDB != nil {
		th.SqlDB.Close()
	}
}

// CreateTestUser creates a test user (for testing purposes - returns mock data)
func (th *TestHelper) CreateTestUser(email, name, password string) model.Users {
	now := time.Now()
	user := model.Users{
		ID:        1, // Mock ID
		UUID:      "test-uuid-" + fmt.Sprintf("%d", time.Now().UnixNano()),
		Email:     email,
		Name:      name,
		Password:  password,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// In a real test, you would set up mock expectations here
	// For now, we just return the mock user data
	return user
}
