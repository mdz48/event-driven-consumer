package main

import (
	"event-driven-consumer/src/core"
	"event-driven-consumer/src/features/orders/application"
	"event-driven-consumer/src/features/orders/infrastructure"
	"event-driven-consumer/src/features/orders/infrastructure/controllers"
	"github.com/gin-gonic/gin"
	"log"
)

type Dependencies struct {
	engine   *gin.Engine
	database *core.Database
	rabbitMQ *core.RabbitMQConnection
	consumer *infrastructure.OrderConsumer
}

func NewDependencies() *Dependencies {
	// Inicializar la base de datos
	database := core.NewDatabase()
	if database == nil {
		log.Fatal("Error al inicializar la base de datos")
	}

	// Inicializar RabbitMQ
	rabbitMQ := core.NewRabbitMQConnection()
	if rabbitMQ == nil {
		log.Fatal("Error al inicializar RabbitMQ")
	}

	// Inicializar el repositorio
	orderRepository := infrastructure.NewMySQL(database.Conn)

	// Inicializar casos de uso
	getOrdersUseCase := application.NewGetOrdersUseCase(orderRepository)
	updateOrderUseCase := application.NewUpdateOrderUseCase(orderRepository)

	// Inicializar el caso de uso ProcessOrderUseCase (necesitas implementar esta clase)
	processOrderUseCase := application.NewProcessOrderUseCase(orderRepository)

	// Inicializar controladores
	getOrdersController := controllers.NewGetOrdersController(getOrdersUseCase)
	updateOrderController := controllers.NewUpdateOrderController(updateOrderUseCase)

	// Inicializar el motor HTTP
	engine := gin.Default()

	// Configurar rutas
	engine.GET("/orders", getOrdersController.GetOrders)
	engine.PUT("/orders", updateOrderController.UpdateOrder)

	// Inicializar el consumidor
	consumer := infrastructure.NewOrderConsumer(rabbitMQ, processOrderUseCase)

	return &Dependencies{
		engine:   engine,
		database: database,
		rabbitMQ: rabbitMQ,
		consumer: consumer,
	}
}

func (d *Dependencies) Run() error {
	// Iniciar el consumidor de mensajes
	err := d.consumer.StartConsuming()
	if err != nil {
		return err
	}

	// Iniciar el servidor HTTP
	return d.engine.Run(":8000")
}