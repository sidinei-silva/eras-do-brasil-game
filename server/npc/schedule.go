package npc

// ScheduleEntry define onde o NPC vai e o que faz durante um período do dia.
type ScheduleEntry struct {
	ZoneID   string
	Activity Activity
}

// Schedule é a agenda diária do NPC: mapeia cada período para uma ScheduleEntry.
//
// Em Go, map[K]V é um tipo embutido de tabela hash — lookup O(1).
// Chaves são os períodos do jogo: "Manha", "Tarde", "Noite", "Madrugada".
//
// Por que não usar uma struct com 4 campos fixos (Manha, Tarde...)?
// O map é mais flexível: se o GDD adicionar um quinto período no futuro,
// o código não muda — só o JSON muda.
type Schedule map[string]ScheduleEntry
