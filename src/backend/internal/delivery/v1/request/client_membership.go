package request

import "github.com/google/uuid"

type ClientMembershipReq struct {
	ID               uuid.UUID `json:"id" binding:"required"`
	StartDate        string    `json:"startdate" binding:"required"`
	EndDate          string    `json:"enddate" binding:"required"`
	MembershipTypeID uuid.UUID `json:"membershiptypeid" binding:"required"`
	ClientID         uuid.UUID `json:"clientid" binding:"required"`
}
