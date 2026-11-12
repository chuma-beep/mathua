package graphtheory

import (
	"math"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 17 tests recompute each answer from the facts the generator publishes.
// The arithmetic is deliberately repeated here: a generator that changes a
// formula without changing the published parameters fails instead of shipping a
// plausible-looking wrong number.

func w17TFact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	raw, ok := p.Facts[key]
	if !ok {
		t.Fatalf("%s: missing fact %q\nfacts: %v", p.Question, key, p.Facts)
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, raw)
	}
	return v
}

func w17TCheck(t *testing.T, p generator.Problem, want int) {
	t.Helper()
	got, err := strconv.Atoi(p.Answer)
	if err != nil {
		t.Fatalf("%s: answer %q is not an integer", p.Question, p.Answer)
	}
	if got != want {
		t.Fatalf("%s: answer %d, want %d\nfacts: %v", p.Question, got, want, p.Facts)
	}
	if p.Facts["answer"] != p.Answer {
		t.Fatalf("%s: fact answer %q disagrees with answer %q", p.Question, p.Facts["answer"], p.Answer)
	}
	if p.Explanation == "" {
		t.Fatalf("%s: empty explanation", p.Question)
	}
	if p.KPOrdinal < 1 || p.KPOrdinal > 3 {
		t.Fatalf("%s: ordinal %d out of range", p.Question, p.KPOrdinal)
	}
}

func w17TC2(n int) int { return n * (n - 1) / 2 }

func w17TCeilDiv(a, b int) int { return (a + b - 1) / b }

func w17TFactorial(n int) int {
	f := 1
	for i := 2; i <= n; i++ {
		f *= i
	}
	return f
}

type w17Case struct {
	name string
	gen  func(int) generator.Problem
	want func(*testing.T, generator.Problem, int) int
}

