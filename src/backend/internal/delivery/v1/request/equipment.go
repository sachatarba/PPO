package request

import "github.com/google/uuid"

type EquipmentReq struct {
	ID          uuid.UUID `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	GymID       uuid.UUID `json:"gymid" binding:"required"`
}
