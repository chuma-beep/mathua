package numtheory

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestBaseConvertExplainsPlaceValue(t *testing.T) {
	reg := generator.NewRegistry()
	Register(reg)

	questionPattern := regexp.MustCompile(`^What is the value of the numeral ([0-9]+) in base ([0-9]+), written in base ten\?$`)
	explanationPattern := regexp.MustCompile(`^The numeral ([0-9]+) in base ([0-9]+) expands as (.+) = ([0-9]+)\. `)

	for i := 0; i < 200; i++ {
		p, err := reg.GenerateContext("nt.base.convert", generator.GeneratorContext{
			Difficulty: float64(i%10) / 10,
			Seed:       int64(i),
		})
		if err != nil {
			t.Fatalf("generate: %v", err)
		}

		q := questionPattern.FindStringSubmatch(p.Question)
		if q == nil {
			t.Fatalf("question does not ask for the base-ten value: %q", p.Question)
		}
		numeral, baseText := q[1], q[2]
		base, err := strconv.Atoi(baseText)
		if err != nil {
			t.Fatalf("base %q: %v", baseText, err)
		}

		want := 0
		for _, ch := range numeral {
			digit := int(ch - '0')
			if digit >= base {
				t.Fatalf("numeral %s has digit %d, which is not valid in base %d", numeral, digit, base)
			}
			want = want*base + digit
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("numeral %s in base %d answers %s, want %d", numeral, base, p.Answer, want)
		}

		e := explanationPattern.FindStringSubmatch(p.Explanation)
		if e == nil {
			t.Fatalf("explanation does not show the expansion: %q", p.Explanation)
		}
		if e[1] != numeral || e[2] != baseText || e[4] != p.Answer {
			t.Fatalf("explanation disagrees with the question or answer: %q", p.Explanation)
		}

		terms := strings.Split(e[3], " + ")
		if len(terms) != len(numeral) {
			t.Fatalf("expansion has %d terms for a %d-digit numeral: %q", len(terms), len(numeral), p.Explanation)
		}
		for j, term := range terms {
			parts := strings.Split(term, "×")
			if len(parts) != 2 {
				t.Fatalf("expansion term %q is not a digit times a place value", term)
			}
			digit, err := strconv.Atoi(parts[0])
			if err != nil || digit != int(numeral[j]-'0') {
				t.Fatalf("expansion term %q does not match digit %c", term, numeral[j])
			}
			place, err := strconv.Atoi(parts[1])
			if err != nil {
				t.Fatalf("place value %q: %v", parts[1], err)
			}
			wantPlace := 1
			for k := 0; k < len(numeral)-1-j; k++ {
				wantPlace *= base
			}
			if place != wantPlace {
				t.Fatalf("expansion term %q uses place value %d, want %d", term, place, wantPlace)
			}
		}
	}
}
