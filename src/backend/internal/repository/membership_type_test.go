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

type MembershipTypeRepoSuite struct {
	suite.Suite

	mock               sqlmock.Sqlmock
	membershipTypeRepo service.IMembershipTypeRepository
}

func (c *MembershipTypeRepoSuite) BeforeEach(t provider.T) {
	t.Title("Init mock db")
	db, mock := mocks.NewMockDB()

	c.membershipTypeRepo = NewMembershipTypeRepo(db)
	c.mock = mock

	t.Tags("fixture", "membershipType", "db")
}

func (s *MembershipTypeRepoSuite) TestRegisterNewMembershipType(t provider.T) {
	t.Title("[Create] Register new membership type")
	t.Tags("repository", "postgres")

	t.WithNewStep("Register new membership type", func(sCtx provider.StepCtx) {
		membershipType := builder.NewMembershipTypeBuilder().
			SetGymID(uuid.New()).
			SetPrice("20.50").
			Build()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`INSERT INTO "membership_types" \("id","type","description","price","days_duration","gym_id"\) VALUES \(.+\)`).
			WithArgs(
				membershipType.ID,
				membershipType.Type,
				membershipType.Description,
				membershipType.Price,
				membershipType.DaysDuration,
				membershipType.GymID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "membershipType", membershipType)

		err := s.membershipTypeRepo.RegisterNewMembershipType(ctx, membershipType)

		sCtx.Assert().NoError(err)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *MembershipTypeRepoSuite) TestChangeMembershipType(t provider.T) {
	t.Title("[Change] Change membership type")
	t.Tags("repository", "postgres")

	t.WithNewStep("Update membership type details", func(sCtx provider.StepCtx) {
		membershipType := builder.NewMembershipTypeBuilder().
			SetID(uuid.New()).
			Build()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`UPDATE "membership_types" SET "type"=\$1,"description"=\$2,"price"=\$3,"days_duration"=\$4,"gym_id"=\$5 WHERE "id" = \$6`).
			WithArgs(
				membershipType.Type,
				membershipType.Description,
				membershipType.Price,
				membershipType.DaysDuration,
				membershipType.GymID,
				membershipType.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "membershipType", membershipType)

		err := s.membershipTypeRepo.ChangeMembershipType(ctx, membershipType)

		sCtx.Assert().NoError(err)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *MembershipTypeRepoSuite) TestDeleteMembershipType(t provider.T) {
	t.Title("[Delete] Delete membership type by ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("Delete membership type by ID", func(sCtx provider.StepCtx) {
		membershipTypeID := uuid.New()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`DELETE FROM "membership_types" WHERE "membership_types"."id" = \$1`).
			WithArgs(membershipTypeID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "membershipTypeID", membershipTypeID)

		err := s.membershipTypeRepo.DeleteMembershipType(ctx, membershipTypeID)

		sCtx.Assert().NoError(err)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *MembershipTypeRepoSuite) TestGetMembershipTypeByID(t provider.T) {
	t.Title("[Get] Get membership type by ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("Retrieve membership type by ID", func(sCtx provider.StepCtx) {
		membershipType := builder.NewMembershipTypeBuilder().
			SetID(uuid.New()).
			Build()

		s.mock.ExpectQuery(`^SELECT \* FROM "membership_types" WHERE "membership_types"."id" = \$1 ORDER BY "membership_types"."id" LIMIT \$2`).
			WithArgs(membershipType.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "type", "description", "price", "days_duration", "gym_id"}).
				AddRow(membershipType.ID, membershipType.Type, membershipType.Description, membershipType.Price, membershipType.DaysDuration, membershipType.GymID))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "membershipTypeID", membershipType.ID)

		result, err := s.membershipTypeRepo.GetMembershipTypeByID(ctx, membershipType.ID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(membershipType.ID, result.ID)
		sCtx.Assert().Equal(membershipType.Type, result.Type)
		sCtx.Assert().Equal(membershipType.Description, result.Description)
		sCtx.Assert().Equal(membershipType.Price, result.Price)
		sCtx.Assert().Equal(membershipType.DaysDuration, result.DaysDuration)
		sCtx.Assert().Equal(membershipType.GymID, result.GymID)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *MembershipTypeRepoSuite) TestListMembershipTypesByGymID(t provider.T) {
	t.Title("[List] List membership types by gym ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("List all membership types for a specific gym", func(sCtx provider.StepCtx) {
		gymID := uuid.New()
		membershipType := builder.NewMembershipTypeBuilder().
			SetGymID(gymID).
			Build()

		s.mock.ExpectQuery(`^SELECT \* FROM "gyms" WHERE "gyms"."id" = \$1 ORDER BY "gyms"."id" LIMIT \$2`).
			WithArgs(gymID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone", "city", "addres", "is_chain"}).
				AddRow(gymID, "Test Gym", "000-000-0000", "Test City", "Test Address", true))

		s.mock.ExpectQuery(`^SELECT \* FROM "membership_types" WHERE "membership_types"."gym_id" = \$1`).
			WithArgs(gymID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "type", "description", "price", "days_duration", "gym_id"}).
				AddRow(membershipType.ID, membershipType.Type, membershipType.Description, membershipType.Price, membershipType.DaysDuration, membershipType.GymID))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "gymID", gymID)

		results, err := s.membershipTypeRepo.ListMembershipTypesByGymID(ctx, gymID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Len(results, 1)
		sCtx.Assert().Equal(membershipType.ID, results[0].ID)
		sCtx.Assert().Equal(membershipType.Type, results[0].Type)
		sCtx.Assert().Equal(membershipType.Description, results[0].Description)
		sCtx.Assert().Equal(membershipType.Price, results[0].Price)
		sCtx.Assert().Equal(membershipType.DaysDuration, results[0].DaysDuration)
		sCtx.Assert().Equal(membershipType.GymID, results[0].GymID)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func TestMembershipTypeSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(MembershipTypeRepoSuite))
}
