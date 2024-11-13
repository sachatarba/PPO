package integration

import (
	"context"
	"errors"
	"os"
	"strconv"
	"testing"
	"time"

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

const (
	clientsCount               = 10
	gymsCount                  = 10
	clientsMembershipPerClient = 3
	MembershipTypesPerGym      = 3
)

type ClientMembershipsServiceSuite struct {
	suite.Suite

	clientMembershipService service.IClientMembershipsService
	db                      *gorm.DB
	// clients                 []entity.Client
	// clientsMembership       []entity.ClientMembership
	// gyms                    []entity.Gym
	// membershipTypes         []entity.MembershipType
}

func (s *ClientMembershipsServiceSuite) BeforeAll(t provider.T) {
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
	repo := repository.NewClientMembershipRepo(db)

	s.clientMembershipService = service.NewClientMembershipService(repo)
}

func (s *ClientMembershipsServiceSuite) AfterAll(t provider.T) {
	tables, err := s.db.Migrator().GetTables()
	t.Assert().NoError(err)

	for _, table := range tables {
		err := s.db.Migrator().DropTable(table)
		t.Assert().NoError(err)
	}
}

func ClientMembershipEqual(sCtx provider.StepCtx,
	expected entity.ClientMembership,
	actual entity.ClientMembership) {

	// expectedStartDate, err := time.Parse(expected.StartDate, time.DateOnly)
	// sCtx.Assert().NoError(err)

	// expectedEndDate, err := time.Parse(expected.EndDate, time.DateOnly)
	// sCtx.Assert().NoError(err)

	// actualStartDate, err := time.Parse(actual.StartDate, time.RFC3339)
	// sCtx.Assert().NoError(err)

	// actualEndDate, err := time.Parse(actual.EndDate, time.RFC3339)
	// sCtx.Assert().NoError(err)

	expectedPrice, err := strconv.ParseFloat(expected.MembershipType.Price, 64)
	sCtx.Assert().NoError(err)

	actualPrice, err := strconv.ParseFloat(actual.MembershipType.Price, 64)
	sCtx.Assert().NoError(err)

	sCtx.Assert().Equal(expected.ID, actual.ID)
	sCtx.Assert().Equal(expected.ClientID, actual.ClientID)
	// sCtx.Assert().Equal(expectedStartDate, actualStartDate.Format(time.DateOnly))
	// sCtx.Assert().Equal(expectedEndDate, actualEndDate.Format(time.DateOnly))
	sCtx.Assert().Equal(expected.MembershipType.ID, actual.MembershipType.ID)
	sCtx.Assert().Equal(expected.MembershipType.DaysDuration, actual.MembershipType.DaysDuration)
	sCtx.Assert().Equal(expected.MembershipType.Description, actual.MembershipType.Description)
	sCtx.Assert().Equal(expected.MembershipType.GymID, actual.MembershipType.GymID)
	sCtx.Assert().InDelta(expectedPrice, actualPrice, 0.0001)
	sCtx.Assert().Equal(expected.MembershipType.Type, actual.MembershipType.Type)
}

func (s *ClientMembershipsServiceSuite) TestCreateNewClientMembership(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[CreateNewClientMembership] Successfully created membership")
	t.Tags("client_membership", "service", "create")
	t.Parallel()

	t.WithNewStep("Correct: successfully created membership", func(sCtx provider.StepCtx) {
		// Подготовка тестовых данных, загрузка их в бд
		ctx := context.TODO()
		gym := builder.NewGymBuilder().
			Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		membershipType := builder.NewMembershipTypeBuilder().
			SetGymID(gym.ID).
			Build()
		membershipTypeOrm := orm.NewMembershipTypeConverter().ConvertFromEntity(membershipType)
		err = s.db.Save(&membershipTypeOrm).Error
		sCtx.Assert().NoError(err)

		client := builder.NewClientBuilder().
			Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)
		err = s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		clientMembership := builder.
			NewClientMembershipBuilder().
			SetClientID(client.ID).
			SetMembershipType(membershipType).
			Build()

		sCtx.WithNewParameters("ctx", ctx, "clientMembership", clientMembership)

		// Вызов метода
		err = s.clientMembershipService.CreateNewClientMembership(ctx, clientMembership)

		// Проверки
		sCtx.Assert().NoError(err)
		actualOrm := orm.ClientMembership{ID: clientMembership.ID}

		err = s.db.Preload("MembershipType").First(&actualOrm).Error

		sCtx.Assert().NoError(err)
		actual := orm.NewClientMembershipConverter().ConvertToEntity(actualOrm)

		ClientMembershipEqual(sCtx, clientMembership, actual)

		// Удаление тестовых данных
		err = s.db.Delete(&orm.ClientMembership{ID: clientMembership.ID}).Error
		sCtx.Assert().NoError(err)

		err = s.db.Delete(&orm.Client{ID: client.ID}).Error
		sCtx.Assert().NoError(err)

		err = s.db.Delete(&orm.MembershipType{ID: membershipType.ID}).Error
		sCtx.Assert().NoError(err)

		err = s.db.Delete(&orm.Gym{ID: gym.ID}).Error
		sCtx.Assert().NoError(err)
	})
}

