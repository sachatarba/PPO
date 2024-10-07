package dto

import "github.com/google/uuid"

type PostLogin struct {
	Id       uuid.UUID `json:"id" binding:"required"`
	Login    string    `json:"login" binding:"required"`
	Password string    `json:"password" binding:"required"`
}
