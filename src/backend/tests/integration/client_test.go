package integration

import (
	"context"
	"errors"
	"os"
	"slices"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/config"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/orm"
	postrgres_adapter "github.com/sachatarba/course-db/internal/postrgres"
	"github.com/sachatarba/course-db/internal/repository"
	"github.com/sachatarba/course-db/internal/service"
	"github.com/sachatarba/course-db/internal/utils/builder"
	"gorm.io/gorm"
)

type ClientServiceSuite struct {
	suite.Suite

	clientService service.IClientService
	db            *gorm.DB
}

func ClientEqual(sCtx provider.StepCtx, expected, actual entity.Client) {
	sCtx.Assert().Equal(expected.ID, actual.ID)
	sCtx.Assert().Equal(expected.Login, actual.Login)
	sCtx.Assert().Equal(expected.Password, actual.Password)
	sCtx.Assert().Equal(expected.Fullname, actual.Fullname)
	sCtx.Assert().Equal(expected.Email, actual.Email)
	sCtx.Assert().Equal(expected.Phone, actual.Phone)
	// sCtx.Assert().Equal(expected.Birthdate, actual.Birthdate)
}

func (s *ClientServiceSuite) BeforeAll(t provider.T) {
	conf := config.NewConfFromEnv()

	postgresConnector := postrgres_adapter.PostgresConnector{
		Conf: conf.PostgresConf,
	}

	db, err := postgresConnector.Connect()
	t.Assert().NoError(err, "Error connection db")

	postgresMigrator := postrgres_adapter.PostgresMigrator{
		DB:     db,
		Tables: orm.TablesORM,
	}

	err = postgresMigrator.Migrate()
	t.Assert().NoError(err, "Error migration db")

	s.db = db
	repo := repository.NewClientRepo(db)

	s.clientService = service.NewClientService(repo)
}

func (s *ClientServiceSuite) AfterAll(t provider.T) {
	tables, err := s.db.Migrator().GetTables()
	t.Assert().NoError(err)

	for _, table := range tables {
		err := s.db.Migrator().DropTable(table)
		t.Assert().NoError(err)
	}
}

func (s *ClientServiceSuite) TestRegisterNewClient(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[RegisterNewClient] Successfully registered a new client")
	t.Tags("client_service", "service", "register")
	t.Parallel()

	t.WithNewStep("Correct: successfully registered new client", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		client := builder.NewClientBuilder().Build()

		sCtx.WithNewParameters("ctx", ctx, "client", client)

		// Вызов метода
		err := s.clientService.RegisterNewClient(ctx, client)

		// Проверка
		sCtx.Assert().NoError(err)
		actualOrm := orm.Client{ID: client.ID}

		err = s.db.First(&actualOrm).Error
		sCtx.Assert().NoError(err)
		actual := orm.NewClientConverter().ConvertToEntity(actualOrm)

		ClientEqual(sCtx, client, actual)

		// Удаление тестовых данных
		err = s.db.Delete(&orm.Client{ID: client.ID}).Error
		sCtx.Assert().NoError(err)
	})
}

func (s *ClientServiceSuite) TestChangeClient(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ChangeClient] Successfully changed client data")
	t.Tags("client_service", "service", "update")
	t.Parallel()

	t.WithNewStep("Correct: successfully updated client", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		client := builder.NewClientBuilder().Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)
		err := s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		// Изменяем данные клиента
		client.Fullname = "Updated Name"

		// Вызов метода
		err = s.clientService.ChangeClient(ctx, client)

		// Проверка
		sCtx.Assert().NoError(err)
		actual := orm.Client{ID: client.ID}
		err = s.db.First(&actual).Error
		sCtx.Assert().NoError(err)
		actualClient := orm.NewClientConverter().ConvertToEntity(actual)

		ClientEqual(sCtx, client, actualClient)

		// Удаление тестовых данных
		s.db.Delete(&orm.Client{ID: client.ID})
	})
}

func (s *ClientServiceSuite) TestDeleteClient(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[DeleteClient] Successfully deleted a client")
	t.Tags("client_service", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted client", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		client := builder.NewClientBuilder().Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)
		err := s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		err = s.clientService.DeleteClient(ctx, client.ID)

		// Проверка
		sCtx.Assert().NoError(err)
		toDelete := &orm.Client{ID: client.ID}
		err = s.db.First(&toDelete).Error
		sCtx.Assert().True(errors.Is(err, gorm.ErrRecordNotFound))
	})
}

func (s *ClientServiceSuite) TestGetClientByID(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[GetClientByID] Successfully retrieved client by ID")
	t.Tags("client_service", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved client by ID", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		client := builder.NewClientBuilder().Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)
		err := s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		actual, err := s.clientService.GetClientByID(ctx, client.ID)

		// Проверка
		sCtx.Assert().NoError(err)
		ClientEqual(sCtx, client, actual)

		// Удаление тестовых данных
		s.db.Delete(&orm.Client{ID: client.ID})
	})
}

func (s *ClientServiceSuite) TestGetClientByLogin(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[GetClientByLogin] Successfully retrieved client by login")
	t.Tags("client_service", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved client by login", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		client := builder.NewClientBuilder().Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)
		err := s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		actual, err := s.clientService.GetClientByLogin(ctx, client.Login)

		// Проверка
		sCtx.Assert().NoError(err)
		ClientEqual(sCtx, client, actual)

		// Удаление тестовых данных
		s.db.Delete(&orm.Client{ID: client.ID})
	})
}

func (s *ClientServiceSuite) TestListClients(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ListClients] Successfully listed all clients")
	t.Tags("client_service", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully listed clients", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		clients := []entity.Client{
			builder.
				NewClientBuilder().
				Build(),
			builder.
				NewClientBuilder().
				Build(),
		}

		for _, client := range clients {
			clientOrm := orm.NewClientConverter().ConvertFromEntity(client)
			err := s.db.Save(&clientOrm).Error
			sCtx.Assert().NoError(err)
		}

		// Вызов метода
		actualClients, err := s.clientService.ListClients(ctx)

		// Проверка
		sCtx.Assert().NoError(err)
		for i, client := range actualClients {
			if slices.ContainsFunc(actualClients, func(c entity.Client) bool {
				return c.ID == client.ID
			}) {
				ClientEqual(sCtx, client, actualClients[i])
			}
		}

		// Удаление тестовых данных

		for _, client := range clients {
			s.db.Delete(&orm.Client{ID: client.ID})
		}
	})
}

func TestClientServiceSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ClientServiceSuite))
}
