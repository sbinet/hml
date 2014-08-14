package main

import (
	"time"
)

type Event struct {
	SubmissionID     int       `hml:"SubmissionId"`
	DateSubmittedUTC time.Time `hml:"DateSubmittedUtc"`
	TeamID           int       `hml:"TeamId"`
	TeamName         string    `hml:"TeamName"`
	UserID           int       `hml:"UserId"`
	UserDisplayName  string    `hml:"UserDisplayName"`
	PublicScore      float64   `hml:"PublicScore"`
	PrivateScore     float64   `hml:"PrivateScore"`
	IsSelected       bool      `hml:"IsSelected"`
	DateRescoredUTC  time.Time `hml:"DateRescoredUtc"`
	PrevPublicScore  float64   `hml:"PreviousPublicScore"`
	PrevPrivateScore float64   `hml:"PreviousPrivateScore"`
}
