package net

import (
	"context"
	"net/http"
	"sync"

	"github.com/sidinei-silva/eras-do-brasil-game/server/net/api"
	"github.com/sidinei-silva/eras-do-brasil-game/server/net/socket"
)

func Start(ctx context.Context, wg *sync.WaitGroup) {
	mux := http.NewServeMux()

	httpServer := api.NewHTTPServer(mux)
	adminSocket := socket.NewAdminSocket()
	clientSocket := socket.NewClientSocket()
	playerHub := socket.NewHub()
	go playerHub.Run(ctx, wg)

	go httpServer.StartHTTPServer(ctx, wg)
	go adminSocket.Start(mux, ctx, wg)
	go clientSocket.Start(mux, ctx, wg)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		playerHub.ServeWS(w, r, wg)
	})

	<-ctx.Done()
	httpServer.StopHTTPServer(ctx)
}
