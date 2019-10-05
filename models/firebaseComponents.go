package models

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"os"
)

var fireAuth *auth.Client
var fireDB *db.Client

func init() {
	if os.Getenv("ENV") != "PROD" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env.production file")
		}
	}

	ctx := context.Background()
	var opt option.ClientOption
	if os.Getenv("ENV") == "PROD" {
		opt = option.WithCredentialsJSON([]byte(os.Getenv("SERVICE_ACCOUNT")))
	} else {
		opt = option.WithCredentialsFile("serviceAccountDEV.json")
	}

	config := &firebase.Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	fireAuth = client

	database, err := app.Database(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	fireDB = database
}

func GetAuth() *auth.Client {
	return fireAuth
}

func GetDB() *db.Client {
	return fireDB
}
