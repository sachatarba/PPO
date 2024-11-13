package integration

import (
	"context"
	"errors"
	"os"
	"slices"
	"testing"

	"github.com/google/uuid"
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

type TrainerServiceSuite struct {
	suite.Suite

	trainerService service.ITrainerService
	db             *gorm.DB
}

func TrainerEqual(sCtx provider.StepCtx, expected, actual entity.Trainer) {
	sCtx.Assert().Equal(expected.ID, actual.ID, "ID should be equal")
	sCtx.Assert().Equal(expected.Fullname, actual.Fullname, "Fullname should be equal")
	sCtx.Assert().Equal(expected.Email, actual.Email, "Email should be equal")
	sCtx.Assert().Equal(expected.Phone, actual.Phone, "Phone should be equal")
	sCtx.Assert().Equal(expected.Qualification, actual.Qualification, "Qualification should be equal")

	sCtx.Assert().InDelta(expected.UnitPrice, actual.UnitPrice, 0.0001, "UnitPrice should be equal")

	sCtx.Assert().ElementsMatch(expected.GymsID, actual.GymsID, "GymsID should contain the same elements")
}

func (s *TrainerServiceSuite) BeforeAll(t provider.T) {
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
	repo := repository.NewTrainerRepo(db)

	s.trainerService = service.NewTrainerService(repo)
}

func (s *TrainerServiceSuite) AfterAll(t provider.T) {
	tables, err := s.db.Migrator().GetTables()
	t.Assert().NoError(err)

	for _, table := range tables {
		err := s.db.Migrator().DropTable(table)
		t.Assert().NoError(err)
	}
}


// Тест изменения данных тренера
func (s *TrainerServiceSuite) TestChangeTrainer(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ChangeTrainer] Successfully changed trainer data")
	t.Tags("trainer_service", "service", "update")
	t.Parallel()

	t.WithNewStep("Correct: successfully updated trainer", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём тренера
		trainer := builder.NewTrainerBuilder().SetGymsID([]uuid.UUID{}).Build()
		trainerOrm := orm.NewTrainerConverter().ConvertFromEntity(trainer)
		err := s.db.Save(&trainerOrm).Error
		sCtx.Assert().NoError(err)

		// Изменяем данные
		trainer.Fullname = "Updated Trainer Name"

		// Вызов метода
		err = s.trainerService.ChangeTrainer(ctx, trainer)

		// Проверка
		sCtx.Assert().NoError(err)
		actual := orm.Trainer{ID: trainer.ID}
		err = s.db.First(&actual).Error
		sCtx.Assert().NoError(err)
		actualTrainer := orm.NewTrainerConverter().ConvertToEntity(actual)

		TrainerEqual(sCtx, trainer, actualTrainer)
	})
}

// Тест удаления тренера
func (s *TrainerServiceSuite) TestDeleteTrainer(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[DeleteTrainer] Successfully deleted a trainer")
	t.Tags("trainer_service", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted trainer", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём тренера
		trainer := builder.NewTrainerBuilder().SetGymsID([]uuid.UUID{}).Build()
		trainerOrm := orm.NewTrainerConverter().ConvertFromEntity(trainer)
		err := s.db.Save(&trainerOrm).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		err = s.trainerService.DeleteTrainer(ctx, trainer.ID)

		// Проверка
		sCtx.Assert().NoError(err)
		toDelete := &orm.Trainer{ID: trainer.ID}
		err = s.db.First(&toDelete).Error
		sCtx.Assert().True(errors.Is(err, gorm.ErrRecordNotFound))
	})
}

// Тест получения тренера по ID
func (s *TrainerServiceSuite) TestGetTrainerByID(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[GetTrainerByID] Successfully retrieved trainer by ID")
	t.Tags("trainer_service", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved trainer by ID", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём тренера
		trainer := builder.NewTrainerBuilder().SetGymsID([]uuid.UUID{}).Build()
		trainerOrm := orm.NewTrainerConverter().ConvertFromEntity(trainer)
		err := s.db.Save(&trainerOrm).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		actual, err := s.trainerService.GetTrainerByID(ctx, trainer.ID)

		// Проверка
		sCtx.Assert().NoError(err)
		TrainerEqual(sCtx, trainer, actual)
	})
}

// Тест списка тренеров
func (s *TrainerServiceSuite) TestListTrainers(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ListTrainers] Successfully listed all trainers")
	t.Tags("trainer_service", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully listed trainers", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём нескольких тренеров
		trainers := []entity.Trainer{
			builder.
				NewTrainerBuilder().
				SetGymsID([]uuid.UUID{}).
				Build(),
			builder.
				NewTrainerBuilder().
				SetGymsID([]uuid.UUID{}).
				Build(),
		}

		for _, trainer := range trainers {
			trainerOrm := orm.NewTrainerConverter().ConvertFromEntity(trainer)
			err := s.db.Save(&trainerOrm).Error
			sCtx.Assert().NoError(err)
		}

		// Вызов метода
		actualTrainers, err := s.trainerService.ListTrainers(ctx)

		// Проверка
		sCtx.Assert().NoError(err)
		for i, trainer := range actualTrainers {
			if slices.ContainsFunc(actualTrainers, func(t entity.Trainer) bool {
				return t.ID == trainer.ID
			}) {
				TrainerEqual(sCtx, trainer, actualTrainers[i])
			}
		}

		// Удаление тестовых данных
		for _, trainer := range trainers {
			s.db.Delete(&orm.Trainer{ID: trainer.ID})
		}
	})
}

// Тест списка тренеров по ID зала
func (s *TrainerServiceSuite) TestListTrainersByGymID(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ListTrainersByGymID] Successfully listed trainers by gym ID")
	t.Tags("trainer_service", "service", "list", "gym")
	t.Parallel()

	t.WithNewStep("Correct: successfully listed trainers by gym ID", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём зал и тренера, связанного с этим залом
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err := s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		trainer := builder.NewTrainerBuilder().SetGymsID([]uuid.UUID{gym.ID}).Build()
		trainerOrm := orm.NewTrainerConverter().ConvertFromEntity(trainer)
		err = s.db.Model(&gymOrm).Association("Trainers").Append(&trainerOrm)
		sCtx.Assert().NoError(err)

		// Вызов метода
		trainers, err := s.trainerService.ListTrainersByGymID(ctx, gym.ID)

		// Проверка
		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(1, len(trainers))
	})
}

func TestTrainerServiceSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(TrainerServiceSuite))
}
