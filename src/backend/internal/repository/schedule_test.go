package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/repository/mocks"
	"github.com/sachatarba/course-db/internal/service"
	"github.com/sachatarba/course-db/internal/utils/builder"
)

type ScheduleRepoSuite struct {
	suite.Suite

	mock         sqlmock.Sqlmock
	scheduleRepo service.IScheduleRepository
}

func (c *ScheduleRepoSuite) BeforeEach(t provider.T) {
	t.Title("Init mock db")
	db, mock := mocks.NewMockDB()

	c.scheduleRepo = NewScheduleRepo(db)
	c.mock = mock

	t.Tags("fixture", "schedule", "db")
}

func (s *ScheduleRepoSuite) TestCreateNewSchedule(t provider.T) {
	t.Title("[Create] Create new schedule")
	t.Tags("repository", "postgres")

	t.WithNewStep("Create new schedule", func(sCtx provider.StepCtx) {
		schedule := builder.NewScheduleBuilder().
			SetClientID(uuid.New()).
			SetTrainingID(uuid.New()).
			SetDayOfTheWeek("Monday").
			SetStartTime("10:00").
			SetEndTime("11:00").
			Build()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`INSERT INTO "schedules" \("id","day_of_the_week","start_time","end_time","client_id","training_id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\)`).
			WithArgs(
				schedule.ID,
				schedule.DayOfTheWeek,
				schedule.StartTime,
				schedule.EndTime,
				schedule.ClientID,
				schedule.TrainingID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "schedule", schedule)

		err := s.scheduleRepo.CreateNewSchedule(ctx, schedule)

		sCtx.Assert().NoError(err)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *ScheduleRepoSuite) TestChangeSchedule(t provider.T) {
	t.Title("[Change] Change schedule")
	t.Tags("repository", "postgres")

	t.WithNewStep("Update schedule details", func(sCtx provider.StepCtx) {
		schedule := builder.NewScheduleBuilder().
			SetID(uuid.New()).
			SetClientID(uuid.New()).
			SetTrainingID(uuid.New()).
			SetDayOfTheWeek("Monday").
			SetStartTime("10:00").
			SetEndTime("11:00").
			Build()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`UPDATE "schedules" SET "day_of_the_week"=\$1,"start_time"=\$2,"end_time"=\$3,"client_id"=\$4,"training_id"=\$5 WHERE "id" = \$6`).
			WithArgs(
				schedule.DayOfTheWeek,
				schedule.StartTime,
				schedule.EndTime,
				schedule.ClientID,
				schedule.TrainingID,
				schedule.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "schedule", schedule)

		err := s.scheduleRepo.ChangeSchedule(ctx, schedule)

		sCtx.Assert().NoError(err)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *ScheduleRepoSuite) TestDeleteSchedule(t provider.T) {
	t.Title("[Delete] Delete schedule by ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("Delete schedule by ID", func(sCtx provider.StepCtx) {
		scheduleID := uuid.New()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`DELETE FROM "schedules" WHERE "schedules"."id" = \$1`).
			WithArgs(scheduleID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "scheduleID", scheduleID)

		err := s.scheduleRepo.DeleteSchedule(ctx, scheduleID)

		sCtx.Assert().NoError(err)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *ScheduleRepoSuite) TestGetScheduleByID(t provider.T) {
	t.Title("[Get] Get schedule by ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("Retrieve schedule by ID", func(sCtx provider.StepCtx) {
		schedule := builder.NewScheduleBuilder().
			SetID(uuid.New()).
			SetClientID(uuid.New()).
			SetTrainingID(uuid.New()).
			SetDayOfTheWeek("Monday").
			SetStartTime("10:00").
			SetEndTime("11:00").
			Build()

		s.mock.ExpectQuery(`^SELECT \* FROM "schedules" WHERE "schedules"."id" = \$1 ORDER BY "schedules"."id" LIMIT \$2`).
			WithArgs(schedule.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "client_id", "training_id", "day_of_the_week", "start_time", "end_time"}).
				AddRow(schedule.ID, schedule.ClientID, schedule.TrainingID, schedule.DayOfTheWeek, schedule.StartTime, schedule.EndTime))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "scheduleID", schedule.ID)

		result, err := s.scheduleRepo.GetScheduleByID(ctx, schedule.ID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(schedule.ID, result.ID)
		sCtx.Assert().Equal(schedule.ClientID, result.ClientID)
		sCtx.Assert().Equal(schedule.TrainingID, result.TrainingID)
		sCtx.Assert().Equal(schedule.DayOfTheWeek, result.DayOfTheWeek)
		sCtx.Assert().Equal(schedule.StartTime, result.StartTime)
		sCtx.Assert().Equal(schedule.EndTime, result.EndTime)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *ScheduleRepoSuite) TestListSchedulesByClientID(t provider.T) {
	t.Title("[List] List schedules by client ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("List all schedules for a specific client", func(sCtx provider.StepCtx) {
		clientID := uuid.New()
		trainingID := uuid.New()
		schedule := builder.NewScheduleBuilder().
			SetClientID(clientID).
			SetTrainingID(trainingID).
			SetDayOfTheWeek("Monday").
			SetStartTime("10:00").
			SetEndTime("11:00").
			Build()

		s.mock.ExpectQuery(`^SELECT \* FROM "schedules" WHERE "schedules"."client_id" = \$1`).
			WithArgs(clientID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "client_id", "training_id", "day_of_the_week", "start_time", "end_time"}).
				AddRow(schedule.ID, schedule.ClientID, schedule.TrainingID, schedule.DayOfTheWeek, schedule.StartTime, schedule.EndTime))

		s.mock.ExpectQuery(`^SELECT \* FROM "trainings" WHERE "trainings"."id" = \$1`).
			WithArgs(trainingID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(trainingID))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "clientID", clientID)

		results, err := s.scheduleRepo.ListSchedulesByClientID(ctx, clientID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Len(results, 1)
		sCtx.Assert().Equal(schedule.ID, results[0].ID)
		sCtx.Assert().Equal(schedule.ClientID, results[0].ClientID)
		sCtx.Assert().Equal(schedule.TrainingID, results[0].TrainingID)
		sCtx.Assert().Equal(schedule.DayOfTheWeek, results[0].DayOfTheWeek)
		sCtx.Assert().Equal(schedule.StartTime, results[0].StartTime)
		sCtx.Assert().Equal(schedule.EndTime, results[0].EndTime)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func TestScheduleSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ScheduleRepoSuite))
}
