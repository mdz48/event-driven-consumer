package application

import "event-driven-consumer/src/features/orders/domain"

type UpdateOrderUseCase struct {
	db domain.IOrder
}

func NewUpdateOrderUseCase(db domain.IOrder) *UpdateOrderUseCase {
	return &UpdateOrderUseCase{db: db}
}

func (uc *UpdateOrderUseCase) Execute(id int, state string) (domain.Order, error) {
	return uc.db.UpdateStatus(id, state)
}