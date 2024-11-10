package orm

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

type MembershipTypeConverter struct {
}

func NewMembershipTypeConverter() entity.IConverter[entity.MembershipType, MembershipType] {
	return &MembershipTypeConverter{}
}

func (conv *MembershipTypeConverter) ConvertFromEntity(membership entity.MembershipType) MembershipType {
	return MembershipType{
		ID:           membership.ID,
		Type:         membership.Type,
		Description:  membership.Description,
		Price:        membership.Price,
		DaysDuration: membership.DaysDuration,
		GymID:        membership.GymID,
	}
}

func (conv *MembershipTypeConverter) ConvertToEntity(membership MembershipType) entity.MembershipType {
	return entity.MembershipType{
		ID:           membership.ID,
		Type:         membership.Type,
		Description:  membership.Description,
		Price:        membership.Price,
		DaysDuration: membership.DaysDuration,
		GymID:        membership.GymID,
	}
}

func (conv *MembershipTypeConverter) ConvertFromEntitySlice(memberships []entity.MembershipType) []MembershipType {
	membershipsOrm := make([]MembershipType, len(memberships))
	for i := 0; i < len(memberships); i++ {
		membershipsOrm[i] = conv.ConvertFromEntity(memberships[i])
	}

	return membershipsOrm
}

func (conv *MembershipTypeConverter) ConvertToEntitySlice(membershipsOrm []MembershipType) []entity.MembershipType {
	memberships := make([]entity.MembershipType, len(membershipsOrm))
	for i := 0; i < len(memberships); i++ {
		memberships[i] = conv.ConvertToEntity(membershipsOrm[i])
	}

	return memberships
}

type ClientMembershipConverter struct {
}

func NewClientMembershipConverter() entity.IConverter[entity.ClientMembership, ClientMembership] {
	return &ClientMembershipConverter{}
}

func (conv *ClientMembershipConverter) ConvertFromEntity(membership entity.ClientMembership) ClientMembership {
	return ClientMembership{
		ID:               membership.ID,
		StartDate:        membership.StartDate,
		EndDate:          membership.EndDate,
		MembershipTypeID: membership.MembershipType.ID,
		MembershipType:   NewMembershipTypeConverter().ConvertFromEntity(membership.MembershipType),
		ClientID:         membership.ClientID,
	}
}

func (c *ClientMembershipConverter) ConvertToEntity(membership ClientMembership) entity.ClientMembership {
	return entity.ClientMembership{
		ID:             membership.ID,
		StartDate:      membership.StartDate,
		EndDate:        membership.EndDate,
		MembershipType: NewMembershipTypeConverter().ConvertToEntity(membership.MembershipType),
		ClientID:       membership.ClientID,
	}
}

func (conv *ClientMembershipConverter) ConvertFromEntitySlice(clientMemberships []entity.ClientMembership) []ClientMembership {
	membershipsOrm := make([]ClientMembership, len(clientMemberships))
	for i := 0; i < len(clientMemberships); i++ {
		membershipsOrm[i] = conv.ConvertFromEntity(clientMemberships[i])
	}

	return membershipsOrm
}

func (conv *ClientMembershipConverter) ConvertToEntitySlice(membershipsOrm []ClientMembership) []entity.ClientMembership {
	memberships := make([]entity.ClientMembership, len(membershipsOrm))
	for i := 0; i < len(memberships); i++ {
		memberships[i] = conv.ConvertToEntity(membershipsOrm[i])
	}

	return memberships
}

type EquipmentConverter struct {
}

func NewEquipmentConverter() entity.IConverter[entity.Equipment, Equipment] {
	return &EquipmentConverter{}
}

func (conv *EquipmentConverter) ConvertFromEntity(e entity.Equipment) Equipment {
	return Equipment{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		GymID:       e.GymID,
	}
}

func (conv *EquipmentConverter) ConvertToEntity(e Equipment) entity.Equipment {
	return entity.Equipment{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		GymID:       e.GymID,
	}
}

func (conv *EquipmentConverter) ConvertFromEntitySlice(eqs []entity.Equipment) []Equipment {
	eqsOrm := make([]Equipment, len(eqs))
	for i := 0; i < len(eqs); i++ {
		eqsOrm[i] = conv.ConvertFromEntity(eqs[i])
	}

	return eqsOrm
}

func (conv *EquipmentConverter) ConvertToEntitySlice(eqsOrm []Equipment) []entity.Equipment {
	eqs := make([]entity.Equipment, len(eqsOrm))
	for i := 0; i < len(eqs); i++ {
		eqs[i] = conv.ConvertToEntity(eqsOrm[i])
	}

	return eqs
}

type TrainingConverter struct {
}

func NewTrainingConverter() entity.IConverter[entity.Training, Training] {
	return &TrainingConverter{}
}

