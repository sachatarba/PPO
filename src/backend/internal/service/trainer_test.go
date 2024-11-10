package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/service/mocks"
	"github.com/sachatarba/course-db/internal/utils/builder"
)

type TrainerServiceSuite struct {
	suite.Suite
}

func (s *TrainerServiceSuite) TestRegisterNewTrainer(t provider.T) {
	t.Title("[RegisterNewTrainer] Successfully registers new trainer")
	t.Tags("trainer", "service", "register")
	t.Parallel()

	t.WithNewStep("Correct: successfully registered new trainer", func(sCtx provider.StepCtx) {
		trainerRepoMock := &mocks.ITrainerRepository{}
		trainerService := &TrainerService{trainerRepoMock}

		ctx := context.TODO()
		trainer := builder.NewTrainerBuilder().Build()

		trainerRepoMock.On("RegisterNewTrainer", ctx, trainer).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "trainer", trainer)

		err := trainerService.RegisterNewTrainer(ctx, trainer)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		trainerRepoMock := &mocks.ITrainerRepository{}
		trainerService := &TrainerService{trainerRepoMock}

		ctx := context.TODO()
		invalidTrainer := builder.NewTrainerBuilder().Invalid().Build()

		err := trainerService.RegisterNewTrainer(ctx, invalidTrainer)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *TrainerServiceSuite) TestChangeTrainer(t provider.T) {
	t.Title("[ChangeTrainer] Successfully changes trainer")
	t.Tags("trainer", "service", "change")
	t.Parallel()

	t.WithNewStep("Correct: successfully changed trainer", func(sCtx provider.StepCtx) {
		trainerRepoMock := &mocks.ITrainerRepository{}
		trainerService := &TrainerService{trainerRepoMock}

		ctx := context.TODO()
		trainer := builder.NewTrainerBuilder().Build()

		trainerRepoMock.On("ChangeTrainer", ctx, trainer).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "trainer", trainer)

		err := trainerService.ChangeTrainer(ctx, trainer)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		trainerRepoMock := &mocks.ITrainerRepository{}
		trainerService := &TrainerService{trainerRepoMock}

		ctx := context.TODO()
		invalidTrainer := builder.NewTrainerBuilder().Invalid().Build()

		err := trainerService.ChangeTrainer(ctx, invalidTrainer)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *TrainerServiceSuite) TestDeleteTrainer(t provider.T) {
	t.Title("[DeleteTrainer] Successfully deletes trainer")
	t.Tags("trainer", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted trainer", func(sCtx provider.StepCtx) {
		trainerRepoMock := &mocks.ITrainerRepository{}
		trainerService := &TrainerService{trainerRepoMock}

		ctx := context.TODO()
		trainerID := uuid.New()

		trainerRepoMock.On("DeleteTrainer", ctx, trainerID).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "trainerID", trainerID)

		err := trainerService.DeleteTrainer(ctx, trainerID)

		sCtx.Assert().NoError(err)
	})
}

func (s *TrainerServiceSuite) TestGetTrainerByID(t provider.T) {
	t.Title("[GetTrainerByID] Successfully retrieves trainer by ID")
	t.Tags("trainer", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved trainer", func(sCtx provider.StepCtx) {
		trainerRepoMock := &mocks.ITrainerRepository{}
		trainerService := &TrainerService{trainerRepoMock}

		ctx := context.TODO()
		trainerID := uuid.New()
		expectedTrainer := builder.NewTrainerBuilder().Build()

		trainerRepoMock.On("GetTrainerByID", ctx, trainerID).Return(expectedTrainer, nil)

		sCtx.WithNewParameters("ctx", ctx, "trainerID", trainerID)

		trainer, err := trainerService.GetTrainerByID(ctx, trainerID)

		sCtx.Assert().Equal(expectedTrainer, trainer)
		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: trainer not found", func(sCtx provider.StepCtx) {
		trainerRepoMock := &mocks.ITrainerRepository{}
		trainerService := &TrainerService{trainerRepoMock}

		ctx := context.TODO()
		trainerID := uuid.New()

		trainerRepoMock.On("GetTrainerByID", ctx, trainerID).Return(entity.Trainer{}, errors.New("trainer not found"))

		sCtx.WithNewParameters("ctx", ctx, "trainerID", trainerID)

		trainer, err := trainerService.GetTrainerByID(ctx, trainerID)

		sCtx.Assert().Equal(entity.Trainer{}, trainer)
		sCtx.Assert().Error(err)
	})
}

func (s *TrainerServiceSuite) TestListTrainers(t provider.T) {
	t.Title("[ListTrainers] Successfully lists trainers")
	t.Tags("trainer", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully listed trainers", func(sCtx provider.StepCtx) {
		trainerRepoMock := &mocks.ITrainerRepository{}
		trainerService := &TrainerService{trainerRepoMock}

		ctx := context.TODO()
		expectedTrainers := []entity.Trainer{
			builder.NewTrainerBuilder().Build(),
			builder.NewTrainerBuilder().Build(),
		}

		trainerRepoMock.On("ListTrainers", ctx).Return(expectedTrainers, nil)

		sCtx.WithNewParameters("ctx", ctx)

		trainers, err := trainerService.ListTrainers(ctx)

		sCtx.Assert().Equal(expectedTrainers, trainers)
		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: failed to list trainers", func(sCtx provider.StepCtx) {
		trainerRepoMock := &mocks.ITrainerRepository{}
		trainerService := &TrainerService{trainerRepoMock}

		ctx := context.TODO()

		trainerRepoMock.On("ListTrainers", ctx).Return([]entity.Trainer{}, errors.New("failed to list trainers"))

		sCtx.WithNewParameters("ctx", ctx)

		trainers, err := trainerService.ListTrainers(ctx)

		sCtx.Assert().Equal([]entity.Trainer{}, trainers)
		sCtx.Assert().Error(err)
	})
}

func TestTrainerServiceSuite(t *testing.T) {
	suite.RunSuite(t, new(TrainerServiceSuite))
}
