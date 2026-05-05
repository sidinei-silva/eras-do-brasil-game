package api

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(mux *http.ServeMux) *HTTPServer {

	Handlers(mux)

	return &HTTPServer{
		server: &http.Server{
			Addr:        ":8080",
			ReadTimeout: 5 * time.Second,
			Handler:     mux,
		},
	}
}

func (httpServer *HTTPServer) StartHTTPServer(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	slog.Info("servidor HTTP iniciado", "addr", httpServer.server.Addr)
	slog.Info("admin dashboard disponível", "url", "http://localhost:8080/admin/")
	slog.Info("admin dashboard v2 disponível", "url", "http://localhost:8080/admin2/")
	if err := httpServer.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("erro no servidor", "err", err)
	}

}

func (httpServer *HTTPServer) StopHTTPServer(ctx context.Context) {
	slog.Info("Iniciando shutdown do servidor...")
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := httpServer.server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Erro durante shutdown do servidor", "err", err)
	} else {
		slog.Info("Servidor shutdown com sucesso")
	}
}
