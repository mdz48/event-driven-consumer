package models

type UpdateOrderRequest struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}
