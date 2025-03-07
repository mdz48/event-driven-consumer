package controllers

import (
	"event-driven-consumer/src/features/orders/application/usecases"
	"github.com/gin-gonic/gin"
)

type GetOrdersController struct {
	GetOrdersUseCase *usecases.GetOrdersUseCase
}

func NewGetOrdersController(getOrdersUseCase *usecases.GetOrdersUseCase) *GetOrdersController {
	return &GetOrdersController{GetOrdersUseCase: getOrdersUseCase}
}

func (c *GetOrdersController) GetOrders(ctx *gin.Context) {
	orders, err := c.GetOrdersUseCase.Execute()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(200, orders)
}
