package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

type AuthorizationService struct {
	// Repo
	sessionRepo ISessionRepository

	// Services
	clientService IClientService
}

func NewAuthorizationService(sessionRepo ISessionRepository, clientService IClientService) IAuthorizationService {
	return &AuthorizationService{
		sessionRepo:   sessionRepo,
		clientService: clientService,
	}
}

func (a *AuthorizationService) IsAuthorize(ctx context.Context, sessionID uuid.UUID) (*entity.Session, error) {
	session, err := a.sessionRepo.GetSessionBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (a *AuthorizationService) Authorize(ctx context.Context, login string, password string) (entity.Session, error) {
	client, err := a.clientService.GetClientByLogin(ctx, login)
	if err != nil {
		return entity.Session{}, err
	}

	if client.Password != password {
		return entity.Session{}, ErrWrongPassword
	}

	session := entity.Session{
		ClientID:  client.ID,
		SessionID: uuid.New(),
		TTL:       time.Now().Add(10 * time.Hour),
	}

	err = a.sessionRepo.CreateNewSession(ctx, session)
	if err != nil {
		return entity.Session{}, err
	}

	return session, nil
}

func (a *AuthorizationService) Register(ctx context.Context, client entity.Client) (entity.Session, error) {
	err := a.clientService.RegisterNewClient(ctx, client)
	if err != nil {
		return entity.Session{}, err
	}

	session := entity.Session{
		ClientID:  client.ID,
		SessionID: uuid.New(),
		TTL:       time.Now().Add(10 * time.Hour),
	}

	err = a.sessionRepo.CreateNewSession(ctx, session)
		if err != nil {
		return entity.Session{}, err
	}

	return session, nil
}

func (a *AuthorizationService) Logout(ctx context.Context, sessionID uuid.UUID) (entity.Session, error) {
	session, err := a.sessionRepo.GetSessionBySessionID(ctx, sessionID)
	if err != nil {
		return entity.Session{}, err
	}

	err = a.sessionRepo.DeleteSession(ctx, sessionID)
	if err != nil {
		return entity.Session{}, err
	}

	session = entity.Session{
		ClientID:  session.ClientID,
		SessionID: session.SessionID,
		TTL:       time.Now().Add(-10 * time.Hour),
	}

	return session, nil
}


func (a *AuthorizationService) DeleteClient(ctx context.Context, clientID uuid.UUID) (entity.Session, error) {
	err := a.clientService.DeleteClient(ctx, clientID)
	if err != nil {
		return entity.Session{}, err
	}

	session := entity.Session{}

	return session, nil
}
