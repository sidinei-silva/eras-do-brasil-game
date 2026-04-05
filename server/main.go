package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"
)


type PeriodOfDay string

const (
	Morning   PeriodOfDay = "Manhã"
	Afternoon PeriodOfDay = "Tarde"
	Night     PeriodOfDay = "Noite"
	MidNight  PeriodOfDay = "Madrugada"
)

type GameTime struct {
	time time.Time
	PeriodOfDay PeriodOfDay
}

type GameLoop struct {
	interval time.Duration
	running  atomic.Bool
	tickCount atomic.Int64
	cancel context.CancelFunc
	lastTickDuration time.Duration
	reactionsForTick []func(gl *GameLoop)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}


func handlers(mux *http.ServeMux) {
	mux.HandleFunc("/", helloHandler)
	adminClientFileServer := http.FileServer(http.Dir("../client/adminClient"))
	mux.Handle("/admin/", http.StripPrefix("/admin/", adminClientFileServer))
}


func NewGameLoop(interval time.Duration, reactions []func(gl *GameLoop)) *GameLoop {
	return &GameLoop{
		interval: interval,
		reactionsForTick: reactions,

	}
}

func (gl *GameLoop) Start(ctx context.Context) {
	if !gl.running.CompareAndSwap(false, true) {
		slog.Warn("Game loop is already running.")
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	gl.cancel = cancel

	slog.Info("Starting game loop")

	ticker := time.NewTicker(gl.interval)
	defer ticker.Stop()

	for {
		select {
			case <- ctx.Done():
				gl.StopGameLoop()
				return
			case <- ticker.C:
				start := time.Now()
				gl.tickCount.Add(1)
				gl.lastTickDuration = time.Since(start)
				for _, reaction := range gl.reactionsForTick {
					reaction(gl)
				}
		}
	}
}

func (gl *GameLoop) StopGameLoop() {
	if !gl.running.CompareAndSwap(true, false) {
		slog.Warn("Game Loop is not running")
		return
	}

	if gl.cancel != nil {
		gl.cancel()
	}

	slog.Info("Game loop stopped")
}

func NewGameTime() GameTime {
	return GameTime{
		time: time.Date(1500, 1, 1, 6, 0, 0, 0, time.UTC),
	}
}

func (gt *GameTime) AdvanceTime(duration time.Duration) {
	gt.time = gt.time.Add(duration)

	hour := gt.time.Hour()

	switch {
	case hour >= 6 && hour < 12:
		gt.PeriodOfDay = Morning
	case hour >= 12 && hour < 18:
		gt.PeriodOfDay = Afternoon
	case hour >= 18 && hour < 24:
		gt.PeriodOfDay = Night
	default:
		gt.PeriodOfDay = MidNight
	}
}


func StartGame(){
	gameTime := NewGameTime()

	gameLoopReactions := []func(gl *GameLoop){
		func(gl *GameLoop) {
			slog.Info("GameTick", "tickCount", gl.tickCount.Load(), "tickDuration", gl.lastTickDuration)
		},

		// Advance game time every tick
		func(gl *GameLoop) {
			gameTime.AdvanceTime(1 * time.Minute)
			slog.Info("Game time advanced", "tickCount", gl.tickCount.Load(), "currentTime", gameTime.time.Format("15:04"), "periodOfDay", gameTime.PeriodOfDay)
		},
	}

	gameLoop := NewGameLoop(1 * time.Second, gameLoopReactions)
	ctx := context.Background()
	go gameLoop.Start(ctx)
}


func StartServer() {
	mux := http.NewServeMux()
	handlers(mux)
	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func main() {
	StartGame()
	StartServer()
}