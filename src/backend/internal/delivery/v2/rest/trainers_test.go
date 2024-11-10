package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/delivery/v2/dto"
	"github.com/sachatarba/course-db/internal/entity"
	"github.com/sachatarba/course-db/internal/service/mocks"
	"github.com/sachatarba/course-db/internal/utils/builder"
	"github.com/stretchr/testify/mock"
)

type TrainerHandlerSuite struct {
	suite.Suite
}

func (s *TrainerHandlerSuite) TestTrainerHandlerGetTrainers(t provider.T) {
	t.Title("[GetTrainers] Successfully retrieves trainers and handles errors")
	t.Tags("trainer", "handler", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieves trainers", func(sCtx provider.StepCtx) {
		trainerServiceMock := &mocks.ITrainerService{}
		trainerHandler := NewTrainerHandler(trainerServiceMock)

		router := gin.Default()
		router.GET("/trainers", trainerHandler.GetTrainers)

		trainers := []entity.Trainer{
			builder.NewTrainerBuilder().Build(),
			builder.NewTrainerBuilder().Build(),
		}

		trainerServiceMock.On("ListTrainers", mock.Anything).Return(trainers, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/trainers", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
		expected, _ := json.Marshal(gin.H{"trainers": trainers})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})

	t.WithNewStep("Incorrect: internal server error", func(sCtx provider.StepCtx) {
		trainerServiceMock := &mocks.ITrainerService{}
		trainerHandler := NewTrainerHandler(trainerServiceMock)

		router := gin.Default()
		router.GET("/trainers", trainerHandler.GetTrainers)

		trainerServiceMock.On("ListTrainers", mock.Anything).Return(nil, fmt.Errorf("internal error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/trainers", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusInternalServerError, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "internal error"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *TrainerHandlerSuite) TestTrainerHandlerGetTrainersByGymID(t provider.T) {
	t.Title("[GetTrainersByGymID] Successfully retrieves trainers by gym ID and handles errors")
	t.Tags("trainer", "handler", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieves trainers by gym ID", func(sCtx provider.StepCtx) {
		trainerServiceMock := &mocks.ITrainerService{}
		trainerHandler := NewTrainerHandler(trainerServiceMock)

		router := gin.Default()
		router.GET("/gyms/:gymId/trainers", trainerHandler.GetTrainersByGymID)

		gymID := uuid.New()
		trainers := []entity.Trainer{
			builder.NewTrainerBuilder().Build(),
			builder.NewTrainerBuilder().Build(),
		}

		trainerServiceMock.On("ListTrainersByGymID", mock.Anything, gymID).Return(trainers, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/gyms/%s/trainers", gymID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
		expected, _ := json.Marshal(gin.H{"trainers": trainers})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})

	t.WithNewStep("Incorrect: invalid gym ID format", func(sCtx provider.StepCtx) {
		trainerServiceMock := &mocks.ITrainerService{}
		trainerHandler := NewTrainerHandler(trainerServiceMock)

		router := gin.Default()
		router.GET("/gyms/:gymId/trainers", trainerHandler.GetTrainersByGymID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/gyms/invalid-uuid/trainers", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid UUID length: 12"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *TrainerHandlerSuite) TestTrainerHandlerPostTrainer(t provider.T) {
	t.Title("[PostTrainer] Successfully posts new trainer and handles errors")
	t.Tags("trainer", "handler", "post")
	t.Parallel()

	t.WithNewStep("Correct: successfully posts new trainer", func(sCtx provider.StepCtx) {
		trainerServiceMock := &mocks.ITrainerService{}
		trainerHandler := NewTrainerHandler(trainerServiceMock)

		router := gin.Default()
		router.POST("/trainers", trainerHandler.PostTrainer)

		newTrainer := builder.NewTrainerBuilder().Build()

		reqBody, _ := json.Marshal(dto.PostTrainer{
			Id:            newTrainer.ID,
			Fullname:      newTrainer.Fullname,
			Email:         newTrainer.Email,
			Phone:         newTrainer.Phone,
			Qualification: newTrainer.Qualification,
			UnitPrice:     newTrainer.UnitPrice,
			GymsID:        newTrainer.GymsID,
		})

		trainerServiceMock.On("RegisterNewTrainer", mock.Anything, mock.AnythingOfType("entity.Trainer")).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/trainers", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid request body", func(sCtx provider.StepCtx) {
		trainerServiceMock := &mocks.ITrainerService{}
		trainerHandler := NewTrainerHandler(trainerServiceMock)

		router := gin.Default()
		router.POST("/trainers", trainerHandler.PostTrainer)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/trainers", nil)
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid request"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *TrainerHandlerSuite) TestTrainerHandlerPutTrainer(t provider.T) {
	t.Title("[PutTrainer] Successfully updates trainer and handles errors")
	t.Tags("trainer", "handler", "put")
	t.Parallel()

	t.WithNewStep("Correct: successfully updates trainer", func(sCtx provider.StepCtx) {
		trainerServiceMock := &mocks.ITrainerService{}
		trainerHandler := NewTrainerHandler(trainerServiceMock)

		router := gin.Default()
		router.PUT("/trainers/:trainerId", trainerHandler.PutTrainer)

		trainerID := uuid.New()
		updatedTrainer := builder.NewTrainerBuilder().
			SetID(trainerID).
			Build()

		reqBody, _ := json.Marshal(dto.PutTrainer{
			Fullname:      updatedTrainer.Fullname,
			Email:         updatedTrainer.Email,
			Phone:         updatedTrainer.Phone,
			Qualification: updatedTrainer.Qualification,
			UnitPrice:     updatedTrainer.UnitPrice,
			GymsID:        updatedTrainer.GymsID,
		})

		trainerServiceMock.On("ChangeTrainer", mock.Anything, mock.AnythingOfType("entity.Trainer")).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/trainers/%s", trainerID.String()), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid trainer ID format", func(sCtx provider.StepCtx) {
		trainerServiceMock := &mocks.ITrainerService{}
		trainerHandler := NewTrainerHandler(trainerServiceMock)

		router := gin.Default()
		router.PUT("/trainers/:trainerId", trainerHandler.PutTrainer)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/trainers/invalid-uuid", nil)
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid request"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *TrainerHandlerSuite) TestTrainerHandlerDeleteTrainer(t provider.T) {
	t.Title("[DeleteTrainer] Successfully deletes trainer and handles errors")
	t.Tags("trainer", "handler", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deletes trainer", func(sCtx provider.StepCtx) {
		trainerServiceMock := &mocks.ITrainerService{}
		trainerHandler := NewTrainerHandler(trainerServiceMock)

		router := gin.Default()
		router.DELETE("/trainers/:trainerId", trainerHandler.DeleteTrainer)

		trainerID := uuid.New()

		trainerServiceMock.On("DeleteTrainer", mock.Anything, trainerID).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/trainers/%s", trainerID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid trainer ID format", func(sCtx provider.StepCtx) {
		trainerServiceMock := &mocks.ITrainerService{}
		trainerHandler := NewTrainerHandler(trainerServiceMock)

		router := gin.Default()
		router.DELETE("/trainers/:trainerId", trainerHandler.DeleteTrainer)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/trainers/invalid-uuid", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid UUID length: 12"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func TestTrainerSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(TrainerHandlerSuite))
}
