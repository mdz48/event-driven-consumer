package main

import (
    "event-driven-consumer/src/core"
    "event-driven-consumer/src/features/orders/application"
    "event-driven-consumer/src/features/orders/infrastructure"
    "event-driven-consumer/src/features/orders/infrastructure/controllers"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "log"
    "net/http"
)

type Dependencies struct {
    engine         *gin.Engine
    database       *core.Database
    rabbitMQ       *core.RabbitMQConnection
    consumer       *infrastructure.OrderConsumer
    wsManager      *core.WebSocketManager
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
    
    // Inicializar WebSocket Manager
    wsManager := core.NewWebSocketManager()

    // Inicializar el repositorio
    orderRepository := infrastructure.NewMySQL(database.Conn)

    // Inicializar casos de uso
    getOrdersUseCase := application.NewGetOrdersUseCase(orderRepository)
    updateOrderUseCase := application.NewUpdateOrderUseCase(orderRepository)

    // Inicializar el caso de uso ProcessOrderUseCase
    processOrderUseCase := application.NewProcessOrderUseCase(orderRepository)

    // Inicializar controladores
    getOrdersController := controllers.NewGetOrdersController(getOrdersUseCase)
    updateOrderController := controllers.NewUpdateOrderController(updateOrderUseCase)

    // Inicializar el motor HTTP
    engine := gin.Default()

    // Configurar rutas
    engine.GET("/orders", getOrdersController.GetOrders)
    engine.PUT("/orders", updateOrderController.UpdateOrder)
    
    // Configurar endpoint WebSocket
    var upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        // Permitir conexiones desde cualquier origen en desarrollo
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }
    
    engine.GET("/ws", func(c *gin.Context) {
        clientID := c.Query("clientId")
        if clientID == "" {
            clientID = "anonymous"
        }
        
        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            log.Printf("Error al establecer conexión WebSocket: %v", err)
            return
        }
        
        // Registrar cliente
        wsManager.AddClient(clientID, conn)
        
        // Gestionar desconexión
        go func() {
            for {
                _, _, err := conn.ReadMessage()
                if err != nil {
                    wsManager.RemoveClient(clientID)
                    conn.Close()
                    break
                }
            }
        }()
    })

    // Inicializar el consumidor (ahora con WebSocket Manager)
    consumer := infrastructure.NewOrderConsumer(rabbitMQ, processOrderUseCase, wsManager)

    return &Dependencies{
        engine:    engine,
        database:  database,
        rabbitMQ:  rabbitMQ,
        consumer:  consumer,
        wsManager: wsManager,
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