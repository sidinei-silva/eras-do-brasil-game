type Hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) Broadcast(snapshot world.Snapshot) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		err := client.WriteJSON(snapshot)
		if err != nil {
			slog.Error("Error broadcasting to client", "error", err)
			client.Close()
			delete(h.clients, client)
		}
	}
}

func (h *Hub) AddClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = true
}

func (h *Hub) RemoveClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[conn]; ok {
		conn.Close()
		delete(h.clients, conn)
	}
}