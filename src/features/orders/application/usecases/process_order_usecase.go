package usecases

import (
	"event-driven-consumer/src/features/orders/application/repositories"
	"event-driven-consumer/src/features/orders/domain"
	"log"
	"time"
)

type ProcessOrderUseCase struct {
	repository   domain.IOrder
	msgPublisher repositories.IMessage
}

func NewProcessOrderUseCase(
	repository domain.IOrder,
	msgPublisher repositories.IMessage,
) *ProcessOrderUseCase {
	return &ProcessOrderUseCase{
		repository:   repository,
		msgPublisher: msgPublisher,
	}
}

func (uc *ProcessOrderUseCase) ReceiveOrder(orderReceived domain.Order) (domain.Order, error) {
	order := domain.Order{
		ID:       orderReceived.ID,
		TableID:  orderReceived.TableID,
		Product:  orderReceived.Product,
		Quantity: orderReceived.Quantity,
		Status:   "recibida",
	}

	// Guardar en repositorio
	createdOrder, err := uc.repository.Create(order)
	if err != nil {
		return domain.Order{}, err
	}

	// Iniciar procesamiento asíncrono
	go uc.processOrderAsync(createdOrder.ID)

	return createdOrder, nil
}

// Método privado para procesar la orden de manera asíncrona
func (uc *ProcessOrderUseCase) processOrderAsync(orderID int) {
	// Esperar un tiempo antes de cambiar el estado
	time.Sleep(2 * time.Second)

	// Cambiar a "procesando"
	updatedOrder, err := uc.StartProcessing(orderID)
	if err != nil {
		log.Printf("Error al actualizar estado a procesando: %v", err)
		return
	}

	// Publicar mensaje de cambio de estado
	if uc.msgPublisher != nil {
		err = uc.msgPublisher.PublishMessage(updatedOrder, "processing")
		if err != nil {
			log.Printf("Error al publicar cambio de estado: %v", err)
		}
	}

	// Esperar otro tiempo antes de completar
	time.Sleep(5 * time.Second)

	// Cambiar a "completada"
	completedOrder, err := uc.FinishProcessing(orderID)
	if err != nil {
		log.Printf("Error al actualizar estado a completado: %v", err)
		return
	}

	// Publicar mensaje de orden completada
	if uc.msgPublisher != nil {
		err = uc.msgPublisher.PublishMessage(completedOrder, "completed")
		if err != nil {
			log.Printf("Error al publicar orden completada: %v", err)
		}
	}
}

// StartProcessing cambia el estado de una orden a "procesando"
func (uc *ProcessOrderUseCase) StartProcessing(orderID int) (domain.Order, error) {
	return uc.repository.UpdateStatus(orderID, "procesando")
}

// FinishProcessing cambia el estado de una orden a "completada"
func (uc *ProcessOrderUseCase) FinishProcessing(orderID int) (domain.Order, error) {
	return uc.repository.UpdateStatus(orderID, "completada")
}
