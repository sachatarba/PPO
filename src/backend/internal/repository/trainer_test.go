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

type TrainerRepoSuite struct {
	suite.Suite

	mock        sqlmock.Sqlmock
	trainerRepo service.ITrainerRepository
}

func (c *TrainerRepoSuite) BeforeEach(t provider.T) {
	t.Title("Init mock db")
	db, mock := mocks.NewMockDB()

	c.trainerRepo = NewTrainerRepo(db)
	c.mock = mock

	t.Tags("fixture", "trainer", "db")
}

func (s *TrainerRepoSuite) TestRegisterNewTrainer(t provider.T) {
	t.Title("Register new trainer")
	t.Tags("repository", "postgres")

	t.WithNewStep("Add a new trainer and associate with gyms", func(sCtx provider.StepCtx) {
		trainerID := uuid.New()
		gymID := uuid.New()
		trainer := builder.NewTrainerBuilder().
			SetID(trainerID).
			SetFullname("John Doe").
			SetEmail("john@example.com").
			SetPhone("+1-800-555-1234").
			SetQualification("Certified Trainer").
			SetUnitPrice(100).
			SetGymsID([]uuid.UUID{gymID}).
			Build()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`^INSERT INTO "trainers"  \("id","fullname","email","phone","qualification","unit_price"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) ON CONFLICT DO NOTHING$`).
			WithArgs(trainer.ID, trainer.Fullname, trainer.Email, trainer.Phone, trainer.Qualification, trainer.UnitPrice).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// s.mock.ExpectQuery(`^SELECT \* FROM "gyms" (.+)$`).
		// 	// WithArgs(gymID).
		// 	WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(gymID))

		s.mock.ExpectExec(`^INSERT INTO "gym_trainers" (.+)$`).
			WithArgs(gymID, trainerID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "trainer", trainer)

		err := s.trainerRepo.RegisterNewTrainer(ctx, trainer)

		sCtx.Assert().NoError(err)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *TrainerRepoSuite) TestChangeTrainer(t provider.T) {
	t.Title("Change trainer details")
	t.Tags("repository", "postgres")

	t.WithNewStep("Update trainer information", func(sCtx provider.StepCtx) {
		trainer := builder.NewTrainerBuilder().
			SetID(uuid.New()).
			SetFullname("Jane Doe").
			SetEmail("jane@example.com").
			SetPhone("+1-800-555-5678").
			SetQualification("Advanced Trainer").
			SetUnitPrice(150).
			Build()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`^UPDATE "trainers" SET (.+)$`).
			WithArgs(trainer.Fullname, trainer.Email, trainer.Phone, trainer.Qualification, trainer.UnitPrice, trainer.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "trainer", trainer)

		err := s.trainerRepo.ChangeTrainer(ctx, trainer)

		sCtx.Assert().NoError(err)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *TrainerRepoSuite) TestDeleteTrainer(t provider.T) {
	t.Title("Delete trainer")
	t.Tags("repository", "postgres")

	t.WithNewStep("Delete a trainer by ID", func(sCtx provider.StepCtx) {
		trainerID := uuid.New()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`DELETE FROM "trainers" WHERE "trainers"."id" = \$1`).
			WithArgs(trainerID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "trainerID", trainerID)

		err := s.trainerRepo.DeleteTrainer(ctx, trainerID)

		sCtx.Assert().NoError(err)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *TrainerRepoSuite) TestGetTrainerByID(t provider.T) {
	t.Title("Get trainer by ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("Fetch a trainer by their ID", func(sCtx provider.StepCtx) {
		trainerID := uuid.New()
		trainer := builder.NewTrainerBuilder().
			SetID(trainerID).
			SetFullname("John Smith").
			SetEmail("john.smith@example.com").
			SetPhone("+1-800-555-0000").
			SetQualification("Expert").
			SetUnitPrice(200).
			Build()

		s.mock.ExpectQuery(`^SELECT \* FROM "trainers" WHERE "trainers"."id" = \$1 ORDER BY "trainers"."id" LIMIT \$2`).
			WithArgs(trainerID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "fullname", "email", "phone", "qualification", "unit_price"}).
				AddRow(trainer.ID, trainer.Fullname, trainer.Email, trainer.Phone, trainer.Qualification, trainer.UnitPrice))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "trainerID", trainerID)

		result, err := s.trainerRepo.GetTrainerByID(ctx, trainerID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(trainer.ID, result.ID)
		sCtx.Assert().Equal(trainer.Fullname, result.Fullname)
		sCtx.Assert().Equal(trainer.Email, result.Email)
		sCtx.Assert().Equal(trainer.Phone, result.Phone)
		sCtx.Assert().Equal(trainer.Qualification, result.Qualification)
		sCtx.Assert().Equal(trainer.UnitPrice, result.UnitPrice)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *TrainerRepoSuite) TestListTrainers(t provider.T) {
	t.Title("List all trainers")
	t.Tags("repository", "postgres")

	t.WithNewStep("Retrieve all trainers", func(sCtx provider.StepCtx) {
		trainer := builder.NewTrainerBuilder().
			SetID(uuid.New()).
			SetFullname("Alice Johnson").
			SetEmail("alice@example.com").
			SetPhone("+1-800-555-1111").
			SetQualification("Professional").
			SetUnitPrice(120).
			Build()

		s.mock.ExpectQuery(`^SELECT \* FROM "trainers"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "fullname", "email", "phone", "qualification", "unit_price"}).
				AddRow(trainer.ID, trainer.Fullname, trainer.Email, trainer.Phone, trainer.Qualification, trainer.UnitPrice))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx)

		results, err := s.trainerRepo.ListTrainers(ctx)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Len(results, 1)
		sCtx.Assert().Equal(trainer.ID, results[0].ID)
		sCtx.Assert().Equal(trainer.Fullname, results[0].Fullname)
		sCtx.Assert().Equal(trainer.Email, results[0].Email)
		sCtx.Assert().Equal(trainer.Phone, results[0].Phone)
		sCtx.Assert().Equal(trainer.Qualification, results[0].Qualification)
		sCtx.Assert().Equal(trainer.UnitPrice, results[0].UnitPrice)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func TestTrainerSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(TrainerRepoSuite))
}
