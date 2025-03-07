package controllers

import (
	"event-driven-consumer/src/features/orders/application/models"
	"event-driven-consumer/src/features/orders/application/usecases"
	"github.com/gin-gonic/gin"
)

type UpdateOrderController struct {
	UpdateOrderUseCase *usecases.UpdateOrderUseCase
}

func NewUpdateOrderController(updateOrderUseCase *usecases.UpdateOrderUseCase) *UpdateOrderController {
	return &UpdateOrderController{UpdateOrderUseCase: updateOrderUseCase}
}

func (c *UpdateOrderController) UpdateOrder(ctx *gin.Context) {
	var request models.UpdateOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": "Datos inválidos"})
		return
	}

	id := request.ID
	order, err := c.UpdateOrderUseCase.Execute(id, request.Status)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error interno del servidor"})
		return
	}

	ctx.JSON(200, order)
}
