package repository

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/entity"
)

type SessionRepoSuite struct {
	mock        redismock.ClientMock
	sessionRepo ISessionRepository
}

func (c *SessionRepoSuite) BeforeEach(t provider.T) {
	t.Title("Init mock redis")

	client, mock := redismock.NewClientMock()
	c.sessionRepo = NewSessionRepo(client)
	c.mock = mock

	t.Tags("fixture", "sessions", "redis")
}

func (s *SessionRepoSuite) TestCreateNewSession(t provider.T) {
	t.Title("Create a new session")
	t.Tags("repository", "redis")

	t.WithNewStep("Create a new session and add to Redis", func(sCtx provider.StepCtx) {
		sessionID := uuid.New()
		clientID := uuid.New()
		session := entity.Session{
			SessionID: sessionID,
			ClientID:  clientID,
			TTL:       time.Now().Add(10 * time.Hour),
		}

		s.mock.ExpectSet("session:"+sessionID.String(), []byte(""), time.Until(session.TTL)).SetVal("OK")
		s.mock.ExpectSAdd("client_session:"+clientID.String(), session.SessionID.String()).SetVal(1)
		s.mock.ExpectExpireAt("client_session:"+clientID.String(), session.TTL).SetVal(true)

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx)

		err := s.sessionRepo.CreateNewSession(ctx, session)

		sCtx.Assert().NoError(err)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *SessionRepoSuite) TestDeleteSession(t provider.T) {
	t.Title("Delete all sessions for a client")
	t.Tags("repository", "redis")

	t.WithNewStep("Delete sessions for a client", func(sCtx provider.StepCtx) {
		clientID := uuid.New()
		sessionIDs := []string{"sessionID1", "sessionID2"}

		s.mock.ExpectSMembers("client_sessions:" + clientID.String()).SetVal(sessionIDs)
		s.mock.ExpectDel("session:" + sessionIDs[0]).SetVal(1)
		s.mock.ExpectDel("session:" + sessionIDs[1]).SetVal(1)
		s.mock.ExpectDel("client_sessions:" + clientID.String()).SetVal(1)

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx)

		err := s.sessionRepo.DeleteSession(ctx, clientID)

		sCtx.Assert().NoError(err)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *SessionRepoSuite) TestDeleteSessionBySessionID(t provider.T) {
	t.Title("Delete session by session ID")
	t.Tags("repository", "redis")

	t.WithNewStep("Delete session by session ID", func(sCtx provider.StepCtx) {
		sessionID := uuid.New()
		clientID := uuid.New()
		session := entity.Session{
			SessionID: sessionID,
			ClientID:  clientID,
			TTL:       time.Now().Add(10 * time.Hour),
		}

		sessionData, _ := json.Marshal(session)
		s.mock.ExpectGet("session:" + sessionID.String()).SetVal(string(sessionData))
		s.mock.ExpectSRem("client_sessions:"+clientID.String(), sessionID.String()).SetVal(1)
		s.mock.ExpectDel("session:" + sessionID.String()).SetVal(1)

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx)

		err := s.sessionRepo.DeleteSessionBySessionID(ctx, sessionID)

		sCtx.Assert().NoError(err)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *SessionRepoSuite) TestGetSessionsByClientID(t provider.T) {
	t.Title("Get all sessions by client ID")
	t.Tags("repository", "redis")

	t.WithNewStep("Get sessions for a client", func(sCtx provider.StepCtx) {
		clientID := uuid.New()
		sessionID1 := uuid.New()
		sessionID2 := uuid.New()
		session := entity.Session{
			SessionID: sessionID1,
			ClientID:  clientID,
			TTL:       time.Now().Add(10 * time.Hour),
		}

		sessionData1, _ := json.Marshal(session)

		s.mock.ExpectSMembers("client_sessions:" + clientID.String()).SetVal([]string{sessionID1.String(), sessionID2.String()})
		s.mock.ExpectGet("session:" + sessionID1.String()).SetVal(string(sessionData1))
		s.mock.ExpectGet("session:" + sessionID2.String()).SetVal(string(sessionData1))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx)

		sessions, err := s.sessionRepo.GetSessionsByClientID(ctx, clientID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Len(sessions, 2)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *SessionRepoSuite) TestGetSessionBySessionID(t provider.T) {
	t.Title("Get session by session ID")
	t.Tags("repository", "redis")

	t.WithNewStep("Get session by session ID", func(sCtx provider.StepCtx) {
		sessionID := uuid.New()
		clientID := uuid.New()
		session := entity.Session{
			SessionID: sessionID,
			ClientID:  clientID,
			TTL:       time.Now().Add(10 * time.Hour),
		}

		sessionData, _ := json.Marshal(session)
		s.mock.ExpectGet("session:" + sessionID.String()).SetVal(string(sessionData))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx)

		resultSession, err := s.sessionRepo.GetSessionBySessionID(ctx, sessionID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(sessionID, resultSession.SessionID)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func TestSessionRepo(t *testing.T) {
	suite.RunSuite(t, new(ScheduleRepoSuite))
}
