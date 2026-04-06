package world

import "time"

type PeriodOfDay string

const (
	Morning   PeriodOfDay = "Manhã"
	Afternoon PeriodOfDay = "Tarde"
	Night     PeriodOfDay = "Noite"
	MidNight  PeriodOfDay = "Madrugada"
)

type GameTime struct {
	Time time.Time
	PeriodOfDay PeriodOfDay
}


func NewGameTime() GameTime {
	return GameTime{
		Time: time.Date(1500, 1, 1, 6, 0, 0, 0, time.UTC),
	}
}

func (gt *GameTime) AdvanceTime(duration time.Duration) {
	gt.Time = gt.Time.Add(duration)

	hour := gt.Time.Hour()

	switch {
	case hour >= 6 && hour < 12:
		gt.PeriodOfDay = Morning
	case hour >= 12 && hour < 18:
		gt.PeriodOfDay = Afternoon
	case hour >= 18 && hour < 24:
		gt.PeriodOfDay = Night
	default:
		gt.PeriodOfDay = MidNight
	}
}
