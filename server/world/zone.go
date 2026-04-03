package world

type ZoneType string

const (
	ZoneTypeVillage  ZoneType = "village"
	ZoneTypeForest   ZoneType = "forest"
	ZoneTypeMountain ZoneType = "mountain"
	ZoneTypeRiver    ZoneType = "river"
	ZoneTypeMine     ZoneType = "mine"
	ZoneTypeCamp     ZoneType = "camp"
	ZoneTypeRuins    ZoneType = "ruins"
	ZoneTypeRupture  ZoneType = "rupture"
)

type Connection struct {
	ZoneID        string
	TravelTimeMin int
	Terrain       string
	Secret        bool
	RequiresItem  string
	Condition     string
}

type Zone struct {
	ID       string
	Name     string
	Type     ZoneType
	LevelMin int
	LevelMax int
	Safe     bool
	Adjacent []Connection
}
