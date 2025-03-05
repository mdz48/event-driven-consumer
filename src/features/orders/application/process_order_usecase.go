package application

import "event-driven-consumer/src/features/orders/domain"

// OrderReceived es la estructura que representa una orden recibida del mensaje
type OrderReceived struct {
    OrderID  int    `json:"order_id"`
    TableID  int    `json:"table_id"`
    Product  string `json:"product"`
    Quantity int    `json:"quantity"`
}

type ProcessOrderUseCase struct {
    db domain.IOrder
}

func NewProcessOrderUseCase(db domain.IOrder) *ProcessOrderUseCase {
    return &ProcessOrderUseCase{db: db}
}

// ReceiveOrder procesa una nueva orden recibida
func (uc *ProcessOrderUseCase) ReceiveOrder(orderReceived OrderReceived) (domain.Order, error) {
    order := domain.Order{
        ID:       orderReceived.OrderID,
        TableID:  orderReceived.TableID,
        Product:  orderReceived.Product,
        Quantity: orderReceived.Quantity,
        Status:   "recibida",
    }

    return uc.db.Create(order)
}

// StartProcessing cambia el estado de una orden a "procesando"
func (uc *ProcessOrderUseCase) StartProcessing(orderID int) (domain.Order, error) {
    return uc.db.UpdateStatus(orderID, "procesando")
}

// FinishProcessing cambia el estado de una orden a "completada"
func (uc *ProcessOrderUseCase) FinishProcessing(orderID int) (domain.Order, error) {
    return uc.db.UpdateStatus(orderID, "completada")
}