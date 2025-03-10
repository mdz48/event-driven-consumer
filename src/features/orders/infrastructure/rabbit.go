package infrastructure

import (
	"encoding/json"
	"event-driven-consumer/src/core"
	"event-driven-consumer/src/features/orders/domain"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	rabbitMQ *core.RabbitMQConnection
}

func NewRabbitMQPublisher(rabbitMQ *core.RabbitMQConnection) *RabbitMQPublisher {
	return &RabbitMQPublisher{rabbitMQ: rabbitMQ}
}

func (p *RabbitMQPublisher) SetupExchangeAndQueue() error {
	if p.rabbitMQ == nil || p.rabbitMQ.Channel == nil {
		return nil
	}

	err := p.rabbitMQ.Channel.ExchangeDeclare(
		"orders.events", 
		"direct",        
		true,            
		false,           
		false,           
		false,          
		nil,             
	)
	if err != nil {
		return err
	}

	queue, err := p.rabbitMQ.DeclareQueue("orders.status_changed")
	if err != nil {
		return err
	}

	
	err = p.rabbitMQ.Channel.QueueBind(
		queue.Name,       
		"status_changed", 
		"orders.events",  
		false,           
		nil,              
	)
	if err != nil {
		return err
	}

	log.Printf("Exchange y cola configurados correctamente")
	return nil
}

func (p *RabbitMQPublisher) PublishMessage(order domain.Order, status string) error {
	if p.rabbitMQ == nil || p.rabbitMQ.Channel == nil {
		log.Println("No se puede enviar mensaje: conexión a RabbitMQ no disponible")
		return nil
	}

	// Crear estructura de mensaje
	message := map[string]interface{}{
		"order":  order,
		"status": status,
		"time":   time.Now().Format(time.RFC3339),
	}

	// Serializar el mensaje
	body, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al serializar mensaje: %s", err)
		return err
	}

	// Publicar mensaje
	err = p.rabbitMQ.Channel.Publish(
		"orders.events",  // exchange
		"status_changed", // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		})

	if err != nil {
		log.Printf("Error al publicar mensaje: %s", err)
		return err
	}

	log.Printf("Mensaje de actualización de estado enviado correctamente")
	return nil
}
