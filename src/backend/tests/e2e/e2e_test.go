package e2e

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sachatarba/course-db/internal/config"
	"github.com/sachatarba/course-db/internal/delivery/v1/request"
	"github.com/sachatarba/course-db/internal/entity"

	"github.com/gavv/httpexpect/v2"

	postrgres_adapter "github.com/sachatarba/course-db/internal/postrgres"
)

type E2ESuite struct {
	suite.Suite

	expect httpexpect.Expect
}

func (s *E2ESuite) BeforeAll(t provider.T) {
	host := os.Getenv("GOLANG_HOST")
	s.expect = *httpexpect.WithConfig(httpexpect.Config{
		Client:   &http.Client{},
		BaseURL:  fmt.Sprintf("http://%s:8080/api/v1/", host),
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

func (s *E2ESuite) AfterAll(t provider.T) {
	conf := config.NewConfFromEnv()

	postgresConnector := postrgres_adapter.PostgresConnector{
		Conf: conf.PostgresConf,
	}

	db, err := postgresConnector.Connect()
	t.Assert().NoError(err, "Error connection db")

	tables, err := db.Migrator().GetTables()
	t.Assert().NoError(err)

	for _, table := range tables {
		err := db.Migrator().DropTable(table)
		t.Assert().NoError(err)
	}
}

type GymsResponse struct {
	Gyms []entity.Gym `json:"gyms"`
}

type ClientMembershipsResponse struct {
	ClientMemberships []entity.ClientMembership `json:"clientMemberships"`
}

func (s *E2ESuite) TestClientMembershipHandlerPostClientMembership(t provider.T) {
	if os.Getenv("SKIP") == "true" {
		t.Skip()
	}
	t.Title("[PostClientMembership] Successfully creates new client membership")
	t.Tags("client_membership", "handler", "post")
	t.Parallel()

	gym := &request.GymReq{
		ID:      uuid.New(),
		Name:    "Качалка",
		Phone:   "+7-999-999-99-99",
		City:    "Москва",
		Addres:  "Улица Боброва, д. 10",
		IsChain: true,
	}

	membershipType := &request.MembershipTypeReq{
		ID:           uuid.New(),
		Type:         "Босс оф гим",
		Description:  "Для настоящих мужчин",
		Price:        "100",
		DaysDuration: 30,
		GymID:        gym.ID,
	}

	client := &request.ClientReq{
		ID:        uuid.New(),
		Fullname:  "Люк Скайвокер",
		Birthdate: "10-10-1999",
		Login:     "best_jedi",
		Password:  "palpatin.net",
		Email:     "sas@mail.ru",
		Phone:     "+7-999-999-99-99",
	}

	clientMembership := &request.ClientMembershipReq{
		ID:               uuid.New(),
		StartDate:        "2024-12-11",
		EndDate:          "2024-12-20",
		MembershipTypeID: membershipType.ID,
		ClientID:         client.ID,
	}

	t.WithNewStep("Create new gym", func(sCtx provider.StepCtx) {
		status := s.expect.POST("/gym/new").
			WithJSON(gym).
			Expect().Raw().Status

		sCtx.Assert().Equal("200 OK", status)
	})

	t.WithNewStep("Get gym", func(sCtx provider.StepCtx) {
		resp := &GymsResponse{}
		s.expect.GET("/gym/all").
			// WithJSON(gym).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Decode(resp)

		// log.Println(resp)
		sCtx.Assert().Equal(gym.ID, resp.Gyms[0].ID)
		sCtx.Assert().Equal(gym.Addres, resp.Gyms[0].Addres)
		sCtx.Assert().Equal(gym.City, resp.Gyms[0].City)
		sCtx.Assert().Equal(gym.Name, resp.Gyms[0].Name)
		sCtx.Assert().Equal(gym.Phone, resp.Gyms[0].Phone)
		sCtx.Assert().Equal(gym.IsChain, resp.Gyms[0].IsChain)
	})

	t.WithNewStep("Create new membership type", func(sCtx provider.StepCtx) {
		status := s.expect.POST("/membershipType/new").
			WithJSON(membershipType).
			Expect().Raw().Status

		sCtx.Assert().Equal("200 OK", status)
	})

	t.WithNewStep("Register client", func(sCtx provider.StepCtx) {
		status := s.expect.POST("/register").
			WithJSON(client).
			Expect().Raw().Status

		sCtx.Assert().Equal("200 OK", status)
	})

	t.WithNewStep("Create new client membership", func(sCtx provider.StepCtx) {
		status := s.expect.POST("/client_membership/new").
			WithJSON(clientMembership).
			Expect().Raw().Status

		sCtx.Assert().Equal("200 OK", status)
	})

	t.WithNewStep("Get client memberships", func(sCtx provider.StepCtx) {
		resp := &ClientMembershipsResponse{}
		s.expect.GET(fmt.Sprintf("client_membership/%s", client.ID)).
			// WithJSON().
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Decode(resp)

		// log.Println(resp)
		sCtx.Assert().Equal(clientMembership.ID, resp.ClientMemberships[0].ID)
		// sCtx.Assert().Equal(clientMembership.StartDate resp.ClientMembership.)
		sCtx.Assert().Equal(clientMembership.MembershipTypeID, resp.ClientMemberships[0].MembershipType.ID)
		sCtx.Assert().Equal(clientMembership.ClientID, resp.ClientMemberships[0].ClientID)
		// sCtx.Assert().Equal(gym.Phone, resp.Gyms[0].Phone)
		// sCtx.Assert().Equal(gym.IsChain, resp.Gyms[0].IsChain)
	})
}


func TestE2e2SuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(E2ESuite))
}
