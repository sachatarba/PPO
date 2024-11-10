package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

type ScheduleService struct {
	scheduleRepo IScheduleRepository
}

func NewScheduleService(repository IScheduleRepository) IScheduleService {
	return &ScheduleService{
		scheduleRepo: repository,
	}
}

func (s *ScheduleService) CreateNewSchedule(ctx context.Context, schedule entity.Schedule) error {
	if !schedule.Validate() {
		return ErrValidation
	}

	err := s.scheduleRepo.CreateNewSchedule(ctx, schedule)

	return err
}

func (s *ScheduleService) ChangeSchedule(ctx context.Context, schedule entity.Schedule) error {
	if !schedule.Validate() {
		return ErrValidation
	}

	err := s.scheduleRepo.ChangeSchedule(ctx, schedule)

	return err
}

func (s *ScheduleService) DeleteSchedule(ctx context.Context, scheduleID uuid.UUID) error {
	err := s.scheduleRepo.DeleteSchedule(ctx, scheduleID)

	return err
}

func (s *ScheduleService) GetScheduleByID(ctx context.Context, scheduleID uuid.UUID) (entity.Schedule, error) {
	schedule, err := s.scheduleRepo.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return entity.Schedule{}, err
	}

	return schedule, nil
}

func (s *ScheduleService) ListSchedulesByClientID(ctx context.Context, clientID uuid.UUID) ([]entity.Schedule, error) {
	schedules, err := s.scheduleRepo.ListSchedulesByClientID(ctx, clientID)
	if err != nil {
		return []entity.Schedule{}, err
	}

	return schedules, nil
}
