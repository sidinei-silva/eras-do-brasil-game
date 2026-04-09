package socket

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type AdminSocket struct {
}

func NewAdminSocket() *AdminSocket {
	return &AdminSocket{}
}

func (sm *AdminSocket) Start(mux *http.ServeMux, ctx context.Context, wg *sync.WaitGroup) {
	mux.HandleFunc("/ws/admin", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Aceitando nova conexão WebSocket Admin")

		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})

		if err != nil {
			slog.Error("Erro ao aceitar conexão WebSocket Admin", "err", err)
			return
		}

		slog.Info("Nova conexão WebSocket Admin estabelecida")

		wg.Add(1)
		defer wg.Done()

		defer func() {
			slog.Info("Fechando conexão WebSocket Admin")
			_, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := conn.CloseNow(); err != nil {
				slog.Error("Erro ao fechar conexão WebSocket Admin", "err", err)
			} else {
				slog.Info("Conexão WebSocket Admin fechada com sucesso")
			}
		}()

		for {
			typeMessage, message, err := conn.Read(r.Context())

			if err != nil {
				slog.Error("Conexão WebSocket Admin encerrada ou erro de leitura", "err", err)
				break
			}

			slog.Info("Mensagem recebida no WebSocket Admin", "message", string(message))

			err = conn.Write(ctx, typeMessage, message)

			if err != nil {
				slog.Error("Erro ao enviar mensagem no WebSocket Admin", "err", err)
				break
			}

			slog.Info("Mensagem enviada no WebSocket Admin", "message", message)
		}
	})

}
