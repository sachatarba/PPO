package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/repository/mocks"
	"github.com/sachatarba/course-db/internal/service"
	"github.com/sachatarba/course-db/internal/utils/builder"
)

type EquipmentRepoSuite struct {
	suite.Suite

	mock          sqlmock.Sqlmock
	equipmentRepo service.IEquipmentRepository
}

func (c *EquipmentRepoSuite) BeforeEach(t provider.T) {
	t.Title("Init mock db")
	db, mock := mocks.NewMockDB()

	c.equipmentRepo = NewEquipmentRepo(db)
	c.mock = mock

	t.Tags("fixture", "equipment", "db")
}

func (s *EquipmentRepoSuite) TestCreateNewEquipment(t provider.T) {
	t.Title("[Create] Create new equipment")
	t.Tags("repository", "postgres")

	t.WithNewStep("Create new equipment", func(sCtx provider.StepCtx) {
		equipment := builder.NewEquipmentBuilder().
			SetGymID(uuid.New()).
			Build()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`INSERT INTO "equipment" (.+) VALUES (.+)`).
			WithArgs(
				equipment.ID,
				equipment.Name,
				equipment.Description,
				equipment.GymID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "equipment", equipment)

		err := s.equipmentRepo.CreateNewEquipment(ctx, equipment)

		sCtx.Assert().NoError(err)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *EquipmentRepoSuite) TestChangeEquipment(t provider.T) {
	t.Title("[Change] Change equipment")
	t.Tags("repository", "postgres")

	t.WithNewStep("Update equipment details", func(sCtx provider.StepCtx) {
		equipment := builder.NewEquipmentBuilder().
			SetID(uuid.New()).
			Build()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`UPDATE "equipment" SET "name"=\$1,"description"=\$2,"gym_id"=\$3 WHERE "id" = \$4`).
			WithArgs(
				equipment.Name,
				equipment.Description,
				equipment.GymID,
				equipment.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "equipment", equipment)

		err := s.equipmentRepo.ChangeEquipment(ctx, equipment)

		sCtx.Assert().NoError(err)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *EquipmentRepoSuite) TestDeleteEquipment(t provider.T) {
	t.Title("[Delete] Delete equipment by ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("Delete equipment by ID", func(sCtx provider.StepCtx) {
		equipmentID := uuid.New()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`DELETE FROM "equipment" WHERE "equipment"."id" = \$1`).
			WithArgs(equipmentID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "equipmentID", equipmentID)

		err := s.equipmentRepo.DeleteEquipment(ctx, equipmentID)

		sCtx.Assert().NoError(err)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *EquipmentRepoSuite) TestGetEquipmentByID(t provider.T) {
	t.Title("[Get] Get equipment by ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("Retrieve equipment by ID", func(sCtx provider.StepCtx) {
		equipment := builder.NewEquipmentBuilder().
			SetID(uuid.New()).
			Build()

		s.mock.ExpectQuery(`^SELECT \* FROM "equipment" WHERE "equipment"."id" = \$1 ORDER BY "equipment"."id" LIMIT \$2`).
			WithArgs(equipment.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "gym_id"}).
				AddRow(equipment.ID, equipment.Name, equipment.Description, equipment.GymID))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "equipmentID", equipment.ID)

		result, err := s.equipmentRepo.GetEquipmentByID(ctx, equipment.ID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(equipment.ID, result.ID)
		sCtx.Assert().Equal(equipment.Name, result.Name)
		sCtx.Assert().Equal(equipment.Description, result.Description)
		sCtx.Assert().Equal(equipment.GymID, result.GymID)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *EquipmentRepoSuite) TestListEquipmentsByGymID(t provider.T) {
	t.Title("[List] List equipments by gym ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("List all equipments for a specific gym", func(sCtx provider.StepCtx) {
		gymID := uuid.New()
		equipment := builder.NewEquipmentBuilder().
			SetGymID(gymID).
			Build()

		s.mock.ExpectQuery(`^SELECT \* FROM "equipment" WHERE "equipment"."gym_id" = \$1`).
			WithArgs(gymID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "gym_id"}).
				AddRow(equipment.ID, equipment.Name, equipment.Description, equipment.GymID))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "gymID", gymID)

		results, err := s.equipmentRepo.ListEquipmentsByGymID(ctx, gymID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Len(results, 1)
		sCtx.Assert().Equal(equipment.ID, results[0].ID)
		sCtx.Assert().Equal(equipment.Name, results[0].Name)
		sCtx.Assert().Equal(equipment.Description, results[0].Description)
		sCtx.Assert().Equal(equipment.GymID, results[0].GymID)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func TestEquipmentSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(EquipmentRepoSuite))
}
