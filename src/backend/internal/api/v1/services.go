package v1

import (
	"github.com/go-redis/redis/v8"
	"github.com/sachatarba/course-db/internal/repository"
	"github.com/sachatarba/course-db/internal/service"
	"gorm.io/gorm"
)

type IApiServicesBuilder interface {
	buildRepos() error
	buildServices() error
	getResult() *service.ApiServices
}

type ApiServicesDirector struct {
	Builder IApiServicesBuilder
}

func (director *ApiServicesDirector) NewServices() (*service.ApiServices, error){
	err := director.Builder.buildRepos()
	if err != nil {
		return nil, err
	}

	err = director.Builder.buildServices()
	if err != nil {
		return nil, err
	}

	return director.Builder.getResult(), nil
}


type ApiServicesBuilder struct {
	Postgres *gorm.DB
	Redis    *redis.Client
	repos    *service.ApiRepos
	services *service.ApiServices
	// Repos service.Repos
}

func (builder *ApiServicesBuilder) buildRepos() error {
	builder.repos = &service.ApiRepos{
		SessionRepo:          repository.NewSessionRepo(builder.Redis),
		ClientMembershipRepo: repository.NewClientMembershipRepo(builder.Postgres),
		ClientRepo:           repository.NewClientRepo(builder.Postgres),
		EquipmentRepo:        repository.NewEquipmentRepo(builder.Postgres),
		GymRepo:              repository.NewGymRepo(builder.Postgres),
		MembershipTypeRepo:   repository.NewMembershipTypeRepo(builder.Postgres),
		ScheduleRepo:         repository.NewScheduleRepo(builder.Postgres),
		TrainerRepo:          repository.NewTrainerRepo(builder.Postgres),
		TrainingRepo:         repository.NewTrainingRepo(builder.Postgres),
	}

	return nil
}

func (builder *ApiServicesBuilder) buildServices() error {
	builder.services = &service.ApiServices{
		AuthorizationService:    service.NewAuthorizationService(builder.repos.SessionRepo, builder.repos.ClientRepo),
		ClientMembershipService: service.NewClientMembershipService(builder.repos.ClientMembershipRepo),
		ClientService:           service.NewClientService(builder.repos.ClientRepo),
		EquipmentService:        service.NewEquipmentService(builder.repos.EquipmentRepo),
		GymService:              service.NewGymService(builder.repos.GymRepo),
		MembershipTypeService:   service.NewMembershipTypeService(builder.repos.MembershipTypeRepo),
		ScheduleService:         service.NewScheduleService(builder.repos.ScheduleRepo),
		TrainerService:          service.NewTrainerService(builder.repos.TrainerRepo),
		TrainingService:         service.NewTrainingService(builder.repos.TrainingRepo),
	}

	return nil
}

func (builder *ApiServicesBuilder) getResult() *service.ApiServices {
	return builder.services
}
