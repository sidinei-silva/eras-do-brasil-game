package npc

type Registry struct {
	npcs map[string]*NPC
}

func NewRegistry() *Registry {
	return &Registry{
		npcs: make(map[string]*NPC),
	}
}

func (r *Registry) ProcessTick(period string) {
	for _, n := range r.npcs {
		n.ProcessTick(period)
	}
}

func (r *Registry) All() []*NPC {
	result := make([]*NPC, 0, len(r.npcs))
	for _, n := range r.npcs {
		result = append(result, n)
	}
	return result
}

func (r *Registry) Get(id string) (*NPC, bool) {
	n, ok := r.npcs[id]
	return n, ok
}

func (r *Registry) Count() int {
	return len(r.npcs)
}

func (r *Registry) AllStates() []NPCState {
	result := make([]NPCState, 0, len(r.npcs))
	for _, n := range r.npcs {
		result = append(result, n.State())
	}
	return result
}

func (r *Registry) State(id string) (NPCState, bool) {
	n, ok := r.npcs[id]
	if !ok {
		return NPCState{}, false
	}
	return n.State(), true
}
