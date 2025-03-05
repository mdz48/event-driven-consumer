package domain

type Order struct {
	ID       int    `json:"id"`
	TableID  int    `json:"table_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}