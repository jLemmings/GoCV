package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *gorm.DB

func init() {
	var database *gorm.DB
	var err error
	if os.Getenv("ENV") == "PROD" {
		fmt.Println("Running in PRODUCTION")
		database, err = gorm.Open("postgres", os.Getenv("DATABASE_URL")+"?sslmode=disable")
	} else {
		fmt.Println("Running in DEVELOPMENT")
		database, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable")
	}

	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	database.AutoMigrate(&User{})

	db = database
}

func InitializeDB(firstName string, lastName string, email string, password string, github string) {
	firstUser := User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		Bio: "EDIT ME",
		GithubProfile: github,
	}
	firstUser.Create()
}

func GetDB() *gorm.DB {
	return db
}
