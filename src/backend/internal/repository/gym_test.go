package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/config"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/orm"
	postrgres_adapter "github.com/sachatarba/course-db/internal/postrgres"
	"github.com/sachatarba/course-db/internal/repository/mocks"
	"github.com/sachatarba/course-db/internal/service"
	"github.com/sachatarba/course-db/internal/utils/builder"
)

type GymRepoSuite struct {
	suite.Suite

	mock    sqlmock.Sqlmock
	gymRepo service.IGymRepository

	ctx  context.Context
	repo service.IGymRepository

	service service.IGymService
}

func (s *GymRepoSuite) BeforeEach(t provider.T) {
	t.Title("Init mock db")
	db, mock := mocks.NewMockDB()

	s.gymRepo = NewGymRepo(db)
	s.mock = mock

	s.service = service.NewGymService(s.gymRepo)

	t.Tags("fixture", "gym", "db")
}

func (s *GymRepoSuite) BeforeAll(t provider.T) {
	s.ctx = context.Background()
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

	s.repo = NewGymRepo(db)
}

func (s *GymRepoSuite) TestRegisterNewGym(t provider.T) {
	t.Title("[Create] Register new gym")
	t.Tags("repository", "postgres")

	t.WithNewStep("Register new gym", func(sCtx provider.StepCtx) {
		gym := builder.NewGymBuilder().Build()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`INSERT INTO "gyms" \("id","name","phone","city","addres","is_chain"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\)`).
			WithArgs(
				gym.ID,
				gym.Name,
				gym.Phone,
				gym.City,
				gym.Addres,
				fmt.Sprintf("%t", gym.IsChain),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "gym", gym)

		err := s.gymRepo.RegisterNewGym(ctx, gym)

		sCtx.Assert().NoError(err)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *GymRepoSuite) TestChangeGym(t provider.T) {
	t.Title("[Change] Change gym details")
	t.Tags("repository", "postgres")

	t.WithNewStep("Update gym details", func(sCtx provider.StepCtx) {
		gym := builder.NewGymBuilder().Build()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`UPDATE "gyms" SET "name"=\$1,"phone"=\$2,"city"=\$3,"addres"=\$4,"is_chain"=\$5 WHERE "id" = \$6`).
			WithArgs(
				gym.Name,
				gym.Phone,
				gym.City,
				gym.Addres,
				fmt.Sprintf("%t", gym.IsChain),
				gym.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "gym", gym)

		err := s.gymRepo.ChangeGym(ctx, gym)

		sCtx.Assert().NoError(err)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *GymRepoSuite) TestDeleteGym(t provider.T) {
	t.Title("[Delete] Delete gym by ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("Delete gym by ID", func(sCtx provider.StepCtx) {
		gymID := uuid.New()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`DELETE FROM "gyms" WHERE "gyms"."id" = \$1`).
			WithArgs(gymID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "gymID", gymID)

		err := s.gymRepo.DeleteGym(ctx, gymID)

		sCtx.Assert().NoError(err)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *GymRepoSuite) TestGetGymByID(t provider.T) {
	t.Title("[Get] Get gym by ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("Retrieve gym by ID", func(sCtx provider.StepCtx) {
		gym := builder.NewGymBuilder().Build()

		s.mock.ExpectQuery(`^SELECT \* FROM "gyms" WHERE "gyms"."id" = \$1 ORDER BY "gyms"."id" LIMIT \$2$`).
			WithArgs(gym.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "addres", "phone"}).
				AddRow(gym.ID, gym.Name, gym.Addres, gym.Phone))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "gymID", gym.ID)

		result, err := s.gymRepo.GetGymByID(ctx, gym.ID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(gym.ID, result.ID)
		sCtx.Assert().Equal(gym.Name, result.Name)
		sCtx.Assert().Equal(gym.Addres, result.Addres)
		sCtx.Assert().Equal(gym.Phone, result.Phone)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *GymRepoSuite) TestListGyms(t provider.T) {
	t.Title("[List] List all gyms")
	t.Tags("repository", "postgres")

	t.WithNewStep("List all gyms", func(sCtx provider.StepCtx) {
		gym1 := builder.NewGymBuilder().Build()
		gym2 := builder.NewGymBuilder().Build()

		s.mock.ExpectQuery(`^SELECT \* FROM "gyms"$`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "addres", "phone"}).
				AddRow(gym1.ID, gym1.Name, gym1.Addres, gym1.Phone).
				AddRow(gym2.ID, gym2.Name, gym2.Addres, gym2.Phone))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx)

		results, err := s.gymRepo.ListGyms(ctx)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Len(results, 2)
		sCtx.Assert().Equal(gym1.ID, results[0].ID)
		sCtx.Assert().Equal(gym2.ID, results[1].ID)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *GymRepoSuite) TestChangeGym1(t provider.T) {
	t.Title("[Change] Change gym details")
	t.Tags("repository", "postgres")

	t.WithNewStep("Update gym details", func(sCtx provider.StepCtx) {
		gym := builder.NewGymBuilder().Build()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`UPDATE "gyms" SET "name"=\$1,"phone"=\$2,"city"=\$3,"addres"=\$4,"is_chain"=\$5 WHERE "id" = \$6`).
			WithArgs(
				gym.Name,
				gym.Phone,
				gym.City,
				gym.Addres,
				fmt.Sprintf("%t", gym.IsChain),
				gym.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "gym", gym)

		err := s.service.ChangeGym(ctx, gym)

		sCtx.Assert().NoError(err)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *GymRepoSuite) TestDeleteGym1(t provider.T) {
	t.Title("[Delete] Delete gym by ID MOOOOOOOOCK")
	t.Tags("repository", "postgres")

	t.WithNewStep("Delete gym by ID", func(sCtx provider.StepCtx) {
		gymID := uuid.New()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`DELETE FROM "gyms" WHERE "gyms"."id" = \$1`).
			WithArgs(gymID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "gymID", gymID)
		err := s.service.DeleteGym(ctx, gymID)

		sCtx.Assert().NoError(err)
		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

func (s *GymRepoSuite) TestGetGymByID1(t provider.T) {
	t.Title("[Get] Get gym by ID")
	t.Tags("repository", "postgres")

	t.WithNewStep("Retrieve gym by ID", func(sCtx provider.StepCtx) {
		gym := builder.NewGymBuilder().Build()

		s.mock.ExpectQuery(`^SELECT \* FROM "gyms" WHERE "gyms"."id" = \$1 ORDER BY "gyms"."id" LIMIT \$2$`).
			WithArgs(gym.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "addres", "phone"}).
				AddRow(gym.ID, gym.Name, gym.Addres, gym.Phone))

		ctx := context.TODO()
		sCtx.WithNewParameters("ctx", ctx, "gymID", gym.ID)

		result, err := s.service.GetGymByID(ctx, gym.ID)

		sCtx.Assert().NoError(err)
		sCtx.Assert().Equal(gym.ID, result.ID)
		sCtx.Assert().Equal(gym.Name, result.Name)
		sCtx.Assert().Equal(gym.Addres, result.Addres)
		sCtx.Assert().Equal(gym.Phone, result.Phone)

		err = s.mock.ExpectationsWereMet()
		sCtx.Assert().NoError(err)
	})
}

type testCase int

const (
	REGISTER = iota
	CHANGE
	DELETE
	GET_BY_ID
)

func (s *GymRepoSuite) TestGymOperations(t provider.T) {
	testCases := []struct {
		name            testCase
		gym             entity.Gym
		expectedName    string
		expectedCity    string
		expectedSuccess bool
	}{
		{
			name: REGISTER,
			gym: entity.Gym{
				ID:      uuid.New(),
				Name:    "Test Gym",
				Phone:   "+7-999-999-99-99",
				City:    "Test City",
				Addres:  "Test Address",
				IsChain: false,
			},
			expectedName:    "Test Gym",
			expectedCity:    "Test City",
			expectedSuccess: true,
		},
		{
			name: CHANGE,
			gym: entity.Gym{
				ID:      uuid.New(),
				Name:    "Test Gym",
				Phone:   "+7-999-999-99-99",
				City:    "Test City",
				Addres:  "Test Address",
				IsChain: false,
			},
			expectedName:    "Updated Gym Name",
			expectedSuccess: true,
		},
		{
			name: DELETE,
			gym: entity.Gym{
				ID:      uuid.New(),
				Name:    "Test Gym",
				Phone:   "+7-999-999-99-99",
				City:    "Test City",
				Addres:  "Test Address",
				IsChain: false,
			},
			expectedSuccess: true,
		},
		{
			name: GET_BY_ID,
			gym: entity.Gym{
				ID:      uuid.New(),
				Name:    "Test Gym",
				Phone:   "+7-999-999-99-99",
				City:    "Test City",
				Addres:  "Test Address",
				IsChain: false,
			},
			expectedName:    "Test Gym",
			expectedCity:    "Test City",
			expectedSuccess: true,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Gym test with data base: %d",tc.name), func(t provider.T) {
			t.WithNewStep("Gym case", func(sCtx provider.StepCtx) {
				switch tc.name {
				case REGISTER:
					err := s.repo.RegisterNewGym(s.ctx, tc.gym)
					t.Assert().Equal(tc.expectedSuccess, err == nil)

				case CHANGE:
					err := s.repo.RegisterNewGym(s.ctx, tc.gym)
					t.Assert().NoError(err)

					tc.gym.Name = "Updated Gym Name"

					err = s.repo.ChangeGym(s.ctx, tc.gym)
					t.Assert().Equal(tc.expectedSuccess, err == nil)

				case DELETE:
					err := s.repo.RegisterNewGym(s.ctx, tc.gym)
					t.Assert().NoError(err)

					err = s.repo.DeleteGym(s.ctx, tc.gym.ID)
					t.Assert().Equal(tc.expectedSuccess, err == nil)

				case GET_BY_ID:
					err := s.repo.RegisterNewGym(s.ctx, tc.gym)
					t.Assert().NoError(err)

					retGym, err := s.repo.GetGymByID(s.ctx, tc.gym.ID)
					t.Assert().Equal(tc.expectedSuccess, err == nil)
					if tc.expectedSuccess {
						t.Assert().Equal(tc.expectedName, retGym.Name)
						t.Assert().Equal(tc.gym.Name, retGym.Name)
						t.Assert().Equal(tc.gym.City, retGym.City)
					}
				}
			})
		})
	}
}

func TestGymSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(GymRepoSuite))
}
