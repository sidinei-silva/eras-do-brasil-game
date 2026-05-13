package npc

import "strings"

// pickMealForHour escolhe a refeição apropriada baseada no horário.
// Função utilitária do pacote, não pertence a Manager nem a NPC.
func pickMealForHour(hour int) Activity {
	switch {
	case hour >= 6 && hour < 11:
		return ActivityBreakfast
	case hour >= 11 && hour < 15:
		return ActivityLunch
	case hour >= 17 && hour < 22:
		return ActivityDinner
	default:
		return ActivitySnack
	}
}

// getPairKey gera uma chave única para um par de NPCs.
// A ordem não importa: getPairKey("a", "b") == getPairKey("b", "a")
// Isso evita processar o mesmo par duas vezes.
func getPairKey(id1, id2 string) string {
	// Colocar em ordem lexicográfica para garantir unicidade
	if strings.Compare(id1, id2) > 0 {
		return id2 + "|" + id1
	}
	return id1 + "|" + id2
}
