package world

// ZoneType classifica o tipo de terreno/ambiente de uma zona.
//
// Em Go, criar um tipo nomeado sobre string (type ZoneType string) tem um benefício:
// o compilador avisa se você tentar passar uma string qualquer onde se espera um ZoneType.
// Ex: uma função que recebe ZoneType não aceita string diretamente sem conversão.
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

// Connection representa uma passagem navegável entre duas zonas.
//
// TravelTimeMin é o custo em minutos de jogo (não minutos reais).
// Condition é opcional: quando preenchido, a passagem só existe sob uma condição
// de mundo ativa (ex: "tide_low" bloqueia acesso na maré alta do Rio das Marés).
// RequiresItem é opcional: ID do item que o personagem precisa carregar para usar a passagem.
type Connection struct {
	ZoneID        string
	TravelTimeMin int
	Terrain       string
	Secret        bool   // passagem oculta — não aparece no mapa sem descoberta prévia
	RequiresItem  string // vazio = sem requisito de item
	Condition     string // vazio = sempre disponível
}

// Zone representa uma zona do mapa: um nó no grafo de exploração.
// NPCs nascem, vivem e transitam entre zonas. Jogadores também se movem entre elas.
//
// Safe=true indica ausência de spawns de inimigos (ex: Vila de São Tomé).
// Adjacent lista todas as conexões disponíveis a partir desta zona.
type Zone struct {
	ID       string
	Name     string
	Type     ZoneType
	LevelMin int
	LevelMax int
	Safe     bool
	Adjacent []Connection
}
