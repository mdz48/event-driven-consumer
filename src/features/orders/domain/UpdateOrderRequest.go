package domain

type UpdateOrderRequest struct {
	ID    string `json:"id"`
	State string `json:"state"`
}