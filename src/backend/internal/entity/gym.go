package entity

import (
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/pkg/validator"
)

type Gym struct {
	ID              uuid.UUID
	Name            string
	Phone           string
	City            string
	Addres          string
	IsChain         bool
	Trainers        []Trainer
	Equipments      []Equipment
	MembershipTypes []MembershipType
}

func (g *Gym) Validate() bool {
	if !validator.IsValidPhoneNumber(g.Phone) {
		return false
	}

	return g.Name != "" && g.City != "" && g.Addres != ""
}
