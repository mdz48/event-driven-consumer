package domain

type IOrderConsumer interface {
	Start() error
	Stop() error
	OnOrderCreated(handler func(order Order) error)
	OnOrderStatusChanged(handler func(order Order) error)
}