func (s *ClientMembershipsServiceSuite) TestChangeClientMembership(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ChangeClientMembership] Successfully changed membership")
	t.Tags("client_membership", "service", "update")
	t.Parallel()

	t.WithNewStep("Correct: successfully changed membership", func(sCtx provider.StepCtx) {
		// Подготовка тестовых данных, загрузка их в бд
		ctx := context.TODO()
		gym := builder.NewGymBuilder().
			Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		membershipType := builder.NewMembershipTypeBuilder().
			SetGymID(gym.ID).
			Build()
		membershipTypeOrm := orm.NewMembershipTypeConverter().ConvertFromEntity(membershipType)
		err = s.db.Save(&membershipTypeOrm).Error
		sCtx.Assert().NoError(err)

		client := builder.NewClientBuilder().
			Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)
		err = s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		clientMembership := builder.
			NewClientMembershipBuilder().
			SetClientID(client.ID).
			SetMembershipType(membershipType).
			Build()
		clientMembershipOrm := orm.NewClientMembershipConverter().ConvertFromEntity(clientMembership)
		err = s.db.Save(&clientMembershipOrm).Error
		sCtx.Assert().NoError(err)

		clientMembership.StartDate = time.Now().Add(time.Hour * 10).Format(time.DateOnly)
		clientMembership.EndDate = time.Now().Add(time.Hour * 124).Format(time.DateOnly)

		// Вызов метода
		err = s.clientMembershipService.ChangeClientMembership(ctx, clientMembership)

		// Проверка
		sCtx.Assert().NoError(err)

		// Проверка изменений
		actual := orm.ClientMembership{ID: clientMembership.ID}
		s.db.Preload("MembershipType").First(&actual)
		ClientMembershipEqual(sCtx, clientMembership, orm.NewClientMembershipConverter().ConvertToEntity(actual))

		// Удаление тестовых данных
		s.db.Delete(&orm.ClientMembership{ID: clientMembership.ID})
		s.db.Delete(&orm.Client{ID: client.ID})
		s.db.Delete(&orm.MembershipType{ID: membershipType.ID})
		s.db.Delete(&orm.Gym{ID: gym.ID})
	})
}

func (s *ClientMembershipsServiceSuite) TestDeleteClientMembership(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[DeleteClientMembership] Successfully deleted membership")
	t.Tags("client_membership", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted membership", func(sCtx provider.StepCtx) {
		// Подготовка тестовых данных, загрузка их в бд
		ctx := context.TODO()
		gym := builder.NewGymBuilder().
			Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		membershipType := builder.NewMembershipTypeBuilder().
			SetGymID(gym.ID).
			Build()
		membershipTypeOrm := orm.NewMembershipTypeConverter().ConvertFromEntity(membershipType)
		err = s.db.Save(&membershipTypeOrm).Error
		sCtx.Assert().NoError(err)

		client := builder.NewClientBuilder().
			Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)
		err = s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		clientMembership := builder.
			NewClientMembershipBuilder().
			SetClientID(client.ID).
			SetMembershipType(membershipType).
			Build()
		clientMembershipOrm := orm.NewClientMembershipConverter().ConvertFromEntity(clientMembership)
		err = s.db.Save(&clientMembershipOrm).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		err = s.clientMembershipService.DeleteClientMembership(ctx, clientMembership.ID)
		sCtx.Assert().NoError(err)

		// Проверка, что запись удалена
		toDelete := &orm.ClientMembership{ID: clientMembership.ID}
		err = s.db.First(&toDelete).Error
		// log.Println("found:", toDelete, "error:", err)
		sCtx.Assert().True(errors.Is(err, gorm.ErrRecordNotFound))

		// Удаление данных
		s.db.Delete(&orm.Client{ID: client.ID})
		s.db.Delete(&orm.MembershipType{ID: membershipType.ID})
		s.db.Delete(&orm.Gym{ID: gym.ID})
	})
}

