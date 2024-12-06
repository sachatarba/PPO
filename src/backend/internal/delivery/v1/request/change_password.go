package request

type ChangePasswordReq struct {
	Login       string `json:"login" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
	Code        string `json:"code" binding:"required"`
}
