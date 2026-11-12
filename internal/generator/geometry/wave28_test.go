package geometry

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func wave28IDs() []string {
	return []string{
		"geo.coord.vector", "geo.coord.dot", "geo.coord.parametric_line",
		"geo.coord.normal_line", "geo.coord.equation_circle",
		"geo.coord.rotated_conic", "geo.coord.polar_curve",
		"geo.coord.area_polygon", "geo.coord.line_intersection",
		"geo.coord.three_d",
		"geo.circle.thales", "geo.circle.arc_sector", "geo.circle.segment_area",
		"geo.circle.common_tangents", "geo.circle.orthogonal",
		"geo.circle.inversion", "geo.circle.nine_point", "geo.circle.euler_line",
		"geo.solid.tetrahedron", "geo.solid.octahedron", "geo.solid.platonic",
		"geo.solid.oblique_prism", "geo.solid.composite_volume",
		"geo.solid.cavalieri", "geo.solid.rotation_volume",
		"geo.solid.truncated_cone",
		"geo.polygon.apothem", "geo.polygon.star",
		"geo.polygon.constructible", "geo.polygon.quadrilateral_area",
		"geo.transform.homothety", "geo.transform.glide",
		"geo.transform.group",
	}
}

func wave28Registry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func wave28Fact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func TestWave28AnswersAreIntegers(t *testing.T) {
	reg := wave28Registry(t)
	for _, id := range wave28IDs() {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("registry: %v", err)
		}
		for kp := 1; kp <= 3; kp++ {
			for i := 0; i < 8; i++ {
				p := gen.Generate(generator.GeneratorContext{KPOrdinal: kp})
				if _, err := strconv.Atoi(p.Answer); err != nil {
					t.Fatalf("%s kp%d: answer %q is not an integer", id, kp, p.Answer)
				}
				if p.Facts["answer"] != p.Answer {
					t.Fatalf("%s kp%d: facts answer %q does not match %q", id, kp, p.Facts["answer"], p.Answer)
				}
				if p.Explanation == "" {
					t.Fatalf("%s kp%d: empty explanation", id, kp)
				}
				if p.KPOrdinal != kp {
					t.Fatalf("%s kp%d: ordinal %d", id, kp, p.KPOrdinal)
				}
			}
		}
	}
}

func TestWave28QuestionVariety(t *testing.T) {
	reg := wave28Registry(t)
	for _, id := range wave28IDs() {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("registry: %v", err)
		}
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 12; i++ {
				seen[gen.Generate(generator.GeneratorContext{KPOrdinal: kp}).Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 12 draws", id, kp, len(seen))
			}
		}
	}
}

type wave28Case struct {
	id   string
	want func(t *testing.T, p generator.Problem, kp int) int
}

