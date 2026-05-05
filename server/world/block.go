package world

type Region string
type BlockType string

const (
	MataCosteira Region = "mata_costeira"
)

const (
	Urban       BlockType = "urban"
	DenseForest BlockType = "dense_forest"
	Mountain    BlockType = "mountain"
	Water       BlockType = "water"
	Ruins       BlockType = "ruins"
	Wasteland   BlockType = "wasteland"
)

type Connection struct {
	Terrain       string // Tipo de terreno (ex.: "distortion_path", "forest_path", "stone").
	ToBlockId     string // ID do bloco de destino desta conexão.
	TravelMinutes int    // Tempo de deslocamento por esta conexão, em minutos.
}

type LevelRange struct {
	Min int // Nível mínimo recomendado para este bloco.
	Max int // Nível máximo recomendado para este bloco.
}

type Block struct {
	Id          string       // Identificador único do bloco.
	Name        string       // Nome do bloco.
	Type        BlockType    // Enum: "urban" | "dense_forest" | "mountain" | "water" | "ruins" | "wasteland".
	Description string       // Descrição do bloco.
	Connections []Connection // Lista de conexões para outros blocos.
	LevelRange  LevelRange   // Faixa de nível recomendada para explorar este bloco.
	Region      Region       // Região à qual o bloco pertence.
	Tags        []string     // Tags de categorização extra (ex.: ["dangerous", "resource-rich"]).
	Pois        []string     // Pontos de interesse dentro do bloco (ex.: "abandoned_house", "hidden_cave").
}

func NewBlock(id string, name string, blockType BlockType, description string, levelRange LevelRange, region Region, tags []string, pois []string) *Block {
	return &Block{
		Id:          id,
		Name:        name,
		Type:        blockType,
		Description: description,
		LevelRange:  levelRange,
		Region:      region,
		Tags:        tags,
		Pois:        pois,
	}
}
