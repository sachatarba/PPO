package integration

import (
	"context"
	"errors"
	"os"
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

type TrainingServiceSuite struct {
	suite.Suite

	trainingService service.ITrainingService
	db              *gorm.DB
}

func TrainingEqual(sCtx provider.StepCtx, expected, actual entity.Training) {
	sCtx.Assert().Equal(expected.ID, actual.ID, "ID should be equal")
	sCtx.Assert().Equal(expected.Title, actual.Title, "Title should be equal")
	sCtx.Assert().Equal(expected.Description, actual.Description, "Description should be equal")
	sCtx.Assert().Equal(expected.TrainingType, actual.TrainingType, "TrainingType should be equal")
	sCtx.Assert().Equal(expected.TrainerID, actual.TrainerID, "TrainerID should be equal")
}

func (s *TrainingServiceSuite) BeforeAll(t provider.T) {
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
	repo := repository.NewTrainingRepo(db)

	s.trainingService = service.NewTrainingService(repo)
}

func (s *TrainingServiceSuite) AfterAll(t provider.T) {
	tables, err := s.db.Migrator().GetTables()
	t.Assert().NoError(err)

	for _, table := range tables {
		err := s.db.Migrator().DropTable(table)
		t.Assert().NoError(err)
	}
}

// Тест создания новой тренировки
func (s *TrainingServiceSuite) TestCreateNewTraining(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[CreateNewTraining] Successfully created a new training")
	t.Tags("training_service", "service", "create")
	t.Parallel()

	t.WithNewStep("Correct: successfully created new training", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём тренера
		trainer := builder.NewTrainerBuilder().SetGymsID([]uuid.UUID{}).Build()
		trainerOrm := orm.NewTrainerConverter().ConvertFromEntity(trainer)
		err := s.db.Save(&trainerOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём тренировку с указанием ID тренера
		training := builder.NewTrainingBuilder().
			SetTrainerID(trainer.ID).
			Build()
		// trainingOrm := orm.NewTrainingConverter().ConvertFromEntity(training)

		// Вызов метода
		err = s.trainingService.CreateNewTraining(ctx, training)

		// Проверка
		sCtx.Assert().NoError(err)
		actual := orm.Training{ID: training.ID}
		err = s.db.First(&actual).Error
		sCtx.Assert().NoError(err)
		actualTraining := orm.NewTrainingConverter().ConvertToEntity(actual)

		TrainingEqual(sCtx, training, actualTraining)

		// Удаление тестовых данных
		s.db.Delete(&trainerOrm)
		s.db.Delete(&actual)
	})
}

// Тест изменения тренировки
func (s *TrainingServiceSuite) TestChangeTraining(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ChangeTraining] Successfully changed training data")
	t.Tags("training_service", "service", "update")
	t.Parallel()

	t.WithNewStep("Correct: successfully updated training", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём тренера
		trainer := builder.NewTrainerBuilder().SetGymsID([]uuid.UUID{}).Build()
		trainerOrm := orm.NewTrainerConverter().ConvertFromEntity(trainer)
		err := s.db.Save(&trainerOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём тренировку
		training := builder.NewTrainingBuilder().
			SetTrainerID(trainer.ID).
			Build()
		trainingOrm := orm.NewTrainingConverter().ConvertFromEntity(training)
		err = s.db.Save(&trainingOrm).Error
		sCtx.Assert().NoError(err)

		// Обновляем данные тренировки
		training.Description = "Updated description"

		// Вызов метода
		err = s.trainingService.ChangeTraining(ctx, training)

		// Проверка
		sCtx.Assert().NoError(err)
		actual := orm.Training{ID: training.ID}
		err = s.db.First(&actual).Error
		sCtx.Assert().NoError(err)
		actualTraining := orm.NewTrainingConverter().ConvertToEntity(actual)

		TrainingEqual(sCtx, training, actualTraining)

		// Удаление тестовых данных
		s.db.Delete(&trainerOrm)
		s.db.Delete(&actual)
	})
}

// Тест удаления тренировки
func (s *TrainingServiceSuite) TestDeleteTraining(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[DeleteTraining] Successfully deleted a training")
	t.Tags("training_service", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted training", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём тренера
		trainer := builder.NewTrainerBuilder().SetGymsID([]uuid.UUID{}).Build()
		trainerOrm := orm.NewTrainerConverter().ConvertFromEntity(trainer)
		err := s.db.Save(&trainerOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём тренировку
		training := builder.NewTrainingBuilder().
			SetTrainerID(trainer.ID).
			Build()
		trainingOrm := orm.NewTrainingConverter().ConvertFromEntity(training)
		err = s.db.Save(&trainingOrm).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		err = s.trainingService.DeleteTraining(ctx, training.ID)

		// Проверка
		sCtx.Assert().NoError(err)
		toDelete := &orm.Training{ID: training.ID}
		err = s.db.First(&toDelete).Error
		sCtx.Assert().True(errors.Is(err, gorm.ErrRecordNotFound))

		// Удаление тестовых данных
		s.db.Delete(&trainerOrm)
	})
}

// Тест получения списка тренировок по ID тренера
func (s *TrainingServiceSuite) TestListTrainingsByTrainerID(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ListTrainingsByTrainerID] Successfully listed trainings by trainer ID")
	t.Tags("training_service", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully listed trainings by trainer ID", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём тренера
		trainer := builder.NewTrainerBuilder().SetGymsID([]uuid.UUID{}).Build()
		trainerOrm := orm.NewTrainerConverter().ConvertFromEntity(trainer)
		err := s.db.Save(&trainerOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём тренировки, связанные с тренером
		trainings := []entity.Training{
			builder.NewTrainingBuilder().SetTrainerID(trainer.ID).Build(),
			builder.NewTrainingBuilder().SetTrainerID(trainer.ID).Build(),
		}
		for _, training := range trainings {
			trainingOrm := orm.NewTrainingConverter().ConvertFromEntity(training)
			err := s.db.Save(&trainingOrm).Error
			sCtx.Assert().NoError(err)
		}

		// Вызов метода
		actualTrainings, err := s.trainingService.ListTrainingsByTrainerID(ctx, trainer.ID)

		// Проверка
		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(len(trainings), len(actualTrainings))

		for i, training := range actualTrainings {
			TrainingEqual(sCtx, trainings[i], training)
		}

		// Удаление тестовых данных
		for _, training := range trainings {
			s.db.Delete(&orm.Training{ID: training.ID})
		}
		s.db.Delete(&trainerOrm)
	})
}

func TestTrainingServiceSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(TrainingServiceSuite))
}
