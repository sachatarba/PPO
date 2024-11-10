package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/service/mocks"
	"github.com/sachatarba/course-db/internal/utils/builder"
	"github.com/stretchr/testify/mock"
)

type AuthServiceSuite struct {
	suite.Suite
}

func (s *AuthServiceSuite) TestAuthServiceRegister(t provider.T) {
	t.Title("[Register] Successfuly registered")
	t.Tags("auth", "service", "register")
	t.Parallel()
	t.WithNewStep("Correct: successfully registered", func(sCtx provider.StepCtx) {
		sessionRepoMock := &mocks.ISessionRepository{}
		clientServiceMock := &mocks.IClientService{}
		authService := NewAuthorizationService(sessionRepoMock, clientServiceMock)

		ctx := context.TODO()
		client := builder.NewClientBuilder().Build()

		sessionRepoMock.On("CreateNewSession", ctx, mock.Anything).Return(nil)
		clientServiceMock.On("RegisterNewClient", ctx, client).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "client", client)

		session, err := authService.Register(ctx, client)

		sCtx.Assert().Equal(client.ID, session.ClientID)
		sCtx.Assert().NoError(err)
	})
}

func (s *AuthServiceSuite) TestAuthServiceErrorRegister(t provider.T) {
	t.Title("[Register] registration failed")
	t.Tags("auth", "service", "register")
	t.Parallel()
	t.WithNewStep("Correct: registration falied", func(sCtx provider.StepCtx) {
		sessionRepoMock := &mocks.ISessionRepository{}
		clientServiceMock := &mocks.IClientService{}
		authService := NewAuthorizationService(sessionRepoMock, clientServiceMock)

		ctx := context.TODO()
		client := builder.NewClientBuilder().Build()

		sessionRepoMock.On("CreateNewSession", ctx, mock.Anything).Return(errors.New("can't create sessions"))
		clientServiceMock.On("RegisterNewClient", ctx, client).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "client", client)

		session, err := authService.Register(ctx, client)

		sCtx.Assert().Equal(entity.Session{}, session)
		sCtx.Assert().Error(err)
	})
}

func (s *AuthServiceSuite) TestAuthServiceLogout(t provider.T) {
	t.Title("[Register] logout ok")
	t.Tags("auth", "service", "logout")
	t.Parallel()
	t.WithNewStep("Correct: successfully logout", func(sCtx provider.StepCtx) {
		sessionRepoMock := &mocks.ISessionRepository{}
		clientServiceMock := &mocks.IClientService{}
		authService := NewAuthorizationService(sessionRepoMock, clientServiceMock)

		ctx := context.TODO()
		client := builder.NewClientBuilder().Build()
		session := builder.NewSessionBuilder().
			SetSessionID(uuid.New()).
			SetClientID(client.ID).
			SetTTL(time.Now().Add(-10 * time.Hour)).
			Build()

		sessionRepoMock.On("GetSessionBySessionID", ctx, session.SessionID).Return(session, nil)
		sessionRepoMock.On("DeleteSession", ctx, session.SessionID).Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "client", session.SessionID)

		realSession, err := authService.Logout(ctx, session.SessionID)

		sCtx.Assert().Equal(realSession.SessionID, session.SessionID)
		sCtx.Assert().Equal(realSession.ClientID, session.ClientID)
		sCtx.Assert().True(realSession.TTL.Before(time.Now()))
		sCtx.Assert().NoError(err)
	})
}

func (s *AuthServiceSuite) TestAuthServiceErrorLogout(t provider.T) {
	t.Title("[Register] logout failed")
	t.Tags("auth", "service", "logout")
	t.Parallel()
	t.WithNewStep("Correct: logout failed", func(sCtx provider.StepCtx) {
		sessionRepoMock := &mocks.ISessionRepository{}
		clientServiceMock := &mocks.IClientService{}
		authService := NewAuthorizationService(sessionRepoMock, clientServiceMock)

		ctx := context.TODO()
		client := builder.NewClientBuilder().Build()
		session := builder.NewSessionBuilder().
			SetSessionID(uuid.New()).
			SetClientID(client.ID).
			SetTTL(time.Now().Add(-10 * time.Hour)).
			Build()

		sessionRepoMock.On("GetSessionBySessionID", ctx, session.SessionID).
			Return(entity.Session{}, errors.New("can't get session"))

		sessionRepoMock.On("DeleteSession", ctx, session.SessionID).
			Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "client", session.SessionID)

		realSession, err := authService.Logout(ctx, session.SessionID)

		sCtx.Assert().Equal(entity.Session{}, realSession)
		sCtx.Assert().Error(err)
	})
}

func (s *AuthServiceSuite) TestAuthServiceDeleteClient(t provider.T) {
	t.Title("[Register] delete client ok")
	t.Tags("auth", "service", "delete client")
	t.Parallel()
	t.WithNewStep("Correct: delete client ok", func(sCtx provider.StepCtx) {
		sessionRepoMock := &mocks.ISessionRepository{}
		clientServiceMock := &mocks.IClientService{}
		authService := NewAuthorizationService(sessionRepoMock, clientServiceMock)

		ctx := context.TODO()
		id := uuid.New()

		clientServiceMock.On("DeleteClient", ctx, id).
			Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "id", id)

		realSession, err := authService.DeleteClient(ctx, id)

		sCtx.Assert().Equal(entity.Session{}, realSession)
		sCtx.Assert().NoError(err)
	})
}

