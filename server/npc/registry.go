package npc

// Registry mantém todos os NPCs ativos no mundo e os processa a cada tick.
//
// Por que map[string]*NPC e não []NPC?
// - Lookup por ID é O(1) com map, seria O(n) com slice.
// - O admin vai precisar inspecionar NPCs por ID.
// - O * (ponteiro) significa que o map guarda referências, não cópias —
//   ProcessTick modifica o NPC real, não uma cópia descartada.
type Registry struct {
	npcs map[string]*NPC
}

// NewRegistry cria um Registry vazio.
func NewRegistry() *Registry {
	return &Registry{
		npcs: make(map[string]*NPC),
	}
}

// ProcessTick avança o estado de todos os NPCs para o tick atual.
// Chamado pelo game loop a cada 500ms, junto com o tick do mundo.
func (r *Registry) ProcessTick(period string) {
	// range sobre map retorna (chave, valor) a cada iteração.
	// _ descarta a chave (ID) — aqui só precisamos do NPC.
	for _, n := range r.npcs {
		n.ProcessTick(period)
	}
}

// All retorna todos os NPCs como slice.
// Útil para iterar (ex: admin, snapshot, fofoca).
// A ordem não é garantida — maps em Go não têm ordem de iteração definida.
func (r *Registry) All() []*NPC {
	result := make([]*NPC, 0, len(r.npcs))
	for _, n := range r.npcs {
		result = append(result, n)
	}
	return result
}

// Get retorna um NPC pelo ID. O bool indica se foi encontrado.
func (r *Registry) Get(id string) (*NPC, bool) {
	n, ok := r.npcs[id]
	return n, ok
}

// Count retorna o total de NPCs registrados.
func (r *Registry) Count() int {
	return len(r.npcs)
}

// AllStates retorna um snapshot do estado de todos os NPCs.
// Usado pelo admin para inspecionar o mundo sem expor ponteiros internos.
func (r *Registry) AllStates() []NPCState {
	result := make([]NPCState, 0, len(r.npcs))
	for _, n := range r.npcs {
		result = append(result, n.State())
	}
	return result
}

// State retorna o estado de um NPC específico pelo ID.
// O bool indica se o NPC foi encontrado.
func (r *Registry) State(id string) (NPCState, bool) {
	n, ok := r.npcs[id]
	if !ok {
		return NPCState{}, false
	}
	return n.State(), true
}
