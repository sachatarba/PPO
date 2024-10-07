package request

import (
	"github.com/google/uuid"
)

type TrainingReq struct {
	ID           uuid.UUID `json:"id" binding:"required"`
	Title        string    `json:"title" binding:"required"`
	Description  string    `json:"description" binding:"required"`
	TrainingType string    `json:"trainingType" binding:"required"`
	TrainerID    uuid.UUID `json:"trainerId" binding:"required"`
}
