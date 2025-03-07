package repositories

import "event-driven-consumer/src/features/orders/domain"

type IMessage interface {
	PublishMessage(order domain.Order, status string) error
}
