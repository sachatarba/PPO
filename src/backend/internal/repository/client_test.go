package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/service"
	"github.com/sachatarba/course-db/internal/utils/builder"
	"github.com/sachatarba/course-db/internal/repository/mocks"
)

type ClientRepoSuite struct {
	suite.Suite

	mock       sqlmock.Sqlmock
	clientRepo service.IClientRepository
}

func (c *ClientRepoSuite) BeforeEach(t provider.T) {
	t.Title("Init mock db")
	db, mock := mocks.NewMockDB()

	c.clientRepo = NewClientRepo(db)
	c.mock = mock

	t.Tags("fixture", "client", "db")
}

func (c *ClientRepoSuite) TestRegisterNewClient(t provider.T) {
	t.Title("[Register] Register new client")
	t.Tags("repository", "postgres")

	t.WithNewStep("Register new client", func(sCtx provider.StepCtx) {
		client := builder.NewClientBuilder().Build()

		c.mock.ExpectBegin()
		c.mock.ExpectExec(`^INSERT INTO "clients" (.+) VALUES (.+)$`).
			WithArgs(client.ID, client.Login, client.Password, client.Fullname, client.Email, client.Phone, client.Birthdate).
			WillReturnResult(sqlmock.NewResult(1, 1))
		c.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "client", client)

		err := c.clientRepo.RegisterNewClient(ctx, client)

		sCtx.Assert().NoError(err)
		err = c.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (c *ClientRepoSuite) TestChangeClient(t provider.T) {
	t.Title("[Change] Change client information")
	t.Tags("repository", "postgres")

	t.WithNewStep("Change client information", func(sCtx provider.StepCtx) {
		client := builder.NewClientBuilder().Build()

		c.mock.ExpectBegin()
		c.mock.ExpectExec(`^UPDATE "clients" SET (.+)$`).
			// WithArgs(client.Login, client.Password, client.Fullname, client.Email, client.Phone, client.Birthdate, client.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		c.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "client", client)

		err := c.clientRepo.ChangeClient(ctx, client)

		sCtx.Assert().NoError(err)
		err = c.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (c *ClientRepoSuite) TestDeleteClient(t provider.T) {
	t.Title("[Delete] Delete client by ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("Delete client by ID", func(sCtx provider.StepCtx) {
		clientID := uuid.New()

		c.mock.ExpectBegin()
		c.mock.ExpectExec(`^DELETE FROM "clients" WHERE (.+)$`).
			WithArgs(clientID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		c.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "clientID", clientID)

		err := c.clientRepo.DeleteClient(ctx, clientID)

		sCtx.Assert().NoError(err)
		err = c.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (c *ClientRepoSuite) TestGetClientByID(t provider.T) {
	t.Title("[Get] Get client by ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("Get client by ID", func(sCtx provider.StepCtx) {
		client := builder.NewClientBuilder().Build()

		c.mock.ExpectQuery(`^SELECT (.+) FROM "clients" WHERE "clients"."id" = \$1 ORDER BY "clients"."id" LIMIT \$2$`).
			WithArgs(client.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "fullname", "email", "phone", "birthdate"}).
				AddRow(client.ID, client.Login, client.Password, client.Fullname, client.Email, client.Phone, client.Birthdate))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "clientID", client.ID)

		result, err := c.clientRepo.GetClientByID(ctx, client.ID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(client.ID, result.ID)
		err = c.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (c *ClientRepoSuite) TestGetClientByLogin(t provider.T) {
	t.Title("[Get] Get client by login")
	t.Tags("repository", "postgres")

	t.WithNewStep("Get client by login", func(sCtx provider.StepCtx) {
		client := builder.NewClientBuilder().Build()

		c.mock.ExpectQuery(`^SELECT (.+) FROM "clients" WHERE login = \$1 (.+) LIMIT \$2`).
			WithArgs(client.Login, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "fullname", "email", "phone", "birthdate"}).
				AddRow(client.ID, client.Login, client.Password, client.Fullname, client.Email, client.Phone, client.Birthdate))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "login", client.Login)

		result, err := c.clientRepo.GetClientByLogin(ctx, client.Login)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(client.Login, result.Login)
		err = c.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (c *ClientRepoSuite) TestListClients(t provider.T) {
	t.Title("[List] List all clients")
	t.Tags("repository", "postgres")

	t.WithNewStep("List all clients", func(sCtx provider.StepCtx) {
		client1 := builder.NewClientBuilder().Build()
		client2 := builder.NewClientBuilder().Build()

		c.mock.ExpectQuery(`^SELECT (.+) FROM "clients"$`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "fullname", "email", "phone", "birthdate"}).
				AddRow(client1.ID, client1.Login, client1.Password, client1.Fullname, client1.Email, client1.Phone, client1.Birthdate).
				AddRow(client2.ID, client2.Login, client2.Password, client2.Fullname, client2.Email, client2.Phone, client2.Birthdate))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx)

		result, err := c.clientRepo.ListClients(ctx)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(2, len(result))
		sCtx.Assert().Equal(client1.ID, result[0].ID)
		sCtx.Assert().Equal(client2.ID, result[1].ID)

		err = c.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func TestClientSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ClientRepoSuite))
}
