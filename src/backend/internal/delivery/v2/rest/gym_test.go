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

type GymHandlerSuite struct {
	suite.Suite
}

func (s *GymHandlerSuite) TestGymHandlerGetGymByID(t provider.T) {
	t.Title("[GetGymByID] Successfully retrieves gym by ID and handles errors")
	t.Tags("gym", "handler", "getByID")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieves gym by ID", func(sCtx provider.StepCtx) {
		gymServiceMock := &mocks.IGymService{}
		gymHandler := NewGymHandler(gymServiceMock)

		router := gin.Default()
		router.GET("/gyms/:gymId", gymHandler.GetGymByID)

		gymID := uuid.New()
		gym := builder.NewGymBuilder().
			SetID(gymID).
			Build()

		gymServiceMock.On("GetGymByID", mock.Anything, gymID).Return(gym, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/gyms/%s", gymID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
		expected, _ := json.Marshal(gin.H{"gym": gym})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})

	t.WithNewStep("Incorrect: invalid gym ID format", func(sCtx provider.StepCtx) {
		gymServiceMock := &mocks.IGymService{}
		gymHandler := NewGymHandler(gymServiceMock)

		router := gin.Default()
		router.GET("/gyms/:gymId", gymHandler.GetGymByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/gyms/invalid-uuid", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid UUID length: 12"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *GymHandlerSuite) TestGymHandlerGetGyms(t provider.T) {
	t.Title("[GetGyms] Successfully lists gyms and handles errors")
	t.Tags("gym", "handler", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully lists gyms", func(sCtx provider.StepCtx) {
		gymServiceMock := &mocks.IGymService{}
		gymHandler := NewGymHandler(gymServiceMock)

		router := gin.Default()
		router.GET("/gyms", gymHandler.GetGyms)

		gyms := []entity.Gym{
			builder.NewGymBuilder().Build(),
			builder.NewGymBuilder().Build(),
		}

		gymServiceMock.On("ListGyms", mock.Anything).Return(gyms, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/gyms", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
		expected, _ := json.Marshal(gin.H{"gyms": gyms})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *GymHandlerSuite) TestGymHandlerPostGym(t provider.T) {
	t.Title("[PostGym] Successfully posts gym and handles errors")
	t.Tags("gym", "handler", "post")
	t.Parallel()

	t.WithNewStep("Correct: successfully posts new gym", func(sCtx provider.StepCtx) {
		gymServiceMock := &mocks.IGymService{}
		gymHandler := NewGymHandler(gymServiceMock)

		router := gin.Default()
		router.POST("/gyms", gymHandler.PostGym)

		newGym := builder.NewGymBuilder().Build()

		reqBody, _ := json.Marshal(dto.PostGym{
			Id:      newGym.ID,
			Name:    newGym.Name,
			Phone:   newGym.Phone,
			City:    newGym.City,
			Addres:  newGym.Addres,
			IsChain: newGym.IsChain,
		})

		gymServiceMock.On("RegisterNewGym", mock.Anything, mock.AnythingOfType("entity.Gym")).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/gyms", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: missing or invalid request body", func(sCtx provider.StepCtx) {
		gymServiceMock := &mocks.IGymService{}
		gymHandler := NewGymHandler(gymServiceMock)

		router := gin.Default()
		router.POST("/gyms", gymHandler.PostGym)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/gyms", nil)
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid request"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *GymHandlerSuite) TestGymHandlerChangeGym(t provider.T) {
	t.Title("[ChangeGym] Successfully updates gym and handles errors")
	t.Tags("gym", "handler", "put")
	t.Parallel()

	t.WithNewStep("Correct: successfully updates gym", func(sCtx provider.StepCtx) {
		gymServiceMock := &mocks.IGymService{}
		gymHandler := NewGymHandler(gymServiceMock)

		router := gin.Default()
		router.PUT("/gyms/:gymId", gymHandler.ChangeGym)

		gymID := uuid.New()
		updatedGym := builder.NewGymBuilder().
			SetID(gymID).
			Build()

		reqBody, _ := json.Marshal(dto.PutGym{
			Name:    updatedGym.Name,
			Phone:   updatedGym.Phone,
			City:    updatedGym.City,
			Addres:  updatedGym.Addres,
			IsChain: updatedGym.IsChain,
		})

		gymServiceMock.On("ChangeGym", mock.Anything, mock.AnythingOfType("entity.Gym")).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/gyms/%s", gymID.String()), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid gym ID format", func(sCtx provider.StepCtx) {
		gymServiceMock := &mocks.IGymService{}
		gymHandler := NewGymHandler(gymServiceMock)

		router := gin.Default()
		router.PUT("/gyms/:gymId", gymHandler.ChangeGym)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/gyms/invalid-uuid", nil)
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid request"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *GymHandlerSuite) TestGymHandlerDeleteGym(t provider.T) {
	t.Title("[DeleteGym] Successfully deletes gym and handles errors")
	t.Tags("gym", "handler", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deletes gym", func(sCtx provider.StepCtx) {
		gymServiceMock := &mocks.IGymService{}
		gymHandler := NewGymHandler(gymServiceMock)

		router := gin.Default()
		router.DELETE("/gyms/:gymId", gymHandler.DeleteGym)

		gymID := uuid.New()
		gymServiceMock.On("DeleteGym", mock.Anything, gymID).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/gyms/%s", gymID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid gym ID format", func(sCtx provider.StepCtx) {
		gymServiceMock := &mocks.IGymService{}
		gymHandler := NewGymHandler(gymServiceMock)

		router := gin.Default()
		router.DELETE("/gyms/:gymId", gymHandler.DeleteGym)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/gyms/invalid-uuid", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid UUID length: 12"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func TestGymHandlerSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(GymHandlerSuite))
}
