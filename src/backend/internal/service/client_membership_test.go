package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/utils/builder"
	"github.com/stretchr/testify/mock"

	"github.com/sachatarba/course-db/internal/service/mocks"
)

type ClientMembershipsServiceSuite struct {
	suite.Suite
}

func (s *ClientMembershipsServiceSuite) TestCreateNewClientMembership(t provider.T) {
	t.Title("[CreateNewClientMembership] Successfully created membership")
	t.Tags("client_membership", "service", "create")
	t.Parallel()

	t.WithNewStep("Correct: successfully created membership", func(sCtx provider.StepCtx) {
		clientMembershipRepoMock := &mocks.IClientMembershipsRepository{}
		clientMembershipService := &ClientMembershipsService{clientMembershipRepoMock}

		ctx := context.TODO()
		clientMembership := builder.NewClientMembershipBuilder().Build()

		clientMembershipRepoMock.On("CreateNewClientMembership", ctx, clientMembership).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "clientMembership", clientMembership)

		err := clientMembershipService.CreateNewClientMembership(ctx, clientMembership)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		clientMembershipRepoMock := &mocks.IClientMembershipsRepository{}
		clientMembershipService := &ClientMembershipsService{clientMembershipRepoMock}

		ctx := context.TODO()
		invalidClientMembership := builder.NewClientMembershipBuilder().Invalid().Build()

		clientMembershipRepoMock.On("CreateNewClientMembership", ctx, mock.Anything).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "invalidClientMembership", invalidClientMembership)

		err := clientMembershipService.CreateNewClientMembership(ctx, invalidClientMembership)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *ClientMembershipsServiceSuite) TestChangeClientMembership(t provider.T) {
	t.Title("[ChangeClientMembership] Successfully changed membership")
	t.Tags("client_membership", "service", "change")
	t.Parallel()

	t.WithNewStep("Correct: successfully changed membership", func(sCtx provider.StepCtx) {
		clientMembershipRepoMock := &mocks.IClientMembershipsRepository{}
		clientMembershipService := &ClientMembershipsService{clientMembershipRepoMock}

		ctx := context.TODO()
		clientMembership := builder.NewClientMembershipBuilder().Build()

		clientMembershipRepoMock.On("ChangeClientMembership", ctx, clientMembership).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "clientMembership", clientMembership)

		err := clientMembershipService.ChangeClientMembership(ctx, clientMembership)

		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: validation failed", func(sCtx provider.StepCtx) {
		clientMembershipRepoMock := &mocks.IClientMembershipsRepository{}
		clientMembershipService := &ClientMembershipsService{clientMembershipRepoMock}

		ctx := context.TODO()
		invalidClientMembership := builder.NewClientMembershipBuilder().Invalid().Build()

		clientMembershipRepoMock.On("ChangeClientMembership", ctx, mock.Anything).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "invalidClientMembership", invalidClientMembership)

		err := clientMembershipService.ChangeClientMembership(ctx, invalidClientMembership)

		sCtx.Assert().Equal(ErrValidation, err)
	})
}

func (s *ClientMembershipsServiceSuite) TestDeleteClientMembership(t provider.T) {
	t.Title("[DeleteClientMembership] Successfully deleted membership")
	t.Tags("client_membership", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted membership", func(sCtx provider.StepCtx) {
		clientMembershipRepoMock := &mocks.IClientMembershipsRepository{}
		clientMembershipService := &ClientMembershipsService{clientMembershipRepoMock}

		ctx := context.TODO()
		clientMembershipID := uuid.New()

		clientMembershipRepoMock.On("DeleteClientMembership", ctx, clientMembershipID).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "clientMembershipID", clientMembershipID)

		err := clientMembershipService.DeleteClientMembership(ctx, clientMembershipID)

		sCtx.Assert().NoError(err)
	})
}

func (s *ClientMembershipsServiceSuite) TestGetClientMembershipByID(t provider.T) {
	t.Title("[GetClientMembershipByID] Successfully retrieved membership by ID")
	t.Tags("client_membership", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved membership", func(sCtx provider.StepCtx) {
		clientMembershipRepoMock := &mocks.IClientMembershipsRepository{}
		clientMembershipService := &ClientMembershipsService{clientMembershipRepoMock}

		ctx := context.TODO()
		clientMembershipID := uuid.New()
		expectedMembership := builder.NewClientMembershipBuilder().Build()

		clientMembershipRepoMock.On("GetClientMembershipByID", ctx, clientMembershipID).Return(expectedMembership, nil)

		sCtx.WithNewParameters("ctx", ctx, "clientMembershipID", clientMembershipID)

		membership, err := clientMembershipService.GetClientMembershipByID(ctx, clientMembershipID)

		sCtx.Assert().Equal(expectedMembership, membership)
		sCtx.Assert().NoError(err)
	})

	t.WithNewStep("Incorrect: membership not found", func(sCtx provider.StepCtx) {
		clientMembershipRepoMock := &mocks.IClientMembershipsRepository{}
		clientMembershipService := &ClientMembershipsService{clientMembershipRepoMock}

		ctx := context.TODO()
		clientMembershipID := uuid.New()

		clientMembershipRepoMock.On("GetClientMembershipByID", ctx, clientMembershipID).Return(entity.ClientMembership{}, errors.New("not found"))

		sCtx.WithNewParameters("ctx", ctx, "clientMembershipID", clientMembershipID)

		membership, err := clientMembershipService.GetClientMembershipByID(ctx, clientMembershipID)

		sCtx.Assert().Equal(entity.ClientMembership{}, membership)
		sCtx.Assert().Error(err)
	})
}

func (s *ClientMembershipsServiceSuite) TestListClientMembershipsByClientID(t provider.T) {
	t.Title("[ListClientMembershipsByClientID] Successfully listed memberships by client ID")
	t.Tags("client_membership", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully listed memberships", func(sCtx provider.StepCtx) {
		clientMembershipRepoMock := &mocks.IClientMembershipsRepository{}
		clientMembershipService := &ClientMembershipsService{clientMembershipRepoMock}

		ctx := context.TODO()
		clientID := uuid.New()
		expectedMemberships := []entity.ClientMembership{
			builder.NewClientMembershipBuilder().Build(),
			builder.NewClientMembershipBuilder().Build(),
		}

		clientMembershipRepoMock.On("ListClientMembershipsByClientID", ctx, clientID).Return(expectedMemberships, nil)

		sCtx.WithNewParameters("ctx", ctx, "clientID", clientID)

		memberships, err := clientMembershipService.ListClientMembershipsByClientID(ctx, clientID)

		sCtx.Assert().Equal(expectedMemberships, memberships)
		sCtx.Assert().NoError(err)
	})
}

func TestClientMembershipsServiceSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ClientMembershipsServiceSuite))
}
