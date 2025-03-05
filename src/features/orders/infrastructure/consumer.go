package infrastructure

import (
	"encoding/json"
	"event-driven-consumer/src/core"
	"event-driven-consumer/src/features/orders/application"
	"log"
	"time"
)

type OrderConsumer struct {
	rabbitMQ            *core.RabbitMQConnection
	processOrderUseCase *application.ProcessOrderUseCase
	wsManager           *core.WebSocketManager
}

func NewOrderConsumer(
	rabbitMQ *core.RabbitMQConnection,
	processOrderUseCase *application.ProcessOrderUseCase,
	wsManager *core.WebSocketManager,
) *OrderConsumer {
	return &OrderConsumer{
		rabbitMQ:            rabbitMQ,
		processOrderUseCase: processOrderUseCase,
		wsManager:           wsManager,
	}
}

// Estructura que coincide con el mensaje publicado por el servicio de pedidos
type OrderMessage struct {
	ID       int    `json:"id"`
	TableID  int    `json:"table_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}

// Estructura para notificaciones WebSocket
type OrderNotification struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Order   interface{} `json:"order"`
}

func (c *OrderConsumer) StartConsuming() error {
	// Verificamos la conexión a RabbitMQ
	if c.rabbitMQ == nil || c.rabbitMQ.Channel == nil {
		return nil
	}

	queueName := "orders.created"

	msgs, err := c.rabbitMQ.Channel.Consume(
		queueName, // nombre de la cola
		"",        // consumidor
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	log.Printf("Cocina: Escuchando nuevas órdenes de la cola %s", queueName)

	go func() {
		for d := range msgs {
			var orderMsg OrderMessage
			if err := json.Unmarshal(d.Body, &orderMsg); err != nil {
				log.Printf("Error al deserializar orden: %s", err)
				d.Nack(false, true) // rechazar y reintentar
				continue
			}

			log.Printf("Nueva orden recibida: #%d - %s (cantidad: %d)",
				orderMsg.ID, orderMsg.Product, orderMsg.Quantity)

			// Convertir el mensaje a estructura esperada por el caso de uso
			orderReceived := application.OrderReceived{
				OrderID:  orderMsg.ID,
				TableID:  orderMsg.TableID,
				Product:  orderMsg.Product,
				Quantity: orderMsg.Quantity,
			}

			// Procesar la orden
			kitchenOrder, err := c.processOrderUseCase.ReceiveOrder(orderReceived)
			if err != nil {
				log.Printf("Error al procesar orden #%d: %s", orderMsg.ID, err)
				d.Nack(false, true) // rechazar y reintentar
				continue
			}

			log.Printf("Orden #%d registrada en cocina con ID interno %d",
				orderMsg.ID, kitchenOrder.ID)

			// Notificar a los clientes WebSocket
			notification := OrderNotification{
				Type:    "NEW_ORDER",
				Message: "Nueva orden recibida",
				Order:   kitchenOrder,
			}
			notificationJson, _ := json.Marshal(notification)
			c.wsManager.BroadcastMessage(notificationJson)

			// Confirmar procesamiento exitoso
			d.Ack(false)

			// Simular inicio automático de procesamiento después de un tiempo
			go func(orderID int) {
				time.Sleep(2 * time.Second)
				updatedOrder, err := c.processOrderUseCase.StartProcessing(orderID)
				if err != nil {
					log.Printf("Error al iniciar procesamiento de orden #%d: %s", orderID, err)
					return
				}

				// Notificar cambio de estado
				statusNotification := OrderNotification{
					Type:    "STATUS_CHANGED",
					Message: "Orden en procesamiento",
					Order:   updatedOrder,
				}
				statusJson, _ := json.Marshal(statusNotification)
				c.wsManager.BroadcastMessage(statusJson)

				// Simular finalización después de otro tiempo
				time.Sleep(5 * time.Second)
				completedOrder, err := c.processOrderUseCase.FinishProcessing(orderID)
				if err != nil {
					log.Printf("Error al finalizar orden #%d: %s", orderID, err)
					return
				}

				// Notificar orden completada
				completedNotification := OrderNotification{
					Type:    "ORDER_COMPLETED",
					Message: "Orden completada",
					Order:   completedOrder,
				}
				completedJson, _ := json.Marshal(completedNotification)
				c.wsManager.BroadcastMessage(completedJson)
			}(kitchenOrder.ID)
		}
	}()

	return nil
}