func (conv *TrainingConverter) ConvertFromEntity(training entity.Training) Training {
	return Training{
		ID:           training.ID,
		Title:        training.Title,
		Description:  training.Description,
		TrainingType: training.TrainingType,
		TrainerID:    training.TrainerID,
	}
}

func (conv *TrainingConverter) ConvertToEntity(training Training) entity.Training {
	return entity.Training{
		ID:           training.ID,
		Title:        training.Title,
		Description:  training.Description,
		TrainingType: training.TrainingType,
		TrainerID:    training.TrainerID,
	}
}

func (conv *TrainingConverter) ConvertFromEntitySlice(training []entity.Training) []Training {
	trainingOrm := make([]Training, len(training))
	for i := 0; i < len(training); i++ {
		trainingOrm[i] = conv.ConvertFromEntity(training[i])
	}

	return trainingOrm
}

func (conv *TrainingConverter) ConvertToEntitySlice(trainingOrm []Training) []entity.Training {
	training := make([]entity.Training, len(trainingOrm))
	for i := 0; i < len(training); i++ {
		training[i] = conv.ConvertToEntity(trainingOrm[i])
	}

	return training
}

type TrainerConverter struct {
}

func NewTrainerConverter() entity.IConverter[entity.Trainer, Trainer] {
	return &TrainerConverter{}
}

func (conv *TrainerConverter) ConvertFromEntity(trainer entity.Trainer) Trainer {
	gyms := make([]*Gym, len(trainer.GymsID))
	for i, id := range trainer.GymsID {
		gyms[i] = &Gym{
			ID: id,
		}
	}

	return Trainer{
		ID:            trainer.ID,
		Fullname:      trainer.Fullname,
		Email:         trainer.Email,
		Phone:         trainer.Phone,
		Qualification: trainer.Qualification,
		UnitPrice:     trainer.UnitPrice,
		Gyms:          gyms,
		Trainings:     NewTrainingConverter().ConvertFromEntitySlice(trainer.Trainings),
	}
}

func (conv *TrainerConverter) ConvertToEntity(trainer Trainer) entity.Trainer {
	gymsID := make([]uuid.UUID, len(trainer.Gyms))
	for i, gym := range trainer.Gyms {
		gymsID[i] = gym.ID
	}

	return entity.Trainer{
		ID:            trainer.ID,
		Fullname:      trainer.Fullname,
		Phone:         trainer.Phone,
		Email:         trainer.Email,
		Qualification: trainer.Qualification,
		UnitPrice:     trainer.UnitPrice,
		GymsID:        gymsID,
		Trainings:     NewTrainingConverter().ConvertToEntitySlice(trainer.Trainings),
	}
}

func (conv *TrainerConverter) ConvertFromEntitySlice(trainers []entity.Trainer) []Trainer {
	trainersOrm := make([]Trainer, len(trainers))
	for i := 0; i < len(trainers); i++ {
		trainersOrm[i] = conv.ConvertFromEntity(trainers[i])
	}

	return trainersOrm
}

func (conv *TrainerConverter) ConvertToEntitySlice(trainersOrm []Trainer) []entity.Trainer {
	trainers := make([]entity.Trainer, len(trainersOrm))
	for i := 0; i < len(trainers); i++ {
		trainers[i] = conv.ConvertToEntity(trainersOrm[i])
	}

	return trainers
}

type GymConverter struct {
}

func NewGymConverter() entity.IConverter[entity.Gym, Gym] {
	return &GymConverter{}
}

func (conv *GymConverter) ConvertFromEntity(gym entity.Gym) Gym {
	trainerConverter := NewTrainerConverter()
	trainersPtr := make([]*Trainer, len(gym.Trainers))
	for i, trainer := range gym.Trainers {
		trainer := trainerConverter.ConvertFromEntity(trainer)
		trainersPtr[i] = &trainer
	}

	return Gym{
		ID:              gym.ID,
		Name:            gym.Name,
		Phone:           gym.Phone,
		City:            gym.City,
		Addres:          gym.Addres,
		IsChain:         strconv.FormatBool(gym.IsChain),
		Trainers:        trainersPtr,
		Equipments:      NewEquipmentConverter().ConvertFromEntitySlice(gym.Equipments),
		MembershipTypes: NewMembershipTypeConverter().ConvertFromEntitySlice(gym.MembershipTypes),
	}
}

func (conv *GymConverter) ConvertToEntity(gym Gym) entity.Gym {
	trainers := make([]entity.Trainer, len(gym.Trainers))
	trainersConverter := NewTrainerConverter()
	for i, trainer := range gym.Trainers {
		if trainer != nil {
			trainers[i] = trainersConverter.ConvertToEntity(*trainer)
		}
	}

	IsChain, _ := strconv.ParseBool(gym.IsChain)

	return entity.Gym{
		ID:              gym.ID,
		Name:            gym.Name,
		Phone:           gym.Phone,
		City:            gym.City,
		Addres:          gym.Addres,
		IsChain:         IsChain,
		Trainers:        trainers,
		Equipments:      NewEquipmentConverter().ConvertToEntitySlice(gym.Equipments),
		MembershipTypes: NewMembershipTypeConverter().ConvertToEntitySlice(gym.MembershipTypes),
	}
}

