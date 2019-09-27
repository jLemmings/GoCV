package models

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"log"
	"os"
)

var fireAuth *auth.Client

func init() {
	ctx := context.Background()
	var opt option.ClientOption
	if os.Getenv("ENV") == "PROD" {
		opt = option.WithCredentialsFile(os.Getenv("SERVICE_ACCOUNT"))
	} else {
		opt = option.WithCredentialsFile("serviceAccountDEV.json")
	}
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	fireAuth = client
}

func GetAuth() *auth.Client {
	return fireAuth
}
