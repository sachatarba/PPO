package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/service"
)

type ISessionRepository interface {
	CreateNewSession(ctx context.Context, session entity.Session) error
	DeleteSession(ctx context.Context, clientID uuid.UUID) error
	DeleteSessionBySessionID(ctx context.Context, sessionID uuid.UUID) error
	GetSessionsByClientID(ctx context.Context, clientID uuid.UUID) ([]entity.Session, error)
	GetSessionBySessionID(ctx context.Context, sessionID uuid.UUID) (entity.Session, error)
}

type SessionRepo struct {
	client *redis.Client
}

func NewSessionRepo(client *redis.Client) ISessionRepository {
	return &SessionRepo{client: client}
}

func (r *SessionRepo) CreateNewSession(ctx context.Context, session entity.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	sessionKey := "session:" + session.SessionID.String()
	err = r.client.Set(ctx, sessionKey, data, time.Until(session.TTL)).Err()
	if err != nil {
		return err
	}

	clientKey := "client_session:" + session.ClientID.String()
	err = r.client.SAdd(ctx, clientKey, session.SessionID.String()).Err()
	if err != nil {
		return err
	}

	r.client.ExpireAt(ctx, clientKey, session.TTL)

	return nil
}

func (r *SessionRepo) DeleteSession(ctx context.Context, clientID uuid.UUID) error {
	clientKey := "client_sessions:" + clientID.String()
	sessionIDs, err := r.client.SMembers(ctx, clientKey).Result()
	if err != nil {
		return err
	}

	for _, sessionID := range sessionIDs {
		r.client.Del(ctx, "session:"+sessionID)
	}

	err = r.client.Del(ctx, clientKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepo) DeleteSessionBySessionID(ctx context.Context, sessionID uuid.UUID) error {
	sessionKey := "session:" + sessionID.String()
	sessionData, err := r.client.Get(ctx, sessionKey).Result()
	if err != nil {
		return err
	}

	var session entity.Session
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		return err
	}

	clientKey := "client_sessions:" + session.ClientID.String()
	err = r.client.SRem(ctx, clientKey, sessionID.String()).Err()
	if err != nil {
		return err
	}

	err = r.client.Del(ctx, sessionKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepo) GetSessionsByClientID(ctx context.Context, clientID uuid.UUID) ([]entity.Session, error) {
	clientKey := "client_sessions:" + clientID.String()
	sessionIDs, err := r.client.SMembers(ctx, clientKey).Result()
	if err != nil {
		return nil, err
	}

	var sessions []entity.Session
	for _, sessionID := range sessionIDs {
		sessionData, err := r.client.Get(ctx, "session:"+sessionID).Result()
		if err != nil {
			return nil, err
		}

		var session entity.Session
		err = json.Unmarshal([]byte(sessionData), &session)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *SessionRepo) GetSessionBySessionID(ctx context.Context, sessionID uuid.UUID) (entity.Session, error) {
	sessionData, err := r.client.Get(ctx, "session:"+sessionID.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			errors.Join(service.ErrSessionNotFound)
		}

		return entity.Session{}, err
	}

	var session entity.Session
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		return entity.Session{}, err
	}

	return session, nil
}
