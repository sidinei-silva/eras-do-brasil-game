package socket

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

type AdminSocket struct {
	conn *websocket.Conn
}

func NewAdminSocket() *AdminSocket {
	return &AdminSocket{}
}

func (sm *AdminSocket) openConnection(w http.ResponseWriter, r *http.Request) {

	slog.Info("Aceitando nova conexão WebSocket Admin")

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		slog.Error("Erro ao aceitar conexão WebSocket Admin", "err", err)
		return
	}
	sm.conn = conn
	slog.Info("Nova conexão WebSocket Admin estabelecida")
}

func (sm *AdminSocket) closeConnection(ctx context.Context) {
	slog.Info("Fechando conexão WebSocket Admin")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sm.conn.CloseNow(); err != nil {
		slog.Error("Erro ao fechar conexão WebSocket Admin", "err", err)
	} else {
		slog.Info("Conexão WebSocket Admin fechada com sucesso")
	}

}

func (sm *AdminSocket) Start(mux *http.ServeMux, ctx context.Context) {
	mux.HandleFunc("/ws/admin", func(w http.ResponseWriter, r *http.Request) {
		sm.openConnection(w, r)
		defer sm.closeConnection(ctx)

		for {
			typeMessage, message, err := sm.conn.Read(r.Context())

			if err != nil {
				slog.Error("Conexão WebSocket Admin encerrada ou erro de leitura", "err", err)
				sm.closeConnection(ctx)
				break
			}

			slog.Info("Mensagem recebida no WebSocket Admin", "message", string(message))

			err = sm.conn.Write(ctx, typeMessage, message)

			if err != nil {
				slog.Error("Erro ao enviar mensagem no WebSocket Admin", "err", err)
				sm.closeConnection(ctx)
				break
			}

			slog.Info("Mensagem enviada no WebSocket Admin", "message", message)
		}
	})

}
