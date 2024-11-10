package builder

import (
	"time"

	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

type ScheduleBuilder struct {
	id           uuid.UUID
	dayOfTheWeek string
	startTime    string
	endTime      string
	clientID     uuid.UUID
	trainingID   uuid.UUID
}

func NewScheduleBuilder() *ScheduleBuilder {
	return &ScheduleBuilder{}
}

func (b *ScheduleBuilder) SetID(id uuid.UUID) *ScheduleBuilder {
	b.id = id
	return b
}

func (b *ScheduleBuilder) SetDayOfTheWeek(day string) *ScheduleBuilder {
	b.dayOfTheWeek = day
	return b
}

func (b *ScheduleBuilder) SetStartTime(start string) *ScheduleBuilder {
	b.startTime = start
	return b
}

func (b *ScheduleBuilder) SetEndTime(end string) *ScheduleBuilder {
	b.endTime = end
	return b
}

func (b *ScheduleBuilder) SetClientID(clientID uuid.UUID) *ScheduleBuilder {
	b.clientID = clientID
	return b
}

func (b *ScheduleBuilder) SetTrainingID(trainingID uuid.UUID) *ScheduleBuilder {
	b.trainingID = trainingID
	return b
}

func (b *ScheduleBuilder) Invalid() *ScheduleBuilder {
	b.id = uuid.Nil
	b.dayOfTheWeek = ""
	b.startTime = "invalid_time"
	b.endTime = "invalid_time"
	b.clientID = uuid.Nil
	b.trainingID = uuid.Nil

	return b
}

func (b *ScheduleBuilder) Build() entity.Schedule {
	if b.id == uuid.Nil {
		b.id = uuid.New()
	}
	if b.dayOfTheWeek == "" {
		temp := time.Now().Format(time.DateOnly)
		b.dayOfTheWeek = temp
	}
	if b.startTime == "" {
		b.startTime = time.Now().Format(time.TimeOnly) 
	}
	if b.endTime == "" {
		b.endTime = time.Now().Add(time.Hour).Format(time.TimeOnly)
	}
	if b.clientID == uuid.Nil {
		b.clientID = uuid.New()
	}
	if b.trainingID == uuid.Nil {
		b.trainingID = uuid.New()
	}

	return entity.Schedule{
		ID:           b.id,
		DayOfTheWeek: b.dayOfTheWeek,
		StartTime:    b.startTime,
		EndTime:      b.endTime,
		ClientID:     b.clientID,
		TrainingID:   b.trainingID,
	}
}
