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

type EquipmentHandlerSuite struct {
	suite.Suite
}

func (s *EquipmentHandlerSuite) TestEquipmentHandlerGetEquipments(t provider.T) {
	t.Title("[GetEquipments] Successfully lists equipment and handles errors")
	t.Tags("equipment", "handler", "list")
	t.Parallel()

	t.WithNewStep("Correct: successfully lists equipment", func(sCtx provider.StepCtx) {
		equipmentServiceMock := &mocks.IEquipmentService{}
		equipmentHandler := NewEquipmentHandler(equipmentServiceMock)

		router := gin.Default()
		router.GET("/gyms/:gymId/equipments", equipmentHandler.GetEquipments)

		gymID := uuid.New()
		equipments := []entity.Equipment{
			builder.NewEquipmentBuilder().Build(),
			builder.NewEquipmentBuilder().Build(),
		}

		equipmentServiceMock.On("ListEquipmentsByGymID", mock.Anything, gymID).Return(equipments, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/gyms/%s/equipments", gymID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
		expected, _ := json.Marshal(gin.H{"equipments": equipments})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})

	t.WithNewStep("Incorrect: invalid gym ID format", func(sCtx provider.StepCtx) {
		equipmentServiceMock := &mocks.IEquipmentService{}
		equipmentHandler := NewEquipmentHandler(equipmentServiceMock)

		router := gin.Default()
		router.GET("/gyms/:gymId/equipments", equipmentHandler.GetEquipments)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/gyms/invalid-uuid/equipments", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid UUID length: 12"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *EquipmentHandlerSuite) TestEquipmentHandlerPostEquipment(t provider.T) {
	t.Title("[PostEquipment] Successfully posts equipment and handles errors")
	t.Tags("equipment", "handler", "post")
	t.Parallel()

	t.WithNewStep("Correct: successfully posts new equipment", func(sCtx provider.StepCtx) {
		equipmentServiceMock := &mocks.IEquipmentService{}
		equipmentHandler := NewEquipmentHandler(equipmentServiceMock)

		router := gin.Default()
		router.POST("/gyms/:gymId/equipments", equipmentHandler.PostEquipment)

		gymID := uuid.New()
		newEquipment := builder.NewEquipmentBuilder().Build()

		reqBody, _ := json.Marshal(dto.PostEquipment{
			Id:          newEquipment.ID,
			Name:        newEquipment.Name,
			Description: newEquipment.Description,
		})

		equipmentServiceMock.On("CreateNewEquipment",
			mock.Anything,
			mock.AnythingOfType("entity.Equipment")).
			Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/gyms/%s/equipments", gymID.String()), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: missing or invalid body", func(sCtx provider.StepCtx) {
		equipmentServiceMock := &mocks.IEquipmentService{}
		equipmentHandler := NewEquipmentHandler(equipmentServiceMock)

		router := gin.Default()
		router.POST("/gyms/:gymId/equipments", equipmentHandler.PostEquipment)

		gymID := uuid.New()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/gyms/%s/equipments", gymID.String()), nil)
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid request"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *EquipmentHandlerSuite) TestEquipmentHandlerPutEquipment(t provider.T) {
	t.Title("[PutEquipment] Successfully updates equipment and handles errors")
	t.Tags("equipment", "handler", "put")
	t.Parallel()

	t.WithNewStep("Correct: successfully updates existing equipment", func(sCtx provider.StepCtx) {
		equipmentServiceMock := &mocks.IEquipmentService{}
		equipmentHandler := NewEquipmentHandler(equipmentServiceMock)

		router := gin.Default()
		router.PUT("/gyms/:gymId/equipments/:equipmentId", equipmentHandler.PutEquipment)

		gymID := uuid.New()
		equipmentID := uuid.New()
		updatedEquipment := builder.NewEquipmentBuilder().
			SetID(equipmentID).
			SetGymID(gymID).
			Build()

		reqBody, _ := json.Marshal(dto.PutEquipment{
			Name:        updatedEquipment.Name,
			Description: updatedEquipment.Description,
		})

		equipmentServiceMock.On("ChangeEquipment",
			mock.Anything,
			mock.AnythingOfType("entity.Equipment")).
			Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/gyms/%s/equipments/%s",
			gymID.String(),
			equipmentID.String()),
			bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid gym or equipment ID", func(sCtx provider.StepCtx) {
		equipmentServiceMock := &mocks.IEquipmentService{}
		equipmentHandler := NewEquipmentHandler(equipmentServiceMock)

		router := gin.Default()
		router.PUT("/gyms/:gymId/equipments/:equipmentId", equipmentHandler.PutEquipment)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/gyms/invalid-uuid/equipments/invalid-uuid", nil)
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid request"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *EquipmentHandlerSuite) TestEquipmentHandlerDeleteEquipment(t provider.T) {
	t.Title("[DeleteEquipment] Successfully deletes equipment and handles errors")
	t.Tags("equipment", "handler", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deletes equipment", func(sCtx provider.StepCtx) {
		equipmentServiceMock := &mocks.IEquipmentService{}
		equipmentHandler := NewEquipmentHandler(equipmentServiceMock)

		router := gin.Default()
		router.DELETE("/equipments/:equipmentId", equipmentHandler.DeleteEquipment)

		equipmentID := uuid.New()
		equipmentServiceMock.On("DeleteEquipment", mock.Anything, equipmentID).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/equipments/%s", equipmentID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid equipment ID format", func(sCtx provider.StepCtx) {
		equipmentServiceMock := &mocks.IEquipmentService{}
		equipmentHandler := NewEquipmentHandler(equipmentServiceMock)

		router := gin.Default()
		router.DELETE("/equipments/:equipmentId", equipmentHandler.DeleteEquipment)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/equipments/invalid-uuid", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid UUID length: 12"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func TestEquipmentHandlerSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(EquipmentHandlerSuite))
}
