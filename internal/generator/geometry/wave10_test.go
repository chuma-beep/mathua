package geometry

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w10GTestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w10GTestPoint(t *testing.T, s string) (int, int) {
	t.Helper()
	s = strings.Trim(strings.TrimSpace(s), "()")
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		t.Fatalf("bad point %q", s)
	}
	x, errX := strconv.Atoi(strings.TrimSpace(parts[0]))
	y, errY := strconv.Atoi(strings.TrimSpace(parts[1]))
	if errX != nil || errY != nil {
		t.Fatalf("bad point %q", s)
	}
	return x, y
}

func TestWave10GeometryAnswerRecomputation(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		for draw := 0; draw < 60; draw++ {
			p := (&w10AltitudeGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			base, height := w10GTestInt(t, p, "base"), w10GTestInt(t, p, "height")
			area, offset := w10GTestInt(t, p, "area"), w10GTestInt(t, p, "offset")
			if area != base*height/2 {
				t.Fatalf("altitude kp%d: area %d does not match base %d and height %d", kp, area, base, height)
			}
			want := height
			if kp == 2 {
				want = area
			} else if kp == 3 {
				want = base + offset
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("altitude kp%d: answer = %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10PowerPointGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			switch kp {
			case 1:
				a, b, c, d := w10GTestInt(t, p, "a"), w10GTestInt(t, p, "b"), w10GTestInt(t, p, "c"), w10GTestInt(t, p, "d")
				if a*b != c*d || p.Answer != strconv.Itoa(d) {
					t.Fatalf("power kp1: a=%d b=%d c=%d d=%d answer=%s", a, b, c, d, p.Answer)
				}
			case 2:
				e, total, tangent := w10GTestInt(t, p, "e"), w10GTestInt(t, p, "total"), w10GTestInt(t, p, "tangent")
				if e*total != tangent*tangent || p.Answer != strconv.Itoa(tangent) {
					t.Fatalf("power kp2: e=%d total=%d tangent=%d answer=%s", e, total, tangent, p.Answer)
				}
			default:
				e1, total1, e2 := w10GTestInt(t, p, "e1"), w10GTestInt(t, p, "total1"), w10GTestInt(t, p, "e2")
				internal2 := w10GTestInt(t, p, "internal2")
				if e1*total1 != e2*(e2+internal2) || p.Answer != strconv.Itoa(internal2) {
					t.Fatalf("power kp3: e1=%d total1=%d e2=%d internal2=%d answer=%s", e1, total1, e2, internal2, p.Answer)
				}
			}

			p = (&w10SectionGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			ax, ay := w10GTestInt(t, p, "ax"), w10GTestInt(t, p, "ay")
			bx, by := w10GTestInt(t, p, "bx"), w10GTestInt(t, p, "by")
			px, py := w10GTestInt(t, p, "px"), w10GTestInt(t, p, "py")
			m, n := w10GTestInt(t, p, "m"), w10GTestInt(t, p, "n")
			gotX, gotY := w10GTestPoint(t, p.Answer)
			if kp == 3 {
				dx, dy := (px-ax)/m, (py-ay)/m
				wantX, wantY := px+n*dx, py+n*dy
				if gotX != wantX || gotY != wantY {
					t.Fatalf("section kp3: answer = %s, want (%d, %d)\nfacts: %v", p.Answer, wantX, wantY, p.Facts)
				}
			} else {
				wantX := (n*ax + m*bx) / (m + n)
				wantY := (n*ay + m*by) / (m + n)
				if gotX != wantX || gotY != wantY {
					t.Fatalf("section kp%d: answer = %s, want (%d, %d)\nfacts: %v", kp, p.Answer, wantX, wantY, p.Facts)
				}
			}
		}
	}
}

func TestWave10GeometryQuestionDiversity(t *testing.T) {
	gens := []struct {
		name string
		gen  generator.Generator
	}{
		{"altitude", &w10AltitudeGen{}},
		{"power", &w10PowerPointGen{}},
		{"section", &w10SectionGen{}},
	}
	for _, tc := range gens {
		for kp := 1; kp <= 3; kp++ {
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
