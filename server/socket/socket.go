package socket

import (
	"net/http"
)

// WsHandler retorna um handler HTTP que delega o upgrade e sessao ao Hub.
// Isso permite injetar dependencias (hub) sem usar variavel global.
func WsHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWS(w, r)
	}
}
