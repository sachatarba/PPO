package builder

import (
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

type EquipmentBuilder struct {
	ID          uuid.UUID
	Name        string
	Description string
	GymID       uuid.UUID
}

func NewEquipmentBuilder() *EquipmentBuilder {
	return &EquipmentBuilder{
		ID:          uuid.New(),
		Name:        "Default Equipment Name",
		Description: "Default Description",
		GymID:       uuid.New(),
	}
}

func (b *EquipmentBuilder) SetID(id uuid.UUID) *EquipmentBuilder {
	b.ID = id
	return b
}

func (b *EquipmentBuilder) SetName(name string) *EquipmentBuilder {
	b.Name = name
	return b
}

func (b *EquipmentBuilder) SetDescription(description string) *EquipmentBuilder {
	b.Description = description
	return b
}

func (b *EquipmentBuilder) SetGymID(gymID uuid.UUID) *EquipmentBuilder {
	b.GymID = gymID
	return b
}

func (b *EquipmentBuilder) Invalid() *EquipmentBuilder {
	b.Name = ""
	b.Description = ""
	b.GymID = uuid.Nil
	return b
}

func (b *EquipmentBuilder) Build() entity.Equipment {
	return entity.Equipment{
		ID:          b.ID,
		Name:        b.Name,
		Description: b.Description,
		GymID:       b.GymID,
	}
}
