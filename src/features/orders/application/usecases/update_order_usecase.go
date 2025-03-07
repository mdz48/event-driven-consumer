package usecases

import (
	"event-driven-consumer/src/features/orders/application/repositories"
	"event-driven-consumer/src/features/orders/domain"
	"log"
)

type UpdateOrderUseCase struct {
	db           domain.IOrder
	msgPublisher repositories.IMessage
}

func NewUpdateOrderUseCase(db domain.IOrder, msgPublisher repositories.IMessage) *UpdateOrderUseCase {
	return &UpdateOrderUseCase{
		db:           db,
		msgPublisher: msgPublisher,
	}
}

func (uc *UpdateOrderUseCase) Execute(id int, status string) (domain.Order, error) {
	// Actualizar en base de datos
	updatedOrder, err := uc.db.UpdateStatus(id, status)
	if err != nil {
		return domain.Order{}, err
	}

	// Publicar mensaje de cambio de estado
	if uc.msgPublisher != nil {
		err = uc.msgPublisher.PublishMessage(updatedOrder, status)
		if err != nil {
			log.Printf("Error al publicar mensaje de cambio de estado: %v", err)
		} else {
			log.Printf("Mensaje de cambio de estado para orden #%d enviado correctamente: %s",
				updatedOrder.ID, status)
		}
	}

	return updatedOrder, nil
}
