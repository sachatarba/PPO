package entity

import (
	"time"

	"github.com/google/uuid"
)

type ClientMembership struct {
	ID             uuid.UUID
	StartDate      string
	EndDate        string
	MembershipType MembershipType
	ClientID       uuid.UUID
}

func (m *ClientMembership) Validate() bool {
	startDate, err := time.Parse(time.DateOnly, m.StartDate)
	if err != nil {
		return false
	}

	m.StartDate = startDate.Format(time.DateOnly)

	endDate, err := time.Parse(time.DateOnly, m.EndDate)
	if err != nil {
		return false
	}
	m.EndDate = endDate.Format(time.DateOnly)

	return startDate.Before(endDate)
}
