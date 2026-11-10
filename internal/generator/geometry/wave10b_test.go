package geometry

import (
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w10gTF(t *testing.T, p generator.Problem, key string) float64 {
	t.Helper()
	v, err := strconv.ParseFloat(p.Facts[key], 64)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not numeric\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w10gNear(got, want float64) bool { return math.Abs(got-want) < 1e-6 }

func TestWave10RestGeometryAnswerRecomputation(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		for draw := 0; draw < 40; draw++ {
			p := (&w10gOrthocenterGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			cx, cy, base := w10GTestInt(t, p, "a"), w10GTestInt(t, p, "b"), w10GTestInt(t, p, "c")
			hx, hy := w10GTestInt(t, p, "hx"), w10GTestInt(t, p, "hy")
			if hx != cx || hy != cx*(base-cx)/cy {
				t.Fatalf("orthocenter kp%d: facts hx=%d hy=%d do not match\nfacts: %v", kp, hx, hy, p.Facts)
			}
			want := hx
			if kp == 2 {
				want = hy
			} else if kp == 3 {
				want = hx + hy
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("orthocenter kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gCircumcenterGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			bLeg, cLeg, radius := w10GTestInt(t, p, "a"), w10GTestInt(t, p, "b"), w10GTestInt(t, p, "c")
			x, y := w10GTestInt(t, p, "x"), w10GTestInt(t, p, "y")
			if x != bLeg/2 || y != cLeg/2 || radius*radius != (bLeg*bLeg+cLeg*cLeg)/4 {
				t.Fatalf("circumcenter kp%d: inconsistent facts\nfacts: %v", kp, p.Facts)
			}
			want = x
			if kp == 2 {
				want = radius
			} else if kp == 3 {
				want = x + y
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("circumcenter kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gIncenterGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			a, b, hyp, r := w10GTestInt(t, p, "a"), w10GTestInt(t, p, "b"), w10GTestInt(t, p, "c"), w10GTestInt(t, p, "r")
			if r != (a+b-hyp)/2 {
				t.Fatalf("incenter kp%d: r=%d is not (a+b-c)/2\nfacts: %v", kp, r, p.Facts)
			}
			want = r
			if kp == 2 {
				want = 2 * r
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("incenter kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gCevaGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			x, y, z := w10GTestInt(t, p, "a"), w10GTestInt(t, p, "b"), w10GTestInt(t, p, "c")
			wantF := float64(x * y * z)
			if kp == 2 {
				wantF = float64(z)
			} else if kp == 3 {
				wantF = 1 / float64(x*y)
			}
			if !w10gNear(w10gTF(t, p, "answer"), wantF) {
				t.Fatalf("ceva kp%d: answer %s, want %v\nfacts: %v", kp, p.Answer, wantF, p.Facts)
			}

			p = (&w10gApolloniusGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			d, h, s := w10GTestInt(t, p, "a"), w10GTestInt(t, p, "b"), w10GTestInt(t, p, "c")
			if d*d+h*h != s*s {
				t.Fatalf("apollonius kp%d: triple is not Pythagorean\nfacts: %v", kp, p.Facts)
			}
			want = h
			if kp == 2 {
				want = s
			} else if kp == 3 {
				want = s*s - d*d
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("apollonius kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gSagittaGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			half, sag, R := w10GTestInt(t, p, "a"), w10GTestInt(t, p, "b"), w10GTestInt(t, p, "c")
			if (R-sag)*(R-sag)+half*half != R*R {
				t.Fatalf("sagitta kp%d: inconsistent triple\nfacts: %v", kp, p.Facts)
			}
			want = R
			if kp == 2 {
				want = sag
			} else if kp == 3 {
				want = 2 * half
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("sagitta kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gRadicalAxisGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			d, r, intersections := w10GTestInt(t, p, "a"), w10GTestInt(t, p, "b"), w10GTestInt(t, p, "c")
			want = d / 2
			if kp == 2 {
				want = d * d
			} else if kp == 3 {
				want = 2
				if d == 2*r {
					want = 1
				} else if d > 2*r {
					want = 0
				}
				if want != intersections {
					t.Fatalf("radical kp3: intersections fact %d, want %d", intersections, want)
				}
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("radical kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gDiagonalGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			n := w10GTestInt(t, p, "a")
			diag, segments := w10GTestInt(t, p, "b"), w10GTestInt(t, p, "c")
			if diag != n*(n-3)/2 || segments != n*(n-1)/2 {
				t.Fatalf("diagonal kp%d: facts do not match formulas\nfacts: %v", kp, p.Facts)
			}
			want = diag
			if kp == 2 {
				want = n
			} else if kp == 3 {
				want = segments
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("diagonal kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gTessellationGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			n = w10GTestInt(t, p, "a")
			sum, angle, used := w10GTestInt(t, p, "b"), w10GTestInt(t, p, "c"), w10GTestInt(t, p, "used")
			if sum != (n-2)*180 || angle != sum/n {
				t.Fatalf("tessellation kp%d: inconsistent angles\nfacts: %v", kp, p.Facts)
			}
			want = sum
			if kp == 2 {
				want = angle
			} else if kp == 3 {
				want = 360 - used*angle
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("tessellation kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gSphereGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			h, k, l := w10GTestInt(t, p, "a"), w10GTestInt(t, p, "b"), w10GTestInt(t, p, "c")
			r = w10GTestInt(t, p, "r")
			px, py, pz := w10GTestInt(t, p, "px"), w10GTestInt(t, p, "py"), w10GTestInt(t, p, "pz")
			dist2 := (px-h)*(px-h) + (py-k)*(py-k) + (pz-l)*(pz-l)
			want = r
			if kp == 2 {
				want = h
			} else if kp == 3 {
				want = dist2
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("sphere kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gCompositionGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			x, y = w10GTestInt(t, p, "a"), w10GTestInt(t, p, "b")
			tx := w10GTestInt(t, p, "tx")
			want = x + y
			if kp == 2 {
				want = -y + tx
			} else if kp == 3 {
				want = y
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("composition kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gIsometryGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			x1, y1 := w10GTestInt(t, p, "a"), w10GTestInt(t, p, "b")
			x2, y2 := w10GTestInt(t, p, "x2"), w10GTestInt(t, p, "y2")
			dist2 = (x2-x1)*(x2-x1) + (y2-y1)*(y2-y1)
			want = dist2
			if kp == 2 {
				want = w10GTestInt(t, p, "area")
			} else if kp == 3 {
				want = -y1
			}
			if p.Answer != strconv.Itoa(want) || (kp == 1 && w10GTestInt(t, p, "c") != dist2) {
				t.Fatalf("isometry kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gHemisphereGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			r = w10GTestInt(t, p, "a")
			want = 2 * r * r * r
			if kp == 2 {
				want = 6 * r * r
			} else if kp == 3 {
				want = 9 * r * r
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("hemisphere kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gDensityGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			mass, volume, density := w10GTestInt(t, p, "a"), w10GTestInt(t, p, "b"), w10GTestInt(t, p, "c")
			if mass != volume*density {
				t.Fatalf("density kp%d: mass is not volume times density\nfacts: %v", kp, p.Facts)
			}
			want = density
			if kp == 2 {
				want = mass
			} else if kp == 3 {
				want = volume
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("density kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}
		}
	}
}

func TestWave10RestGeometryDiversity(t *testing.T) {
	gens := []struct {
		name string
		gen  generator.Generator
	}{
		{"orthocenter", &w10gOrthocenterGen{}}, {"circumcenter", &w10gCircumcenterGen{}}, {"incenter", &w10gIncenterGen{}},
		{"ceva", &w10gCevaGen{}}, {"apollonius", &w10gApolloniusGen{}}, {"sagitta", &w10gSagittaGen{}},
		{"radical-axis", &w10gRadicalAxisGen{}}, {"diagonal", &w10gDiagonalGen{}}, {"tessellation", &w10gTessellationGen{}},
		{"sphere", &w10gSphereGen{}}, {"composition", &w10gCompositionGen{}}, {"isometry", &w10gIsometryGen{}},
		{"hemisphere", &w10gHemisphereGen{}}, {"density", &w10gDensityGen{}},
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
