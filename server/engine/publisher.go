package engine

import "github.com/sidinei-silva/eras-do-brasil-game/server/world"

type MultiPublisher struct {
	publishers []SnapshotPublisher
}

func NewMultiPublisher(publishers ...SnapshotPublisher) *MultiPublisher {
	return &MultiPublisher{publishers: publishers}
}

func (m *MultiPublisher) Publish(snap world.Snapshot) {
	for _, p := range m.publishers {
		p.Publish(snap)
	}
}
