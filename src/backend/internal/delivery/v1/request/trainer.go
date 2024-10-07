package request

import "github.com/google/uuid"

type TrainerReq struct {
	ID            uuid.UUID   `json:"id" binding:"required"`
	Fullname      string      `json:"fullname" binding:"required"`
	Email         string      `json:"email" binding:"required"`
	Phone         string      `json:"phone" binding:"required"`
	Qualification string      `json:"qualification" binding:"required"`
	UnitPrice     float64     `json:"unitprice" binding:"required"`
	GymsID        []uuid.UUID `json:"gymsid" binding:"required"`
}