func TestWave28AnswersRecompute(t *testing.T) {
	reg := wave28Registry(t)
	cases := []wave28Case{
		{"geo.coord.vector", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "sq")
			}
			if kp == 2 {
				return wave28Fact(t, p, "a")
			}
			return wave28Fact(t, p, "b")
		}},
		{"geo.coord.dot", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "dot")
			}
			if kp == 2 {
				return wave28Fact(t, p, "a")
			}
			return wave28Fact(t, p, "b")
		}},
		{"geo.coord.parametric_line", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "x")
			}
			if kp == 2 {
				return wave28Fact(t, p, "dx")
			}
			return wave28Fact(t, p, "t")
		}},
		{"geo.coord.normal_line", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "sq")
			}
			if kp == 2 {
				return wave28Fact(t, p, "a")
			}
			return wave28Fact(t, p, "c")
		}},
		{"geo.coord.equation_circle", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "r2")
			}
			if kp == 2 {
				return wave28Fact(t, p, "r")
			}
			return wave28Fact(t, p, "h") + wave28Fact(t, p, "k")
		}},
		{"geo.coord.rotated_conic", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "trace")
			}
			if kp == 2 {
				return wave28Fact(t, p, "a")
			}
			return wave28Fact(t, p, "c")
		}},
		{"geo.coord.polar_curve", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "r")
			}
			if kp == 2 {
				return wave28Fact(t, p, "twoa")
			}
			return 2
		}},
		{"geo.coord.area_polygon", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "twice")
			}
			if kp == 2 {
				return wave28Fact(t, p, "a")
			}
			return wave28Fact(t, p, "b")
		}},
		{"geo.coord.line_intersection", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "det")
			}
			if kp == 2 {
				return wave28Fact(t, p, "slope")
			}
			return 1
		}},
		{"geo.coord.three_d", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "sq")
			}
			if kp == 2 {
				return wave28Fact(t, p, "dx")
			}
			return wave28Fact(t, p, "dy")
		}},
		{"geo.circle.thales", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "r")
			}
			if kp == 2 {
				return wave28Fact(t, p, "d")
			}
			return 2
		}},
		{"geo.circle.arc_sector", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "arc")
			}
			if kp == 2 {
				return wave28Fact(t, p, "sector")
			}
			return wave28Fact(t, p, "r")
		}},
		{"geo.circle.segment_area", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "sector")
			}
			if kp == 2 {
				return wave28Fact(t, p, "r")
			}
			return wave28Fact(t, p, "theta")
		}},
		{"geo.circle.common_tangents", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "sum")
			}
			if kp == 2 {
				return wave28Fact(t, p, "tangents")
			}
			return wave28Fact(t, p, "d")
		}},
		{"geo.circle.orthogonal", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "d2")
			}
			if kp == 2 {
				return wave28Fact(t, p, "sum")
			}
			return wave28Fact(t, p, "r1")
		}},
		{"geo.circle.inversion", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "image")
			}
			if kp == 2 {
				return wave28Fact(t, p, "r2")
			}
			return wave28Fact(t, p, "R")
		}},
		{"geo.circle.nine_point", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "r9")
			}
			if kp == 2 {
				return wave28Fact(t, p, "R")
			}
			return 2
		}},
		{"geo.circle.euler_line", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "ratio")
			}
			if kp == 2 {
				return wave28Fact(t, p, "a")
			}
			return wave28Fact(t, p, "b")
		}},
		{"geo.solid.tetrahedron", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "faces")
			}
			if kp == 2 {
				return wave28Fact(t, p, "a")
			}
			return wave28Fact(t, p, "edges")
		}},
		{"geo.solid.octahedron", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "e")
			}
			if kp == 2 {
				return wave28Fact(t, p, "v")
			}
			return wave28Fact(t, p, "f")
		}},
		{"geo.solid.platonic", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "f")
			}
			if kp == 2 {
				return 5
			}
			return 1
		}},
		{"geo.solid.oblique_prism", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "volume")
			}
			if kp == 2 {
				return wave28Fact(t, p, "base")
			}
			return wave28Fact(t, p, "height")
		}},
		{"geo.solid.composite_volume", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "total")
			}
			if kp == 2 {
				return wave28Fact(t, p, "v1")
			}
			return wave28Fact(t, p, "v2") + wave28Fact(t, p, "v3")
		}},
		{"geo.solid.cavalieri", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "count")
			}
			if kp == 2 {
				return 0
			}
			return 1
		}},
		{"geo.solid.rotation_volume", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "r2h")
			}
			if kp == 2 {
				return wave28Fact(t, p, "r2")
			}
			return wave28Fact(t, p, "h")
		}},
		{"geo.solid.truncated_cone", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "factor")
			}
			if kp == 2 {
				return wave28Fact(t, p, "volume")
			}
			return wave28Fact(t, p, "h")
		}},
		{"geo.polygon.apothem", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "a")
			}
			if kp == 2 {
				return wave28Fact(t, p, "n")
			}
			return wave28Fact(t, p, "area")
		}},
		{"geo.polygon.star", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "k")
			}
			if kp == 2 {
				return wave28Fact(t, p, "n")
			}
			return wave28Fact(t, p, "g")
		}},
		{"geo.polygon.constructible", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 5
			}
			if kp == 2 {
				return wave28Fact(t, p, "n")
			}
			return wave28Fact(t, p, "primes")
		}},
		{"geo.polygon.quadrilateral_area", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "area")
			}
			if kp == 2 {
				return wave28Fact(t, p, "a")
			}
			return wave28Fact(t, p, "b")
		}},
		{"geo.transform.homothety", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "image")
			}
			if kp == 2 {
				return wave28Fact(t, p, "k")
			}
			return wave28Fact(t, p, "d")
		}},
		{"geo.transform.glide", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave28Fact(t, p, "d")
			}
			if kp == 2 {
				return wave28Fact(t, p, "n")
			}
			return 2
		}},
		{"geo.transform.group", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave28Fact(t, p, "n")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return 2 * n
			}
			return 2
		}},
	}

	if len(cases) != len(wave28IDs()) {
		t.Fatalf("case table has %d entries, want %d", len(cases), len(wave28IDs()))
	}
	for _, tc := range cases {
		gen, err := reg.Get(tc.id)
		if err != nil {
			t.Fatalf("%s: %v", tc.id, err)
		}
		for kp := 1; kp <= 3; kp++ {
			for i := 0; i < 6; i++ {
				p := gen.Generate(generator.GeneratorContext{KPOrdinal: kp})
				want := tc.want(t, p, kp)
				got, err := strconv.Atoi(p.Answer)
				if err != nil {
					t.Fatalf("%s kp%d: answer %q", tc.id, kp, p.Answer)
				}
				if got != want {
					t.Fatalf("%s kp%d: answer %d, want %d\nquestion: %s\nfacts: %v", tc.id, kp, got, want, p.Question, p.Facts)
				}
			}
		}
	}
}
