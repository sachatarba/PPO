package request

import "github.com/google/uuid"

type ScheduleReq struct {
	ID           uuid.UUID `json:"id" binding:"required"`
	DayOfTheWeek string    `json:"dayoftheweek" binding:"required"`
	StartTime    string    `json:"starttime" binding:"required"`
	EndTime      string    `json:"endtime" binding:"required"`
	ClientID     uuid.UUID `json:"clientid" binding:"required"`
	TrainingID   uuid.UUID `json:"trainingid" binding:"required"`
}
