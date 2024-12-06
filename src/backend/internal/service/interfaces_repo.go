package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

type (
	//go:generate mockery --name IEquipmentRepository
	IEquipmentRepository interface {
		CreateNewEquipment(ctx context.Context, equipment entity.Equipment) error
		ChangeEquipment(ctx context.Context, equipment entity.Equipment) error
		DeleteEquipment(ctx context.Context, equipmentID uuid.UUID) error
		GetEquipmentByID(ctx context.Context, equipmentID uuid.UUID) (entity.Equipment, error)
		ListEquipmentsByGymID(ctx context.Context, gymID uuid.UUID) ([]entity.Equipment, error)
	}

	//go:generate mockery --name IGymRepository
	IGymRepository interface {
		RegisterNewGym(ctx context.Context, gym entity.Gym) error
		ChangeGym(ctx context.Context, gym entity.Gym) error
		DeleteGym(ctx context.Context, gymID uuid.UUID) error
		GetGymByID(ctx context.Context, gymID uuid.UUID) (entity.Gym, error)
		ListGyms(ctx context.Context) ([]entity.Gym, error)
	}

	//go:generate mockery --name IMembershipTypeRepository
	IMembershipTypeRepository interface {
		RegisterNewMembershipType(ctx context.Context, membershipType entity.MembershipType) error
		ChangeMembershipType(ctx context.Context, membershipType entity.MembershipType) error
		DeleteMembershipType(ctx context.Context, membershipTypeID uuid.UUID) error
		GetMembershipTypeByID(ctx context.Context, membershipTypeID uuid.UUID) (entity.MembershipType, error)
		ListMembershipTypesByGymID(ctx context.Context, gymID uuid.UUID) ([]entity.MembershipType, error)
	}

	//go:generate mockery --name ITrainerRepository
	ITrainerRepository interface {
		RegisterNewTrainer(ctx context.Context, trainer entity.Trainer) error
		ChangeTrainer(ctx context.Context, trainer entity.Trainer) error
		DeleteTrainer(ctx context.Context, trainerID uuid.UUID) error
		GetTrainerByID(ctx context.Context, trainerID uuid.UUID) (entity.Trainer, error)
		ListTrainers(ctx context.Context) ([]entity.Trainer, error)
		ListTrainersByGymID(ctx context.Context, gymID uuid.UUID) ([]entity.Trainer, error)
	}

	//go:generate mockery --name IClientRepository
	IClientRepository interface {
		RegisterNewClient(ctx context.Context, client entity.Client) error
		ChangeClient(ctx context.Context, client entity.Client) error
		DeleteClient(ctx context.Context, clientID uuid.UUID) error
		GetClientByID(ctx context.Context, clientID uuid.UUID) (entity.Client, error)
		GetClientByLogin(ctx context.Context, login string) (entity.Client, error)
		ListClients(ctx context.Context) ([]entity.Client, error)
	}

	//go:generate mockery --name IClientMembershipsRepository
	IClientMembershipsRepository interface {
		CreateNewClientMembership(ctx context.Context, clientMembership entity.ClientMembership) error
		ChangeClientMembership(ctx context.Context, clientMembership entity.ClientMembership) error
		DeleteClientMembership(ctx context.Context, clientMembershipID uuid.UUID) error
		GetClientMembershipByID(ctx context.Context, clientMembershipID uuid.UUID) (entity.ClientMembership, error)
		ListClientMembershipsByClientID(ctx context.Context, clientID uuid.UUID) ([]entity.ClientMembership, error)
	}

	//go:generate mockery --name IScheduleRepository
	IScheduleRepository interface {
		CreateNewSchedule(ctx context.Context, shedule entity.Schedule) error
		ChangeSchedule(ctx context.Context, shedule entity.Schedule) error
		DeleteSchedule(ctx context.Context, scheduleID uuid.UUID) error
		GetScheduleByID(ctx context.Context, sheduleID uuid.UUID) (entity.Schedule, error)
		ListSchedulesByClientID(ctx context.Context, clientID uuid.UUID) ([]entity.Schedule, error)
	}

	//go:generate mockery --name ITrainingRepository
	ITrainingRepository interface {
		CreateNewTraining(ctx context.Context, training entity.Training) error
		ChangeTraining(ctx context.Context, training entity.Training) error
		DeleteTraining(ctx context.Context, trainingID uuid.UUID) error
		ListTrainingsByTrainerID(ctx context.Context, trainerID uuid.UUID) ([]entity.Training, error)
	}

	//go:generate mockery --name ISessionRepository
	ISessionRepository interface {
		CreateNewSession(ctx context.Context, session entity.Session) error
		DeleteSession(ctx context.Context, clientID uuid.UUID) error
		DeleteSessionBySessionID(ctx context.Context, sessionID uuid.UUID) error
		GetSessionsByClientID(ctx context.Context, clientID uuid.UUID) ([]entity.Session, error)
		GetSessionBySessionID(ctx context.Context, sessionID uuid.UUID) (entity.Session, error)
	}

	ICodeRepository interface {
		SaveCode(ctx context.Context, code string, clienID uuid.UUID) error
		GetCodeByClientID(ctx context.Context, clientID uuid.UUID) (string, error)
	}
)
