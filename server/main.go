package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sidinei-silva/eras-do-brasil-game/server/admin"
	"github.com/sidinei-silva/eras-do-brasil-game/server/engine"
	"github.com/sidinei-silva/eras-do-brasil-game/server/socket"
	"github.com/sidinei-silva/eras-do-brasil-game/server/world"
)

// ServerError é um tipo de erro customizado (não obrigatório, mas bom para debug).
type ServerError struct {
	Message string
	Err     error
}

func (se *ServerError) Error() string {
	return se.Message
}

func (se *ServerError) Unwrap() error {
	return se.Err
}

// hubPublisher adapta o player Hub para implementar engine.SnapshotPublisher.
type hubPublisher struct {
	hub *socket.Hub
}

func (p *hubPublisher) Publish(snap world.Snapshot) {
	p.hub.Broadcast(socket.OutboundEvent{
		Type: "world_snapshot",
		Data: snap,
	})
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// --- World ---
	mundo := world.NewWorld()

	// --- Player Hub ---
	playerHub := socket.NewHub()
	go playerHub.Run(ctx)

	// --- Game Loop ---
	// O admin hub precisa do gameLoop, mas o gameLoop precisa do publisher.
	// Criamos o gameLoop primeiro com publisher nil, depois setamos.
	loop := engine.NewGameLoop(500*time.Millisecond, mundo, nil)

	// --- Admin Hub ---
	adminHub := admin.NewHub(loop, mundo, playerHub.OnlineCount)

	// O player hub notifica o admin de eventos (join, leave, chat, rename).
	playerHub.SetObserver(adminHub)

	go adminHub.Run(ctx)

	// --- Wiring do publisher ---
	// MultiPublisher envia snapshots para player hub E admin hub.
	multi := engine.NewMultiPublisher(
		&hubPublisher{hub: playerHub},
		adminHub,
	)
	loop.SetPublisher(multi)

	go loop.Start(ctx)

	// --- HTTP Routes ---
	http.HandleFunc("/admin/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"game_loop": loop.Status(),
			"world":     mundo.Snapshot(),
			"lobby": map[string]any{
				"online": playerHub.OnlineCount(),
			},
		})
	})

	http.HandleFunc("/ws", socket.WsHandler(playerHub))
	http.HandleFunc("/ws/admin", func(w http.ResponseWriter, r *http.Request) {
		adminHub.ServeWS(w, r)
	})

	// Serve arquivos estáticos do admin client.
	adminFS := http.FileServer(http.Dir("../client/adminClient"))
	http.Handle("/admin/", http.StripPrefix("/admin/", adminFS))

	server := &http.Server{
		Addr:        ":8080",
		ReadTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("server iniciado", "addr", server.Addr)
		slog.Info("admin dashboard disponivel em http://localhost:8080/admin/")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("erro no servidor", "err", err)
			stop()
		}
	}()

	<-ctx.Done()

	loop.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("erro no shutdown", "err", err)
		os.Exit(1)
	}

	slog.Info("programa finalizado")
}
