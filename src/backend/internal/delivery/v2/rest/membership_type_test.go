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

type MembershipTypeHandlerSuite struct {
	suite.Suite
}

func (s *MembershipTypeHandlerSuite) TestMembershipTypeHandlerGetMembershipTypeByGymID(t provider.T) {
	t.Title("[GetMembershipTypeByGymID] Successfully retrieves membership types by gym ID and handles errors")
	t.Tags("membershipType", "handler", "getByGymID")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieves membership types by gym ID", func(sCtx provider.StepCtx) {
		membershipTypeServiceMock := &mocks.IMembershipTypeService{}
		membershipTypeHandler := NewMembershipTypeHandler(membershipTypeServiceMock)

		router := gin.Default()
		router.GET("/gyms/:gymId/membershipTypes", membershipTypeHandler.GetMembershipTypeByGymID)

		gymID := uuid.New()
		membershipTypes := []entity.MembershipType{
			builder.NewMembershipTypeBuilder().SetGymID(gymID).Build(),
			builder.NewMembershipTypeBuilder().SetGymID(gymID).Build(),
		}

		membershipTypeServiceMock.On("ListMembershipTypesByGymID", mock.Anything, gymID).Return(membershipTypes, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/gyms/%s/membershipTypes", gymID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
		expected, _ := json.Marshal(gin.H{"membershipTypes": membershipTypes})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})

	t.WithNewStep("Incorrect: invalid gym ID format", func(sCtx provider.StepCtx) {
		membershipTypeServiceMock := &mocks.IMembershipTypeService{}
		membershipTypeHandler := NewMembershipTypeHandler(membershipTypeServiceMock)

		router := gin.Default()
		router.GET("/gyms/:gymId/membershipTypes", membershipTypeHandler.GetMembershipTypeByGymID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/gyms/invalid-uuid/membershipTypes", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid UUID length: 12"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *MembershipTypeHandlerSuite) TestMembershipTypeHandlerPostMembershipType(t provider.T) {
	t.Title("[PostMembershipType] Successfully posts membership type and handles errors")
	t.Tags("membershipType", "handler", "post")
	t.Parallel()

	t.WithNewStep("Correct: successfully posts new membership type", func(sCtx provider.StepCtx) {
		membershipTypeServiceMock := &mocks.IMembershipTypeService{}
		membershipTypeHandler := NewMembershipTypeHandler(membershipTypeServiceMock)

		router := gin.Default()
		router.POST("/gyms/:gymId/membershipTypes", membershipTypeHandler.PostMembershipType)

		newMembershipType := builder.NewMembershipTypeBuilder().Build()

		reqBody, _ := json.Marshal(dto.PostMembershipType{
			Id:           newMembershipType.ID,
			Type:         newMembershipType.Type,
			Description:  newMembershipType.Description,
			Price:        newMembershipType.Price,
			DaysDuration: newMembershipType.DaysDuration,
		})

		gymID := uuid.New()
		membershipTypeServiceMock.On("RegisterNewMembershipType", mock.Anything, mock.AnythingOfType("entity.MembershipType")).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/gyms/%s/membershipTypes", gymID.String()), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid request body", func(sCtx provider.StepCtx) {
		membershipTypeServiceMock := &mocks.IMembershipTypeService{}
		membershipTypeHandler := NewMembershipTypeHandler(membershipTypeServiceMock)

		router := gin.Default()
		router.POST("/gyms/:gymId/membershipTypes", membershipTypeHandler.PostMembershipType)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/gyms/invalid-uuid/membershipTypes", nil)
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid request"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *MembershipTypeHandlerSuite) TestMembershipTypeHandlerPutMembershipType(t provider.T) {
	t.Title("[PutMembershipType] Successfully updates membership type and handles errors")
	t.Tags("membershipType", "handler", "put")
	t.Parallel()

	t.WithNewStep("Correct: successfully updates membership type", func(sCtx provider.StepCtx) {
		membershipTypeServiceMock := &mocks.IMembershipTypeService{}
		membershipTypeHandler := NewMembershipTypeHandler(membershipTypeServiceMock)

		router := gin.Default()
		router.PUT("/gyms/:gymId/membershipTypes/:membershipTypeId", membershipTypeHandler.PutMembershipType)

		gymID := uuid.New()
		membershipTypeID := uuid.New()
		updatedMembershipType := builder.NewMembershipTypeBuilder().
			SetID(membershipTypeID).
			SetGymID(gymID).
			Build()

		reqBody, _ := json.Marshal(dto.PutMembershipType{
			Type:         updatedMembershipType.Type,
			Price:        updatedMembershipType.Price,
			DaysDuration: updatedMembershipType.DaysDuration,
			Description:  updatedMembershipType.Description,
		})

		membershipTypeServiceMock.
			On("ChangeMembershipType",
				mock.Anything,
				mock.AnythingOfType("entity.MembershipType")).
			Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT",
			fmt.Sprintf("/gyms/%s/membershipTypes/%s",
				gymID.String(),
				membershipTypeID.String()),
			bytes.NewBuffer(reqBody))

		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid gym ID format", func(sCtx provider.StepCtx) {
		membershipTypeServiceMock := &mocks.IMembershipTypeService{}
		membershipTypeHandler := NewMembershipTypeHandler(membershipTypeServiceMock)

		router := gin.Default()
		router.PUT("/gyms/:gymId/membershipTypes/:membershipTypeId", membershipTypeHandler.PutMembershipType)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/gyms/invalid-uuid/membershipTypes/invalid-uuid", nil)
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid request"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *MembershipTypeHandlerSuite) TestMembershipTypeHandlerDeleteMembershipType(t provider.T) {
	t.Title("[DeleteMembershipType] Successfully deletes membership type and handles errors")
	t.Tags("membershipType", "handler", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deletes membership type", func(sCtx provider.StepCtx) {
		membershipTypeServiceMock := &mocks.IMembershipTypeService{}
		membershipTypeHandler := NewMembershipTypeHandler(membershipTypeServiceMock)

		router := gin.Default()
		router.DELETE("/gyms/:gymId/membershipTypes/:membershipTypeId", membershipTypeHandler.DeleteMembershipType)

		membershipTypeID := uuid.New()
		membershipTypeServiceMock.On("DeleteMembershipType", mock.Anything, membershipTypeID).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/gyms/%s/membershipTypes/%s", uuid.New().String(), membershipTypeID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid membership type ID format", func(sCtx provider.StepCtx) {
		membershipTypeServiceMock := &mocks.IMembershipTypeService{}
		membershipTypeHandler := NewMembershipTypeHandler(membershipTypeServiceMock)

		router := gin.Default()
		router.DELETE("/gyms/:gymId/membershipTypes/:membershipTypeId", membershipTypeHandler.DeleteMembershipType)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/gyms/invalid-uuid/membershipTypes/invalid-uuid", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid UUID length: 12"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func TestMembershipTypeHandlerSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(MembershipTypeHandlerSuite))
}
