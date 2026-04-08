package net

import (
	"context"
	"net/http"
	"sync"

	"github.com/sidinei-silva/eras-do-brasil-game/server/net/api"
	"github.com/sidinei-silva/eras-do-brasil-game/server/net/socket"
)

func Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	mux := http.NewServeMux()

	httpServer := api.NewHTTPServer(mux)
	adminSocket := socket.NewAdminSocket()
	clientSocket := socket.NewClientSocket()

	go httpServer.StartHTTPServer(ctx, wg)
	go adminSocket.Start(mux, ctx)
	go clientSocket.Start(mux, ctx)

	<-ctx.Done()
	httpServer.StopHTTPServer(ctx)
}
