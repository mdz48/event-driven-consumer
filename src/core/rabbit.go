package core

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

type RabbitMQConnection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQConnection() *RabbitMQConnection {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		log.Println("Falta en .env")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("Error conectando a RabbitMQ: %s", err)
		return nil
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Printf("Error abriendo canal de RabbitMQ: %s", err)
		return nil
	}

	log.Println("Conexión establecida con RabbitMQ")

	return &RabbitMQConnection{
		Conn:    conn,
		Channel: ch,
	}
}

func (r *RabbitMQConnection) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
}

func (r *RabbitMQConnection) DeclareQueue(name string) (amqp.Queue, error) {
    return r.Channel.QueueDeclare(
        name,  // nombre
        true,  // durable
        false, // delete when unused
        false, // exclusive
        false, // no-wait
        nil,   // arguments
    )
}