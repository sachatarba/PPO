package integration

import (
	"context"
	"errors"
	"os"
	"slices"
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

type GymServiceSuite struct {
	suite.Suite

	gymService service.IGymService
	db         *gorm.DB
}

func GymEqual(sCtx provider.StepCtx, expected, actual entity.Gym) {
	sCtx.Assert().Equal(expected.ID, actual.ID, "ID should be equal")
	sCtx.Assert().Equal(expected.Name, actual.Name, "Name should be equal")
	sCtx.Assert().Equal(expected.Addres, actual.Addres, "Address should be equal")
	sCtx.Assert().Equal(expected.City, actual.City, "City should be equal")
	sCtx.Assert().Equal(expected.IsChain, actual.IsChain, "IsChain should be equal")
	sCtx.Assert().Equal(expected.Phone, actual.Phone, "City should be equal")
}

func (s *GymServiceSuite) BeforeAll(t provider.T) {
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
	repo := repository.NewGymRepo(db)

	s.gymService = service.NewGymService(repo)
}

func (s *GymServiceSuite) AfterAll(t provider.T) {
	tables, err := s.db.Migrator().GetTables()
	t.Assert().NoError(err)

	for _, table := range tables {
		err := s.db.Migrator().DropTable(table)
		t.Assert().NoError(err)
	}
}

func (s *GymServiceSuite) TestRegisterNewGym(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[RegisterNewGym] Successfully registered a new gym")
	t.Tags("gym_service", "service", "register")
	t.Parallel()

	t.WithNewStep("Correct: successfully registered new gym", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		gym := builder.NewGymBuilder().Build()

		sCtx.WithNewParameters("ctx", ctx, "gym", gym)

		// Вызов метода
		err := s.gymService.RegisterNewGym(ctx, gym)

		// Проверка
		sCtx.Assert().NoError(err)
		actualOrm := orm.Gym{ID: gym.ID}

		err = s.db.First(&actualOrm).Error
		sCtx.Assert().NoError(err)
		actual := orm.NewGymConverter().ConvertToEntity(actualOrm)

		GymEqual(sCtx, gym, actual)

		// Удаление тестовых данных
		err = s.db.Delete(&orm.Gym{ID: gym.ID}).Error
		sCtx.Assert().NoError(err)
	})
}

func (s *GymServiceSuite) TestChangeGym(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ChangeGym] Successfully changed gym data")
	t.Tags("gym_service", "service", "update")
	t.Parallel()

	t.WithNewStep("Correct: successfully updated gym", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		// Изменяем данные
		gym.Name = "Updated Gym Name"

		// Вызов метода
		err = s.gymService.ChangeGym(ctx, gym)

		// Проверка
		sCtx.Assert().NoError(err)
		actual := orm.Gym{ID: gym.ID}
		err = s.db.First(&actual).Error
		sCtx.Assert().NoError(err)
		actualGym := orm.NewGymConverter().ConvertToEntity(actual)

		GymEqual(sCtx, gym, actualGym)

		// Удаление тестовых данных
		s.db.Delete(&orm.Gym{ID: gym.ID})
	})
}

func (s *GymServiceSuite) TestDeleteGym(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[DeleteGym] Successfully deleted a gym")
	t.Tags("gym_service", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted gym", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		err = s.gymService.DeleteGym(ctx, gym.ID)

		// Проверка
		sCtx.Assert().NoError(err)
		toDelete := &orm.Gym{ID: gym.ID}
		err = s.db.First(&toDelete).Error
		sCtx.Assert().True(errors.Is(err, gorm.ErrRecordNotFound))
	})
}

func (s *GymServiceSuite) TestGetGymByID(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[GetGymByID] Successfully retrieved gym by ID")
	t.Tags("gym_service", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved gym by ID", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		actual, err := s.gymService.GetGymByID(ctx, gym.ID)

		// Проверка
		sCtx.Assert().NoError(err)
		GymEqual(sCtx, gym, actual)

		// Удаление тестовых данных
		s.db.Delete(&orm.Gym{ID: gym.ID})
	})
}

func (s *GymServiceSuite) TestListGyms(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ListGyms] Successfully listed all gyms")
	t.Tags("gym_service", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully listed gyms", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		gyms := []entity.Gym{
			builder.NewGymBuilder().Build(),
			builder.NewGymBuilder().Build(),
		}

		for _, gym := range gyms {
			gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
			err := s.db.Save(&gymOrm).Error
			sCtx.Assert().NoError(err)
		}

		// Вызов метода
		actualGyms, err := s.gymService.ListGyms(ctx)

		// Проверка
		sCtx.Assert().NoError(err)
		for i, gym := range actualGyms {
			if slices.ContainsFunc(actualGyms, func(g entity.Gym) bool {
				return g.ID == gym.ID
			}) {
				GymEqual(sCtx, gym, actualGyms[i])
			}
		}

		// Удаление тестовых данных
		for _, gym := range gyms {
			s.db.Delete(&orm.Gym{ID: gym.ID})
		}
	})
}

func TestGymServiceSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(GymServiceSuite))
}
