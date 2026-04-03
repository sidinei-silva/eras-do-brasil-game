package world

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed data/mata-costeira.json
var mataCosteiraSrc []byte

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

type WorldMap struct {
	RegionID   string
	RegionName string
	zones      map[string]*Zone
}

func (m *WorldMap) Zone(id string) (*Zone, bool) {
	z, ok := m.zones[id]
	return z, ok
}

func (m *WorldMap) Zones() []*Zone {
	result := make([]*Zone, 0, len(m.zones))
	for _, z := range m.zones {
		result = append(result, z)
	}
	return result
}

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

func LoadMataCosteira() (*WorldMap, error) {
	var region regionJSON
	if err := json.Unmarshal(mataCosteiraSrc, &region); err != nil {
		return nil, fmt.Errorf("loading mata-costeira.json: %w", err)
	}

	m := &WorldMap{
		RegionID:   region.ID,
		RegionName: region.Name,
		zones:      make(map[string]*Zone, len(region.Zones)),
	}

	for _, zj := range region.Zones {
		zone := &Zone{
			ID:       zj.ID,
			Name:     zj.Name,
			Type:     ZoneType(zj.Type),
			LevelMin: zj.LevelMin,
			LevelMax: zj.LevelMax,
			Safe:     zj.Safe,
		}

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
