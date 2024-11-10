package builder

import (
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

type MembershipTypeBuilder struct {
	membershipType entity.MembershipType
}

func NewMembershipTypeBuilder() *MembershipTypeBuilder {
	return &MembershipTypeBuilder{
		membershipType: entity.MembershipType{
			ID:           uuid.New(),              
			Type:         "Standard",              
			Description:  "Basic membership type", 
			Price:        "100.00",                
			DaysDuration: 30,                      
			GymID:        uuid.New(),              
		},
	}
}

func (b *MembershipTypeBuilder) SetID(id uuid.UUID) *MembershipTypeBuilder {
	b.membershipType.ID = id
	return b
}

func (b *MembershipTypeBuilder) SetType(membershipType string) *MembershipTypeBuilder {
	b.membershipType.Type = membershipType
	return b
}

func (b *MembershipTypeBuilder) SetDescription(description string) *MembershipTypeBuilder {
	b.membershipType.Description = description
	return b
}

func (b *MembershipTypeBuilder) SetPrice(price string) *MembershipTypeBuilder {
	b.membershipType.Price = price
	return b
}

func (b *MembershipTypeBuilder) SetDaysDuration(daysDuration int) *MembershipTypeBuilder {
	b.membershipType.DaysDuration = daysDuration
	return b
}

func (b *MembershipTypeBuilder) SetGymID(gymID uuid.UUID) *MembershipTypeBuilder {
	b.membershipType.GymID = gymID
	return b
}

func (b *MembershipTypeBuilder) Invalid() *MembershipTypeBuilder {
	b.membershipType.ID = uuid.New()
	b.membershipType.Type = ""
	b.membershipType.Description = ""
	b.membershipType.Price = ""
	b.membershipType.DaysDuration = -1
	b.membershipType.GymID = uuid.Nil

	return b
}

func (b *MembershipTypeBuilder) Build() entity.MembershipType {
	return b.membershipType
}
