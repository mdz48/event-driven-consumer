package domain

type IOrder interface {
	GetAll() ([]Order, error)
	GetOrder(id int) (Order, error)
	Create(order Order) (Order, error)
	UpdateStatus(id int, status string) (Order, error)
	Delete(id int) error
}