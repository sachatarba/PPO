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

type ScheduleHandlerSuite struct {
	suite.Suite
}

func (s *ScheduleHandlerSuite) TestScheduleHandlerGetSchedules(t provider.T) {
	t.Title("[GetSchedules] Successfully retrieves schedules by client ID and handles errors")
	t.Tags("schedule", "handler", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieves schedules by client ID", func(sCtx provider.StepCtx) {
		scheduleServiceMock := &mocks.IScheduleService{}
		scheduleHandler := NewScheduleHandler(scheduleServiceMock)

		router := gin.Default()
		router.GET("/clients/:clientId/schedules", scheduleHandler.GetSchedules)

		clientID := uuid.New()
		schedules := []entity.Schedule{
			builder.NewScheduleBuilder().
				SetClientID(clientID).
				Build(),
			builder.NewScheduleBuilder().
				SetClientID(clientID).
				Build(),
		}

		scheduleServiceMock.On("ListSchedulesByClientID", mock.Anything, clientID).Return(schedules, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/clients/%s/schedules", clientID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
		expected, _ := json.Marshal(gin.H{"schedules": schedules})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})

	t.WithNewStep("Incorrect: invalid client ID format", func(sCtx provider.StepCtx) {
		scheduleServiceMock := &mocks.IScheduleService{}
		scheduleHandler := NewScheduleHandler(scheduleServiceMock)

		router := gin.Default()
		router.GET("/clients/:clientId/schedules", scheduleHandler.GetSchedules)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/clients/invalid-uuid/schedules", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid UUID length: 12"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *ScheduleHandlerSuite) TestScheduleHandlerPostSchedule(t provider.T) {
	t.Title("[PostSchedule] Successfully posts new schedule and handles errors")
	t.Tags("schedule", "handler", "post")
	t.Parallel()

	t.WithNewStep("Correct: successfully posts new schedule", func(sCtx provider.StepCtx) {
		scheduleServiceMock := &mocks.IScheduleService{}
		scheduleHandler := NewScheduleHandler(scheduleServiceMock)

		router := gin.Default()
		router.POST("/clients/:clientId/schedules", scheduleHandler.PostSchedule)

		newSchedule := builder.NewScheduleBuilder().Build()

		reqBody, _ := json.Marshal(dto.PostSchedule{
			Id:           newSchedule.ID,
			DayOfTheWeek: newSchedule.DayOfTheWeek,
			StartTime:    newSchedule.StartTime,
			EndTime:      newSchedule.EndTime,
			TrainingID:   newSchedule.TrainingID,
		})

		clientID := uuid.New()
		scheduleServiceMock.On("CreateNewSchedule", mock.Anything, mock.AnythingOfType("entity.Schedule")).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/clients/%s/schedules", clientID.String()), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid request body", func(sCtx provider.StepCtx) {
		scheduleServiceMock := &mocks.IScheduleService{}
		scheduleHandler := NewScheduleHandler(scheduleServiceMock)

		router := gin.Default()
		router.POST("/clients/:clientId/schedules", scheduleHandler.PostSchedule)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/clients/invalid-uuid/schedules", nil)
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid request"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *ScheduleHandlerSuite) TestScheduleHandlerPutSchedule(t provider.T) {
	t.Title("[PutSchedule] Successfully updates schedule and handles errors")
	t.Tags("schedule", "handler", "put")
	t.Parallel()

	t.WithNewStep("Correct: successfully updates schedule", func(sCtx provider.StepCtx) {
		scheduleServiceMock := &mocks.IScheduleService{}
		scheduleHandler := NewScheduleHandler(scheduleServiceMock)

		router := gin.Default()
		router.PUT("/clients/:clientId/schedules/:scheduleId", scheduleHandler.PutSchedule)

		clientID := uuid.New()
		scheduleID := uuid.New()
		updatedSchedule := builder.NewScheduleBuilder().
			SetID(scheduleID).
			SetClientID(clientID).
			Build()

		reqBody, _ := json.Marshal(dto.PutSchedule{
			DayOfTheWeek: updatedSchedule.DayOfTheWeek,
			StartTime:    updatedSchedule.StartTime,
			EndTime:      updatedSchedule.EndTime,
			TrainingID:   updatedSchedule.TrainingID,
		})

		scheduleServiceMock.On("ChangeSchedule", mock.Anything, mock.AnythingOfType("entity.Schedule")).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/clients/%s/schedules/%s", clientID.String(), scheduleID.String()), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid client or schedule ID format", func(sCtx provider.StepCtx) {
		scheduleServiceMock := &mocks.IScheduleService{}
		scheduleHandler := NewScheduleHandler(scheduleServiceMock)

		router := gin.Default()
		router.PUT("/clients/:clientId/schedules/:scheduleId", scheduleHandler.PutSchedule)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/clients/invalid-uuid/schedules/invalid-uuid", nil)
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid request"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *ScheduleHandlerSuite) TestScheduleHandlerDeleteSchedule(t provider.T) {
	t.Title("[DeleteSchedule] Successfully deletes schedule and handles errors")
	t.Tags("schedule", "handler", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deletes schedule", func(sCtx provider.StepCtx) {
		scheduleServiceMock := &mocks.IScheduleService{}
		scheduleHandler := NewScheduleHandler(scheduleServiceMock)

		router := gin.Default()
		router.DELETE("/clients/:clientId/schedules/:scheduleId", scheduleHandler.DeleteSchedule)

		scheduleID := uuid.New()
		scheduleServiceMock.On("DeleteSchedule", mock.Anything, scheduleID).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/clients/%s/schedules/%s", uuid.New().String(), scheduleID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid schedule ID format", func(sCtx provider.StepCtx) {
		scheduleServiceMock := &mocks.IScheduleService{}
		scheduleHandler := NewScheduleHandler(scheduleServiceMock)

		router := gin.Default()
		router.DELETE("/clients/:clientId/schedules/:scheduleId", scheduleHandler.DeleteSchedule)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/clients/invalid-uuid/schedules/invalid-uuid", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid UUID length: 12"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func TestScheduleSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ScheduleHandlerSuite))
}
