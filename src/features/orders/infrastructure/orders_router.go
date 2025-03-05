package infrastructure

import (
	"event-driven-consumer/src/features/orders/infrastructure/controllers"
	"github.com/gin-gonic/gin"
)

type OrdersRouter struct {
	engine                *gin.Engine
	getOrdersController   *controllers.GetOrdersController
	updateOrderController *controllers.UpdateOrderController
}

func NewOrdersRouter(
	engine *gin.Engine,
	getOrdersController *controllers.GetOrdersController,
	updateOrderController *controllers.UpdateOrderController,
) *OrdersRouter {
	return &OrdersRouter{
		engine:                engine,
		getOrdersController:   getOrdersController,
		updateOrderController: updateOrderController,
	}
}

func (s *OrdersRouter) SetupRoutes() {
	r := s.engine.Group("/orders/consumer")
	{
		r.GET("/", s.getOrdersController.GetOrders)
		r.PUT("/", s.updateOrderController.UpdateOrder)
	}
}