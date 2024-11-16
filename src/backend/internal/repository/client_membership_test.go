package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/utils/builder"
	"github.com/sachatarba/course-db/internal/utils/fabric"
	"github.com/sachatarba/course-db/internal/repository/mocks"
	"github.com/sachatarba/course-db/internal/service"
)

type ClientMembershipRepoSuite struct {
	suite.Suite

	mock                 sqlmock.Sqlmock
	clientMembershipRepo service.IClientMembershipsRepository
}

func (c *ClientMembershipRepoSuite) BeforeAll(t provider.T) {
	t.Title("Init mock db")
	db, mock := mocks.NewMockDB()

	c.clientMembershipRepo = NewClientMembershipRepo(db)
	c.mock = mock

	t.Tags("fixture", "client membership", "db")
}

func (c *ClientMembershipRepoSuite) TestCreateNewClientMembership(t provider.T) {
	t.Title("[Create] Create client membership")
	t.Tags("repository", "postgres")

	clientMembership := fabric.DeafaultClientMembership()

	c.mock.ExpectBegin()
	c.mock.ExpectExec(`INSERT INTO "client_memberships"`).
		WithArgs(
			clientMembership.ID,
			clientMembership.StartDate,
			clientMembership.EndDate,
			clientMembership.MembershipType.ID,
			clientMembership.ClientID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	c.mock.ExpectCommit()

	t.WithNewStep("Create client membership", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		sCtx.WithNewParameters("ctx", ctx, "client membership", clientMembership)

		err := c.clientMembershipRepo.CreateNewClientMembership(ctx, clientMembership)

		sCtx.Assert().NoError(err)
	})
}

func (c *ClientMembershipRepoSuite) TestChangeClientMembership(t provider.T) {
	t.Title("[Update] Change client membership")
	t.Tags("repository", "postgres")

	t.WithNewStep("Change client membership", func(sCtx provider.StepCtx) {
		clientMembership := builder.NewClientMembershipBuilder().
			SetStartDate(time.Now().Format(time.DateOnly)).
			SetEndDate(time.Now().Format(time.DateOnly)).
			SetID(uuid.New()).
			Build()

		c.mock.ExpectBegin()

		c.mock.ExpectExec(`UPDATE "client_memberships" SET`).
			WithArgs(
				clientMembership.StartDate,
				clientMembership.EndDate,
				clientMembership.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		c.mock.ExpectCommit()
		ctx := context.TODO()

		sCtx.WithNewParameters("ctx", ctx, "client membership", clientMembership)

		err := c.clientMembershipRepo.ChangeClientMembership(ctx, clientMembership)

		sCtx.Assert().NoError(err)
	})
}

func (c *ClientMembershipRepoSuite) TestDeleteClientMembership(t provider.T) {
	t.Title("[Delete] Delete client membership")
	t.Tags("repository", "postgres")

	t.WithNewStep("Delete client membership", func(sCtx provider.StepCtx) {
		clientMembershipID := uuid.New()

		c.mock.ExpectBegin()
		c.mock.ExpectExec(`DELETE FROM "client_memberships"`).
			WithArgs(clientMembershipID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		c.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "clientMembershipID", clientMembershipID)

		err := c.clientMembershipRepo.DeleteClientMembership(ctx, clientMembershipID)

		sCtx.Assert().NoError(err)

		err = c.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (c *ClientMembershipRepoSuite) TestGetClientMembershipByID(t provider.T) {
	t.Title("[Get] Get client membership by ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("Get client membership by ID", func(sCtx provider.StepCtx) {
		clientMembershipID := uuid.New()
		clientMembership := builder.NewClientMembershipBuilder().
			SetMembershipType(builder.NewMembershipTypeBuilder().Build()).
			Build()

		c.mock.ExpectQuery(`^SELECT (.+) FROM "client_memberships" WHERE (.+)$`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "start_date", "end_date", "membership_type_id"}).
				AddRow(clientMembership.ID, clientMembership.StartDate, clientMembership.EndDate, clientMembership.MembershipType.ID))

		c.mock.ExpectQuery(`^SELECT (.+) FROM "membership_types" WHERE (.+)$`).
			WithArgs(clientMembership.MembershipType.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "type", "description", "price", "days_duration", "gym_id"}).
				AddRow(clientMembership.MembershipType.ID,
					clientMembership.MembershipType.Type,
					clientMembership.MembershipType.Description,
					clientMembership.MembershipType.Price,
					clientMembership.MembershipType.DaysDuration,
					clientMembership.MembershipType.GymID))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "clientMembershipID", clientMembershipID)

		result, err := c.clientMembershipRepo.GetClientMembershipByID(ctx, clientMembershipID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(clientMembership.ID, result.ID)
		sCtx.Assert().Equal(clientMembership.StartDate, result.StartDate)
		sCtx.Assert().Equal(clientMembership.EndDate, result.EndDate)
		sCtx.Assert().Equal(clientMembership.MembershipType.ID, result.MembershipType.ID)

		err = c.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (c *ClientMembershipRepoSuite) TestListClientMembershipsByClientID(t provider.T) {
	t.Title("[List] List client memberships by client ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("List client memberships by client ID", func(sCtx provider.StepCtx) {
		clientID := uuid.New()
		membershipTypeID := uuid.New()

		clientMemberships := []entity.ClientMembership{
			builder.NewClientMembershipBuilder().
				SetClientID(clientID).
				SetMembershipType(
					builder.NewMembershipTypeBuilder().
						SetID(membershipTypeID).
						Build(),
				).
				Build(),
			builder.NewClientMembershipBuilder().
				SetClientID(clientID).
				SetMembershipType(
					builder.NewMembershipTypeBuilder().
						SetID(membershipTypeID).
						Build(),
				).
				Build(),
		}

		c.mock.ExpectQuery(`^SELECT (.+) FROM "client_memberships" WHERE (.+)$`).
			WithArgs(clientID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "start_date", "end_date", "membership_type_id"}).
				AddRow(clientMemberships[0].ID,
					clientMemberships[0].StartDate,
					clientMemberships[0].EndDate,
					clientMemberships[0].MembershipType.ID).
				AddRow(clientMemberships[1].ID,
					clientMemberships[1].StartDate,
					clientMemberships[1].EndDate,
					clientMemberships[1].MembershipType.ID))

		c.mock.ExpectQuery(`^SELECT (.+) FROM "membership_types" WHERE (.+)$`).
			WithArgs(membershipTypeID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "type", "description", "price", "days_duration", "gym_id"}).
				AddRow(clientMemberships[0].MembershipType.ID,
					clientMemberships[0].MembershipType.Type,
					clientMemberships[0].MembershipType.Description,
					clientMemberships[0].MembershipType.Price,
					clientMemberships[0].MembershipType.DaysDuration,
					clientMemberships[0].MembershipType.GymID))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "clientID", clientID)

		result, err := c.clientMembershipRepo.ListClientMembershipsByClientID(ctx, clientID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(len(clientMemberships), len(result))
		for i, membership := range clientMemberships {
			sCtx.Assert().Equal(membership.ID, result[i].ID)
			sCtx.Assert().Equal(membership.StartDate, result[i].StartDate)
			sCtx.Assert().Equal(membership.EndDate, result[i].EndDate)
			sCtx.Assert().Equal(membership.MembershipType.ID, result[i].MembershipType.ID)
		}

		err = c.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ClientMembershipRepoSuite))
}
