package repository

// import (
// 	"context"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/sachatarba/course-db/internal/config"
// 	"github.com/sachatarba/course-db/internal/entity"
// 	"github.com/sachatarba/course-db/internal/orm"
// 	postrgres_adapter "github.com/sachatarba/course-db/internal/postrgres"
// 	"github.com/stretchr/testify/assert"
// )

// type testCase int

// const (
// 	REGISTER = iota
// 	CHANGE
// 	DELETE
// 	GET_BY_ID
// )

// func TestGymOperations(t *testing.T) {
// 	testCases := []struct {
// 		name            testCase
// 		gym             entity.Gym
// 		expectedName    string
// 		expectedCity    string
// 		expectedSuccess bool
// 	}{
// 		{
// 			name: REGISTER,
// 			gym: entity.Gym{
// 				ID:      uuid.New(),
// 				Name:    "Test Gym",
// 				Phone:   "+7-999-999-99-99",
// 				City:    "Test City",
// 				Addres:  "Test Address",
// 				IsChain: false,
// 			},
// 			expectedName:    "Test Gym",
// 			expectedCity:    "Test City",
// 			expectedSuccess: true,
// 		},
// 		{
// 			name: CHANGE,
// 			gym: entity.Gym{
// 				ID:      uuid.New(),
// 				Name:    "Test Gym",
// 				Phone:   "+7-999-999-99-99",
// 				City:    "Test City",
// 				Addres:  "Test Address",
// 				IsChain: false,
// 			},
// 			expectedName:    "Updated Gym Name",
// 			expectedSuccess: true,
// 		},
// 		{
// 			name: DELETE,
// 			gym: entity.Gym{
// 				ID:      uuid.New(),
// 				Name:    "Test Gym",
// 				Phone:   "+7-999-999-99-99",
// 				City:    "Test City",
// 				Addres:  "Test Address",
// 				IsChain: false,
// 			},
// 			expectedSuccess: true,
// 		},
// 		{
// 			name: GET_BY_ID,
// 			gym: entity.Gym{
// 				ID:      uuid.New(),
// 				Name:    "Test Gym",
// 				Phone:   "+7-999-999-99-99",
// 				City:    "Test City",
// 				Addres:  "Test Address",
// 				IsChain: false,
// 			},
// 			expectedName:    "Test Gym",
// 			expectedCity:    "Test City",
// 			expectedSuccess: true,
// 		},
// 	}

// 	ctx := context.Background()
// 	conf := config.NewConfFromEnv()
// 	// Подключение к тестовой базе данных PostgreSQL
// 	postgresConnector := postrgres_adapter.PostgresConnector{
// 		Conf: conf.PostgresConf,
// 	}

// 	db, err := postgresConnector.Connect()
// 	assert.NoError(t, err, "Error connection db")

// 	postgresMigrator := postrgres_adapter.PostgresMigrator{
// 		DB:     db,
// 		Tables: orm.TablesORM,
// 	}

// 	err = postgresMigrator.Migrate()
// 	assert.NoError(t, err, "Error migration db")

// 	repo := NewGymRepo(db)

// 	for _, tc := range testCases {
// 		switch tc.name {
// 		case REGISTER:
// 			err := repo.RegisterNewGym(ctx, tc.gym)
// 			assert.Equal(t, tc.expectedSuccess, err == nil)

// 		case CHANGE:
// 			err := repo.RegisterNewGym(ctx, tc.gym)
// 			assert.NoError(t, err)

// 			tc.gym.Name = "Updated Gym Name"

// 			err = repo.ChangeGym(ctx, tc.gym)
// 			assert.Equal(t, tc.expectedSuccess, err == nil)

// 		case DELETE:
// 			err := repo.RegisterNewGym(ctx, tc.gym)
// 			assert.NoError(t, err)

// 			err = repo.DeleteGym(ctx, tc.gym.ID)
// 			assert.Equal(t, tc.expectedSuccess, err == nil)

// 		case GET_BY_ID:
// 			err := repo.RegisterNewGym(ctx, tc.gym)
// 			assert.NoError(t, err)

// 			retGym, err := repo.GetGymByID(ctx, tc.gym.ID)
// 			assert.Equal(t, tc.expectedSuccess, err == nil)
// 			if tc.expectedSuccess {
// 				assert.Equal(t, tc.expectedName, retGym.Name)
// 				assert.Equal(t, tc.gym.Name, retGym.Name)
// 				assert.Equal(t, tc.gym.City, retGym.City)
// 			}
// 		}
// 	}
// }
