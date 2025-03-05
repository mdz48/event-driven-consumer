package core

import (
    "log"
    "sync"

    "github.com/gorilla/websocket"
)

// WebSocketManager gestiona las conexiones WebSocket y el envío de mensajes
type WebSocketManager struct {
    // Mapa de conexiones activas (cliente ID -> conexión)
    clients map[string]*websocket.Conn
    // Mutex para proteger el acceso al mapa de clientes
    mutex sync.RWMutex
}

// NewWebSocketManager crea un nuevo gestor de WebSockets
func NewWebSocketManager() *WebSocketManager {
    return &WebSocketManager{
        clients: make(map[string]*websocket.Conn),
    }
}

// AddClient añade un nuevo cliente al gestor
func (m *WebSocketManager) AddClient(clientID string, conn *websocket.Conn) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    m.clients[clientID] = conn
    log.Printf("Cliente WebSocket conectado: %s", clientID)
}

// RemoveClient elimina un cliente del gestor
func (m *WebSocketManager) RemoveClient(clientID string) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    if _, exists := m.clients[clientID]; exists {
        delete(m.clients, clientID)
        log.Printf("Cliente WebSocket desconectado: %s", clientID)
    }
}

// BroadcastMessage envía un mensaje a todos los clientes conectados
func (m *WebSocketManager) BroadcastMessage(message []byte) {
    m.mutex.RLock()
    defer m.mutex.RUnlock()

    for clientID, conn := range m.clients {
        err := conn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            log.Printf("Error al enviar mensaje a cliente %s: %v", clientID, err)
            // No eliminamos el cliente aquí para evitar modificar el mapa durante la iteración
        }
    }
}

// SendToClient envía un mensaje a un cliente específico
func (m *WebSocketManager) SendToClient(clientID string, message []byte) error {
    m.mutex.RLock()
    defer m.mutex.RUnlock()

    if conn, exists := m.clients[clientID]; exists {
        return conn.WriteMessage(websocket.TextMessage, message)
    }
    return nil // El cliente no existe, no es un error crítico
}