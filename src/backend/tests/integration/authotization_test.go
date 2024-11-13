package integration

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/config"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/orm"
	postrgres_adapter "github.com/sachatarba/course-db/internal/postrgres"
	redis_adapter "github.com/sachatarba/course-db/internal/redis"
	"github.com/sachatarba/course-db/internal/repository"
	"github.com/sachatarba/course-db/internal/service"
	"github.com/sachatarba/course-db/internal/utils/builder"
	"gorm.io/gorm"
)

type AuthorizationServiceSuite struct {
	suite.Suite

	// clientService service.IClientService
	authService service.IAuthorizationService

	db     *gorm.DB
	client *redis.Client
	// clients                 []entity.Client
	// clientsMembership       []entity.ClientMembership
	// gyms                    []entity.Gym
	// membershipTypes         []entity.MembershipType
}

func (s *AuthorizationServiceSuite) BeforeAll(t provider.T) {
	conf := config.NewConfFromEnv()

	postgresConnector := postrgres_adapter.PostgresConnector{
		Conf: conf.PostgresConf,
	}

	redisConnector := redis_adapter.RedisConnector{
		Conf: conf.RedisConf,
	}

	s.client = redisConnector.Connect()
	t.Assert().NotNil(s.client, "Error connetcion redis")

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

	clientService := service.NewClientService(repo)
	sessionRepo := repository.NewSessionRepo(s.client)

	s.authService = service.NewAuthorizationService(sessionRepo, clientService)
}

func (s *AuthorizationServiceSuite) AfterAll(t provider.T) {
	tables, err := s.db.Migrator().GetTables()
	t.Assert().NoError(err)

	for _, table := range tables {
		err := s.db.Migrator().DropTable(table)
		t.Assert().NoError(err)
	}
}

// Тест для проверки авторизации с использованием Redis
func (s *AuthorizationServiceSuite) TestIsAuthorize(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[IsAuthorize] Successfully checked if session is authorized")
	t.Tags("authorization_service", "service", "authorize")
	t.Parallel()

	t.WithNewStep("Correct: successfully checked if session is authorized", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём клиента и преобразуем его в orm.Client
		client := builder.NewClientBuilder().Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)

		// Сохраняем клиента в базе данных (передаем по указателю)
		err := s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём сессию напрямую в Redis
		session := entity.Session{
			ClientID:  client.ID,
			SessionID: uuid.New(),
			TTL:       time.Now().Add(10 * time.Hour),
		}
		data, err := json.Marshal(session)
		sCtx.Assert().NoError(err)

		sessionKey := "session:" + session.SessionID.String()
		err = s.client.Set(ctx, sessionKey, data, time.Until(session.TTL)).Err()
		sCtx.Assert().NoError(err)

		clientKey := "client_session:" + session.ClientID.String()
		err = s.client.SAdd(ctx, clientKey, session.SessionID.String()).Err()
		sCtx.Assert().NoError(err)

		// Вызов метода
		actualSession, err := s.authService.IsAuthorize(ctx, session.SessionID)

		// Проверка
		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(session.SessionID, actualSession.SessionID)

		// Удаление тестовых данных
		s.db.Delete(&clientOrm)
		s.client.Del(ctx, "session:"+session.SessionID.String())
		s.client.Del(ctx, "client_session:"+session.ClientID.String())
	})
}

// Тест для авторизации с использованием Redis
func (s *AuthorizationServiceSuite) TestAuthorize(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[Authorize] Successfully authorized client")
	t.Tags("authorization_service", "service", "authorize")
	t.Parallel()

	t.WithNewStep("Correct: successfully authorized client", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём клиента и преобразуем его в orm.Client
		client := builder.NewClientBuilder().Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)

		// Сохраняем клиента в базе данных
		err := s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		// Вызов метода для авторизации
		session, err := s.authService.Authorize(ctx, client.Login, client.Password)

		// Проверка
		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(session.SessionID)

		// Сохраняем сессию напрямую в Redis
		data, err := json.Marshal(session)
		sCtx.Assert().NoError(err)

		sessionKey := "session:" + session.SessionID.String()
		err = s.client.Set(ctx, sessionKey, data, time.Until(session.TTL)).Err()
		sCtx.Assert().NoError(err)

		clientKey := "client_session:" + session.ClientID.String()
		err = s.client.SAdd(ctx, clientKey, session.SessionID.String()).Err()
		sCtx.Assert().NoError(err)

		// Удаление тестовых данных
		s.db.Delete(&clientOrm)
		s.client.Del(ctx, "session:"+session.SessionID.String())
		s.client.Del(ctx, "client_session:"+session.ClientID.String())
	})
}

