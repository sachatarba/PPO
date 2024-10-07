package dto

import "github.com/google/uuid"

type PostTrainer struct {
	Id            uuid.UUID   `json:"id" binding:"required"`
	Fullname      string      `json:"fullname" binding:"required"`
	Email         string      `json:"email" binding:"required"`
	Phone         string      `json:"phone" binding:"required"`
	Qualification string      `json:"qualification" binding:"required"`
	UnitPrice     float64     `json:"unit_price" binding:"required"`
	GymsID        []uuid.UUID `json:"gyms_id" binding:"required"`
}
