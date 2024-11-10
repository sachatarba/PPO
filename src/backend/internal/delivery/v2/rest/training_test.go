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

type TrainingHandlerSuite struct {
	suite.Suite
}

func (s *TrainingHandlerSuite) TestTrainingHandlerGetTrainingsByTrainerID(t provider.T) {
	t.Title("[GetTrainingsByTrainerID] Successfully retrieves trainings by trainer ID and handles errors")
	t.Tags("training", "handler", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieves trainings by trainer ID", func(sCtx provider.StepCtx) {
		trainingServiceMock := &mocks.ITrainingService{}
		trainingHandler := NewTrainingHandler(trainingServiceMock)

		router := gin.Default()
		router.GET("/trainers/:trainerId/trainings", trainingHandler.GetTrainingsByTrainerID)

		trainerID := uuid.New()
		trainings := []entity.Training{
			builder.NewTrainingBuilder().Build(),
			builder.NewTrainingBuilder().Build(),
		}

		trainingServiceMock.On("ListTrainingsByTrainerID", mock.Anything, trainerID).Return(trainings, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/trainers/%s/trainings", trainerID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
		expected, _ := json.Marshal(gin.H{"trainings": trainings})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})

	t.WithNewStep("Incorrect: invalid trainer ID format", func(sCtx provider.StepCtx) {
		trainingServiceMock := &mocks.ITrainingService{}
		trainingHandler := NewTrainingHandler(trainingServiceMock)

		router := gin.Default()
		router.GET("/trainers/:trainerId/trainings", trainingHandler.GetTrainingsByTrainerID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/trainers/invalid-uuid/trainings", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid UUID length: 12"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *TrainingHandlerSuite) TestTrainingHandlerPostTraining(t provider.T) {
	t.Title("[PostTraining] Successfully posts new training and handles errors")
	t.Tags("training", "handler", "post")
	t.Parallel()

	t.WithNewStep("Correct: successfully posts new training", func(sCtx provider.StepCtx) {
		trainingServiceMock := &mocks.ITrainingService{}
		trainingHandler := NewTrainingHandler(trainingServiceMock)

		router := gin.Default()
		router.POST("/trainers/:trainerId/trainings", trainingHandler.PostTraining)

		trainerID := uuid.New()
		newTraining := builder.NewTrainingBuilder().Build()

		reqBody, _ := json.Marshal(dto.PostTraining{
			Id:           newTraining.ID,
			Title:        newTraining.Title,
			Description:  newTraining.Description,
			TrainingType: newTraining.TrainingType,
		})

		trainingServiceMock.
			On("CreateNewTraining",
				mock.Anything,
				mock.AnythingOfType("entity.Training")).
			Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/trainers/%s/trainings",
			trainerID.String()),
			bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid request body", func(sCtx provider.StepCtx) {
		trainingServiceMock := &mocks.ITrainingService{}
		trainingHandler := NewTrainingHandler(trainingServiceMock)

		router := gin.Default()
		router.POST("/trainers/:trainerId/trainings", trainingHandler.PostTraining)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/trainers/invalid-uuid/trainings", nil)
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid request"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *TrainingHandlerSuite) TestTrainingHandlerPutTraining(t provider.T) {
	t.Title("[PutTraining] Successfully updates training and handles errors")
	t.Tags("training", "handler", "put")
	t.Parallel()

	t.WithNewStep("Correct: successfully updates training", func(sCtx provider.StepCtx) {
		trainingServiceMock := &mocks.ITrainingService{}
		trainingHandler := NewTrainingHandler(trainingServiceMock)

		router := gin.Default()
		router.PUT("/trainers/:trainerId/trainings/:trainingId", trainingHandler.PutTraining)

		trainerID := uuid.New()
		trainingID := uuid.New()
		updatedTraining := builder.NewTrainingBuilder().
			SetID(trainingID).
			SetTrainerID(trainerID).Build()

		reqBody, _ := json.Marshal(dto.PutTraining{
			Title:        updatedTraining.Title,
			Description:  updatedTraining.Description,
			TrainingType: updatedTraining.TrainingType,
		})

		trainingServiceMock.On("ChangeTraining", mock.Anything, mock.AnythingOfType("entity.Training")).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT",
			fmt.Sprintf("/trainers/%s/trainings/%s",
				trainerID.String(),
				trainingID.String()),
			bytes.NewBuffer(reqBody))

		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid training ID format", func(sCtx provider.StepCtx) {
		trainingServiceMock := &mocks.ITrainingService{}
		trainingHandler := NewTrainingHandler(trainingServiceMock)

		router := gin.Default()
		router.PUT("/trainers/:trainerId/trainings/:trainingId", trainingHandler.PutTraining)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/trainers/invalid-uuid/trainings/invalid-uuid", nil)
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid request"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *TrainingHandlerSuite) TestTrainingHandlerDeleteTraining(t provider.T) {
	t.Title("[DeleteTraining] Successfully deletes training and handles errors")
	t.Tags("training", "handler", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deletes training", func(sCtx provider.StepCtx) {
		trainingServiceMock := &mocks.ITrainingService{}
		trainingHandler := NewTrainingHandler(trainingServiceMock)

		router := gin.Default()
		router.DELETE("/trainers/:trainerId/trainings/:trainingId", trainingHandler.DeleteTraining)

		trainingID := uuid.New()

		trainingServiceMock.On("DeleteTraining", mock.Anything, trainingID).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE",
			fmt.Sprintf("/trainers/invalid-uuid/trainings/%s",
				trainingID.String()),
			nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid training ID format", func(sCtx provider.StepCtx) {
		trainingServiceMock := &mocks.ITrainingService{}
		trainingHandler := NewTrainingHandler(trainingServiceMock)

		router := gin.Default()
		router.DELETE("/trainers/:trainerId/trainings/:trainingId", trainingHandler.DeleteTraining)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/trainers/invalid-uuid/trainings/invalid-uuid", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid UUID length: 12"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func TestTrainingSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(TrainingHandlerSuite))
}
