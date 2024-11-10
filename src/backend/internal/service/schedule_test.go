package service

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/service/mocks"
	"github.com/sachatarba/course-db/internal/utils/builder"
)

type ScheduleServiceSuite struct {
	suite.Suite
}

var ErrNotFound = errors.New("bot found")

func (s *ScheduleServiceSuite) TestScheduleServiceCreateNewSchedule(t provider.T) {
	t.Title("[CreateNewSchedule] Successfully creates a new schedule")
	t.Tags("schedule", "service", "create")
	t.Parallel()

	t.WithNewStep("Correct: successfully created new schedule", func(sCtx provider.StepCtx) {
		scheduleRepoMock := &mocks.IScheduleRepository{}
		scheduleService := &ScheduleService{scheduleRepoMock}

		ctx := context.TODO()
		schedule := builder.NewScheduleBuilder().Build()

		validatedSchedule := schedule
		validatedSchedule.Validate()

		scheduleRepoMock.On("CreateNewSchedule", ctx, validatedSchedule).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "schedule", schedule)

		log.Println("bebebab")
		err := scheduleService.CreateNewSchedule(ctx, schedule)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		scheduleRepoMock := &mocks.IScheduleRepository{}
		scheduleService := &ScheduleService{scheduleRepoMock}

		ctx := context.TODO()
		invalidSchedule := builder.NewScheduleBuilder().Invalid().Build()

		err := scheduleService.CreateNewSchedule(ctx, invalidSchedule)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *ScheduleServiceSuite) TestScheduleServiceChangeSchedule(t provider.T) {
	t.Title("[ChangeSchedule] Successfully changes schedule")
	t.Tags("schedule", "service", "change")
	t.Parallel()

	t.WithNewStep("Correct: successfully changed schedule", func(sCtx provider.StepCtx) {
		scheduleRepoMock := &mocks.IScheduleRepository{}
		scheduleService := &ScheduleService{scheduleRepoMock}

		ctx := context.TODO()
		schedule := builder.NewScheduleBuilder().Build()

		validatedSchedule := schedule
		validatedSchedule.Validate()

		scheduleRepoMock.On("ChangeSchedule", ctx, validatedSchedule).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "schedule", schedule)

		err := scheduleService.ChangeSchedule(ctx, schedule)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		scheduleRepoMock := &mocks.IScheduleRepository{}
		scheduleService := &ScheduleService{scheduleRepoMock}

		ctx := context.TODO()
		invalidSchedule := builder.NewScheduleBuilder().Invalid().Build()

		err := scheduleService.ChangeSchedule(ctx, invalidSchedule)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *ScheduleServiceSuite) TestScheduleServiceDeleteSchedule(t provider.T) {
	t.Title("[DeleteSchedule] Successfully deletes a schedule")
	t.Tags("schedule", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted schedule", func(sCtx provider.StepCtx) {
		scheduleRepoMock := &mocks.IScheduleRepository{}
		scheduleService := &ScheduleService{scheduleRepoMock}

		ctx := context.TODO()
		scheduleID := uuid.New()

		scheduleRepoMock.On("DeleteSchedule", ctx, scheduleID).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "scheduleID", scheduleID)

		err := scheduleService.DeleteSchedule(ctx, scheduleID)

		sCtx.Assert().NoError(err)
	})
}

func (s *ScheduleServiceSuite) TestScheduleServiceGetScheduleByID(t provider.T) {
	t.Title("[GetScheduleByID] Successfully retrieves a schedule")
	t.Tags("schedule", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved schedule", func(sCtx provider.StepCtx) {
		scheduleRepoMock := &mocks.IScheduleRepository{}
		scheduleService := &ScheduleService{scheduleRepoMock}

		ctx := context.TODO()
		scheduleID := uuid.New()
		expectedSchedule := builder.NewScheduleBuilder().Build()

		scheduleRepoMock.On("GetScheduleByID", ctx, scheduleID).Return(expectedSchedule, nil)

		sCtx.WithNewParameters("ctx", ctx, "scheduleID", scheduleID)

		schedule, err := scheduleService.GetScheduleByID(ctx, scheduleID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(expectedSchedule, schedule)
	})

	t.WithNewStep("Incorrect: error retrieving schedule", func(sCtx provider.StepCtx) {
		scheduleRepoMock := &mocks.IScheduleRepository{}
		scheduleService := &ScheduleService{scheduleRepoMock}

		ctx := context.TODO()
		scheduleID := uuid.New()

		scheduleRepoMock.On("GetScheduleByID", ctx, scheduleID).Return(entity.Schedule{}, ErrNotFound)

		sCtx.WithNewParameters("ctx", ctx, "scheduleID", scheduleID)

		schedule, err := scheduleService.GetScheduleByID(ctx, scheduleID)

		sCtx.Assert().Equal(ErrNotFound, err)
		sCtx.Assert().Equal(entity.Schedule{}, schedule)
	})
}

func (s *ScheduleServiceSuite) TestScheduleServiceListSchedulesByClientID(t provider.T) {
	t.Title("[ListSchedulesByClientID] Successfully retrieves schedules by client ID")
	t.Tags("schedule", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved schedules", func(sCtx provider.StepCtx) {
		scheduleRepoMock := &mocks.IScheduleRepository{}
		scheduleService := &ScheduleService{scheduleRepoMock}

		ctx := context.TODO()
		clientID := uuid.New()
		expectedSchedules := []entity.Schedule{builder.NewScheduleBuilder().Build()}

		scheduleRepoMock.On("ListSchedulesByClientID", ctx, clientID).Return(expectedSchedules, nil)

		sCtx.WithNewParameters("ctx", ctx, "clientID", clientID)

		schedules, err := scheduleService.ListSchedulesByClientID(ctx, clientID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(expectedSchedules, schedules)
	})

	t.WithNewStep("Incorrect: error retrieving schedules", func(sCtx provider.StepCtx) {
		scheduleRepoMock := &mocks.IScheduleRepository{}
		scheduleService := &ScheduleService{scheduleRepoMock}

		ctx := context.TODO()
		clientID := uuid.New()

		scheduleRepoMock.On("ListSchedulesByClientID", ctx, clientID).Return([]entity.Schedule{}, ErrNotFound)

		sCtx.WithNewParameters("ctx", ctx, "clientID", clientID)

		schedules, err := scheduleService.ListSchedulesByClientID(ctx, clientID)

		sCtx.Assert().Equal(ErrNotFound, err)
		sCtx.Assert().Equal([]entity.Schedule{}, schedules)
	})
}

func TestScheduleServiceSuite(t *testing.T) {
	suite.RunSuite(t, new(ScheduleServiceSuite))
}