func (s *AuthServiceSuite) TestAuthServiceErrorDeleteClient(t provider.T) {
	t.Title("[Register] delete client failed")
	t.Tags("auth", "service", "delete client")
	t.Parallel()
	t.WithNewStep("Correct: delete client failed", func(sCtx provider.StepCtx) {
		sessionRepoMock := &mocks.ISessionRepository{}
		clientServiceMock := &mocks.IClientService{}
		authService := NewAuthorizationService(sessionRepoMock, clientServiceMock)

		ctx := context.TODO()
		id := uuid.New()

		clientServiceMock.On("DeleteClient", ctx, id).
			Return(errors.New("can't delete session"))

		sCtx.WithNewParameters("ctx", ctx, "id", id)

		realSession, err := authService.DeleteClient(ctx, id)

		sCtx.Assert().Equal(entity.Session{}, realSession)
		sCtx.Assert().Error(err)
	})
}

func (s *AuthServiceSuite) TestAuthorizationServiceAuthorize(t provider.T) {
	t.Title("[Register] authorize ok")
	t.Tags("auth", "service", "authorize")
	t.Parallel()
	t.WithNewStep("Correct: authorize", func(sCtx provider.StepCtx) {
		sessionRepoMock := &mocks.ISessionRepository{}
		clientServiceMock := &mocks.IClientService{}
		authService := NewAuthorizationService(sessionRepoMock, clientServiceMock)

		ctx := context.TODO()
		client := builder.NewClientBuilder().Build()
		session := builder.NewSessionBuilder().
			SetClientID(client.ID).
			Build()

		clientServiceMock.On("GetClientByLogin", ctx, client.Login).
			Return(client, nil)
		sessionRepoMock.On("CreateNewSession", ctx, mock.Anything).
			Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "login", client.Login, "password", client.Password)

		realSession, err := authService.Authorize(ctx, client.Login, client.Password)

		sCtx.Assert().Equal(session.ClientID, realSession.ClientID)
		sCtx.Assert().NoError(err)
	})
}

func (s *AuthServiceSuite) TestAuthorizationServiceWrongPasswordAuthorize(t provider.T) {
	t.Title("[Register] authorize ok")
	t.Tags("auth", "service", "authorize")
	t.Parallel()
	t.WithNewStep("Correct: authorize", func(sCtx provider.StepCtx) {
		sessionRepoMock := &mocks.ISessionRepository{}
		clientServiceMock := &mocks.IClientService{}
		authService := NewAuthorizationService(sessionRepoMock, clientServiceMock)

		ctx := context.TODO()
		client := builder.NewClientBuilder().Build()

		clientServiceMock.On("GetClientByLogin", ctx, client.Login).
			Return(client, nil)
		sessionRepoMock.On("CreateNewSession", ctx, mock.Anything).
			Return(nil)

		sCtx.WithNewParameters("ctx", ctx, "login", client.Login, "password", client.Password+"bebra")

		realSession, err := authService.Authorize(ctx, client.Login, client.Password+"berba")

		sCtx.Assert().Equal(entity.Session{}, realSession)
		sCtx.Assert().ErrorIs(err, ErrWrongPassword)
	})
}

func (s *AuthServiceSuite) TestAuthorizationServiceIsAuthorize(t provider.T) {
	t.Title("[Register] is authorize")
	t.Tags("auth", "service", "authorize")
	t.Parallel()
	t.WithNewStep("Correct: authorize", func(sCtx provider.StepCtx) {
		sessionRepoMock := &mocks.ISessionRepository{}
		clientServiceMock := &mocks.IClientService{}
		authService := NewAuthorizationService(sessionRepoMock, clientServiceMock)

		ctx := context.TODO()
		id := uuid.New()
		session := builder.NewSessionBuilder().Build()

		sessionRepoMock.On("GetSessionBySessionID", ctx, id).
			Return(session, nil)

		sCtx.WithNewParameters("ctx", ctx, "id", id)

		realSession, err := authService.IsAuthorize(ctx, id)

		sCtx.Assert().Equal(session, *realSession)
		sCtx.Assert().NoError(err)
	})
}

func (s *AuthServiceSuite) TestAuthorizationServiceErorIsAuthorize(t provider.T) {
	t.Title("[Register] is authorize")
	t.Tags("auth", "service", "authorize")
	t.Parallel()
	t.WithNewStep("Correct: is authorize failed", func(sCtx provider.StepCtx) {
		sessionRepoMock := &mocks.ISessionRepository{}
		clientServiceMock := &mocks.IClientService{}
		authService := NewAuthorizationService(sessionRepoMock, clientServiceMock)

		ctx := context.TODO()
		id := uuid.New()

		sessionRepoMock.On("GetSessionBySessionID", ctx, id).
			Return(entity.Session{}, ErrSessionNotFound)

		sCtx.WithNewParameters("ctx", ctx, "id", id)

		realSession, err := authService.IsAuthorize(ctx, id)

		sCtx.Assert().Nil(realSession)
		sCtx.Assert().ErrorIs(err, ErrSessionNotFound)
	})
}

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(AuthServiceSuite))
}
