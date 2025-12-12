package core

import (
	"fmt"
	"log"
	"time"

	"github.com/OtchereDev/hrms-go/internal/config"
	"github.com/OtchereDev/hrms-go/internal/core/models/base"
	"github.com/OtchereDev/hrms-go/internal/core/models/hr"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database holds the database connection
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := cfg.Database.GetDSN()

	// Configure GORM logger
	var gormLogger logger.Interface
	if cfg.App.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true, // Better performance
		PrepareStmt:            true, // Prepare statements for better performance
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdle)

	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(cfg.Database.MaxConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully")

	return &Database{db}, nil
}

// AutoMigrate runs database migrations
func (db *Database) AutoMigrate() error {
	log.Println("Running database migrations...")

	// Migrate base models
	if err := db.DB.AutoMigrate(
		&base.User{},
		&base.Role{},
		&base.Permission{},
		&base.Session{},
	); err != nil {
		return fmt.Errorf("failed to migrate base models: %w", err)
	}

	// Migrate HR organization models
	if err := db.DB.AutoMigrate(
		&hr.Company{},
		&hr.Department{},
		&hr.DepartmentApprover{},
		&hr.Branch{},
		&hr.Designation{},
		&hr.EmploymentType{},
		&hr.EmployeeGrade{},
		&hr.HolidayList{},
		&hr.Holiday{},
	); err != nil {
		return fmt.Errorf("failed to migrate HR organization models: %w", err)
	}

	// Migrate HR employee models
	if err := db.DB.AutoMigrate(
		&hr.Employee{},
		&hr.EmployeeEducation{},
		&hr.EmployeeExternalWorkHistory{},
		&hr.EmployeeInternalWorkHistory{},
		&hr.EmployeeSkill{},
	); err != nil {
		return fmt.Errorf("failed to migrate HR employee models: %w", err)
	}

	// TODO: Add more model migrations as we create them
	// &hr.Attendance{},
	// &hr.LeaveApplication{},
	// &payroll.SalarySlip{},
	// etc.

	log.Println("Database migrations completed successfully")
	return nil
}

// SeedDefaultData seeds default data into the database
func (db *Database) SeedDefaultData() error {
	log.Println("Seeding default data...")

	// Create default roles if they don't exist
	roles := base.DefaultRoles()
	for _, role := range roles {
		var existingRole base.Role
		result := db.Where("name = ?", role.Name).First(&existingRole)
		if result.RowsAffected == 0 {
			if err := db.Create(&role).Error; err != nil {
				return fmt.Errorf("failed to create role %s: %w", role.Name, err)
			}
			log.Printf("Created role: %s\n", role.Name)
		}
	}

	// Create default admin user if it doesn't exist
	var adminUser base.User
	result := db.Where("email = ?", "admin@hrms.local").First(&adminUser)
	if result.RowsAffected == 0 {
		adminUser = base.User{
			Email:     "admin@hrms.local",
			Username:  "admin",
			Password:  "admin123", // Will be hashed by BeforeCreate hook
			FirstName: "System",
			LastName:  "Administrator",
			Enabled:   true,
		}

		if err := db.Create(&adminUser).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}

		// Assign System Manager role
		var systemManagerRole base.Role
		db.Where("name = ?", "System Manager").First(&systemManagerRole)
		if systemManagerRole.ID > 0 {
			db.Model(&adminUser).Association("Roles").Append(&systemManagerRole)
		}

		log.Println("Created default admin user: admin@hrms.local / admin123")
	}

	log.Println("Default data seeding completed")
	return nil
}

// Close closes the database connection
func (db *Database) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Transaction executes a function within a database transaction
func (db *Database) Transaction(fn func(*gorm.DB) error) error {
	return db.DB.Transaction(fn)
}

// Ping checks if the database connection is alive
func (db *Database) Ping() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
