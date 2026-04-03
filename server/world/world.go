package world

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/sidinei-silva/eras-do-brasil-game/server/npc"
)

type Snapshot struct {
	Tick      uint64    `json:"tick"`
	GameTime  time.Time `json:"game_time"`
	Period    string    `json:"period"`
	UpdatedAt time.Time `json:"updated_at"`
}

type World struct {
	mu sync.RWMutex

	tick     uint64
	gameTime time.Time

	WorldMap *WorldMap
	NPCs     *npc.Registry
}

func NewWorld() *World {
	slog.Info("Creating new world")

	gameMap, err := LoadMataCosteira()
	if err != nil {
		panic(fmt.Sprintf("failed to load game map: %v", err))
	}
	slog.Info("Game map loaded", slog.String("region", gameMap.RegionName), slog.Int("zones", len(gameMap.Zones())))

	npcs, err := npc.LoadMataCosteira()
	if err != nil {
		panic(fmt.Sprintf("failed to load NPCs: %v", err))
	}
	slog.Info("NPCs loaded", slog.Int("count", npcs.Count()))

	return &World{
		gameTime: time.Date(1500, 1, 1, 6, 0, 0, 0, time.UTC),
		WorldMap: gameMap,
		NPCs:     npcs,
	}
}

func (w *World) ProcessTick() Snapshot {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.tick++
	w.gameTime = w.gameTime.Add(1 * time.Minute)
	period := getPeriod(w.gameTime.Hour())

	w.NPCs.ProcessTick(period)

	slog.Info("Processing tick", slog.Uint64("tick", w.tick))

	return Snapshot{
		Tick:      w.tick,
		GameTime:  w.gameTime,
		Period:    period,
		UpdatedAt: time.Now().UTC(),
	}
}

func (w *World) Snapshot() Snapshot {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return Snapshot{
		Tick:      w.tick,
		GameTime:  w.gameTime,
		Period:    getPeriod(w.gameTime.Hour()),
		UpdatedAt: time.Now().UTC(),
	}
}

const (
	PeriodoManha     = "Manha"
	PeriodoTarde     = "Tarde"
	PeriodoNoite     = "Noite"
	PeriodoMadrugada = "Madrugada"
)

func getPeriod(h int) string {
	switch {
	case h >= 6 && h < 12:
		return PeriodoManha
	case h >= 12 && h < 18:
		return PeriodoTarde
	case h >= 18 && h < 22:
		return PeriodoNoite
	default:
		return PeriodoMadrugada
	}
}
