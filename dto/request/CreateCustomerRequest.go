package dto

type CreateCustomerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
