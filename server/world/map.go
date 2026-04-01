package world

import (
	_ "embed" // necessário para ativar a diretiva //go:embed
	"encoding/json"
	"fmt"
)

// A diretiva go:embed instrui o compilador a incluir o conteúdo do arquivo
// diretamente no binário em tempo de compilação. Vantagem: nenhum arquivo
// externo precisa existir em produção — os dados ficam dentro do executável.
//
// A variável precisa ser do tipo []byte ou string, e o caminho é relativo
// ao arquivo .go que contém a diretiva.

//go:embed data/mata-costeira.json
var mataCosteiraSrc []byte

// --- Tipos internos para leitura do JSON ---
//
// Esses tipos existem apenas para decodificar o JSON. Eles são privados (letra minúscula)
// e não ficam expostos fora do pacote. Depois de carregar, convertemos para os tipos
// públicos Zone e Connection.
//
// As tags `json:"..."` dizem ao decoder qual chave JSON mapeia para qual campo.
// Ex: `json:"travel_time_min"` lê a chave "travel_time_min" do JSON e coloca em TravelTimeMin.
// `json:",omitempty"` significa: ao serializar, omite o campo se ele for o valor zero.

type connectionJSON struct {
	ZoneID        string `json:"zone_id"`
	TravelTimeMin int    `json:"travel_time_min"`
	Terrain       string `json:"terrain"`
	Secret        bool   `json:"secret,omitempty"`
	RequiresItem  string `json:"requires_item,omitempty"`
	Condition     string `json:"condition,omitempty"`
}

type zoneJSON struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Type     string           `json:"type"`
	LevelMin int              `json:"level_min"`
	LevelMax int              `json:"level_max"`
	Safe     bool             `json:"safe"`
	Adjacent []connectionJSON `json:"adjacent"`
}

type regionJSON struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	LevelMin int        `json:"level_min"`
	LevelMax int        `json:"level_max"`
	Zones    []zoneJSON `json:"zones"`
}

// --- WorldMap ---

// WorldMap é o grafo de zonas carregado em memória.
// Internamente usa map[string]*Zone para lookup O(1) por ID.
// O campo zones é privado (minúsculo) — acesso só via métodos públicos.
type WorldMap struct {
	RegionID   string
	RegionName string
	zones      map[string]*Zone
}

// Zone retorna a zona pelo ID. O segundo valor (bool) indica se foi encontrada.
// Esse padrão de retorno (valor, ok) é idiomático em Go — evita panic por nil pointer.
//
// Uso:
//
//	zone, ok := m.Zone("vila-sao-tome")
//	if !ok { /* zona não existe */ }
func (m *WorldMap) Zone(id string) (*Zone, bool) {
	z, ok := m.zones[id]
	return z, ok
}

// Zones retorna todas as zonas do mapa como slice.
// Útil para listar/iterar todas as zonas (ex: admin, snapshot).
//
// Nota: a ordem não é garantida — maps em Go não têm ordem de iteração definida.
func (m *WorldMap) Zones() []*Zone {
	// make(tipo, tamanho) cria um slice com capacidade pré-alocada.
	// Pré-alocar evita realocações de memória durante o append.
	result := make([]*Zone, 0, len(m.zones))
	for _, z := range m.zones {
		result = append(result, z)
	}
	return result
}

// Adjacent retorna as zonas diretamente alcançáveis a partir da zona com o ID dado.
// Passagens secretas e condicionais são incluídas — quem chama decide o que exibir.
// Se a zona não existir, retorna nil.
func (m *WorldMap) Adjacent(id string) []*Zone {
	zone, ok := m.zones[id]
	if !ok {
		return nil
	}

	result := make([]*Zone, 0, len(zone.Adjacent))
	for _, conn := range zone.Adjacent {
		if z, ok := m.zones[conn.ZoneID]; ok {
			result = append(result, z)
		}
	}
	return result
}

// Connection retorna os dados da conexão entre duas zonas (custo, terreno, etc).
// Retorna nil se a conexão não existir.
func (m *WorldMap) Connection(fromID, toID string) *Connection {
	zone, ok := m.zones[fromID]
	if !ok {
		return nil
	}
	for i := range zone.Adjacent {
		if zone.Adjacent[i].ZoneID == toID {
			return &zone.Adjacent[i]
		}
	}
	return nil
}

// --- Carregamento ---

// LoadMataCosteira lê o JSON embutido e constrói o WorldMap da Região 1.
// Retorna erro se o JSON for inválido — o chamador decide como tratar (panic, log, fallback).
//
// Conceito Go: fmt.Errorf com %w "envolve" (wrap) o erro original.
// Isso permite que quem receber o erro use errors.Is/errors.As para inspecionar
// a causa raiz, sem perder a mensagem contextual.
func LoadMataCosteira() (*WorldMap, error) {
	var region regionJSON

	// json.Unmarshal converte []byte (JSON) para uma struct Go.
	// O segundo argumento é um ponteiro (&region) — necessário para que
	// Unmarshal consiga modificar a variável, não uma cópia dela.
	if err := json.Unmarshal(mataCosteiraSrc, &region); err != nil {
		return nil, fmt.Errorf("loading mata-costeira.json: %w", err)
	}

	m := &WorldMap{
		RegionID:   region.ID,
		RegionName: region.Name,
		// make(map[K]V, hint) cria um map com capacidade inicial hint.
		// Não é obrigatório, mas evita rehashing quando o map cresce.
		zones: make(map[string]*Zone, len(region.Zones)),
	}

	for _, zj := range region.Zones {
		zone := &Zone{
			ID:       zj.ID,
			Name:     zj.Name,
			Type:     ZoneType(zj.Type), // conversão explícita de string → ZoneType
			LevelMin: zj.LevelMin,
			LevelMax: zj.LevelMax,
			Safe:     zj.Safe,
		}

		// Converte cada connectionJSON para Connection.
		// Pré-alocamos o slice com a capacidade exata para evitar realocações.
		zone.Adjacent = make([]Connection, 0, len(zj.Adjacent))
		for _, cj := range zj.Adjacent {
			zone.Adjacent = append(zone.Adjacent, Connection{
				ZoneID:        cj.ZoneID,
				TravelTimeMin: cj.TravelTimeMin,
				Terrain:       cj.Terrain,
				Secret:        cj.Secret,
				RequiresItem:  cj.RequiresItem,
				Condition:     cj.Condition,
			})
		}

		m.zones[zone.ID] = zone
	}

	return m, nil
}
