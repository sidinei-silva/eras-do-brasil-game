package engine

import "github.com/sidinei-silva/eras-do-brasil-game/server/world"

// MultiPublisher distribui um snapshot para múltiplos destinos.
// Usado para enviar o mesmo tick tanto para o hub de jogadores quanto para o admin.
type MultiPublisher struct {
	publishers []SnapshotPublisher
}

// NewMultiPublisher cria um publisher composto a partir de N publishers.
func NewMultiPublisher(publishers ...SnapshotPublisher) *MultiPublisher {
	return &MultiPublisher{publishers: publishers}
}

// Publish envia o snapshot para cada publisher registrado.
func (m *MultiPublisher) Publish(snap world.Snapshot) {
	for _, p := range m.publishers {
		p.Publish(snap)
	}
}
