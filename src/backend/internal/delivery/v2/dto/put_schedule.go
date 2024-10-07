package dto

import "github.com/google/uuid"

type PutSchedule struct {
	DayOfTheWeek string    `json:"day_of_the_week" binding:"required"`
	StartTime    string    `json:"start_time" binding:"required"`
	EndTime      string    `json:"end_time" binding:"required"`
	TrainingID   uuid.UUID `json:"training_id" binding:"required"`
}
