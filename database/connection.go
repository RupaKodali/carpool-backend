package database

import (
	"carpool-backend/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance
var DB *gorm.DB

// ConnectDb establishes a connection to the database and runs migrations
func ConnectDb(migrateDB bool) (*gorm.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Format the DSN (Data Source Name) correctly
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Configure GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Enable color
		},
	)

	// Open the GORM DB connection
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	log.Println("Connected to the database successfully!")

	// Run migrations
	if migrateDB {
		if err := RunMigrations(DB); err != nil {
			log.Printf("Migration failed: %v", err)
			return nil, err
		}
	}

	// Set connection pool settings
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return DB, nil
}

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
	log.Println("Running migrations...")

	models := []interface{}{
		&models.User{},
		&models.Ride{},
		&models.Conversation{},
		&models.RequiredRide{},
		&models.Booking{},
		&models.Rating{},
		&models.Message{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
