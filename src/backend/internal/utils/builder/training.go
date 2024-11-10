package builder

import (
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

type TrainingBuilder struct {
	training entity.Training
}

func NewTrainingBuilder() *TrainingBuilder {
	return &TrainingBuilder{
		training: entity.Training{
			ID:           uuid.New(),
			Title:        "Default Training",
			Description:  "Default description",
			TrainingType: entity.TrainingType[entity.Aerobic],
			TrainerID:    uuid.New(),
		},
	}
}
func (b *TrainingBuilder) SetID(id uuid.UUID) *TrainingBuilder {
	b.training.ID = id
	return b
}

func (b *TrainingBuilder) SetTitle(title string) *TrainingBuilder {
	b.training.Title = title
	return b
}

func (b *TrainingBuilder) SetDescription(description string) *TrainingBuilder {
	b.training.Description = description
	return b
}

func (b *TrainingBuilder) SetTrainingType(trainingType string) *TrainingBuilder {
	b.training.TrainingType = trainingType
	return b
}

func (b *TrainingBuilder) SetTrainerID(trainerID uuid.UUID) *TrainingBuilder {
	b.training.TrainerID = trainerID
	return b
}

func (b *TrainingBuilder) Invalid() *TrainingBuilder {
	b.training.ID = uuid.Nil
	b.training.Title = ""
	b.training.Description = ""
	b.training.TrainingType = ""
	b.training.TrainerID = uuid.Nil
	return b
}

func (b *TrainingBuilder) Build() entity.Training {
	return b.training
}
