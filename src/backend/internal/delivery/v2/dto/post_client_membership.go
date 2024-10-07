package dto

import "github.com/google/uuid"

type PostClientMembership struct {
	Id               uuid.UUID `json:"id" binding:"required"`
	StartDate        string    `json:"start_date" binding:"required"`
	EndDate          string    `json:"end_date" binding:"required"`
	MembershipTypeID uuid.UUID `json:"membership_type_id" binding:"required"`
}
