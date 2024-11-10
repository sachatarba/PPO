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
	"github.com/stretchr/testify/mock"
)

type ClientHandlerSuite struct {
	suite.Suite
}

func (s *ClientHandlerSuite) TestClientHandlerGetClients(t provider.T) {
	t.Title("[GetClients] Successfully lists clients")
	t.Tags("client", "handler", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully lists clients", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientService{}
		clientHandler := NewClientHandler(clientServiceMock)

		router := gin.Default()
		router.GET("/clients", clientHandler.GetClients)

		clients := []entity.Client{
			{ID: uuid.New(), Login: "client1", Fullname: "Client One"},
			{ID: uuid.New(), Login: "client2", Fullname: "Client Two"},
		}

		clientServiceMock.On("ListClients", mock.Anything).Return(clients, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/clients", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
		expected, _ := json.Marshal(gin.H{"clients": clients})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})

	t.WithNewStep("Incorrect: service error", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientService{}
		clientHandler := NewClientHandler(clientServiceMock)

		router := gin.Default()
		router.GET("/clients", clientHandler.GetClients)

		clientServiceMock.On("ListClients", mock.Anything).Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/clients", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusInternalServerError, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "service error"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *ClientHandlerSuite) TestClientHandlerGetClientByLogin(t provider.T) {
	t.Title("[GetClientByLogin] Successfully retrieves client by login")
	t.Tags("client", "handler", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieves client by login", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientService{}
		clientHandler := NewClientHandler(clientServiceMock)

		router := gin.Default()
		router.GET("/clients/login/:login", clientHandler.GetClientByLogin)

		login := "client1"
		client := entity.Client{ID: uuid.New(), Login: login, Fullname: "Client One"}

		clientServiceMock.On("GetClientByLogin", mock.Anything, login).Return(client, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/clients/login/%s", login), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
		expected, _ := json.Marshal(gin.H{"client": client})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})

	t.WithNewStep("Incorrect: service error", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientService{}
		clientHandler := NewClientHandler(clientServiceMock)

		router := gin.Default()
		router.GET("/clients/login/:login", clientHandler.GetClientByLogin)

		login := "nonexistent"
		clientServiceMock.On("GetClientByLogin", mock.Anything, login).Return(entity.Client{}, errors.New("service error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/clients/login/%s", login), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusInternalServerError, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "service error"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func (s *ClientHandlerSuite) TestClientHandlerGetClientByID(t provider.T) {
	t.Title("[GetClientByID] Successfully retrieves client by ID")
	t.Tags("client", "handler", "get")
	t.Parallel()

	t.WithNewStep("Correct: successfully retrieves client by ID", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientService{}
		clientHandler := NewClientHandler(clientServiceMock)

		router := gin.Default()
		router.GET("/clients/:clientId", clientHandler.GetClientByID)

		clientID := uuid.New()
		client := entity.Client{ID: clientID, Login: "client1", Fullname: "Client One"}

		clientServiceMock.On("GetClientByID", mock.Anything, clientID).Return(client, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/clients/%s", clientID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
		expected, _ := json.Marshal(gin.H{"client": client})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})

	t.WithNewStep("Incorrect: invalid client ID format", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientService{}
		clientHandler := NewClientHandler(clientServiceMock)

		router := gin.Default()
		router.GET("/clients/:clientId", clientHandler.GetClientByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/clients/invalid-uuid", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
	})
}

func (s *ClientHandlerSuite) TestClientHandlerPutClient(t provider.T) {
	t.Title("[PutClient] Successfully updates client information")
	t.Tags("client", "handler", "put")
	t.Parallel()

	t.WithNewStep("Correct: successfully updates client", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientService{}
		clientHandler := NewClientHandler(clientServiceMock)

		router := gin.Default()
		router.PUT("/clients/:clientId", clientHandler.PutClient)

		clientID := uuid.New()
		reqBody := dto.PutClient{
			Login:     "updatedLogin",
			Password:  "newPassword",
			Fullname:  "Updated Name",
			Email:     "newemail@example.com",
			Phone:     "+7-999-999-99-99",
			Birthdate: time.Now().Format(time.DateOnly),
		}

		clientServiceMock.On("ChangeClient", mock.Anything, mock.AnythingOfType("entity.Client")).
			Return(nil)

		body, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/clients/%s", clientID.String()), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid client ID format", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientService{}
		clientHandler := NewClientHandler(clientServiceMock)

		router := gin.Default()
		router.PUT("/clients/:clientId", clientHandler.PutClient)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/clients/invalid-uuid", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
	})
}

func (s *ClientHandlerSuite) TestClientHandlerDeleteClient(t provider.T) {
	t.Title("[DeleteClient] Successfully deletes client and handles errors")
	t.Tags("client", "handler", "delete")
	t.Parallel()

	t.WithNewStep("Correct: successfully deletes client", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientService{}
		clientHandler := NewClientHandler(clientServiceMock)

		router := gin.Default()
		router.DELETE("/clients/:clientId", clientHandler.DeleteClient)

		clientID := uuid.New()
		clientServiceMock.On("DeleteClient", mock.Anything, clientID).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/clients/%s", clientID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusOK, w.Code)
	})

	t.WithNewStep("Incorrect: invalid client ID format", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientService{}
		clientHandler := NewClientHandler(clientServiceMock)

		router := gin.Default()
		router.DELETE("/clients/:clientId", clientHandler.DeleteClient)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/clients/invalid-uuid", nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusBadRequest, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "invalid UUID length: 12"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})

	t.WithNewStep("Incorrect: client not found", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientService{}
		clientHandler := NewClientHandler(clientServiceMock)

		router := gin.Default()
		router.DELETE("/clients/:clientId", clientHandler.DeleteClient)

		clientID := uuid.New()
		clientServiceMock.On("DeleteClient", mock.Anything, clientID).Return(errors.New("client not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/clients/%s", clientID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusInternalServerError, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "client not found"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})

	t.WithNewStep("Incorrect: service error", func(sCtx provider.StepCtx) {
		clientServiceMock := &mocks.IClientService{}
		clientHandler := NewClientHandler(clientServiceMock)

		router := gin.Default()
		router.DELETE("/clients/:clientId", clientHandler.DeleteClient)

		clientID := uuid.New()
		clientServiceMock.On("DeleteClient", mock.Anything, clientID).Return(errors.New("service error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/clients/%s", clientID.String()), nil)

		sCtx.WithNewParameters("w", w, "req", req)
		router.ServeHTTP(w, req)

		sCtx.Assert().Equal(http.StatusInternalServerError, w.Code)
		expected, _ := json.Marshal(gin.H{"err": "service error"})
		sCtx.Assert().Equal(string(expected), w.Body.String())
	})
}

func TestClientHandlerSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ClientHandlerSuite))
}
