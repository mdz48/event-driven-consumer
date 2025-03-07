package main

import (
	"event-driven-consumer/src/core"
	"event-driven-consumer/src/features/orders/application/usecases"
	"event-driven-consumer/src/features/orders/infrastructure"
	"event-driven-consumer/src/features/orders/infrastructure/controllers"
	"log"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	engine   *gin.Engine
	database *core.Database
	rabbitMQ *core.RabbitMQConnection
}

func NewDependencies() *Dependencies {
	// Inicializar la base de datos
	database := core.NewDatabase()
	if database == nil {
		log.Fatal("Error al inicializar la base de datos")
	}

	// Inicializar RabbitMQ (solo para envío de mensajes)
	rabbitMQ := core.NewRabbitMQConnection()
	if rabbitMQ == nil {
		log.Fatal("Error al inicializar RabbitMQ")
	}

	// Inicializar el repositorio de órdenes (adaptador secundario)
	orderRepository := infrastructure.NewMySQL(database.Conn)

	// Inicializar el publicador de mensajes (adaptador secundario)
	messagePublisher := infrastructure.NewRabbitMQPublisher(rabbitMQ)

	// Inicializar casos de uso (capa de aplicación)
	getOrdersUseCase := usecases.NewGetOrdersUseCase(orderRepository)
	updateOrderUseCase := usecases.NewUpdateOrderUseCase(orderRepository)
	processOrderUseCase := usecases.NewProcessOrderUseCase(orderRepository, messagePublisher)

	// Inicializar controladores (adaptadores primarios)
	getOrdersController := controllers.NewGetOrdersController(getOrdersUseCase)
	updateOrderController := controllers.NewUpdateOrderController(updateOrderUseCase)
	processOrderController := controllers.NewProcessOrderController(processOrderUseCase)

	// Inicializar el motor HTTP
	engine := gin.Default()

	// Configurar rutas
	engine.GET("/orders", getOrdersController.GetOrders)
	engine.PUT("/orders/:id", updateOrderController.UpdateOrder)
	engine.POST("/orders/process", processOrderController.ProcessOrder)

	return &Dependencies{
		engine:   engine,
		database: database,
		rabbitMQ: rabbitMQ,
	}
}

func (d *Dependencies) Run() error {
	// Iniciar el servidor HTTP
	return d.engine.Run(":8000")
}
