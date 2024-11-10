package builder

import (
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

type TrainerBuilder struct {
	trainer entity.Trainer
}

func NewTrainerBuilder() *TrainerBuilder {
	return &TrainerBuilder{
		trainer: entity.Trainer{
			ID:            uuid.New(),
			Fullname:      "Default Name",
			Email:         "default@example.com",
			Phone:         "+7-999-999-99-99",
			Qualification: "Default Qualification",
			UnitPrice:     1.0,
			GymsID: []uuid.UUID{
				uuid.New(),
			},
		},
	}
}

func (b *TrainerBuilder) SetID(id uuid.UUID) *TrainerBuilder {
	b.trainer.ID = id
	return b
}

func (b *TrainerBuilder) SetFullname(name string) *TrainerBuilder {
	b.trainer.Fullname = name
	return b
}

func (b *TrainerBuilder) SetEmail(email string) *TrainerBuilder {
	b.trainer.Email = email
	return b
}

func (b *TrainerBuilder) SetPhone(phone string) *TrainerBuilder {
	b.trainer.Phone = phone
	return b
}

func (b *TrainerBuilder) SetQualification(qualification string) *TrainerBuilder {
	b.trainer.Qualification = qualification
	return b
}

func (b *TrainerBuilder) SetUnitPrice(price float64) *TrainerBuilder {
	b.trainer.UnitPrice = price
	return b
}

func (b *TrainerBuilder) SetGymsID(gymsID []uuid.UUID) *TrainerBuilder {
	b.trainer.GymsID = gymsID
	return b
}

func (b *TrainerBuilder) SetTrainings(trainings []entity.Training) *TrainerBuilder {
	b.trainer.Trainings = trainings
	return b
}

func (b *TrainerBuilder) Invalid() *TrainerBuilder {
	b.trainer = entity.Trainer{
		ID:            uuid.New(),
		Fullname:      "",
		Email:         "invalidemail",
		Phone:         "+1-000-000-0000",
		Qualification: "",
		UnitPrice:     -1.0,
		GymsID:        []uuid.UUID{},
	}
	return b
}

func (b *TrainerBuilder) Build() entity.Trainer {
	return b.trainer
}
