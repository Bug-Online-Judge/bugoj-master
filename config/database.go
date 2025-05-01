package config

import (
	"bugoj-master/model"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initRoles(db *gorm.DB) error {
	var defaultRoles = []model.Role{
		{Name: "user"},
		{Name: "admin"},
	}
	for _, role := range defaultRoles {
		if err := (db.FirstOrCreate(&model.Role{}, model.Role{Name: role.Name})).
			Error; err != nil {
			return err
		}
	}
	return nil
}

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to PostgreSQL:", err)
	}

	err = model.MigrateAll(db)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	DB = db

	// Initialize roles in the database
	if err := initRoles(db); err != nil {
		log.Fatalf("Failed to initialize roles: %v", err)
	}
}
