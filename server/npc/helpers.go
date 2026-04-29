package npc

// pickMealForHour escolhe a refeição apropriada baseada no horário.
// Função utilitária do pacote, não pertence a Manager nem a NPC.
func pickMealForHour(hour int) Activity {
	switch {
	case hour >= 6 && hour < 10:
		return ActivityBreakfast
	case hour >= 11 && hour < 14:
		return ActivityLunch
	case hour >= 18 && hour < 21:
		return ActivityDinner
	default:
		return ActivitySnack
	}
}
