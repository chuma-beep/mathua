package realanalysis

import (
	"math"
	"strconv"
	"testing"
)

// TestRAWave2HCantorRecomputes checks the three Cantor-set shapes from their
// published stage number. The length answers use the same six-decimal form the
// generator publishes, so a change to either the formula or the rounding shows
// up here rather than as a wrong answer in a lesson.
func TestRAWave2HCantorRecomputes(t *testing.T) {
	reg := raRegistry(t)
	for kp := 1; kp <= 3; kp++ {
		for i := 0; i < 100; i++ {
			p := raSampleKP(t, reg, "ra.measure.cantor_set", kp, i)
			k := raFactInt(t, p, "k")
			if k < 1 {
				t.Fatalf("cantor KP%d: stage k=%d", kp, k)
			}
			switch p.Facts["kind"] {
			case "remaining_intervals":
				if want := strconv.Itoa(1 << k); p.Answer != want {
					t.Fatalf("cantor KP%d: k=%d answers %s, want %s", kp, k, p.Answer, want)
				}
			case "removed_intervals":
				if want := strconv.Itoa(1 << (k - 1)); p.Answer != want {
					t.Fatalf("cantor KP%d: k=%d answers %s, want %s", kp, k, p.Answer, want)
				}
			case "removed_length":
				want := decimal6(1 - math.Pow(2.0/3.0, float64(k)))
				if p.Answer != want {
					t.Fatalf("cantor KP%d: k=%d answers %s, want %s", kp, k, p.Answer, want)
				}
			case "remaining_length":
				want := decimal6(math.Pow(2.0/3.0, float64(k)))
				if p.Answer != want {
					t.Fatalf("cantor KP%d: k=%d answers %s, want %s", kp, k, p.Answer, want)
				}
			default:
				t.Fatalf("cantor KP%d: unknown kind %q", kp, p.Facts["kind"])
			}
		}
	}
}