func w17Cases() []w17Case {
	return []w17Case{
		{
			"handshake",
			w17Handshake,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return 2 * w17TFact(t, p, "m")
				case 2:
					return w17TFact(t, p, "sum") / 2
				default:
					return w17TFact(t, p, "n") * w17TFact(t, p, "avg") / 2
				}
			},
		},
		{
			"dirac",
			w17Dirac,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					n := w17TFact(t, p, "n")
					return (n + 1) / 2
				case 2:
					n := w17TFact(t, p, "n")
					deg := w17TFact(t, p, "deg")
					total := n * deg
					if total%2 != 0 {
						total++
					}
					return total / 2
				default:
					return w17TFact(t, p, "k") - 1
				}
			},
		},
		{
			"ore",
			w17Ore,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "n")
				case 2:
					return w17TFact(t, p, "d1") + w17TFact(t, p, "d2")
				default:
					return (w17TFact(t, p, "n") + 1) / 2
				}
			},
		},
		{
			"erdos_stone",
			w17ErdosStone,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "r") - 1
				case 2:
					n := w17TFact(t, p, "n")
					r := w17TFact(t, p, "r")
					parts := make([]int, r)
					for i := 0; i < n; i++ {
						parts[i%r]++
					}
					edges := 0
					for i := 0; i < r; i++ {
						for j := i + 1; j < r; j++ {
							edges += parts[i] * parts[j]
						}
					}
					return edges
				default:
					return w17TFact(t, p, "r") + 1
				}
			},
		},
		{
			"cage",
			w17Cage,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					r := w17TFact(t, p, "r")
					g := w17TFact(t, p, "g")
					levels := (g - 1) / 2
					total := 1
					term := r
					for i := 1; i <= levels; i++ {
						total += term
						term *= r - 1
					}
					return total
				case 2:
					return w17TFact(t, p, "v") * w17TFact(t, p, "r") / 2
				default:
					return w17TFact(t, p, "v")
				}
			},
		},
		{
			"zarankiewicz",
			w17Zarankiewicz,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					n := w17TFact(t, p, "n")
					return int(math.Floor(float64(n) / 4 * (1 + math.Sqrt(float64(4*n-3)))))
				case 2:
					return w17TFact(t, p, "s") * w17TFact(t, p, "t")
				default:
					n := w17TFact(t, p, "n")
					return w17TC2(n) * w17TC2(n)
				}
			},
		},
		{
			"critical",
			w17Critical,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "k") - 1
				case 2:
					n := w17TFact(t, p, "n")
					total := n * (w17TFact(t, p, "k") - 1)
					if total%2 != 0 {
						total++
					}
					return total / 2
				default:
					return w17TC2(w17TFact(t, p, "k"))
				}
			},
		},
		{
			"goodman",
			w17Goodman,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					n := w17TFact(t, p, "n")
					return n * (n - 1) * (n - 2) / 6
				case 2:
					n := w17TFact(t, p, "n")
					d := w17TFact(t, p, "d")
					return d * (n - 1 - d)
				default:
					n := w17TFact(t, p, "n")
					mixed := w17TFact(t, p, "mixed")
					return n*(n-1)*(n-2)/6 - mixed/2
				}
			},
		},
		{
			"friendship",
			w17Friendship,
			func(t *testing.T, p generator.Problem, kp int) int {
				k := w17TFact(t, p, "k")
				switch kp {
				case 1:
					return 2*k + 1
				case 2:
					return 3 * k
				default:
					return k
				}
			},
		},
		{
			"hajnal",
			w17Hajnal,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "delta") + 1
				case 2:
					return w17TCeilDiv(w17TFact(t, p, "n"), w17TFact(t, p, "k"))
				default:
					return w17TFact(t, p, "n") / w17TFact(t, p, "k")
				}
			},
		},
		{
			"hall",
			w17Hall,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "n")
				case 2:
					return w17TFact(t, p, "n") - w17TFact(t, p, "best")
				default:
					return w17TFactorial(w17TFact(t, p, "n"))
				}
			},
		},
		{
			"konig",
			w17Konig,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "n") / 2
				case 2:
					return w17TFact(t, p, "n") / 2
				default:
					m := w17TFact(t, p, "m")
					n := w17TFact(t, p, "n")
					if n < m {
						return n
					}
					return m
				}
			},
		},
		{
			"min_cost",
			w17MinCost,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "f") * w17TFact(t, p, "c")
				case 2:
					return w17TFact(t, p, "f1")*w17TFact(t, p, "c1") +
						w17TFact(t, p, "f2")*w17TFact(t, p, "c2")
				default:
					return w17TFact(t, p, "f") * w17TFact(t, p, "c")
				}
			},
		},
		{
			"circulation",
			w17Circulation,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "in")
				case 2:
					return w17TFact(t, p, "k") * w17TFact(t, p, "l")
				default:
					return w17TFact(t, p, "k") * w17TFact(t, p, "c")
				}
			},
		},
		{
			"gomory_hu",
			w17GomoryHu,
			func(t *testing.T, p generator.Problem, kp int) int {
				n := w17TFact(t, p, "n")
				if kp == 2 {
					return w17TC2(n)
				}
				return n - 1
			},
		},
		{
			"multiway_cut",
			w17MultiwayCut,
			func(t *testing.T, p generator.Problem, kp int) int {
				k := w17TFact(t, p, "k")
				switch kp {
				case 1:
					return k
				case 2:
					return 2*k - 2
				default:
					return w17TC2(k)
				}
			},
		},
		{
			"flow_integrality",
			w17FlowIntegrality,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "c")
				case 2:
					return w17TFact(t, p, "k") * w17TFact(t, p, "c")
				default:
					return w17TFact(t, p, "d") * w17TFact(t, p, "c")
				}
			},
		},
		{
			"edge_disjoint",
			w17EdgeDisjoint,
			func(t *testing.T, p generator.Problem, kp int) int {
				if kp == 1 {
					return w17TFact(t, p, "n") - 1
				}
				if kp == 2 {
					return 2
				}
				return 1
			},
		},
		{
			"dinic",
			w17Dinic,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "n") - 1
				case 2:
					return int(math.Floor(math.Sqrt(float64(w17TFact(t, p, "m")))))
				default:
					n := w17TFact(t, p, "n")
					return n * n * w17TFact(t, p, "m")
				}
			},
		},
		{
			"push_relabel",
			w17PushRelabel,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return 2*w17TFact(t, p, "n") - 1
				case 2:
					return w17TFact(t, p, "d") * w17TFact(t, p, "c")
				default:
					return w17TFact(t, p, "l")
				}
			},
		},
		{
			"choosability",
			w17Choosability,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "n")
				case 2:
					return 2
				default:
					return 3
				}
			},
		},
		{
			"grundy",
			w17Grundy,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "n")
				case 2:
					return 2
				default:
					return w17TFact(t, p, "n") - 1
				}
			},
		},
		{
			"achromatic",
			w17Achromatic,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "n")
				case 2:
					return (w17TFact(t, p, "n") + 1) / 2
				default:
					return 2
				}
			},
		},
		{
			"fractional",
			w17Fractional,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "n")
				case 2:
					return w17TFact(t, p, "k")
				default:
					return 2
				}
			},
		},
		{
			"circular",
			w17Circular,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return 2*w17TFact(t, p, "k") + 1
				case 2:
					return w17TFact(t, p, "n")
				default:
					return 2
				}
			},
		},
		{
			"defective",
			w17Defective,
			func(t *testing.T, p generator.Problem, kp int) int {
				d := w17TFact(t, p, "d")
				switch kp {
				case 1:
					return w17TCeilDiv(w17TFact(t, p, "delta")+1, d+1)
				case 2:
					return w17TFact(t, p, "n") * d / 2
				default:
					return w17TFact(t, p, "k")*(d+1) - 1
				}
			},
		},
		{
			"acyclic",
			w17Acyclic,
			func(t *testing.T, p generator.Problem, kp int) int {
				if kp == 2 {
					return w17TFact(t, p, "n")
				}
				if kp == 1 {
					return 2
				}
				return 3
			},
		},
		{
			"total",
			w17Total,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "delta") + 1
				case 2:
					n := w17TFact(t, p, "n")
					if n%2 == 0 {
						return n + 1
					}
					return n
				default:
					return 3
				}
			},
		},
		{
			"hadwiger",
			w17Hadwiger,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "n")
				case 2:
					return 3
				default:
					return 2
				}
			},
		},
		{
			"complete_coloring",
			w17CompleteColoring,
			func(t *testing.T, p generator.Problem, kp int) int {
				if kp == 2 {
					m := w17TFact(t, p, "m")
					k := 0
					for (k+1)*k/2 <= m {
						k++
					}
					return k
				}
				return w17TC2(w17TFact(t, p, "k"))
			},
		},
		{
			"incidence",
			w17Incidence,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "n") * w17TFact(t, p, "m")
				case 2:
					return 2 * w17TFact(t, p, "m")
				default:
					return w17TFact(t, p, "n") - 1
				}
			},
		},
		{
			"normalized",
			w17Normalized,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "n")
				case 2:
					return w17TFact(t, p, "c")
				default:
					return w17TFact(t, p, "n")
				}
			},
		},
		{
			"signless",
			w17Signless,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return 2 * w17TFact(t, p, "m")
				case 2:
					return w17TFact(t, p, "d")
				default:
					return 1
				}
			},
		},
		{
			"interlacing",
			w17Interlacing,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "n")
				case 2:
					return w17TFact(t, p, "n") - w17TFact(t, p, "k")
				default:
					return w17TFact(t, p, "k") + 1
				}
			},
		},
		{
			"closed_walks",
			w17ClosedWalks,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return 2 * w17TFact(t, p, "m")
				case 2:
					return 6 * w17TFact(t, p, "t")
				default:
					return w17TFact(t, p, "d")
				}
			},
		},
		{
			"expander",
			w17Expander,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "d") * w17TFact(t, p, "n") / 4
				case 2:
					return w17TFact(t, p, "lambda") * w17TFact(t, p, "s")
				default:
					return w17TFact(t, p, "d") / 4
				}
			},
		},
		{
			"cheeger",
			w17Cheeger,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "n")
				case 2:
					return w17TFact(t, p, "n") / 2
				default:
					return w17TFact(t, p, "k") * w17TFact(t, p, "s")
				}
			},
		},
		{
			"heat",
			w17Heat,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TFact(t, p, "n")
				case 2:
					return w17TFact(t, p, "h")
				default:
					return w17TFact(t, p, "c")
				}
			},
		},
		{
			"distance",
			w17Distance,
			func(t *testing.T, p generator.Problem, kp int) int {
				n := w17TFact(t, p, "n")
				switch kp {
				case 1:
					return n * n
				case 2:
					return n * (n*n - 1) / 6
				default:
					return n - 1
				}
			},
		},
		{
			"energy",
			w17Energy,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return 2 * (w17TFact(t, p, "n") - 1)
				case 2:
					return 2 * w17TFact(t, p, "a") * w17TFact(t, p, "b")
				default:
					return 2*(w17TFact(t, p, "a")-1) + 2*(w17TFact(t, p, "b")-1)
				}
			},
		},
		{
			"rotation",
			w17Rotation,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return 2 * w17TFact(t, p, "m")
				case 2:
					return w17TFact(t, p, "m") - w17TFact(t, p, "n") + 2
				default:
					return w17TFact(t, p, "d")
				}
			},
		},
		{
			"nonorientable",
			w17Nonorientable,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return 2 - w17TFact(t, p, "n") + w17TFact(t, p, "m") - w17TFact(t, p, "f")
				case 2:
					return w17TFact(t, p, "k")
				default:
					return w17TFact(t, p, "m") - w17TFact(t, p, "n") + 2 - w17TFact(t, p, "k")
				}
			},
		},
		{
			"toroidal",
			w17Toroidal,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return 3 * w17TFact(t, p, "n")
				case 2:
					return 3 * w17TFact(t, p, "n")
				default:
					return 2 * w17TFact(t, p, "n")
				}
			},
		},
		{
			"one_planar",
			w17OnePlanar,
			func(t *testing.T, p generator.Problem, kp int) int {
				n := w17TFact(t, p, "n")
				switch kp {
				case 1:
					return 4*n - 8
				case 2:
					return n - 2
				default:
					return w17TC2(n)
				}
			},
		},
		{
			"book",
			w17Book,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return w17TCeilDiv(w17TFact(t, p, "n"), 2)
				case 2:
					return 1
				default:
					return w17TFact(t, p, "k")
				}
			},
		},
		{
			"page",
			w17Page,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return 1
				case 2:
					return 4
				default:
					return w17TCeilDiv(w17TFact(t, p, "n"), 2)
				}
			},
		},
		{
			"queue",
			w17Queue,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1, 2:
					return 1
				default:
					return w17TFact(t, p, "n") / 2
				}
			},
		},
		{
			"bandwidth",
			w17Bandwidth,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					return 1
				case 2:
					return 2
				default:
					return w17TFact(t, p, "n") - 1
				}
			},
		},
		{
			"separator",
			w17Separator,
			func(t *testing.T, p generator.Problem, kp int) int {
				switch kp {
				case 1:
					n := w17TFact(t, p, "n")
					return int(math.Ceil(2 * math.Sqrt(2*float64(n))))
				case 2:
					return 2 * w17TFact(t, p, "n") / 3
				default:
					return w17TFact(t, p, "k")
				}
			},
		},
		{
			"apex",
			w17Apex,
			func(t *testing.T, p generator.Problem, kp int) int {
				n := w17TFact(t, p, "n")
				switch kp {
				case 1:
					return n - 4
				case 2:
					return w17TC2(n - 1)
				default:
					return 3*(n-1) - 6
				}
			},
		},
	}
}

func TestWave17AnswersRecompute(t *testing.T) {
	for _, tc := range w17Cases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for kp := 1; kp <= 3; kp++ {
				for draw := 0; draw < 8; draw++ {
					p := tc.gen(kp)
					if p.KPOrdinal != kp {
						t.Fatalf("%s kp%d: ordinal %d", tc.name, kp, p.KPOrdinal)
					}
					w17TCheck(t, p, tc.want(t, p, kp))
				}
			}
		})
	}
}
