package dto

type PutGym struct {
	Name    string    `json:"name" binding:"required"`
	Phone   string    `json:"phone" binding:"required"`
	City    string    `json:"city" binding:"required"`
	Addres  string    `json:"addres" binding:"required"`
	IsChain bool      `json:"is_chain" binding:"required"`
}
