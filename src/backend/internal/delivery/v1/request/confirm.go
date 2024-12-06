package request

import "github.com/google/uuid"

type ConfirmReq struct {
	ClientID uuid.UUID `json:"client_id" binding:"required"`
	Code     string    `json:"code" binding:"required"`
}
