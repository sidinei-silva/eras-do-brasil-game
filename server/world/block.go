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
	Terrain           string // The type of terrain, e.g., "distortion_path", "forest_path", "stone", etc.
	ToBlockId         string // The ID of the block this connection leads to
	TravelTimeMinutes int    // The time it takes to travel through this connection, in minutes
}

type Block struct {
	Id          string       // Unique identifier for the block
	Name        string       // Name of the block
	Type        string       // Enum: "urban" | "dense_forest" | "mountain" | "water" | "ruins" | "wasteland"
	Description string       // Description of the block
	Connections []Connection // List of connections to other blocks
	LevelRange  [2]int       // The recommended level range for this block, e.g., [1, 10]
	Region      Region       // The region this block belongs to, e.g., "Forest", "Desert", etc.
	Tags        []string     // List of tags for additional categorization, e.g., ["dangerous", "resource-rich"]
}

func NewBlock(id, name, blockType, description string, levelRange [2]int, region Region, tags []string) *Block {
	return &Block{
		Id:          id,
		Name:        name,
		Type:        blockType,
		Description: description,
		LevelRange:  levelRange,
		Region:      region,
		Tags:        tags,
	}
}
