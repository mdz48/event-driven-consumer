package usecases

import "event-driven-consumer/src/features/orders/domain"

type GetOrdersUseCase struct {
	db domain.IOrder
}

func NewGetOrdersUseCase(db domain.IOrder) *GetOrdersUseCase {
	return &GetOrdersUseCase{db: db}
}

func (uc *GetOrdersUseCase) Execute() ([]domain.Order, error) {
	orders, err := uc.db.GetAll()
	if err != nil {
		return nil, err
	}

	return orders, nil
}
