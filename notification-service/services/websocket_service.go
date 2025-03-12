package services

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	upgrader  websocket.Upgrader
	clients   map[*websocket.Conn]bool
	broadcast chan string
	mu        sync.Mutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan string),
	}
}

// Maneja conexiones WebSocket
func (ws *WebSocketService) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar conexión:", err)
		return
	}
	defer conn.Close()

	ws.mu.Lock()
	ws.clients[conn] = true
	ws.mu.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Cliente desconectado")
			ws.mu.Lock()
			delete(ws.clients, conn)
			ws.mu.Unlock()
			break
		}
	}
}

// Maneja el envío de mensajes a los clientes
func (ws *WebSocketService) HandleMessages() {
	for {
		message := <-ws.broadcast
		ws.mu.Lock()
		for client := range ws.clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("Error al enviar mensaje:", err)
				client.Close()
				delete(ws.clients, client)
			}
		}
		ws.mu.Unlock()
	}
}

// Envía un mensaje a todos los clientes WebSocket
func (ws *WebSocketService) SendMessage(message string) {
	ws.broadcast <- message
}
