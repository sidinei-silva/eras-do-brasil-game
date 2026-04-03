package socket

import (
	"net/http"
)

func WsHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWS(w, r)
	}
}
