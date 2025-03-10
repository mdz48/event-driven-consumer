package main

import (
	"event-driven-consumer/src/core"
	"event-driven-consumer/src/features/orders/application/usecases"
	"event-driven-consumer/src/features/orders/infrastructure"
	"event-driven-consumer/src/features/orders/infrastructure/controllers"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	engine   *gin.Engine
	database *core.Database
	rabbitMQ *core.RabbitMQConnection
}

func NewDependencies() *Dependencies {
	return &Dependencies{
		engine:   gin.Default(),
		database: core.NewDatabase(),
		rabbitMQ: core.NewRabbitMQConnection(),
	}
}

func (d *Dependencies) Run() {
	database := core.NewDatabase()
	rabbit := core.NewRabbitMQConnection()

	ordersDataBase := infrastructure.NewMySQL(database.Conn)
	ordersRabbit := infrastructure.NewRabbitMQPublisher(rabbit)
	ordersUpdateUseCase := usecases.NewUpdateOrderUseCase(ordersDataBase, ordersRabbit)
	ordersController := controllers.NewUpdateOrderController(ordersUpdateUseCase)

	getOrdersUseCase := usecases.NewGetOrdersUseCase(ordersDataBase)
	getOrdersController := controllers.NewGetOrdersController(getOrdersUseCase)

	ordersRouter := infrastructure.NewOrdersRouter(d.engine, getOrdersController, ordersController)
	ordersRouter.SetupRoutes()

	d.engine.Run(":8080")
}