func (conv *GymConverter) ConvertFromEntitySlice(gyms []entity.Gym) []Gym {
	gymsOrm := make([]Gym, len(gyms))
	for i, gym := range gyms {
		gymsOrm[i] = conv.ConvertFromEntity(gym)
	}

	return gymsOrm
}

func (conv *GymConverter) ConvertToEntitySlice(gymsOrm []Gym) []entity.Gym {
	gyms := make([]entity.Gym, len(gymsOrm))
	for i, gym := range gymsOrm {
		gyms[i] = conv.ConvertToEntity(gym)
	}

	return gyms
}

type ScheduleConverter struct{}

func NewScheduleConverter() entity.IConverter[entity.Schedule, Schedule] {
	return &ScheduleConverter{}
}

func (conv *ScheduleConverter) ConvertFromEntity(schedule entity.Schedule) Schedule {
	return Schedule{
		ID:           schedule.ID,
		DayOfTheWeek: schedule.DayOfTheWeek,
		StartTime:    schedule.StartTime,
		EndTime:      schedule.EndTime,
		ClientID:     schedule.ClientID,
		TrainingID:   schedule.TrainingID,
		Training:     NewTrainingConverter().ConvertFromEntity(schedule.Training),
	}
}

func (conv *ScheduleConverter) ConvertToEntity(scheduleOrm Schedule) entity.Schedule {
	// start := time.Parse(time., scheduleOrm.StartTime)

	return entity.Schedule{
		ID:           scheduleOrm.ID,
		DayOfTheWeek: scheduleOrm.DayOfTheWeek,
		StartTime:    scheduleOrm.StartTime,
		EndTime:      scheduleOrm.EndTime,
		ClientID:     scheduleOrm.ClientID,
		TrainingID:   scheduleOrm.TrainingID,
		Training:     NewTrainingConverter().ConvertToEntity(scheduleOrm.Training),
	}
}

func (conv *ScheduleConverter) ConvertFromEntitySlice(schedules []entity.Schedule) []Schedule {
	schedulesOrm := make([]Schedule, len(schedules))
	for i, schedule := range schedules {
		schedulesOrm[i] = conv.ConvertFromEntity(schedule)
	}
	return schedulesOrm
}

func (conv *ScheduleConverter) ConvertToEntitySlice(schedulesOrm []Schedule) []entity.Schedule {
	schedules := make([]entity.Schedule, len(schedulesOrm))
	for i, schedule := range schedulesOrm {
		schedules[i] = conv.ConvertToEntity(schedule)
	}
	return schedules
}

type ClientConverter struct {
}

func NewClientConverter() entity.IConverter[entity.Client, Client] {
	return &ClientConverter{}
}

func (conv ClientConverter) ConvertFromEntity(client entity.Client) Client {
	return Client{
		ID:                client.ID,
		Login:             client.Login,
		Password:          client.Password,
		Fullname:          client.Fullname,
		Email:             client.Email,
		Phone:             client.Phone,
		Birthdate:         client.Birthdate,
		ClientMemberships: NewClientMembershipConverter().ConvertFromEntitySlice(client.ClientMemberships),
		Schedules:         NewScheduleConverter().ConvertFromEntitySlice(client.Schedules),
	}
}

func (conv ClientConverter) ConvertToEntity(clientOrm Client) entity.Client {
	return entity.Client{
		ID:                clientOrm.ID,
		Login:             clientOrm.Login,
		Password:          clientOrm.Password,
		Fullname:          clientOrm.Fullname,
		Email:             clientOrm.Email,
		Phone:             clientOrm.Phone,
		Birthdate:         clientOrm.Birthdate,
		ClientMemberships: NewClientMembershipConverter().ConvertToEntitySlice(clientOrm.ClientMemberships),
		Schedules:         NewScheduleConverter().ConvertToEntitySlice(clientOrm.Schedules),
	}
}

func (conv ClientConverter) ConvertFromEntitySlice(clients []entity.Client) []Client {
	clientOrm := make([]Client, len(clients))
	for i, entity := range clients {
		clientOrm[i] = conv.ConvertFromEntity(entity)
	}
	return clientOrm
}

func (conv ClientConverter) ConvertToEntitySlice(clientsOrm []Client) []entity.Client {
	clients := make([]entity.Client, len(clientsOrm))
	for i, client := range clientsOrm {
		clients[i] = conv.ConvertToEntity(client)
	}
	return clients
}
