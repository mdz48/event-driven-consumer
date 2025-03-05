package infrastructure

import (
	"encoding/json"
	"event-driven-consumer/src/core"
	"event-driven-consumer/src/features/orders/application"
	// "event-driven-consumer/src/features/orders/domain"
	"log"
	"time"

	// amqp "github.com/rabbitmq/amqp091-go"
)

type OrderConsumer struct {
	rabbitMQ           *core.RabbitMQConnection
	processOrderUseCase *application.ProcessOrderUseCase
}

func NewOrderConsumer(
	rabbitMQ *core.RabbitMQConnection,
	processOrderUseCase *application.ProcessOrderUseCase,
) *OrderConsumer {
	return &OrderConsumer{
		rabbitMQ:           rabbitMQ,
		processOrderUseCase: processOrderUseCase,
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

			// Confirmar procesamiento exitoso
			d.Ack(false)

			// Simular inicio automático de procesamiento después de un tiempo
			go func(orderID int) {
				time.Sleep(2 * time.Second)
				_, err := c.processOrderUseCase.StartProcessing(orderID)
				if err != nil {
					log.Printf("Error al iniciar procesamiento de orden #%d: %s", orderID, err)
				}
			}(kitchenOrder.ID)
		}
	}()

	return nil
}