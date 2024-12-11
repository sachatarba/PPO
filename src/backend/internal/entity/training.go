package entity

import "github.com/google/uuid"

type TrainingTypeID int

const (
	Aerobic TrainingTypeID = iota
	Anaerobic
	Flexibility
	Strength
)

var TrainingType map[TrainingTypeID]string = map[TrainingTypeID]string{
	Aerobic:     "aerobic",
	Anaerobic:   "anaerobic",
	Flexibility: "flexibility",
	Strength:    "strength",
}

type Training struct {
	ID           uuid.UUID
	Title        string
	Description  string
	TrainingType string
	TrainerID    uuid.UUID
}

func (t *Training) Validate() bool {
	found := false
	for id := range TrainingType {
		if value := TrainingType[id]; value == t.TrainingType {
			found = true
		}
	}

	if !found {
		return false
	}

	return t.Title != "" && t.Description != ""
}
