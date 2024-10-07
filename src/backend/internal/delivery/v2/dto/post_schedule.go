package dto

import "github.com/google/uuid"

type PostSchedule struct {
	Id           uuid.UUID `json:"id" binding:"required"`
	DayOfTheWeek string    `json:"day_of_the_week" binding:"required"`
	StartTime    string    `json:"start_time" binding:"required"`
	EndTime      string    `json:"end_time" binding:"required"`
	TrainingID   uuid.UUID `json:"training_id" binding:"required"`
}
