package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

type TrainingService struct {
	trainingRepo ITrainingRepository
}

func NewTrainingService(repository ITrainingRepository) ITrainingService {
	return &TrainingService{
		trainingRepo: repository,
	}
}

func (s *TrainingService) CreateNewTraining(ctx context.Context, training entity.Training) error {
	if !training.Validate() {
		return ErrValidation
	}

	err := s.trainingRepo.CreateNewTraining(ctx, training)

	return err
}

func (s *TrainingService) ChangeTraining(ctx context.Context, training entity.Training) error {
	if !training.Validate() {
		return ErrValidation
	}

	err := s.trainingRepo.ChangeTraining(ctx, training)

	return err
}

func (s *TrainingService) DeleteTraining(ctx context.Context, trainingID uuid.UUID) error {
	err := s.trainingRepo.DeleteTraining(ctx, trainingID)

	return err
}

func (s *TrainingService) ListTrainingsByTrainerID(ctx context.Context, trainerID uuid.UUID) ([]entity.Training, error) {
	trainings, err := s.trainingRepo.ListTrainingsByTrainerID(ctx, trainerID)
	log.Print("service: ", trainings, err)
	if err != nil {
		return []entity.Training{}, err
	}

	return trainings, nil
}
