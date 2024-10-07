package dto

import "github.com/google/uuid"

type PostMembershipType struct {
	Id           uuid.UUID `json:"id" binding:"required"`
	Type         string    `json:"type" binding:"required"`
	Description  string    `json:"description" binding:"required"`
	Price        string    `json:"price" binding:"required"`
	DaysDuration int       `json:"days_duration" binding:"required"`
}
