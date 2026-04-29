package npc

import (
	"time"
)

const (
	ActivityEating      Activity = "eating"
	ActivityWalking     Activity = "walking"
	ActivityIdle        Activity = "idle"
	ActivitySleeping    Activity = "sleeping"
	ActivityWorking     Activity = "working"
	ActivityBreakfast   Activity = "breakfast"
	ActivityLunch       Activity = "lunch"
	ActivityDinner      Activity = "dinner"
	ActivitySnack       Activity = "snack"
	ActivitySocializing Activity = "socializing"
)

type ActivityEffect struct {
	HungerDelta     float64
	FatigueDelta    float64
	LonelinessDelta float64
}

var activityIntensity = map[Activity]float64{
	ActivityIdle:        0.7,
	ActivityWalking:     0.9,
	ActivityWorking:     1.2,
	ActivityBreakfast:   0.5,
	ActivityLunch:       0.5,
	ActivityDinner:      0.5,
	ActivitySnack:       0.6,
	ActivitySleeping:    0.3,
	ActivitySocializing: 0.7,
}

var activityEffects = map[Activity]ActivityEffect{
	ActivityBreakfast:   {HungerDelta: -50},
	ActivityLunch:       {HungerDelta: -60},
	ActivityDinner:      {HungerDelta: -60},
	ActivitySnack:       {HungerDelta: -20},
	ActivitySleeping:    {FatigueDelta: -80},
	ActivitySocializing: {LonelinessDelta: -30},
}

var activityDuration = map[Activity]time.Duration{
	ActivityBreakfast:   1 * time.Hour,
	ActivityLunch:       1 * time.Hour,
	ActivityDinner:      1 * time.Hour,
	ActivitySnack:       30 * time.Minute,
	ActivitySleeping:    8 * time.Hour,
	ActivitySocializing: 1 * time.Hour,
}

func getActivityIntensity(activity Activity) float64 {
	if intensity, exists := activityIntensity[activity]; exists {
		return intensity
	}
	return 1.0
}

func getActivityEffects(activity Activity) ActivityEffect {
	if effects, exists := activityEffects[activity]; exists {
		return effects
	}
	return ActivityEffect{}
}

func getActivityDuration(activity Activity) time.Duration {
	if duration, exists := activityDuration[activity]; exists {
		return duration
	}
	return 1 * time.Hour
}
