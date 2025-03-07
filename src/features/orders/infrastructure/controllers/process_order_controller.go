package controllers

import (
	"event-driven-consumer/src/features/orders/application/usecases"
	"event-driven-consumer/src/features/orders/domain"
	"github.com/gin-gonic/gin"
)

type ProcessOrderController struct {
	ProcessOrderUseCase *usecases.ProcessOrderUseCase
}

func NewProcessOrderController(processOrderUseCase *usecases.ProcessOrderUseCase) *ProcessOrderController {
	return &ProcessOrderController{
		ProcessOrderUseCase: processOrderUseCase,
	}
}

// ProcessOrder procesa una orden recibida desde el consumidor externo
func (c *ProcessOrderController) ProcessOrder(ctx *gin.Context) {
	var request domain.Order
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": "Formato de datos inválido"})
		return
	}

	// Procesar la orden (caso de uso)
	order, err := c.ProcessOrderUseCase.ReceiveOrder(request)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error al procesar la orden"})
		return
	}

	// Devolver la respuesta
	ctx.JSON(200, order)
}
