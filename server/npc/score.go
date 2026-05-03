package npc

type NeedWeight struct {
	Hunger     float64
	Fatigue    float64
	Loneliness float64
	Schedule   float64
}

func (npc *Npc) CalculateScores(hour int) map[string]float64 {
	scores := make(map[string]float64)

	scores["Hunger"] = npc.NeedsWeight.Hunger * npc.Needs.Hunger
	scores["Fatigue"] = npc.NeedsWeight.Fatigue * npc.Needs.Fatigue
	scores["Loneliness"] = npc.NeedsWeight.Loneliness * npc.Needs.Loneliness

	if _, found := npc.ActiveScheduleAt(hour); found {
		scores["Schedule"] = npc.NeedsWeight.Schedule * 100
	} else {
		scores["Schedule"] = 0
	}

	return scores
}

func (npc *Npc) PickWinnerScore(scores map[string]float64, currentActivity Activity) (string, float64) {
	maxScore := -1.0
	winner := ""

	for name, score := range scores {
		if score > maxScore {
			maxScore = score
			winner = name
		}
	}

	return winner, maxScore
}
