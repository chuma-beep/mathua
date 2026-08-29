package scheduler

import "math"

type SM2 struct {
	Repetitions int
	EFactor     float64
	Interval    int
}

func DefaultSM2() SM2 {
	return SM2{EFactor: 2.5}
}

func ComputeSM2(prev SM2, quality int) SM2 {
	return ComputeSM2WithSpeed(prev, quality, 1.0)
}

func ComputeSM2WithSpeed(prev SM2, quality int, learningSpeed float64) SM2 {
	if learningSpeed < 0.5 {
		learningSpeed = 0.5
	}
	if learningSpeed > 2.0 {
		learningSpeed = 2.0
	}
	if quality < 3 {
		return SM2{Repetitions: 0, EFactor: prev.EFactor, Interval: 1}
	}
	rep := prev.Repetitions + 1
	ef := prev.EFactor + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
	ef = math.Max(ef, 1.3)
	var interval int
	switch prev.Repetitions {
	case 0:
		interval = int(math.Round(1 * learningSpeed))
		if interval < 1 {
			interval = 1
		}
	case 1:
		interval = int(math.Round(6 * learningSpeed))
		if interval < 1 {
			interval = 1
		}
	default:
		interval = int(math.Round(float64(prev.Interval) * prev.EFactor * learningSpeed))
		if interval < 1 {
			interval = 1
		}
	}
	return SM2{Repetitions: rep, EFactor: ef, Interval: interval}
}

// UpdateLearningSpeed adjusts per-topic speed from timeRatio + streak.
// 0.5 slow .. 1.0 normal .. 2.0 fast.
func UpdateLearningSpeed(prev float64, correct bool, timeRatio float64, streak int) float64 {
	if prev == 0 {
		prev = 1.0
	}
	if !correct {
		prev -= 0.1
		if prev < 0.5 {
			prev = 0.5
		}
		return math.Round(prev*20) / 20
	}
	if timeRatio <= 0.7 && streak >= 2 {
		prev += 0.05
	} else if timeRatio <= 1.0 {
		// normal — small drift up
		prev += 0.02
	} else {
		prev -= 0.05
	}
	if prev < 0.5 {
		prev = 0.5
	}
	if prev > 2.0 {
		prev = 2.0
	}
	return math.Round(prev*20) / 20
}
