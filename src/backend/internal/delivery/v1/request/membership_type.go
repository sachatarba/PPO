package request

import "github.com/google/uuid"

type MembershipTypeReq struct {
	ID           uuid.UUID `json:"id" binding:"required"`
	Type         string    `json:"type" binding:"required"`
	Description  string    `json:"description" binding:"required"`
	Price        string    `json:"price" binding:"required"`
	DaysDuration int       `json:"daysduration" binding:"required"`
	GymID        uuid.UUID `json:"gymid" binding:"required"`
}