func (s *ClientMembershipsServiceSuite) TestGetClientMembershipByID(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[GetClientMembershipByID] Successfully retrieved membership by ID")
	t.Tags("client_membership", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved membership by ID", func(sCtx provider.StepCtx) {
		// Подготовка тестовых данных, загрузка их в бд
		ctx := context.TODO()
		gym := builder.NewGymBuilder().
			Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		membershipType := builder.NewMembershipTypeBuilder().
			SetGymID(gym.ID).
			Build()
		membershipTypeOrm := orm.NewMembershipTypeConverter().ConvertFromEntity(membershipType)
		err = s.db.Save(&membershipTypeOrm).Error
		sCtx.Assert().NoError(err)

		client := builder.NewClientBuilder().
			Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)
		err = s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		clientMemberships := []entity.ClientMembership{
			builder.
				NewClientMembershipBuilder().
				SetClientID(client.ID).
				SetMembershipType(membershipType).
				Build(),
			builder.
				NewClientMembershipBuilder().
				SetClientID(client.ID).
				SetMembershipType(membershipType).
				Build(),
		}
		for _, membership := range clientMemberships {
			clientMembershipOrm := orm.NewClientMembershipConverter().ConvertFromEntity(membership)
			err = s.db.Save(&clientMembershipOrm).Error
			sCtx.Assert().NoError(err)
		}

		// Вызов метода
		actual, err := s.clientMembershipService.GetClientMembershipByID(ctx, clientMemberships[1].ID)

		// Проверка
		sCtx.Assert().NoError(err)
		ClientMembershipEqual(sCtx, clientMemberships[1], actual)

		// Удаление
		for _, mem := range clientMemberships {
			s.db.Delete(&orm.ClientMembership{ID: mem.ID})
		}
		s.db.Delete(&orm.Client{ID: client.ID})
		s.db.Delete(&orm.MembershipType{ID: membershipType.ID})
		s.db.Delete(&orm.Gym{ID: gym.ID})
	})
}

func (s *ClientMembershipsServiceSuite) TestListClientMembershipsByClientID(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ListClientMembershipsByClientID] Successfully listed memberships by client ID")
	t.Tags("client_membership", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully listed memberships by client ID", func(sCtx provider.StepCtx) {
		// Подготовка тестовых данных, загрузка их в бд
		ctx := context.TODO()
		gym := builder.NewGymBuilder().
			Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		membershipType := builder.NewMembershipTypeBuilder().
			SetGymID(gym.ID).
			Build()
		membershipTypeOrm := orm.NewMembershipTypeConverter().ConvertFromEntity(membershipType)
		err = s.db.Save(&membershipTypeOrm).Error
		sCtx.Assert().NoError(err)

		client := builder.NewClientBuilder().
			Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)
		err = s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		clientMemberships := []entity.ClientMembership{
			builder.
				NewClientMembershipBuilder().
				SetClientID(client.ID).
				SetMembershipType(membershipType).
				Build(),
			builder.
				NewClientMembershipBuilder().
				SetClientID(client.ID).
				SetMembershipType(membershipType).
				Build(),
		}
		for _, membership := range clientMemberships {
			clientMembershipOrm := orm.NewClientMembershipConverter().ConvertFromEntity(membership)
			err = s.db.Save(&clientMembershipOrm).Error
			sCtx.Assert().NoError(err)
		}

		// Вызов метода
		actualMemberships, err := s.clientMembershipService.ListClientMembershipsByClientID(ctx, client.ID)

		// Проверка
		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(len(clientMemberships), len(actualMemberships))
		for i, membership := range clientMemberships {
			ClientMembershipEqual(sCtx, membership, actualMemberships[i])
		}

		// Удаление
		for _, mem := range clientMemberships {
			s.db.Delete(&orm.ClientMembership{ID: mem.ID})
		}
		s.db.Delete(&orm.Client{ID: client.ID})
		s.db.Delete(&orm.MembershipType{ID: membershipType.ID})
		s.db.Delete(&orm.Gym{ID: gym.ID})
	})
}

func TestClientMembershipsServiceSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ClientMembershipsServiceSuite))
}
