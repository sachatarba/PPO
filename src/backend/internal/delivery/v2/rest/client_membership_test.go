package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

type ClientMembershipHandlerSuite struct {
	suite.Suite
}

func (s *ClientMembershipHandlerSuite) TestClientMembershipHandlerGetClientMemberships(t provider.T) {
	t.Title("[GetClientMemberships] Successfully get clients membership by client ID")
	t.Tags("training", "service", "list")
	t.Parallel()

	t.WithNewStep("Correct: get clients membership by client ID", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientMembershipsService{}
		clientMembershipHandler := NewClientMembershipHandler(clientServiceMock)

		router := gin.Default()
		router.GET("/clients/:clientId/client_memberships", clientMembershipHandler.GetClientMemberships)

		clientID := uuid.New()
		clientMemberships := []entity.ClientMembership{
			builder.NewClientMembershipBuilder().Build(),
			builder.NewClientMembershipBuilder().Build(),
		}

		clientServiceMock.On("ListClientMembershipsByClientID", mock.Anything, clientID).
			Return(clientMemberships, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/clients/%s/client_memberships", clientID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
		expected, _ := json.Marshal(gin.H{"clientMemberships": clientMemberships})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})

	t.WithNewStep("Incorrect: get clients membership by client ID", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientMembershipsService{}
		clientMembershipHandler := NewClientMembershipHandler(clientServiceMock)

		router := gin.Default()
		router.GET("/clients/:clientId/client_memberships", clientMembershipHandler.GetClientMemberships)

		clientID := uuid.New()
		clientMemberships := []entity.ClientMembership{}

		err := errors.New("service error")

		clientServiceMock.On("ListClientMembershipsByClientID", mock.Anything, clientID).
			Return(clientMemberships, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/clients/%s/client_memberships", clientID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusInternalServerError, w.Code)
		expected, _ := json.Marshal(gin.H{"err": err.Error()})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *ClientMembershipHandlerSuite) TestClientMembershipHandlerPostClientMembership(t provider.T) {
	t.Title("[PostClientMembership] Successfully creates new client membership")
	t.Tags("client_membership", "handler", "post")
	t.Parallel()

	t.WithNewStep("Correct: successfully creates client membership", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientMembershipsService{}
		clientMembershipHandler := NewClientMembershipHandler(clientServiceMock)

		router := gin.Default()
		router.POST("/clients/:clientId/client_memberships", clientMembershipHandler.PostClientMembership)

		clientID := uuid.New()
		reqBody := dto.PostClientMembership{
			Id:               uuid.New(),
			StartDate:        time.Now().Format(time.DateOnly),
			EndDate:          time.Now().AddDate(0, 1, 0).Format(time.DateOnly),
			MembershipTypeID: uuid.New(),
		}

		clientServiceMock.On("CreateNewClientMembership", mock.Anything, mock.AnythingOfType("entity.ClientMembership")).
			Return(nil)

		body, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/clients/%s/client_memberships", clientID.String()), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid client ID format", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientMembershipsService{}
		clientMembershipHandler := NewClientMembershipHandler(clientServiceMock)

		router := gin.Default()
		router.POST("/clients/:clientId/client_memberships", clientMembershipHandler.PostClientMembership)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/clients/invalid-uuid/client_memberships", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
	})
}

func (s *ClientMembershipHandlerSuite) TestClientMembershipHandlerPutClientMembership(t provider.T) {
	t.Title("[PutClientMembership] Successfully updates client membership")
	t.Tags("client_membership", "handler", "put")
	t.Parallel()

	t.WithNewStep("Correct: successfully updates client membership", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientMembershipsService{}
		clientMembershipHandler := NewClientMembershipHandler(clientServiceMock)

		router := gin.Default()
		router.PUT("/clients/:clientId/client_memberships/:clientMembershipId", clientMembershipHandler.PutClientMembership)

		clientID := uuid.New()
		clientMembershipID := uuid.New()
		reqBody := dto.PutClientMembership{
			StartDate:        time.Now().Format(time.DateOnly),
			EndDate:          time.Now().AddDate(0, 1, 0).Format(time.DateOnly),
			MembershipTypeID: uuid.New(),
		}

		clientServiceMock.On("ChangeClientMembership", mock.Anything, mock.AnythingOfType("entity.ClientMembership")).
			Return(nil)

		body, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/clients/%s/client_memberships/%s", clientID.String(), clientMembershipID.String()), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid client membership ID format", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientMembershipsService{}
		clientMembershipHandler := NewClientMembershipHandler(clientServiceMock)

		router := gin.Default()
		router.PUT("/clients/:clientId/client_memberships/:clientMembershipId", clientMembershipHandler.PutClientMembership)

		clientID := uuid.New()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/clients/%s/client_memberships/invalid-uuid", clientID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
	})
}

func (s *ClientMembershipHandlerSuite) TestClientMembershipHandlerDeleteClientMembership(t provider.T) {
	t.Title("[DeleteClientMembership] Successfully deletes client membership")
	t.Tags("client_membership", "handler", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deletes client membership", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientMembershipsService{}
		clientMembershipHandler := NewClientMembershipHandler(clientServiceMock)

		router := gin.Default()
		router.DELETE("/clients/:clientId/client_memberships/:clientMembershipId", clientMembershipHandler.DeleteClientMembership)

		clientMembershipID := uuid.New()
		clientServiceMock.On("DeleteClientMembership", mock.Anything, clientMembershipID).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/clients/%s/client_memberships/%s", uuid.New().String(), clientMembershipID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid client membership ID format", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientMembershipsService{}
		clientMembershipHandler := NewClientMembershipHandler(clientServiceMock)

		router := gin.Default()
		router.DELETE("/clients/:clientId/client_memberships/:clientMembershipId", clientMembershipHandler.DeleteClientMembership)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/clients/client-id/client_memberships/invalid-uuid", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
	})
}

func TestClientMembershipsHandlerSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ClientMembershipHandlerSuite))
}
