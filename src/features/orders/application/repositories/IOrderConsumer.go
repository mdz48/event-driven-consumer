package repositories

import "event-driven-consumer/src/features/orders/domain"

type IOrderConsumer interface {
	Start() error
	Stop() error
	OnOrderCreated(handler func(order domain.Order) error)
	OnOrderStatusChanged(handler func(order domain.Order) error)
}
