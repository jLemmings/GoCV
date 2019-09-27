package models

import (
	"context"
	"fmt"
	"github.com/google/go-github/v28/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"log"
	"os"
)

var gitClient *github.Client

func init() {
	ctx := context.Background()
	if os.Getenv("ENV") != "PROD" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env.production file")
		}
	}
	gitKey := os.Getenv("GITHUB_KEY")
	fmt.Println(gitKey)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_KEY")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	gitClient = client
}

func GetGitClient() *github.Client {
	return gitClient
}
