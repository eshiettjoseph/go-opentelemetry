package initializers

import (
	"go-rest-api/models"
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
  
var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := os.ExpandEnv("host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} port=5432 sslmode=${SSL_MODE}")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	db.AutoMigrate(&models.User{})
	DB = db
}

