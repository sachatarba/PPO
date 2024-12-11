package orm

import (
	"github.com/google/uuid"
)

type ClientMembership struct {
	ID               uuid.UUID `gorm:"type:UUID;primaryKey"`
	StartDate        string    `gorm:"type:DATE"`
	EndDate          string    `gorm:"type:DATE"`
	MembershipTypeID uuid.UUID
	MembershipType   MembershipType
	ClientID         uuid.UUID
}
