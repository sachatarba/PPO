package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/sachatarba/course-db/internal/delivery/v2/rest"
	"github.com/sachatarba/course-db/internal/repository"
	"github.com/sachatarba/course-db/internal/service"
	"gorm.io/gorm"
)

type ApiHandlers struct {
	Postgres *gorm.DB
}

func (api *ApiHandlers) InitHandlers(router gin.IRouter) {
	clientMembershipRepo := repository.NewClientMembershipRepo(api.Postgres)
	clientRepo := repository.NewClientRepo(api.Postgres)
	equipmentRepo := repository.NewEquipmentRepo(api.Postgres)
	gymRepo := repository.NewGymRepo(api.Postgres)
	membershipTypeRepo := repository.NewMembershipTypeRepo(api.Postgres)
	scheduleRepo := repository.NewScheduleRepo(api.Postgres)
	trainerRepo := repository.NewTrainerRepo(api.Postgres)
	trainingRepo := repository.NewTrainingRepo(api.Postgres)

	clientMembershipService := service.NewClientMembershipService(clientMembershipRepo)
	clientService := service.NewClientService(clientRepo)
	equipmentService := service.NewEquipmentService(equipmentRepo)
	gymService := service.NewGymService(gymRepo)
	membershipTypeService := service.NewMembershipTypeService(membershipTypeRepo)
	scheduleService := service.NewScheduleService(scheduleRepo)
	trainerService := service.NewTrainerService(trainerRepo)
	trainingService := service.NewTrainingService(trainingRepo)

	clientMembershipHandler := rest.NewClientMembershipHandler(clientMembershipService)
	clientHandler := rest.NewClientHandler(clientService)
	equipmentHandler := rest.NewEquipmentHandler(equipmentService)
	gymHandler := rest.NewGymHandler(gymService)
	membershipTypeHandler := rest.NewMembershipTypeHandler(membershipTypeService)
	scheduleHandler := rest.NewScheduleHandler(scheduleService)
	trainerHandler := rest.NewTrainerHandler(trainerService)
	trainingHandler := rest.NewTrainingHandler(trainingService)

	apiV2 := router.Group("/api")
	{
		v2 := apiV2.Group("/v2")
		{
			gymGroup := v2.Group("/gyms")
			{
				gymGroup.POST("", gymHandler.PostGym)
				gymGroup.GET("", gymHandler.GetGyms)
				gymGroup.GET("/:gymId", gymHandler.GetGymByID)
				gymGroup.PUT("/:gymId", gymHandler.ChangeGym)
				gymGroup.DELETE("/:gymId", gymHandler.DeleteGym)

				equipmentGroup := gymGroup.Group("/:gymId/equipments")
				{
					equipmentGroup.POST("", equipmentHandler.PostEquipment)
					equipmentGroup.GET("", equipmentHandler.GetEquipments)
					equipmentGroup.PUT("/:equipmentId", equipmentHandler.PutEquipment)
					equipmentGroup.DELETE("/:equipmentId", equipmentHandler.DeleteEquipment)
				}

				membershipTypeGroup := gymGroup.Group("/:gymId/membership_types")
				{
					membershipTypeGroup.POST("", membershipTypeHandler.PostMembershipType)
					membershipTypeGroup.GET("", membershipTypeHandler.GetMembershipTypeByGymID)
					membershipTypeGroup.PUT("/:membershipTypeId", membershipTypeHandler.PutMembershipType)
					membershipTypeGroup.DELETE("/:membershipTypeId", membershipTypeHandler.DeleteMembershipType)
				}
			}

			clientGroup := v2.Group("/clients")
			{
				// clientGroup.POST("", )
				clientGroup.GET("", clientHandler.GetClients)
				clientGroup.GET("/:clientId", clientHandler.GetClientByID)
				clientGroup.PUT("/:clientId", clientHandler.PutClient)
				clientGroup.DELETE("/:clientId", clientHandler.DeleteClient)

				clientMembershipGroup := clientGroup.Group("/:clientId/client_memberships")
				{
					clientMembershipGroup.POST("", clientMembershipHandler.PostClientMembership)
					clientMembershipGroup.GET("", clientMembershipHandler.GetClientMemberships)
					clientMembershipGroup.PUT("/:clientMembershipId", clientMembershipHandler.PutClientMembership)
					clientMembershipGroup.DELETE("/:clientMembershipId", clientMembershipHandler.DeleteClientMembership)
				}

				scheduleGroup := clientGroup.Group("/:clientId/schedules")
				{
					scheduleGroup.POST("", scheduleHandler.PostSchedule)
					scheduleGroup.GET("", scheduleHandler.GetSchedules)
					scheduleGroup.PUT("/:scheduleId", scheduleHandler.PutSchedule)
					scheduleGroup.DELETE("/:scheduleId", scheduleHandler.DeleteSchedule)
				}
			}

			trainer := v2.Group("/trainers")
			{
				trainer.POST("", trainerHandler.PostTrainer)
				trainer.PUT("", trainerHandler.PutTrainer)
				trainer.DELETE("", trainerHandler.DeleteTrainer)
				trainer.GET("", trainerHandler.GetTrainers)
				trainer.GET("/:trainerId", trainerHandler.GetTrainersByGymID)

				training := trainer.Group("/:trainerId/trainings")
				{
					training.POST("", trainingHandler.PostTraining)
					training.PUT("/:trainingId", trainingHandler.PutTraining)
					training.DELETE("/:trainingId", trainingHandler.DeleteTraining)
					training.GET("", trainingHandler.GetTrainingsByTrainerID)
				}
			}
		}
	}
}
