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

type ScheduleServiceSuite struct {
	suite.Suite

	scheduleService service.IScheduleService
	db              *gorm.DB
}

func ScheduleEqual(sCtx provider.StepCtx, expected, actual entity.Schedule) {
	sCtx.Assert().Equal(expected.ID, actual.ID, "ID should be equal")
	// sCtx.Assert().Equal(expected.DayOfTheWeek, actual.DayOfTheWeek, "DayOfTheWeek should be equal")
	// sCtx.Assert().Equal(expected.StartTime, actual.StartTime, "StartTime should be equal")
	// sCtx.Assert().Equal(expected.EndTime, actual.EndTime, "EndTime should be equal")
	sCtx.Assert().Equal(expected.ClientID, actual.ClientID, "ClientID should be equal")
	sCtx.Assert().Equal(expected.TrainingID, actual.TrainingID, "TrainingID should be equal")
	// sCtx.Assert().Equal(expected.Training, actual.Training, "Training should be equal")
}

func (s *ScheduleServiceSuite) BeforeAll(t provider.T) {
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
	repo := repository.NewScheduleRepo(db)

	s.scheduleService = service.NewScheduleService(repo)
}

func (s *ScheduleServiceSuite) AfterAll(t provider.T) {
	tables, err := s.db.Migrator().GetTables()
	t.Assert().NoError(err)

	for _, table := range tables {
		err := s.db.Migrator().DropTable(table)
		t.Assert().NoError(err)
	}
}

func (s *ScheduleServiceSuite) TestCreateNewSchedule(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[CreateNewSchedule] Successfully created a new schedule")
	t.Tags("schedule_service", "service", "create")
	t.Parallel()

	t.WithNewStep("Correct: successfully created new schedule", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём необходимые зависимости для Schedule
		// Создаём Client
		client := builder.NewClientBuilder().Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)
		err := s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём Gym
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err = s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём Trainer
		trainer := builder.NewTrainerBuilder().SetGymsID([]uuid.UUID{gym.ID}).Build()
		trainerOrm := orm.NewTrainerConverter().ConvertFromEntity(trainer)
		err = s.db.Model(&gymOrm).Association("Trainers").Append(&trainerOrm)
		sCtx.Assert().NoError(err)

		// Создаём Training с зависимостью от Gym и Trainer
		training := builder.NewTrainingBuilder().
			SetTrainerID(trainer.ID).
			Build()
		trainingOrm := orm.NewTrainingConverter().ConvertFromEntity(training)
		err = s.db.Save(&trainingOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём объект Schedule с зависимостями от Client и Training
		schedule := builder.NewScheduleBuilder().
			SetClientID(client.ID).
			SetTrainingID(training.ID).
			Build()

		sCtx.WithNewParameters("ctx", ctx, "schedule", schedule)

		// Вызов метода
		err = s.scheduleService.CreateNewSchedule(ctx, schedule)

		// Проверка
		sCtx.Assert().NoError(err)
		actualOrm := orm.Schedule{ID: schedule.ID}

		err = s.db.First(&actualOrm).Error
		sCtx.Assert().NoError(err)
		actual := orm.NewScheduleConverter().ConvertToEntity(actualOrm)

		ScheduleEqual(sCtx, schedule, actual)
	})
}

func (s *ScheduleServiceSuite) TestChangeSchedule(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[ChangeSchedule] Successfully changed schedule data")
	t.Tags("schedule_service", "service", "update")
	t.Parallel()

	t.WithNewStep("Correct: successfully updated schedule", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём необходимые зависимости для Schedule
		// Создаём Client
		client := builder.NewClientBuilder().Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)
		err := s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём Gym
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err = s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём Trainer
		trainer := builder.NewTrainerBuilder().SetGymsID([]uuid.UUID{gym.ID}).Build()
		trainerOrm := orm.NewTrainerConverter().ConvertFromEntity(trainer)
		err = s.db.Model(&gymOrm).Association("Trainers").Append(&trainerOrm)
		sCtx.Assert().NoError(err)

		// Создаём Training с зависимостью от Gym и Trainer
		training := builder.NewTrainingBuilder().
			SetTrainerID(trainer.ID).
			Build()
		trainingOrm := orm.NewTrainingConverter().ConvertFromEntity(training)
		err = s.db.Save(&trainingOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём объект Schedule
		schedule := builder.NewScheduleBuilder().
			SetClientID(client.ID).
			SetTrainingID(training.ID).
			Build()
		scheduleValidated := schedule
		scheduleValidated.Validate()
		scheduleOrm := orm.NewScheduleConverter().ConvertFromEntity(scheduleValidated)
		// scheduleOrm.StartTime
		err = s.db.Save(&scheduleOrm).Error
		sCtx.Assert().NoError(err)

		// Изменяем данные
		schedule.StartTime = "09:00:00"

		// Вызов метода
		err = s.scheduleService.ChangeSchedule(ctx, schedule)

		// Проверка
		sCtx.Assert().NoError(err)
		actual := orm.Schedule{ID: schedule.ID}
		err = s.db.First(&actual).Error
		sCtx.Assert().NoError(err)
		actualSchedule := orm.NewScheduleConverter().ConvertToEntity(actual)

		ScheduleEqual(sCtx, schedule, actualSchedule)

	})
}

