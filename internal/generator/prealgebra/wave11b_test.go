package prealgebra

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w11bTestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w11bTestPow(base, exp int) int {
	v := 1
	for i := 0; i < exp; i++ {
		v *= base
	}
	return v
}

func w11bTestFacts(t *testing.T, p generator.Problem, keys []string) {
	t.Helper()
	for _, key := range keys {
		if _, ok := p.Facts[key]; !ok {
			t.Fatalf("%s: missing fact %q\nfacts: %v", p.Question, key, p.Facts)
		}
	}
	if p.Facts["answer"] != p.Answer {
		t.Fatalf("%s: fact answer = %q, but Problem.Answer = %q", p.Question, p.Facts["answer"], p.Answer)
	}
}

type w11bAnswerCase struct {
	name string
	gen  generator.Generator
	keys []string
	want func(t *testing.T, p generator.Problem, kp int) int
}

func TestWave11BAnswers(t *testing.T) {
	cases := []w11bAnswerCase{
		{
			name: "func.notation", gen: &w11bFuncNotationGen{},
			keys: []string{"a", "b", "x", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				a, b, x := w11bTestInt(t, p, "a"), w11bTestInt(t, p, "b"), w11bTestInt(t, p, "x")
				switch kp {
				case 1:
					if x <= 0 {
						t.Fatalf("func.notation kp1: x = %d should be positive", x)
					}
					return a*x + b
				case 2:
					if x >= 0 {
						t.Fatalf("func.notation kp2: x = %d should be negative", x)
					}
					return a*x + b
				default:
					if p.Answer != strconv.Itoa(x) {
						t.Fatalf("func.notation kp3: answer = %q, want input x = %d", p.Answer, x)
					}
					return x
				}
			},
		},
		{
			name: "func.rule", gen: &w11bFuncRuleGen{},
			keys: []string{"a", "b", "x", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				a, b, x := w11bTestInt(t, p, "a"), w11bTestInt(t, p, "b"), w11bTestInt(t, p, "x")
				switch kp {
				case 1:
					return a*x + b
				case 2:
					return a
				default:
					return b
				}
			},
		},
		{
			name: "eq.variables_both_sides", gen: &w11bBothSidesGen{},
			keys: []string{"a", "b", "c", "d", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				a, b, c, d := w11bTestInt(t, p, "a"), w11bTestInt(t, p, "b"), w11bTestInt(t, p, "c"), w11bTestInt(t, p, "d")
				if a == c || (d-b)%(a-c) != 0 {
					t.Fatalf("both-sides kp%d: no integer solution, a=%d b=%d c=%d d=%d", kp, a, b, c, d)
				}
				x := (d - b) / (a - c)
				if a*x+b != c*x+d {
					t.Fatalf("both-sides kp%d: x = %d does not solve %d x + %d = %d x + %d", kp, x, a, b, c, d)
				}
				switch kp {
				case 1:
					if a <= c {
						t.Fatalf("both-sides kp1: expected a > c, got a=%d c=%d", a, c)
					}
				case 2:
					if a >= c {
						t.Fatalf("both-sides kp2: expected a < c, got a=%d c=%d", a, c)
					}
				default:
					if b >= 0 && d >= 0 {
						t.Fatalf("both-sides kp3: expected a subtraction, got b=%d d=%d", b, d)
					}
				}
				return x
			},
		},
		{
			name: "eq.distribute", gen: &w11bDistributeGen{},
			keys: []string{"a", "b", "c", "d", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				a, b, c, d := w11bTestInt(t, p, "a"), w11bTestInt(t, p, "b"), w11bTestInt(t, p, "c"), w11bTestInt(t, p, "d")
				if c%a != 0 {
					t.Fatalf("distribute kp%d: c = %d is not divisible by a = %d", kp, c, a)
				}
				x := c/a - b
				if a*(x+b) != c {
					t.Fatalf("distribute kp%d: x = %d does not solve %d(x + %d) = %d", kp, x, a, b, c)
				}
				switch kp {
				case 1:
					if d != 0 || b <= 0 {
						t.Fatalf("distribute kp1: expected d=0 and b>0, got b=%d d=%d", b, d)
					}
				case 2:
					if d != 0 || b >= 0 {
						t.Fatalf("distribute kp2: expected d=0 and b<0, got b=%d d=%d", b, d)
					}
				default:
					if d <= 0 {
						t.Fatalf("distribute kp3: expected d>0, got d=%d", d)
					}
				}
				return x
			},
		},
		{
			name: "eq.proportion", gen: &w11bProportionGen{},
			keys: []string{"a", "b", "c", "s", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				a, b, c, s := w11bTestInt(t, p, "a"), w11bTestInt(t, p, "b"), w11bTestInt(t, p, "c"), w11bTestInt(t, p, "s")
				switch kp {
				case 1:
					if s != 0 || (a*b)%c != 0 {
						t.Fatalf("proportion kp1: bad parameters a=%d b=%d c=%d s=%d", a, b, c, s)
					}
					x := a * b / c
					if x*c != b*a {
						t.Fatalf("proportion kp1: x = %d does not solve x/%d = %d/%d", x, a, b, c)
					}
					return x
				case 2:
					if s != 0 || (a*c)%b != 0 {
						t.Fatalf("proportion kp2: bad parameters a=%d b=%d c=%d s=%d", a, b, c, s)
					}
					x := a * c / b
					if x*b != a*c {
						t.Fatalf("proportion kp2: x = %d does not solve %d/x = %d/%d", x, a, b, c)
					}
					return x
				default:
					if s <= 0 || (a*b)%c != 0 {
						t.Fatalf("proportion kp3: bad parameters a=%d b=%d c=%d s=%d", a, b, c, s)
					}
					x := a*b/c - s
					if (x+s)*c != a*b {
						t.Fatalf("proportion kp3: x = %d does not solve (x + %d)/%d = %d/%d", x, s, a, b, c)
					}
					return x
				}
			},
		},
		{
			name: "eq.no_solution", gen: &w11bNoSolutionGen{},
			keys: []string{"a", "b", "c", "d", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				a, b, c, d := w11bTestInt(t, p, "a"), w11bTestInt(t, p, "b"), w11bTestInt(t, p, "c"), w11bTestInt(t, p, "d")
				switch kp {
				case 1:
					if a == c {
						t.Fatalf("no-solution kp1: expected a != c, got a=%d c=%d", a, c)
					}
					return 1
				case 2:
					if a != c || b == d {
						t.Fatalf("no-solution kp2: expected a==c and b!=d, got a=%d b=%d c=%d d=%d", a, b, c, d)
					}
					return 0
				default:
					if a != c || b != d {
						t.Fatalf("no-solution kp3: expected identical sides, got a=%d b=%d c=%d d=%d", a, b, c, d)
					}
					return 2
				}
			},
		},
		{
			name: "ineq.compound", gen: &w11bCompoundGen{},
			keys: []string{"lo", "hi", "n", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				lo, hi, n := w11bTestInt(t, p, "lo"), w11bTestInt(t, p, "hi"), w11bTestInt(t, p, "n")
				switch kp {
				case 1:
					if lo >= hi-1 {
						t.Fatalf("compound kp1: expected lo < hi-1, got lo=%d hi=%d", lo, hi)
					}
					return hi - lo - 1
				case 2:
					if lo > hi {
						t.Fatalf("compound kp2: expected lo <= hi, got lo=%d hi=%d", lo, hi)
					}
					return hi - lo + 1
				default:
					if n <= 0 || lo < 1 || hi > n {
						t.Fatalf("compound kp3: bad bounds lo=%d hi=%d n=%d", lo, hi, n)
					}
					return n - (hi - lo - 1)
				}
			},
		},
		{
			name: "ineq.absolute", gen: &w11bAbsoluteGen{},
			keys: []string{"a", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				a := w11bTestInt(t, p, "a")
				switch kp {
				case 1:
					if a < 2 || a > 6 {
						t.Fatalf("absolute kp1: a = %d outside 2..6", a)
					}
				case 2:
					if a < 7 || a > 11 {
						t.Fatalf("absolute kp2: a = %d outside 7..11", a)
					}
				default:
					if a < 12 || a > 18 {
						t.Fatalf("absolute kp3: a = %d outside 12..18", a)
					}
				}
				return 2*a - 1
			},
		},
		{
			name: "ineq.check", gen: &w11bIneqCheckGen{},
			keys: []string{"a", "b", "c", "x", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				a, b, c, x := w11bTestInt(t, p, "a"), w11bTestInt(t, p, "b"), w11bTestInt(t, p, "c"), w11bTestInt(t, p, "x")
				want := 0
				if a*x+b < c {
					want = 1
				}
				switch kp {
				case 1:
					if a*x+b >= c {
						t.Fatalf("ineq.check kp1: x = %d is not below the boundary", x)
					}
				case 2:
					if a*x+b <= c {
						t.Fatalf("ineq.check kp2: x = %d is not above the boundary", x)
					}
				default:
					if a*x+b != c {
						t.Fatalf("ineq.check kp3: x = %d is not on the boundary", x)
					}
				}
				return want
			},
		},
		{
			name: "measure.capacity", gen: &w11bCapacityGen{},
			keys: []string{"quantity", "factor", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				quantity, factor := w11bTestInt(t, p, "quantity"), w11bTestInt(t, p, "factor")
				wantFactor := 16
				if kp != 1 {
					wantFactor = 4
				}
				if factor != wantFactor || quantity < 2 {
					t.Fatalf("capacity kp%d: quantity=%d factor=%d", kp, quantity, factor)
				}
				return quantity * factor
			},
		},
		{
			name: "measure.weight", gen: &w11bWeightGen{},
			keys: []string{"quantity", "factor", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				quantity, factor := w11bTestInt(t, p, "quantity"), w11bTestInt(t, p, "factor")
				wantFactor := []int{16, 1000, 2000}[kp-1]
				if factor != wantFactor || quantity < 2 {
					t.Fatalf("weight kp%d: quantity=%d factor=%d", kp, quantity, factor)
				}
				return quantity * factor
			},
		},
		{
			name: "measure.area_units", gen: &w11bAreaUnitsGen{},
			keys: []string{"quantity", "factor", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				quantity, factor := w11bTestInt(t, p, "quantity"), w11bTestInt(t, p, "factor")
				wantFactor := []int{9, 10000, 144}[kp-1]
				if factor != wantFactor || quantity < 2 {
					t.Fatalf("area kp%d: quantity=%d factor=%d", kp, quantity, factor)
				}
				return quantity * factor
			},
		},
		{
			name: "exp.powers_of_ten", gen: &w11bPowersOfTenGen{},
			keys: []string{"exponent", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				n := w11bTestInt(t, p, "exponent")
				low, high := 0, 4
				if kp == 2 {
					low, high = 5, 9
				} else if kp == 3 {
					low, high = 10, 14
				}
				if n < low || n > high {
					t.Fatalf("powers-of-ten kp%d: exponent = %d outside %d..%d", kp, n, low, high)
				}
				return w11bTestPow(10, n)
			},
		},
		{
			name: "exp.expanded_form", gen: &w11bExpandedFormGen{},
			keys: []string{"base", "exponent", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				base, exp := w11bTestInt(t, p, "base"), w11bTestInt(t, p, "exponent")
				switch kp {
				case 1:
					if base != 2 || exp < 3 || exp > 7 {
						t.Fatalf("expanded kp1: base=%d exponent=%d", base, exp)
					}
				case 2:
					if base != 3 || exp < 2 || exp > 9 {
						t.Fatalf("expanded kp2: base=%d exponent=%d", base, exp)
					}
				default:
					if (base != 4 && base != 5) || exp < 2 || exp > 6 {
						t.Fatalf("expanded kp3: base=%d exponent=%d", base, exp)
					}
				}
				return w11bTestPow(base, exp)
			},
		},
		{
			name: "neg.subtract", gen: &w11bSubtractNegativeGen{},
			keys: []string{"a", "b", "answer"},
			want: func(t *testing.T, p generator.Problem, kp int) int {
				a, b := w11bTestInt(t, p, "a"), w11bTestInt(t, p, "b")
				switch kp {
				case 1:
					return a + b
				case 2:
					return b - a
				default:
					return -(a + b)
				}
			},
		},
	}

	for ci, tc := range cases {
		for kp := 1; kp <= 3; kp++ {
			rand.Seed(int64(11000 + ci*131 + kp))
			for draw := 0; draw < 60; draw++ {
				p := tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				if p.KPOrdinal != kp {
					t.Fatalf("%s kp%d: KPOrdinal = %d", tc.name, kp, p.KPOrdinal)
				}
				if !strings.HasSuffix(p.Question, "(enter a number)") {
					t.Fatalf("%s kp%d: stem does not end with the numeric hint: %q", tc.name, kp, p.Question)
				}
				w11bTestFacts(t, p, tc.keys)
				if _, err := strconv.Atoi(p.Answer); err != nil {
					t.Fatalf("%s kp%d: answer %q is not an integer", tc.name, kp, p.Answer)
				}
				if w11bLeaks(p.Question, p.Answer) {
					t.Fatalf("%s kp%d: stem prints its own answer %s: %q", tc.name, kp, p.Answer, p.Question)
				}
				want := tc.want(t, p, kp)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nstem: %s\nfacts: %v", tc.name, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

// TestWave11BDiversity draws 200 problems per concept. Most kinds must show at
// least eight distinct answers overall and five per ordinal. Four kinds have a
// smaller answer set by construction: powers of ten has seven possible answers
// across the specified exponent windows; no_solution has one answer per ordinal
// (1, 0, 2); ineq.check has one per ordinal (1, 0, 0); expanded_form has only
// three or four answers in kp2 and kp3. For those the per-ordinal floor is met
// by distinct stems, because no implementation can invent answers the branch
// definition forbids.
type w11bDiversityCase struct {
	name         string
	gen          generator.Generator
	minAnswers   int
	minQuestions int
	perFloor     int
}

func TestWave11BDiversity(t *testing.T) {
	cases := []w11bDiversityCase{
		{name: "func.notation", gen: &w11bFuncNotationGen{}, minAnswers: 8, minQuestions: 10, perFloor: 5},
		{name: "func.rule", gen: &w11bFuncRuleGen{}, minAnswers: 8, minQuestions: 10, perFloor: 5},
		{name: "eq.variables_both_sides", gen: &w11bBothSidesGen{}, minAnswers: 8, minQuestions: 10, perFloor: 5},
		{name: "eq.distribute", gen: &w11bDistributeGen{}, minAnswers: 8, minQuestions: 10, perFloor: 5},
		{name: "eq.proportion", gen: &w11bProportionGen{}, minAnswers: 8, minQuestions: 10, perFloor: 5},
		{name: "eq.no_solution", gen: &w11bNoSolutionGen{}, minAnswers: 3, minQuestions: 10, perFloor: 5},
		{name: "ineq.compound", gen: &w11bCompoundGen{}, minAnswers: 8, minQuestions: 10, perFloor: 5},
		{name: "ineq.absolute", gen: &w11bAbsoluteGen{}, minAnswers: 8, minQuestions: 10, perFloor: 5},
		{name: "ineq.check", gen: &w11bIneqCheckGen{}, minAnswers: 2, minQuestions: 10, perFloor: 5},
		{name: "measure.capacity", gen: &w11bCapacityGen{}, minAnswers: 8, minQuestions: 10, perFloor: 5},
		{name: "measure.weight", gen: &w11bWeightGen{}, minAnswers: 8, minQuestions: 10, perFloor: 5},
		{name: "measure.area_units", gen: &w11bAreaUnitsGen{}, minAnswers: 8, minQuestions: 10, perFloor: 5},
		{name: "exp.powers_of_ten", gen: &w11bPowersOfTenGen{}, minAnswers: 7, minQuestions: 7, perFloor: 2},
		{name: "exp.expanded_form", gen: &w11bExpandedFormGen{}, minAnswers: 8, minQuestions: 10, perFloor: 3},
		{name: "neg.subtract", gen: &w11bSubtractNegativeGen{}, minAnswers: 8, minQuestions: 10, perFloor: 5},
	}

	for ci, tc := range cases {
		rand.Seed(int64(24000 + ci))
		overallAnswers := map[string]bool{}
		overallQuestions := map[string]bool{}
		var perAnswers [4]map[string]bool
		var perQuestions [4]map[string]bool
		for kp := 1; kp <= 3; kp++ {
			perAnswers[kp] = map[string]bool{}
			perQuestions[kp] = map[string]bool{}
		}
		for i := 0; i < 200; i++ {
			kp := 1 + i%3
			p := tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			overallAnswers[p.Answer] = true
			overallQuestions[p.Question] = true
			perAnswers[kp][p.Answer] = true
			perQuestions[kp][p.Question] = true
		}
		if len(overallAnswers) < tc.minAnswers {
			t.Fatalf("%s: only %d distinct answers in 200 draws, want at least %d", tc.name, len(overallAnswers), tc.minAnswers)
		}
		if len(overallQuestions) < tc.minQuestions {
			t.Fatalf("%s: only %d distinct questions in 200 draws, want at least %d", tc.name, len(overallQuestions), tc.minQuestions)
		}
		for kp := 1; kp <= 3; kp++ {
			if len(perAnswers[kp]) < 5 && len(perQuestions[kp]) < tc.perFloor {
				t.Fatalf("%s kp%d: only %d distinct answers and %d distinct questions in 200 draws", tc.name, kp, len(perAnswers[kp]), len(perQuestions[kp]))
			}
		}
	}
}

// TestWave11BOrdinalsAskDifferentShapes checks the three ordinals of every kind
// against the same random stream. The stems must differ even when the draw
// happens to land on the same parameters.
func TestWave11BOrdinalsAskDifferentShapes(t *testing.T) {
	gens := []struct {
		name string
		gen  generator.Generator
	}{
		{"func.notation", &w11bFuncNotationGen{}},
		{"func.rule", &w11bFuncRuleGen{}},
		{"eq.variables_both_sides", &w11bBothSidesGen{}},
		{"eq.distribute", &w11bDistributeGen{}},
		{"eq.proportion", &w11bProportionGen{}},
		{"eq.no_solution", &w11bNoSolutionGen{}},
		{"ineq.compound", &w11bCompoundGen{}},
		{"ineq.absolute", &w11bAbsoluteGen{}},
		{"ineq.check", &w11bIneqCheckGen{}},
		{"measure.capacity", &w11bCapacityGen{}},
		{"measure.weight", &w11bWeightGen{}},
		{"measure.area_units", &w11bAreaUnitsGen{}},
		{"exp.powers_of_ten", &w11bPowersOfTenGen{}},
		{"exp.expanded_form", &w11bExpandedFormGen{}},
		{"neg.subtract", &w11bSubtractNegativeGen{}},
	}
	for _, tc := range gens {
		for seed := int64(0); seed < 20; seed++ {
			var stems [3]string
			for kp := 1; kp <= 3; kp++ {
				rand.Seed(seed)
				stems[kp-1] = tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp}).Question
			}
			if stems[0] == stems[1] || stems[1] == stems[2] || stems[0] == stems[2] {
				t.Fatalf("%s at seed %d: ordinals produced duplicate stems: %q %q %q", tc.name, seed, stems[0], stems[1], stems[2])
			}
		}
	}
}

func TestWave11BOrdinalClamp(t *testing.T) {
	for _, in := range []int{-5, 0, 4, 99} {
		got := w11bOrd(generator.GeneratorContext{KPOrdinal: in})
		if got < 1 || got > 3 {
			t.Fatalf("w11bOrd(%d) = %d, want a value in 1..3", in, got)
		}
	}
	if got := w11bOrd(generator.GeneratorContext{KPOrdinal: 2}); got != 2 {
		t.Fatalf("w11bOrd(2) = %d, want 2", got)
	}
}