// Тест для регистрации клиента с использованием Redis
func (s *AuthorizationServiceSuite) TestRegister(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[Register] Successfully registered new client")
	t.Tags("authorization_service", "service", "register")
	t.Parallel()

	t.WithNewStep("Correct: successfully registered new client", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём нового клиента и преобразуем в orm.Client
		client := builder.NewClientBuilder().Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)

		// Вызов метода для регистрации
		session, err := s.authService.Register(ctx, client)

		// Проверка
		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(session.SessionID)

		// Сохраняем сессию напрямую в Redis
		data, err := json.Marshal(session)
		sCtx.Assert().NoError(err)

		sessionKey := "session:" + session.SessionID.String()
		err = s.client.Set(ctx, sessionKey, data, time.Until(session.TTL)).Err()
		sCtx.Assert().NoError(err)

		clientKey := "client_session:" + session.ClientID.String()
		err = s.client.SAdd(ctx, clientKey, session.SessionID.String()).Err()
		sCtx.Assert().NoError(err)

		// Удаление тестовых данных
		s.db.Delete(&clientOrm)
		s.client.Del(ctx, "session:"+session.SessionID.String())
		s.client.Del(ctx, "client_session:"+session.ClientID.String())
	})
}

// Тест для выхода из системы с использованием Redis
func (s *AuthorizationServiceSuite) TestLogout(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[Logout] Successfully logged out client")
	t.Tags("authorization_service", "service", "logout")
	t.Parallel()

	t.WithNewStep("Correct: successfully logged out client", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём клиента и преобразуем его в orm.Client
		client := builder.NewClientBuilder().Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)

		// Сохраняем клиента в базе данных
		err := s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём сессию
		session := entity.Session{
			ClientID:  client.ID,
			SessionID: uuid.New(),
			TTL:       time.Now().Add(10 * time.Hour),
		}
		data, err := json.Marshal(session)
		sCtx.Assert().NoError(err)

		sessionKey := "session:" + session.SessionID.String()
		err = s.client.Set(ctx, sessionKey, data, time.Until(session.TTL)).Err()
		sCtx.Assert().NoError(err)

		clientKey := "client_session:" + session.ClientID.String()
		err = s.client.SAdd(ctx, clientKey, session.SessionID.String()).Err()
		sCtx.Assert().NoError(err)

		// Вызов метода для выхода
		loggedOutSession, err := s.authService.Logout(ctx, session.SessionID)

		// Проверка
		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(session.SessionID, loggedOutSession.SessionID)
		sCtx.Assert().True(loggedOutSession.TTL.Before(time.Now()))

		// Удаление тестовых данных
		s.db.Delete(&clientOrm)
		s.client.Del(ctx, "session:"+session.SessionID.String())
		s.client.Del(ctx, "client_session:"+session.ClientID.String())
	})
}

// Тест для удаления клиента с использованием Redis
func (s *AuthorizationServiceSuite) TestDeleteClient(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	log.Println(os.Getenv("SKIP"))
	t.Title("[DeleteClient] Successfully deleted client")
	t.Tags("authorization_service", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted client", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём клиента и преобразуем в orm.Client
		client := builder.NewClientBuilder().Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)

		// Сохраняем клиента в базе данных
		err := s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём сессию
		session := entity.Session{
			ClientID:  client.ID,
			SessionID: uuid.New(),
			TTL:       time.Now().Add(10 * time.Hour),
		}
		data, err := json.Marshal(session)
		sCtx.Assert().NoError(err)

		sessionKey := "session:" + session.SessionID.String()
		err = s.client.Set(ctx, sessionKey, data, time.Until(session.TTL)).Err()
		sCtx.Assert().NoError(err)

		clientKey := "client_session:" + session.ClientID.String()
		err = s.client.SAdd(ctx, clientKey, session.SessionID.String()).Err()
		sCtx.Assert().NoError(err)

		// Вызов метода для удаления клиента
		session, err = s.authService.DeleteClient(ctx, client.ID)

		// Проверка
		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(session.ClientID, uuid.Nil)

		// Удаление тестовых данных
		s.db.Delete(&clientOrm)
		s.client.Del(ctx, "session:"+session.SessionID.String())
		s.client.Del(ctx, "client_session:"+session.ClientID.String())
	})
}

func TestAuthorizationServiceSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(AuthorizationServiceSuite))
}
