package dto

import "github.com/google/uuid"

type PostEquipment struct {
	Id          uuid.UUID `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
}
