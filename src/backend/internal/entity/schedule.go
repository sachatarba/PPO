package entity

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID           uuid.UUID
	DayOfTheWeek string
	StartTime    string
	EndTime      string
	ClientID     uuid.UUID
	TrainingID   uuid.UUID
	Training     Training
}

func (sh *Schedule) Validate() bool {
	date, err := time.Parse(time.DateOnly, sh.DayOfTheWeek)
	if err != nil {
		return false
	}
	sh.DayOfTheWeek = date.Format(time.DateOnly)

	startTime, err := time.Parse(time.TimeOnly, sh.StartTime)
	if err != nil {
		return false
	}
	sh.StartTime = time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		startTime.Hour(),
		startTime.Minute(),
		startTime.Second(),
		0,
		time.UTC,
	).
		Format(time.RFC3339)

	endTime, err := time.Parse(time.TimeOnly, sh.EndTime)
	if err != nil {
		return false
	}
	sh.EndTime = time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		endTime.Hour(),
		endTime.Minute(),
		endTime.Second(),
		0,
		time.UTC,
	).
		Format(time.RFC3339)

	// если кто-то ночью тренируется, то не очень работает:D
	// return startTime.Before(endTime)

	return true
}
