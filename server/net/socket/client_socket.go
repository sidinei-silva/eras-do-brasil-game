package socket

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type ClientSocket struct {
}

func NewClientSocket() *ClientSocket {
	return &ClientSocket{}
}

func (sm *ClientSocket) Start(mux *http.ServeMux, ctx context.Context, wg *sync.WaitGroup) {
	mux.HandleFunc("/ws_old", func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Aceitando nova conexão WebSocket Client")

		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			slog.Error("Erro ao aceitar conexão WebSocket Client", "err", err)
			return
		}

		slog.Info("Nova conexão WebSocket Client estabelecida")

		wg.Add(1)
		defer wg.Done()
		defer func() {
			slog.Info("Fechando conexão WebSocket Client")
			_, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := conn.CloseNow(); err != nil {
				slog.Error("Erro ao fechar conexão WebSocket Client", "err", err)
			} else {
				slog.Info("Conexão WebSocket Client fechada com sucesso")
			}
		}()

		for {
			typeMessage, message, err := conn.Read(r.Context())

			if err != nil {
				slog.Error("Conexão WebSocket Client encerrada ou erro de leitura", "err", err)
				break
			}

			slog.Info("Mensagem recebida no WebSocket Client", "message", string(message))

			err = conn.Write(ctx, typeMessage, message)

			if err != nil {
				slog.Error("Erro ao enviar mensagem no WebSocket Client", "err", err)
				break
			}

			slog.Info("Mensagem enviada no WebSocket Client", "message", message)
		}
	})

}
