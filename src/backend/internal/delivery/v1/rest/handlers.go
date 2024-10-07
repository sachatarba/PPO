package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sachatarba/course-db/internal/config"
	"github.com/sachatarba/course-db/internal/service"
)

type PaymentHandler struct {
	ApiKey string
	ShopID string
}

func (h *PaymentHandler) InitPayment(router gin.IRouter) {
	router.POST("/create-payment", h.CreatePayment)
}

func NewPaymentHandler(conf *config.PaymentApiConfig) *PaymentHandler {
	return &PaymentHandler{
		ApiKey:    conf.ApiKey,
		ShopID: conf.ShopID,
	}
}

type Handler struct {
	authorizationService    service.IAuthorizationService
	clientMembershipService service.IClientMembershipsService
	clientService           service.IClientService
	equipmentService        service.IEquipmentService
	gymService              service.IGymService
	membershipTypeService   service.IMembershipTypeService
	scheduleService         service.IScheduleService
	trainerService          service.ITrainerService
	trainingService         service.ITrainingService
}

func NewHandler(services service.ApiServices) *Handler {
	return &Handler{
		authorizationService:    services.AuthorizationService,
		clientMembershipService: services.ClientMembershipService,
		clientService:           services.ClientService,
		equipmentService:        services.EquipmentService,
		gymService:              services.GymService,
		membershipTypeService:   services.MembershipTypeService,
		scheduleService:         services.ScheduleService,
		trainerService:          services.TrainerService,
		trainingService:         services.TrainingService,
	}
}

func (h *Handler) InitApi(router gin.IRouter) {
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/register", h.RegisterNewUser)
			v1.POST("/login", h.Login)
			v1.GET("/logout", h.Logout)

			v1.GET("/isauthorize", h.IsAuthorize)

			clientMembership := v1.Group("/client_membership")
			{
				clientMembership.POST("/new", h.CreateNewClientMembership)
				clientMembership.POST("/change", h.ChangeClientMembership)
				clientMembership.POST("/delete", h.DeleteClientMembership)
				clientMembership.GET("/:id", h.ListClientMembershipsByClientID)
			}

			client := v1.Group("/client")
			{
				client.GET("/login/:login", h.GetClientByLogin)
				client.POST("/delete", h.DeleteClient)
				client.POST("/change", h.ChangeClient)
				client.GET("/all", h.ListClients)
				client.GET("/:id", h.GetClientByID)
			}

			equipment := v1.Group("/equipment")
			{
				equipment.POST("/new", h.CreateNewEquipment)
				equipment.POST("/change", h.ChangeEquipment)
				equipment.POST("/delete", h.DeleteEquipment)
				equipment.GET("/:id", h.ListEquipmentsByGymID)
			}

			gym := v1.Group("/gym")
			{
				gym.POST("/new", h.CreateNewGym)
				gym.POST("/change", h.ChangeGym)
				gym.POST("/delete", h.DeleteGym)
				gym.GET("/all", h.ListGyms)
				gym.GET("/:id", h.GetGymByID)
			}

			membershipType := v1.Group("/membershipType")
			{
				membershipType.POST("/new", h.CreateNewMembershipType)
				membershipType.POST("/change", h.ChangeMembershipType)
				membershipType.POST("/delete", h.DeleteMembershipType)
				membershipType.GET("/:id", h.ListMembershipTypeByGymID)
			}

			schedule := v1.Group("/schedule")
			{
				schedule.POST("/new", h.CreateNewSchedule)
				schedule.POST("/change", h.ChangeSchedule)
				schedule.POST("/delete", h.DeleteSchedule)
				schedule.GET("/:id", h.ListSchedulesByClientID)
			}

			trainer := v1.Group("/trainer")
			{
				trainer.POST("/new", h.CreateNewTrainer)
				trainer.POST("/change", h.ChangeTrainer)
				trainer.POST("/delete", h.DeleteTrainer)
				trainer.GET("/all", h.GetListTrainers)
				trainer.GET("/:id", h.GetListTrainersByGymID)
			}

			training := v1.Group("/training")
			{
				training.POST("/new", h.CreateNewTraining)
				training.POST("/change", h.ChangeTraining)
				training.POST("/delete", h.DeleteTraining)
				training.GET("/:id", h.ListTrainingsByTrainerID)
			}
		}
	}
}
