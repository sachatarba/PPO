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

type TrainingServiceSuite struct {
	suite.Suite
}

var ErrDatabase = errors.New("error: database")

func (s *TrainingServiceSuite) TestCreateNewTraining(t provider.T) {
	t.Title("[CreateNewTraining] Successfully creates new training")
	t.Tags("training", "service", "create")
	t.Parallel()

	t.WithNewStep("Correct: successfully created new training", func(sCtx provider.StepCtx) {
		trainingRepoMock := &mocks.ITrainingRepository{}
		trainingService := &TrainingService{trainingRepoMock}

		ctx := context.TODO()
		training := builder.NewTrainingBuilder().Build()

		trainingRepoMock.On("CreateNewTraining", ctx, training).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "training", training)

		err := trainingService.CreateNewTraining(ctx, training)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		trainingRepoMock := &mocks.ITrainingRepository{}
		trainingService := &TrainingService{trainingRepoMock}

		ctx := context.TODO()
		invalidTraining := builder.NewTrainingBuilder().Invalid().Build()

		err := trainingService.CreateNewTraining(ctx, invalidTraining)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *TrainingServiceSuite) TestChangeTraining(t provider.T) {
	t.Title("[ChangeTraining] Successfully changes training")
	t.Tags("training", "service", "change")
	t.Parallel()

	t.WithNewStep("Correct: successfully changed training", func(sCtx provider.StepCtx) {
		trainingRepoMock := &mocks.ITrainingRepository{}
		trainingService := &TrainingService{trainingRepoMock}

		ctx := context.TODO()
		training := builder.NewTrainingBuilder().Build()

		trainingRepoMock.On("ChangeTraining", ctx, training).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "training", training)

		err := trainingService.ChangeTraining(ctx, training)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		trainingRepoMock := &mocks.ITrainingRepository{}
		trainingService := &TrainingService{trainingRepoMock}

		ctx := context.TODO()
		invalidTraining := builder.NewTrainingBuilder().Invalid().Build()

		err := trainingService.ChangeTraining(ctx, invalidTraining)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *TrainingServiceSuite) TestDeleteTraining(t provider.T) {
	t.Title("[DeleteTraining] Successfully deletes training")
	t.Tags("training", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted training", func(sCtx provider.StepCtx) {
		trainingRepoMock := &mocks.ITrainingRepository{}
		trainingService := &TrainingService{trainingRepoMock}

		ctx := context.TODO()
		trainingID := uuid.New()

		trainingRepoMock.On("DeleteTraining", ctx, trainingID).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "trainingID", trainingID)

		err := trainingService.DeleteTraining(ctx, trainingID)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: delete training fails", func(sCtx provider.StepCtx) {
		trainingRepoMock := &mocks.ITrainingRepository{}
		trainingService := &TrainingService{trainingRepoMock}

		ctx := context.TODO()
		trainingID := uuid.New()

		trainingRepoMock.On("DeleteTraining", ctx, trainingID).Return(ErrNotFound)

		sCtx.WithNewParameters("ctx", ctx, "trainingID", trainingID)

		err := trainingService.DeleteTraining(ctx, trainingID)

		sCtx.Assert().Equal(ErrNotFound, err)
	})
}

func (s *TrainingServiceSuite) TestListTrainingsByTrainerID(t provider.T) {
	t.Title("[ListTrainingsByTrainerID] Successfully lists trainings by trainer ID")
	t.Tags("training", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully lists trainings by trainer ID", func(sCtx provider.StepCtx) {
		trainingRepoMock := &mocks.ITrainingRepository{}
		trainingService := &TrainingService{trainingRepoMock}

		ctx := context.TODO()
		trainerID := uuid.New()
		trainings := []entity.Training{
			{ID: uuid.New(), TrainerID: trainerID},
		}

		trainingRepoMock.On("ListTrainingsByTrainerID", ctx, trainerID).Return(trainings, nil)

		sCtx.WithNewParameters("ctx", ctx, "trainerID", trainerID)

		result, err := trainingService.ListTrainingsByTrainerID(ctx, trainerID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(len(trainings), len(result))
	})

	t.WithNewStep("Incorrect: repository call fails", func(sCtx provider.StepCtx) {
		trainingRepoMock := &mocks.ITrainingRepository{}
		trainingService := &TrainingService{trainingRepoMock}

		ctx := context.TODO()
		trainerID := uuid.New()

		trainingRepoMock.On("ListTrainingsByTrainerID", ctx, trainerID).Return(nil, ErrDatabase)

		sCtx.WithNewParameters("ctx", ctx, "trainerID", trainerID)

		result, err := trainingService.ListTrainingsByTrainerID(ctx, trainerID)

		sCtx.Assert().Equal(ErrDatabase, err)
		sCtx.Assert().Empty(result)
	})
}

func TestTrainingServiceSuite(t *testing.T) {
	suite.RunSuite(t, new(TrainingServiceSuite))
}
