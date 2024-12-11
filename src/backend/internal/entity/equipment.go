package entity

import "github.com/google/uuid"

type Equipment struct {
	ID          uuid.UUID
	Name        string
	Description string
	GymID       uuid.UUID
}

func (e *Equipment) Validate() bool {
	return e.Name != "" && e.Description != ""
}