func (s *ScheduleServiceSuite) TestDeleteSchedule(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[DeleteSchedule] Successfully deleted a schedule")
	t.Tags("schedule_service", "service", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deleted schedule", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём необходимые зависимости для Schedule
		// Создаём Client
		client := builder.NewClientBuilder().Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)
		err := s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём Gym
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err = s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём Trainer
		trainer := builder.NewTrainerBuilder().SetGymsID([]uuid.UUID{gym.ID}).Build()
		trainerOrm := orm.NewTrainerConverter().ConvertFromEntity(trainer)
		err = s.db.Model(&gymOrm).Association("Trainers").Append(&trainerOrm)
		sCtx.Assert().NoError(err)

		// Создаём Training с зависимостью от Gym и Trainer
		training := builder.NewTrainingBuilder().
			SetTrainerID(trainer.ID).
			Build()
		trainingOrm := orm.NewTrainingConverter().ConvertFromEntity(training)
		err = s.db.Save(&trainingOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём объект Schedule
		schedule := builder.NewScheduleBuilder().
			SetClientID(client.ID).
			SetTrainingID(training.ID).
			Build()
		scheduleValidated := schedule
		scheduleValidated.Validate()
		scheduleOrm := orm.NewScheduleConverter().ConvertFromEntity(scheduleValidated)
		err = s.db.Save(&scheduleOrm).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		err = s.scheduleService.DeleteSchedule(ctx, schedule.ID)

		// Проверка
		sCtx.Assert().NoError(err)
		toDelete := &orm.Schedule{ID: schedule.ID}
		err = s.db.First(&toDelete).Error
		sCtx.Assert().True(errors.Is(err, gorm.ErrRecordNotFound))

	})
}

func (s *ScheduleServiceSuite) TestGetScheduleByID(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[GetScheduleByID] Successfully retrieved schedule by ID")
	t.Tags("schedule_service", "service", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieved schedule by ID", func(sCtx provider.StepCtx) {
		ctx := context.TODO()

		// Создаём необходимые зависимости для Schedule
		// Создаём Client
		client := builder.NewClientBuilder().Build()
		clientOrm := orm.NewClientConverter().ConvertFromEntity(client)
		err := s.db.Save(&clientOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём Gym
		gym := builder.NewGymBuilder().Build()
		gymOrm := orm.NewGymConverter().ConvertFromEntity(gym)
		err = s.db.Save(&gymOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём Trainer
		trainer := builder.NewTrainerBuilder().SetGymsID([]uuid.UUID{gym.ID}).Build()
		trainerOrm := orm.NewTrainerConverter().ConvertFromEntity(trainer)
		err = s.db.Model(&gymOrm).Association("Trainers").Append(&trainerOrm)
		sCtx.Assert().NoError(err)

		// Создаём Training с зависимостью от Gym и Trainer
		training := builder.NewTrainingBuilder().
			SetTrainerID(trainer.ID).
			Build()

		trainingOrm := orm.NewTrainingConverter().ConvertFromEntity(training)
		err = s.db.Save(&trainingOrm).Error
		sCtx.Assert().NoError(err)

		// Создаём объект Schedule
		schedule := builder.NewScheduleBuilder().
			SetClientID(client.ID).
			SetTrainingID(training.ID).
			Build()

		scheduleValidated := schedule
		scheduleValidated.Validate()
		scheduleOrm := orm.NewScheduleConverter().ConvertFromEntity(scheduleValidated)
		err = s.db.Save(&scheduleOrm).Error
		sCtx.Assert().NoError(err)

		// Вызов метода
		actual, err := s.scheduleService.GetScheduleByID(ctx, schedule.ID)

		// Проверка
		sCtx.Assert().NoError(err)
		ScheduleEqual(sCtx, schedule, actual)

	})
}

func TestScheduleServiceSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ScheduleServiceSuite))
}
