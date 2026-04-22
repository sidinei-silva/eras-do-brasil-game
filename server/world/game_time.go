package world

import "time"

type PeriodOfDay string

const (
	Morning   PeriodOfDay = "manhã"
	Afternoon PeriodOfDay = "tarde"
	Night     PeriodOfDay = "noite"
	MidNight  PeriodOfDay = "madrugada"
)

type GameTime struct {
	Time        time.Time
	PeriodOfDay PeriodOfDay
}

func NewGameTime() GameTime {
	time := time.Date(1500, 1, 1, 6, 0, 0, 0, time.UTC)
	periodOfDay := CalculatePeriodOfDay(time)
	return GameTime{
		Time:        time,
		PeriodOfDay: periodOfDay,
	}
}

func (gt *GameTime) AdvanceTime(duration time.Duration) {
	gt.Time = gt.Time.Add(duration)
	gt.PeriodOfDay = CalculatePeriodOfDay(gt.Time)
}

func CalculatePeriodOfDay(time time.Time) PeriodOfDay {
	hour := time.Hour()

	switch {
	case hour >= 6 && hour < 12:
		return Morning
	case hour >= 12 && hour < 18:
		return Afternoon
	case hour >= 18 && hour < 24:
		return Night
	default:
		return MidNight
	}
}
