package dto

type PutMembershipType struct {
	Type         string    `json:"type" binding:"required"`
	Description  string    `json:"description" binding:"required"`
	Price        string    `json:"price" binding:"required"`
	DaysDuration int       `json:"days_duration" binding:"required"`
}
