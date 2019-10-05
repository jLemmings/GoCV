package models

import "time"

type Education struct {
	Title     string
	Institute string
	From      time.Time
	To        time.Time
}
