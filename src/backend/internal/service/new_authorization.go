package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	pb "github.com/sachatarba/course-db/internal/api/grpc"
	"github.com/sachatarba/course-db/internal/entity"
)

type AuthorizationNewService struct {
	// Repo
	sessionRepo ISessionRepository
	codeRepo    ICodeRepository

	smtpService ISmtpService

	// Services
	clientService pb.ClientServiceClient
}

func NewAuthorizationNewService(sessionRepo ISessionRepository,
	clientService pb.ClientServiceClient,
	smtpService ISmtpService, codeRepo ICodeRepository) IAuthorizationService {
	return &AuthorizationNewService{
		sessionRepo:   sessionRepo,
		clientService: clientService,
		smtpService:   smtpService,
		codeRepo:      codeRepo,
	}
}

func (a *AuthorizationNewService) IsAuthorize(ctx context.Context, sessionID uuid.UUID) (*entity.Session, error) {
	session, err := a.sessionRepo.GetSessionBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (a *AuthorizationNewService) Authorize(ctx context.Context, login string, password string) (entity.Session, error) {
	client, err := a.checkAuth(ctx, login, password)
	if err != nil {
		log.Println("Error", login, password, err)
		return entity.Session{}, err
	}

	err = a.sendAndSaveCode(ctx, client.ID, client.Email, "[Качалка] Код подтверждения авторизации на сайте")
	// log.Println("Error: ", login, password, err)
	if err != nil {
		return entity.Session{}, err
	}

	return entity.Session{}, nil
}

func (a *AuthorizationNewService) ChangePassword(ctx context.Context, login string, newPassword string, code string) error {
	client, err := a.clientService.GetClientByLogin(ctx, &pb.LoginRequest{Login: login})
	if err != nil {
		return err
	}

	codeFound, err := a.codeRepo.GetCodeByClientID(ctx, uuid.MustParse(client.Client.Id))
	if err != nil {
		return fmt.Errorf("can't get code: %w", err)
	}

	if code != codeFound {
		return fmt.Errorf("invalid code")
	}

	birthdate, err := time.Parse(time.RFC3339, client.GetClient().Birthdate)
	if err != nil {
		return fmt.Errorf("can't parse birthdate client: %w", err)
	}

	client.GetClient().Password = newPassword
	client.GetClient().Birthdate = birthdate.Format(time.DateOnly)

	_, err = a.clientService.ChangeClient(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("can't change client: %w", err)
	}

	return nil
}

func (a *AuthorizationNewService) Confirm2FA(ctx context.Context, clientID uuid.UUID, code string) (entity.Session, error) {
	codeFound, err := a.codeRepo.GetCodeByClientID(ctx, clientID)
	if err != nil {
		return entity.Session{}, fmt.Errorf("can't get code: %w", err)
	}

	if code != codeFound {
		return entity.Session{}, fmt.Errorf("invalid code")
	}

	session, err := a.createSession(ctx, clientID)
	if err != nil {
		return entity.Session{}, err
	}

	return session, nil
}

func (a *AuthorizationNewService) Register(ctx context.Context, client entity.Client) (entity.Session, error) {
	_, err := a.clientService.RegisterNewClient(ctx, &pb.Client{
		Id:        client.ID.String(),
		Fullname:  client.Fullname,
		Login:     client.Login,
		Email:     client.Email,
		Phone:     client.Phone,
		Birthdate: client.Birthdate,
		Password:  client.Password,
	})
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

func (a *AuthorizationNewService) Logout(ctx context.Context, sessionID uuid.UUID) (entity.Session, error) {
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

func (a *AuthorizationNewService) DeleteClient(ctx context.Context, clientID uuid.UUID) (entity.Session, error) {
	_, err := a.clientService.DeleteClient(ctx, &pb.UUID{Value: clientID.String()})
	if err != nil {
		return entity.Session{}, err
	}

	session := entity.Session{}

	return session, nil
}

func generate2FACode() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b)[:6], nil
}

func (a *AuthorizationNewService) checkAuth(ctx context.Context, login string, password string) (entity.Client, error) {
	client, err := a.clientService.GetClientByLogin(ctx, &pb.LoginRequest{Login: login})
	if err != nil {
		return entity.Client{}, err
	}

	if client.GetClient().Password != password {
		return entity.Client{}, ErrWrongPassword
	}

	return entity.Client{
		ID:        uuid.MustParse(client.Client.Id),
		Fullname:  client.Client.Fullname,
		Login:     client.Client.Login,
		Email:     client.Client.Email,
		Phone:     client.Client.Phone,
		Birthdate: client.Client.Birthdate,
		Password:  client.Client.Password,
	}, nil
}

func (a *AuthorizationNewService) sendAndSaveCode(ctx context.Context, id uuid.UUID, email string, subject string) error {
	code, err := generate2FACode()
	if err != nil {
		return fmt.Errorf("can't generate 2fa code: %w", err)
	}
	err = a.codeRepo.SaveCode(ctx, code, id)
	if err != nil {
		return fmt.Errorf("can't save code %w", err)
	}

	index := strings.IndexAny(email, "@")
	message := fmt.Sprintf("Привет, %s!\r\nВ ваш аккаунт был выполнен вход.\r\nДля подтверждения авториазции введите код: %s\r\nЕсли это были не вы, проигнорируйте это письмо.", email[:index], code)
	err = a.smtpService.SendMail(message, email, subject)
	if err != nil {
		return fmt.Errorf("can't send email: %w", err)
	} 

	return nil
}

func (a *AuthorizationNewService) createSession(ctx context.Context, clientID uuid.UUID) (entity.Session, error) {
	session := entity.Session{
		ClientID:  clientID,
		SessionID: uuid.New(),
		TTL:       time.Now().Add(10 * time.Hour),
	}

	err := a.sessionRepo.CreateNewSession(ctx, session)
	if err != nil {
		return entity.Session{}, err
	}

	return session, nil
}
