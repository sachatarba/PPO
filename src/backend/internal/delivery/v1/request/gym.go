package request

import "github.com/google/uuid"

type GymReq struct {
	ID      uuid.UUID `json:"id" binding:"required"`
	Name    string    `json:"name" binding:"required"`
	Phone   string    `json:"phone" binding:"required"`
	City    string    `json:"city" binding:"required"`
	Addres  string    `json:"addres" binding:"required"`
	IsChain bool      `json:"ischain" binding:"required"`
}
