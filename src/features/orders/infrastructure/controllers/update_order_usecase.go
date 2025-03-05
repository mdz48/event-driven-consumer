package controllers

import (
	"event-driven-consumer/src/features/orders/application"
	"event-driven-consumer/src/features/orders/domain"
	"github.com/gin-gonic/gin"
	"strconv"
)

type UpdateOrderController struct {
	UpdateOrderUseCase *application.UpdateOrderUseCase
}

func NewUpdateOrderController(updateOrderUseCase *application.UpdateOrderUseCase) *UpdateOrderController {
	return &UpdateOrderController{UpdateOrderUseCase: updateOrderUseCase}
}

func (c *UpdateOrderController) UpdateOrder(ctx *gin.Context) {
	var request domain.UpdateOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": "Datos inválidos"})
		return
	}

	id, err := strconv.Atoi(request.ID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "ID inválido"})
		return
	}

	order, err := c.UpdateOrderUseCase.Execute(id, request.State)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error interno del servidor"})
		return
	}

	ctx.JSON(200, order)
}