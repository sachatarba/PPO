package integration

import (
	"context"
	"errors"
	"os"
	"strconv"
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

type MembershipTypeServiceSuite struct {
	suite.Suite

	membershipTypeService service.IMembershipTypeService
	db                    *gorm.DB
}

func MembershipTypeEqual(sCtx provider.StepCtx, expected, actual entity.MembershipType) {
	sCtx.Assert().Equal(expected.ID, actual.ID, "ID should be equal")
	sCtx.Assert().Equal(expected.Type, actual.Type, "Type should be equal")
	sCtx.Assert().Equal(expected.Description, actual.Description, "Description should be equal")
	sCtx.Assert().Equal(expected.DaysDuration, actual.DaysDuration, "DaysDuration should be equal")
	sCtx.Assert().Equal(expected.GymID, actual.GymID, "GymID should be equal")

	expectedPrice, err := strconv.ParseFloat(expected.Price, 64)
	sCtx.Assert().NoError(err)

	actualPrice, err := strconv.ParseFloat(actual.Price, 64)
	sCtx.Assert().NoError(err)

	sCtx.Assert().InDelta(expectedPrice, actualPrice, 0.0001, "Price should be equal")
}

func (s *MembershipTypeServiceSuite) BeforeAll(t provider.T) {
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
	repo := repository.NewMembershipTypeRepo(db)

	s.membershipTypeService = service.NewMembershipTypeService(repo)
}

func (s *MembershipTypeServiceSuite) AfterAll(t provider.T) {
	tables, err := s.db.Migrator().GetTables()
	t.Assert().NoError(err)

	for _, table := range tables {
		err := s.db.Migrator().DropTable(table)
		t.Assert().NoError(err)
	}
}

func (s *MembershipTypeServiceSuite) TestRegisterNewMembershipType(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[RegisterNewMembershipType] Successfully registered new membership type")
	t.Tags("membership_type", "service", "create")
	t.Parallel()

	t.WithNewStep("Correct: successfully registered new membership type", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём зал для привязки к абонементу
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём новый тип абонемента с привязкой к залу
		membershipType := builder.NewMembershipTypeBuilder().SetGymID(gym.ID).Build()

		// Вызов метода
		err = s.membershipTypeService.RegisterNewMembershipType(ctx, membershipType)

		// Проверка
		sCtx.Assert().NoError(err)
		actualOrm := orm.MembershipType{ID: membershipType.ID}
		err = s.db.First(&actualOrm).Error
		sCtx.Assert().NoError(err)
		actual := orm.NewMembershipTypeConverter().ConvertToEntity(actualOrm)

		MembershipTypeEqual(sCtx, membershipType, actual)

		// Удаление тестовых данных
		err = s.db.Delete(&orm.MembershipType{ID: membershipType.ID}).Error
		sCtx.Assert().NoError(err)
		err = s.db.Delete(&orm.Gym{ID: gym.ID}).Error
		sCtx.Assert().NoError(err)
	})
}

func (s *MembershipTypeServiceSuite) TestChangeMembershipType(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ChangeMembershipType] Successfully changed membership type")
	t.Tags("membership_type", "service", "update")
	t.Parallel()

	t.WithNewStep("Correct: successfully changed membership type", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём зал для привязки к абонементу
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём новый тип абонемента
		membershipType := builder.NewMembershipTypeBuilder().SetGymID(gym.ID).Build()
		err = s.db.Save(&membershipType).Error
		sCtx.Assert().NoError(err)

		// Изменяем данные абонемента
		membershipType.Description = "Access to all facilities and classes"

		// Вызов метода
		err = s.membershipTypeService.ChangeMembershipType(ctx, membershipType)
		sCtx.Assert().NoError(err)

		// Проверка
		actualOrm := orm.MembershipType{ID: membershipType.ID}
		err = s.db.First(&actualOrm).Error
		sCtx.Assert().NoError(err)
		actual := orm.NewMembershipTypeConverter().ConvertToEntity(actualOrm)

		MembershipTypeEqual(sCtx, membershipType, actual)

		// Удаление тестовых данных
		err = s.db.Delete(&orm.MembershipType{ID: membershipType.ID}).Error
		sCtx.Assert().NoError(err)
		err = s.db.Delete(&orm.Gym{ID: gym.ID}).Error
		sCtx.Assert().NoError(err)
	})
}

func (s *MembershipTypeServiceSuite) TestDeleteMembershipType(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[DeleteMembershipType] Successfully deleted membership type")
	t.Tags("membership_type", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted membership type", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём зал для привязки к абонементу
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём новый тип абонемента
		membershipType := builder.NewMembershipTypeBuilder().SetGymID(gym.ID).Build()
		err = s.db.Save(&membershipType).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		err = s.membershipTypeService.DeleteMembershipType(ctx, membershipType.ID)
		sCtx.Assert().NoError(err)

		// Проверка, что запись удалена
		toDelete := &orm.MembershipType{ID: membershipType.ID}
		err = s.db.First(&toDelete).Error
		sCtx.Assert().True(errors.Is(err, gorm.ErrRecordNotFound))

		// Удаление данных
		s.db.Delete(&orm.Gym{ID: gym.ID})
	})
}

func (s *MembershipTypeServiceSuite) TestGetMembershipTypeByID(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[GetMembershipTypeByID] Successfully retrieved membership type by ID")
	t.Tags("membership_type", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved membership type by ID", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём зал для привязки к абонементу
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём новый тип абонемента
		membershipType := builder.NewMembershipTypeBuilder().SetGymID(gym.ID).Build()
		err = s.db.Save(&membershipType).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		actual, err := s.membershipTypeService.GetMembershipTypeByID(ctx, membershipType.ID)
		sCtx.Assert().NoError(err)

		// Проверка
		MembershipTypeEqual(sCtx, membershipType, actual)

		// Удаление данных
		s.db.Delete(&orm.MembershipType{ID: membershipType.ID})
		s.db.Delete(&orm.Gym{ID: gym.ID})
	})
}

func (s *MembershipTypeServiceSuite) TestListMembershipTypesByGymID(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ListMembershipTypesByGymID] Successfully listed membership types by gym ID")
	t.Tags("membership_type", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully listed membership types by gym ID", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём зал для привязки к абонементу
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём несколько типов абонементов
		membershipTypes := []entity.MembershipType{
			builder.NewMembershipTypeBuilder().SetGymID(gym.ID).Build(),
			builder.NewMembershipTypeBuilder().SetGymID(gym.ID).Build(),
		}
		for _, membershipType := range membershipTypes {
			err := s.db.Save(&membershipType).Error
			sCtx.Assert().NoError(err)
		}

		// Вызов метода
		actualMembershipTypes, err := s.membershipTypeService.ListMembershipTypesByGymID(ctx, gym.ID)
		sCtx.Assert().NoError(err)

		// Проверка
		sCtx.Assert().Equal(len(membershipTypes), len(actualMembershipTypes))
		for i, membershipType := range membershipTypes {
			MembershipTypeEqual(sCtx, membershipType, actualMembershipTypes[i])
		}

		// Удаление данных
		for _, membershipType := range membershipTypes {
			s.db.Delete(&orm.MembershipType{ID: membershipType.ID})
		}
		s.db.Delete(&orm.Gym{ID: gym.ID})
	})
}

func TestMembershipTypeServiceSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(MembershipTypeServiceSuite))
}
