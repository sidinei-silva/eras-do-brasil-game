package npc

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed data/mata-costeira-npcs.json
var mataCosteiraNPCsSrc []byte

// Tipos privados usados apenas para decodificar o JSON.
// Separar os tipos de parsing dos tipos de runtime é uma boa prática:
// o JSON pode ter uma estrutura conveniente para quem escreve o arquivo,
// enquanto o código usa a estrutura mais conveniente para o jogo.

type scheduleEntryJSON struct {
	ZoneID   string `json:"zone_id"`
	Activity string `json:"activity"`
}

type npcJSON struct {
	ID        string                       `json:"id"`
	Name      string                       `json:"name"`
	SpawnZone string                       `json:"spawn_zone"`
	// map[string] aqui porque o JSON tem chaves dinâmicas ("Manha", "Tarde"...).
	// Se fossem campos fixos, usaríamos uma struct com campos nomeados.
	Schedule  map[string]scheduleEntryJSON `json:"schedule"`
}

type npcDataJSON struct {
	NPCs []npcJSON `json:"npcs"`
}

// LoadMataCosteira lê o JSON embutido e retorna um Registry populado.
//
// Padrão de retorno (*Registry, error) é idiomático em Go para funções
// que podem falhar. O chamador decide o que fazer com o erro.
func LoadMataCosteira() (*Registry, error) {
	var data npcDataJSON
	if err := json.Unmarshal(mataCosteiraNPCsSrc, &data); err != nil {
		return nil, fmt.Errorf("loading mata-costeira-npcs.json: %w", err)
	}

	r := NewRegistry()

	for _, nj := range data.NPCs {
		schedule := make(Schedule, len(nj.Schedule))
		for period, ej := range nj.Schedule {
			schedule[period] = ScheduleEntry{
				ZoneID:   ej.ZoneID,
				Activity: Activity(ej.Activity),
			}
		}

		r.npcs[nj.ID] = &NPC{
			ID:              nj.ID,
			Name:            nj.Name,
			CurrentZoneID:   nj.SpawnZone,
			Schedule:        schedule,
			currentActivity: ActivityIdle,
			// Needs começa zerada (zero value de float64 é 0.0) — NPC nasce satisfeito.
		}
	}

	return r, nil
}
