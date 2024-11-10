package builder

import (
	"time"

	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

type ClientMembershipBuilder struct {
	membership entity.ClientMembership
}

func NewClientMembershipBuilder() *ClientMembershipBuilder {
	return &ClientMembershipBuilder{
		membership: entity.ClientMembership{
			ID:        uuid.New(),
			StartDate: time.Now().Format(time.DateOnly),
			EndDate:   time.Now().Add(time.Hour * 24).Format(time.DateOnly),
		},
	}
}

func (b *ClientMembershipBuilder) Invalid() *ClientMembershipBuilder {
	b.membership.ID = uuid.Nil
	b.membership.StartDate = ""
	b.membership.EndDate = ""

	return b
}

func (b *ClientMembershipBuilder) SetID(id uuid.UUID) *ClientMembershipBuilder {
	b.membership.ID = id
	return b
}

func (b *ClientMembershipBuilder) SetStartDate(startDate string) *ClientMembershipBuilder {
	b.membership.StartDate = startDate
	return b
}

func (b *ClientMembershipBuilder) SetEndDate(endDate string) *ClientMembershipBuilder {
	b.membership.EndDate = endDate
	return b
}

func (b *ClientMembershipBuilder) SetMembershipType(membershipType entity.MembershipType) *ClientMembershipBuilder {
	b.membership.MembershipType = membershipType
	return b
}

func (b *ClientMembershipBuilder) SetClientID(clientID uuid.UUID) *ClientMembershipBuilder {
	b.membership.ClientID = clientID
	return b
}

func (b *ClientMembershipBuilder) Build() entity.ClientMembership {
	return b.membership
}
