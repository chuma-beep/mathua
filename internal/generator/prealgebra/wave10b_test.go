package prealgebra

import (
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w10bTF(t *testing.T, p generator.Problem, key string) float64 {
	t.Helper()
	v, err := strconv.ParseFloat(p.Facts[key], 64)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not numeric\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w10bNear(got, want float64) bool { return math.Abs(got-want) < 1e-6 }

func TestWave10RestPrealgebraAnswerRecomputation(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		for draw := 0; draw < 40; draw++ {
			p := (&w10bTapeGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			a, b, c, unit, total := w10TestInt(t, p, "a"), w10TestInt(t, p, "b"), w10TestInt(t, p, "c"), w10TestInt(t, p, "unit"), w10TestInt(t, p, "total")
			want := total
			if kp == 2 {
				want = unit
			} else if kp == 3 {
				want = c
				total = (a + b + c) * unit
			}
			if p.Answer != strconv.Itoa(want) || (kp == 1 && total != (a+b+c)*unit) {
				t.Fatalf("tape kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10bConstantGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			k, x, x2 := w10TestInt(t, p, "a"), w10TestInt(t, p, "b"), w10TestInt(t, p, "x2")
			want = k
			if kp == 2 {
				want = k * x2
			} else if kp == 3 {
				want = x2
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("constant kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}
			if w10TestInt(t, p, "c") != k*x {
				t.Fatalf("constant kp%d: table pair is not consistent", kp)
			}

			p = (&w10bPercentPointsGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			old, newer := w10TestInt(t, p, "a"), w10TestInt(t, p, "b")
			want = newer - old
			if kp == 2 {
				want = (newer - old) * 100 / old
			} else if kp == 3 {
				want = old
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("percent points kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10bSuccessiveGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			base, p1, p2 := w10TestInt(t, p, "a"), w10TestInt(t, p, "b"), w10TestInt(t, p, "c")
			after1 := base * (100 + p1) / 100
			want = after1 * (100 + p2) / 100
			if kp == 2 {
				want = after1 * (100 - p2) / 100
			} else if kp == 3 {
				want = ((100+p1)*(100+p2) - 10000) / 100
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("successive kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10bReverseGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			pct, part, whole := w10TestInt(t, p, "a"), w10TestInt(t, p, "b"), w10TestInt(t, p, "c")
			want = whole
			if kp == 2 {
				want = part
			} else if kp == 3 {
				want = pct
			}
			if p.Answer != strconv.Itoa(want) || part != whole*pct/100 {
				t.Fatalf("reverse kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10bDivDecimalGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			divisor, dividend, q := w10bTF(t, p, "divisor"), w10bTF(t, p, "dividend"), w10TestInt(t, p, "b")
			wantF := float64(q)
			if kp == 2 {
				wantF = dividend * 10
			} else if kp == 3 {
				wantF = dividend
			}
			if !w10bNear(w10bTF(t, p, "answer"), wantF) || !w10bNear(divisor*float64(q), dividend) {
				t.Fatalf("div decimal kp%d: answer %s, want %v\nfacts: %v", kp, p.Answer, wantF, p.Facts)
			}

			p = (&w10bEstimateGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			aF, bF := w10bTF(t, p, "a"), w10bTF(t, p, "b")
			ea, eb := math.Floor(aF+0.5), math.Floor(bF+0.5)
			wantF = ea + eb
			if kp == 2 {
				wantF = ea * eb
			} else if kp == 3 {
				wantF = aF + bF
			}
			if !w10bNear(w10bTF(t, p, "answer"), wantF) {
				t.Fatalf("estimate kp%d: answer %s, want %v\nfacts: %v", kp, p.Answer, wantF, p.Facts)
			}

			p = (&w10bExpandedGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			h, tens, ones := w10TestInt(t, p, "a"), w10TestInt(t, p, "b"), w10TestInt(t, p, "c")
			tenths, hundredths := w10TestInt(t, p, "tenths"), w10TestInt(t, p, "hundredths")
			wantF = float64(h*100+tens*10+ones) + float64(tenths)/10 + float64(hundredths)/100
			if kp == 1 {
				wantF = float64(tenths) / 10
				if p.Facts["asked"] == "hundredths" {
					wantF = float64(hundredths) / 100
				}
			} else if kp == 3 {
				wantF = float64(hundredths) / 100
			}
			if !w10bNear(w10bTF(t, p, "answer"), wantF) {
				t.Fatalf("expanded kp%d: answer %s, want %v\nfacts: %v", kp, p.Answer, wantF, p.Facts)
			}

			p = (&w10bTwoVariablesGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			a, b, c = w10TestInt(t, p, "a"), w10TestInt(t, p, "b"), w10TestInt(t, p, "c")
			x, y := w10TestInt(t, p, "x"), w10TestInt(t, p, "y")
			want = a*x + b*y
			if kp == 2 {
				want = a*x + b*y + c
			} else if kp == 3 {
				want = x
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("two variables kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10bCoefficientGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			a, b, c = w10TestInt(t, p, "a"), w10TestInt(t, p, "b"), w10TestInt(t, p, "c")
			want = a + b
			if kp == 2 {
				want = c
			} else if kp == 3 {
				want = 2
				if c == 0 {
					want = 1
				}
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("coefficient kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10bDistributeSimplifyGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			a, b, c = w10TestInt(t, p, "a"), w10TestInt(t, p, "b"), w10TestInt(t, p, "c")
			x = w10TestInt(t, p, "x")
			want = a
			if kp == 2 {
				want = a*b + c
			} else if kp == 3 {
				want = a*(x+b) + c
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("distribute kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10bCheckGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			a, b, c = w10TestInt(t, p, "a"), w10TestInt(t, p, "b"), w10TestInt(t, p, "c")
			x, proposed := w10TestInt(t, p, "x"), w10TestInt(t, p, "proposed")
			lhs := a*proposed + b
			want = lhs
			if kp == 2 {
				want = lhs - c
			} else if kp == 3 {
				want = x
			}
			if p.Answer != strconv.Itoa(want) || a*x+b != c {
				t.Fatalf("check kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10bTimeGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			h, m := w10TestInt(t, p, "a"), w10TestInt(t, p, "b")
			want = 60*h + m
			if kp == 2 {
				want = 60 * (60*h + m)
			} else if kp == 3 {
				want = 60*h + m + 45
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("time kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10bPrecisionGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			aF, bF = w10bTF(t, p, "a"), w10bTF(t, p, "b")
			if kp == 1 {
				s := w10bDec(aF)
				digits, seen := 0, false
				for _, r := range s {
					if r >= '1' && r <= '9' {
						seen = true
					}
					if seen && r >= '0' && r <= '9' {
						digits++
					}
				}
				if digits < 1 {
					digits = 1
				}
				if p.Answer != strconv.Itoa(digits) {
					t.Fatalf("precision kp1: answer %s, want %d\nfacts: %v", p.Answer, digits, p.Facts)
				}
			} else {
				wantF = math.Floor(aF*10+0.5) / 10
				if kp == 3 {
					wantF = math.Abs(aF - bF)
				}
				if !w10bNear(w10bTF(t, p, "answer"), wantF) {
					t.Fatalf("precision kp%d: answer %s, want %v\nfacts: %v", kp, p.Answer, wantF, p.Facts)
				}
			}
		}
	}
}

func TestWave10RestPrealgebraDiversity(t *testing.T) {
	gens := []struct {
		name string
		gen  generator.Generator
	}{
		{"tape", &w10bTapeGen{}}, {"constant", &w10bConstantGen{}}, {"points", &w10bPercentPointsGen{}},
		{"successive", &w10bSuccessiveGen{}}, {"reverse", &w10bReverseGen{}}, {"div-decimal", &w10bDivDecimalGen{}},
		{"estimate", &w10bEstimateGen{}}, {"expanded", &w10bExpandedGen{}}, {"two-variables", &w10bTwoVariablesGen{}},
		{"coefficient", &w10bCoefficientGen{}}, {"distribute", &w10bDistributeSimplifyGen{}}, {"check", &w10bCheckGen{}},
		{"time", &w10bTimeGen{}}, {"precision", &w10bPrecisionGen{}},
	}
	for _, tc := range gens {
		for kp := 1; kp <= 3; kp++ {
			rand.Seed(int64(20261110 + kp))
			seen := map[string]bool{}
			for i := 0; i < 200; i++ {
				seen[tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp}).Question] = true
			}
			if len(seen) < 10 {
				t.Fatalf("%s kp%d: only %d distinct questions in 200 draws", tc.name, kp, len(seen))
			}
			short := map[string]bool{}
			for i := 0; i < 12; i++ {
				short[tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp}).Question] = true
			}
			if len(short) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 12 draws", tc.name, kp, len(short))
			}
		}
	}
}
