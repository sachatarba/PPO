package dto

import "github.com/google/uuid"

type PostTraining struct {
	Id           uuid.UUID `json:"id" binding:"required"`
	Title        string    `json:"title" binding:"required"`
	Description  string    `json:"description" binding:"required"`
	TrainingType string    `json:"trainingType" binding:"required"`
}
