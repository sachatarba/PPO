package builder

import (
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

type GymBuilder struct {
	gym entity.Gym
}

func NewGymBuilder() *GymBuilder {
	return &GymBuilder{
		gym: entity.Gym{
			ID:      uuid.New(),
			Name:    "Default Gym Name",
			Addres:  "Default Address",
			Phone:   "+7-985-985-98-98",
			City:    "Deafult City",
			IsChain: true,
		},
	}
}

func (b *GymBuilder) SetID(id uuid.UUID) *GymBuilder {
	b.gym.ID = id
	return b
}

func (b *GymBuilder) SetName(name string) *GymBuilder {
	b.gym.Name = name
	return b
}

func (b *GymBuilder) SetAddress(address string) *GymBuilder {
	b.gym.Addres = address
	return b
}

func (b *GymBuilder) SetPhone(phone string) *GymBuilder {
	b.gym.Phone = phone
	return b
}

func (b *GymBuilder) SetCity(city string) *GymBuilder {
	b.gym.Phone = city
	return b
}

func (b *GymBuilder) SetIsChain(isChain bool) *GymBuilder {
	b.gym.IsChain = isChain
	return b
}

func (b *GymBuilder) Build() entity.Gym {
	return b.gym
}

func (b *GymBuilder) Invalid() *GymBuilder {
	b.gym.Name = ""         
	b.gym.Addres = ""       
	b.gym.Phone = "invalid" 
	b.gym.City = ""         
	b.gym.IsChain = false
	return b
}


