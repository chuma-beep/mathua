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
	if quality < 3 {
		return SM2{Repetitions: 0, EFactor: prev.EFactor, Interval: 1}
	}
	rep := prev.Repetitions + 1
	ef := prev.EFactor + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
	ef = math.Max(ef, 1.3)
	var interval int
	switch prev.Repetitions {
	case 0:
		interval = 1
	case 1:
		interval = 6
	default:
		interval = int(math.Round(float64(prev.Interval) * prev.EFactor))
	}
	return SM2{Repetitions: rep, EFactor: ef, Interval: interval}
}
