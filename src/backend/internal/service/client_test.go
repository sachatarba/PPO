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

type ClientServiceSuite struct {
	suite.Suite
}

func (s *ClientServiceSuite) TestRegisterNewClient(t provider.T) {
	t.Title("[RegisterNewClient] Successfully registers new client")
	t.Tags("client", "service", "register")
	t.Parallel()

	t.WithNewStep("Correct: successfully registered new client", func(sCtx provider.StepCtx) {
		clientRepoMock := &mocks.IClientRepository{}
		clientService := &ClientService{clientRepoMock}

		ctx := context.TODO()
		client := builder.NewClientBuilder().Build()

		clientRepoMock.On("RegisterNewClient", ctx, client).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "client", client)

		err := clientService.RegisterNewClient(ctx, client)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		clientRepoMock := &mocks.IClientRepository{}
		clientService := &ClientService{clientRepoMock}

		ctx := context.TODO()
		invalidClient := builder.NewClientBuilder().Invalid().Build()

		err := clientService.RegisterNewClient(ctx, invalidClient)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *ClientServiceSuite) TestChangeClient(t provider.T) {
	t.Title("[ChangeClient] Successfully changes client")
	t.Tags("client", "service", "change")
	t.Parallel()

	t.WithNewStep("Correct: successfully changed client", func(sCtx provider.StepCtx) {
		clientRepoMock := &mocks.IClientRepository{}
		clientService := &ClientService{clientRepoMock}

		ctx := context.TODO()
		client := builder.NewClientBuilder().Build()

		clientRepoMock.On("ChangeClient", ctx, client).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "client", client)

		err := clientService.ChangeClient(ctx, client)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		clientRepoMock := &mocks.IClientRepository{}
		clientService := &ClientService{clientRepoMock}

		ctx := context.TODO()
		invalidClient := builder.NewClientBuilder().Invalid().Build()

		err := clientService.ChangeClient(ctx, invalidClient)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *ClientServiceSuite) TestDeleteClient(t provider.T) {
	t.Title("[DeleteClient] Successfully deletes client")
	t.Tags("client", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted client", func(sCtx provider.StepCtx) {
		clientRepoMock := &mocks.IClientRepository{}
		clientService := &ClientService{clientRepoMock}

		ctx := context.TODO()
		clientID := uuid.New()

		clientRepoMock.On("DeleteClient", ctx, clientID).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "clientID", clientID)

		err := clientService.DeleteClient(ctx, clientID)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: client not found", func(sCtx provider.StepCtx) {
		clientRepoMock := &mocks.IClientRepository{}
		clientService := &ClientService{clientRepoMock}

		ctx := context.TODO()
		clientID := uuid.New()

		clientRepoMock.On("DeleteClient", ctx, clientID).Return(errors.New("not found"))

		err := clientService.DeleteClient(ctx, clientID)

		sCtx.Assert().Error(err)
	})
}

func (s *ClientServiceSuite) TestGetClientByID(t provider.T) {
	t.Title("[GetClientByID] Successfully retrieves client by ID")
	t.Tags("client", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved client", func(sCtx provider.StepCtx) {
		clientRepoMock := &mocks.IClientRepository{}
		clientService := &ClientService{clientRepoMock}

		ctx := context.TODO()
		clientID := uuid.New()
		expectedClient := builder.NewClientBuilder().Build()

		clientRepoMock.On("GetClientByID", ctx, clientID).Return(expectedClient, nil)

		sCtx.WithNewParameters("ctx", ctx, "clientID", clientID)

		client, err := clientService.GetClientByID(ctx, clientID)

		sCtx.Assert().Equal(expectedClient, client)
		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: client not found", func(sCtx provider.StepCtx) {
		clientRepoMock := &mocks.IClientRepository{}
		clientService := &ClientService{clientRepoMock}

		ctx := context.TODO()
		clientID := uuid.New()

		clientRepoMock.On("GetClientByID", ctx, clientID).Return(entity.Client{}, errors.New("not found"))

		client, err := clientService.GetClientByID(ctx, clientID)

		sCtx.Assert().Equal(entity.Client{}, client)
		sCtx.Assert().Error(err)
	})
}

func (s *ClientServiceSuite) TestGetClientByLogin(t provider.T) {
	t.Title("[GetClientByLogin] Successfully retrieves client by login")
	t.Tags("client", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved client by login", func(sCtx provider.StepCtx) {
		clientRepoMock := &mocks.IClientRepository{}
		clientService := &ClientService{clientRepoMock}

		ctx := context.TODO()
		login := "testlogin"
		expectedClient := builder.NewClientBuilder().SetLogin(login).Build()

		clientRepoMock.On("GetClientByLogin", ctx, login).Return(expectedClient, nil)

		sCtx.WithNewParameters("ctx", ctx, "login", login)

		client, err := clientService.GetClientByLogin(ctx, login)

		sCtx.Assert().Equal(expectedClient, client)
		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: client not found by login", func(sCtx provider.StepCtx) {
		clientRepoMock := &mocks.IClientRepository{}
		clientService := &ClientService{clientRepoMock}

		ctx := context.TODO()
		login := "nonexistentlogin"

		clientRepoMock.On("GetClientByLogin", ctx, login).Return(entity.Client{}, errors.New("not found"))

		client, err := clientService.GetClientByLogin(ctx, login)

		sCtx.Assert().Equal(entity.Client{}, client)
		sCtx.Assert().Error(err)
	})
}

func (s *ClientServiceSuite) TestListClients(t provider.T) {
	t.Title("[ListClients] Successfully lists clients")
	t.Tags("client", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully listed clients", func(sCtx provider.StepCtx) {
		clientRepoMock := &mocks.IClientRepository{}
		clientService := &ClientService{clientRepoMock}

		ctx := context.TODO()
		expectedClients := []entity.Client{
			builder.NewClientBuilder().Build(),
			builder.NewClientBuilder().Build(),
		}

		clientRepoMock.On("ListClients", ctx).Return(expectedClients, nil)

		sCtx.WithNewParameters("ctx", ctx)

		clients, err := clientService.ListClients(ctx)

		sCtx.Assert().Equal(expectedClients, clients)
		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: error while listing clients", func(sCtx provider.StepCtx) {
		clientRepoMock := &mocks.IClientRepository{}
		clientService := &ClientService{clientRepoMock}

		ctx := context.TODO()

		clientRepoMock.On("ListClients", ctx).Return([]entity.Client{}, errors.New("database error"))

		clients, err := clientService.ListClients(ctx)

		sCtx.Assert().Empty(clients)
		sCtx.Assert().Error(err)
	})
}

func TestClientServiceSuite(t *testing.T) {
	suite.RunSuite(t, new(ClientServiceSuite))
}
