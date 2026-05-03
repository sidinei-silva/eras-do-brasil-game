package npc

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
