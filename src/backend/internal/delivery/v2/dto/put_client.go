package dto

type PutClient struct {
	Fullname  string    `json:"fullname" binding:"required"`
	Login     string    `json:"login" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	Birthdate string    `json:"birthdate" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	Phone     string    `json:"phone" binding:"required"`
}
