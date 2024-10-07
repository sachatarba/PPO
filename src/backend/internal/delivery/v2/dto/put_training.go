package dto

type PutTraining struct {
	Title        string    `json:"title" binding:"required"`
	Description  string    `json:"description" binding:"required"`
	TrainingType string    `json:"trainingType" binding:"required"`
}
