package fabric

import (
	"time"

	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

func DeafaultClientMembership() entity.ClientMembership {
	return entity.ClientMembership{
		ID:        uuid.New(),
		StartDate: time.Now().Format(time.DateOnly),
		EndDate:   time.Now().Format(time.DateOnly),
		MembershipType: entity.MembershipType{
			ID: uuid.New(),
		},
		ClientID: uuid.New(),
	}
}
