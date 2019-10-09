package models

import (
	"github.com/google/go-github/v28/github"
)

type Project struct {
	Name        string
	Description string
	Language    string
	URL         string
	Stack       []string
	LastUpdate  github.Timestamp
}
