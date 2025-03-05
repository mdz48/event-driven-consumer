package controllers

import (
	"event-driven-consumer/src/features/orders/application"
	"github.com/gin-gonic/gin"
)

type GetOrdersController struct {
	GetOrdersUseCase *application.GetOrdersUseCase
}

func NewGetOrdersController(getOrdersUseCase *application.GetOrdersUseCase) *GetOrdersController {
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