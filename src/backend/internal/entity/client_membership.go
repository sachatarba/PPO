package entity

import (
	"log"
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
		log.Println("start:", startDate, err)
		return false
	}

	m.StartDate = startDate.Format(time.DateOnly)

	endDate, err := time.Parse(time.DateOnly, m.EndDate)
	if err != nil {
		log.Println("end:", startDate, err)
		return false
	}
	m.EndDate = endDate.Format(time.DateOnly)

	if (!startDate.Before(endDate)) {
		log.Println("not before:")
	}

	return startDate.Before(endDate)
}
