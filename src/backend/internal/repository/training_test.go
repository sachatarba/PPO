package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/utils/builder"
	"github.com/sachatarba/course-db/internal/repository/mocks"
	"github.com/sachatarba/course-db/internal/service"
)

type TrainingRepoSuite struct {
	suite.Suite

	mock         sqlmock.Sqlmock
	trainingRepo service.ITrainingRepository
}

func (c *TrainingRepoSuite) BeforeEach(t provider.T) {
	t.Title("Init mock db")
	db, mock := mocks.NewMockDB()

	c.trainingRepo = NewTrainingRepo(db)
	c.mock = mock

	t.Tags("fixture", "training", "db")
}

func (s *TrainingRepoSuite) TestCreateNewTraining(t provider.T) {
	t.Title("Create new training")
	t.Tags("repository", "postgres")

	t.WithNewStep("Create a new training", func(sCtx provider.StepCtx) {
		trainingID := uuid.New()
		trainerID := uuid.New()
		training := builder.NewTrainingBuilder().
			SetID(trainingID).
			SetTitle("Yoga Class").
			SetDescription("A relaxing yoga session").
			SetTrainingType("flexibility").
			SetTrainerID(trainerID).
			Build()

		s.mock.ExpectBegin()

		s.mock.ExpectExec(`^INSERT INTO "trainings" (.+)$`).
			WillReturnResult(sqlmock.NewResult(1, 1))

		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "training", training)

		err := s.trainingRepo.CreateNewTraining(ctx, training)

		sCtx.Assert().NoError(err)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *TrainingRepoSuite) TestChangeTraining(t provider.T) {
	t.Title("Change training details")
	t.Tags("repository", "postgres")

	t.WithNewStep("Change an existing training", func(sCtx provider.StepCtx) {
		trainingID := uuid.New()
		trainerID := uuid.New()
		training := builder.NewTrainingBuilder().
			SetID(trainingID).
			SetTitle("Advanced Yoga").
			SetDescription("A challenging yoga session").
			SetTrainingType("flexibility").
			SetTrainerID(trainerID).
			Build()

		s.mock.ExpectBegin()

		s.mock.ExpectExec(`^UPDATE "trainings" SET "title"=\$1,"description"=\$2,"training_type"=\$3,"trainer_id"=\$4 WHERE "id" = \$5$`).
			WithArgs(training.Title, training.Description, training.TrainingType, training.TrainerID, training.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "training", training)

		err := s.trainingRepo.ChangeTraining(ctx, training)

		sCtx.Assert().NoError(err)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *TrainingRepoSuite) TestDeleteTraining(t provider.T) {
	t.Title("Delete training")
	t.Tags("repository", "postgres")

	t.WithNewStep("Delete a training by ID", func(sCtx provider.StepCtx) {
		trainingID := uuid.New()

		s.mock.ExpectBegin()

		s.mock.ExpectExec(`^DELETE FROM "trainings" WHERE "trainings"."id" = \$1`).
			WithArgs(trainingID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "trainingID", trainingID)

		err := s.trainingRepo.DeleteTraining(ctx, trainingID)

		sCtx.Assert().NoError(err)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *TrainingRepoSuite) TestListTrainingsByTrainerID(t provider.T) {
	t.Title("List trainings by trainer ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("List all trainings for a specific trainer", func(sCtx provider.StepCtx) {
		trainerID := uuid.New()
		training1 := builder.NewTrainingBuilder().
			SetID(uuid.New()).
			SetTitle("Yoga Basics").
			SetDescription("A basic yoga session").
			SetTrainingType("flexibility").
			SetTrainerID(trainerID).
			Build()
		training2 := builder.NewTrainingBuilder().
			SetID(uuid.New()).
			SetTitle("Strength Training").
			SetDescription("A high-intensity strength session").
			SetTrainingType("strength").
			SetTrainerID(trainerID).
			Build()

		s.mock.ExpectQuery(`^SELECT \* FROM "trainings" WHERE "trainings"."trainer_id" = \$1`).
			WithArgs(trainerID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "training_type", "trainer_id"}).
				AddRow(training1.ID, training1.Title, training1.Description, training1.TrainingType, training1.TrainerID).
				AddRow(training2.ID, training2.Title, training2.Description, training2.TrainingType, training2.TrainerID))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "trainerID", trainerID)

		results, err := s.trainingRepo.ListTrainingsByTrainerID(ctx, trainerID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Len(results, 2)
		sCtx.Assert().Equal(training1.ID, results[0].ID)
		sCtx.Assert().Equal(training2.ID, results[1].ID)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func TestTrainingSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(TrainingRepoSuite))
}
