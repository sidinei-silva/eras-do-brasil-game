package socket

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

type ClientSocket struct {
	conn *websocket.Conn
}

func NewClientSocket() *ClientSocket {
	return &ClientSocket{}
}

func (sm *ClientSocket) openConnection(w http.ResponseWriter, r *http.Request) {

	slog.Info("Aceitando nova conexão WebSocket Client")

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		slog.Error("Erro ao aceitar conexão WebSocket Client", "err", err)
		return
	}
	sm.conn = conn
	slog.Info("Nova conexão WebSocket Client estabelecida")
}

func (sm *ClientSocket) closeConnection(ctx context.Context) {
	slog.Info("Fechando conexão WebSocket Client")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sm.conn.CloseNow(); err != nil {
		slog.Error("Erro ao fechar conexão WebSocket Client", "err", err)
	} else {
		slog.Info("Conexão WebSocket Client fechada com sucesso")
	}

}

func (sm *ClientSocket) Start(mux *http.ServeMux, ctx context.Context) {
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		sm.openConnection(w, r)
		defer sm.closeConnection(ctx)

		for {
			typeMessage, message, err := sm.conn.Read(r.Context())

			if err != nil {
				slog.Error("Conexão WebSocket Client encerrada ou erro de leitura", "err", err)
				sm.closeConnection(ctx)
				break
			}

			slog.Info("Mensagem recebida no WebSocket Client", "message", string(message))

			err = sm.conn.Write(ctx, typeMessage, message)

			if err != nil {
				slog.Error("Erro ao enviar mensagem no WebSocket Client", "err", err)
				sm.closeConnection(ctx)
				break
			}

			slog.Info("Mensagem enviada no WebSocket Client", "message", message)
		}
	})

}
