package fractions

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func fuzzGen(t *testing.T, gen generator.Generator) {
	t.Helper()
	for i := 0; i < 100; i++ {
		d := rand.Float64()
		p := gen.Generate(generator.GeneratorContext{Difficulty: d})
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field at difficulty=%.2f: q=%q a=%q e=%q", d, p.Question, p.Answer, p.Explanation)
		}
	}
}

func TestFracConceptGen(t *testing.T)      { fuzzGen(t, &fracConceptGen{}) }
func TestFracPartsGen(t *testing.T)        { fuzzGen(t, &fracPartsGen{}) }
func TestFracNumberLineGen(t *testing.T)   { fuzzGen(t, &fracNumberLineGen{}) }
func TestFracEquivalentGen(t *testing.T)   { fuzzGen(t, &fracEquivalentGen{}) }
func TestFracSimplifyGen(t *testing.T)     { fuzzGen(t, &fracSimplifyGen{}) }
func TestFracCompareGen(t *testing.T)      { fuzzGen(t, &fracCompareGen{}) }
func TestFracBenchmarkGen(t *testing.T)    { fuzzGen(t, &fracBenchmarkGen{}) }
func TestFracToDecimalGen(t *testing.T)    { fuzzGen(t, &fracToDecimalGen{}) }
func TestFracOpSameDenGen(t *testing.T)    { fuzzGen(t, &fracOpSameDenGen{}) }
func TestFracOpDiffDenGen(t *testing.T)    { fuzzGen(t, &fracOpDiffDenGen{}) }
func TestFracWordGen(t *testing.T)         { fuzzGen(t, &fracWordGen{}) }
func TestFracMultGen(t *testing.T)         { fuzzGen(t, &fracMultGen{}) }
func TestFracDivGen(t *testing.T)          { fuzzGen(t, &fracDivGen{}) }
func TestFracMixedConvertGen(t *testing.T) { fuzzGen(t, &fracMixedConvertGen{}) }
func TestFracMixedOpGen(t *testing.T)      { fuzzGen(t, &fracMixedOpGen{}) }

func TestFracBenchmarkHasMoreThanThreePairs(t *testing.T) {
	g := &fracBenchmarkGen{}
	seen := map[string]bool{}
	for i := 0; i < 500; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64()})
		seen[p.Question] = true
	}
	if len(seen) < 8 {
		t.Fatalf("expected a wider benchmark pool, got %d distinct questions", len(seen))
	}
}

func TestOfUnitNamesTheMeasuringUnit(t *testing.T) {
	g := &ofUnitGen{}
	for i := 0; i < 200; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64()})
		if strings.Contains(p.Question, "of its parts") {
			t.Fatalf("question still uses the vague 'parts' wording: %q", p.Question)
		}
		named := false
		for _, unit := range []string{"minute", "cent", "hour", "centimetre"} {
			if strings.Contains(p.Question, unit) {
				named = true
				break
			}
		}
		if !named {
			t.Fatalf("question does not name a measurable unit: %q", p.Question)
		}
		if !strings.Contains(p.Answer, "/") {
			t.Fatalf("answer is not a fraction: %q", p.Answer)
		}
	}
}

func TestFracMixedOpNamesTheNotation(t *testing.T) {
	g := &fracMixedOpGen{}
	for i := 0; i < 100; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64()})
		if !strings.Contains(p.Question, "mixed numbers") {
			t.Fatalf("expected the question to name mixed numbers, got %q", p.Question)
		}
		if !strings.Contains(p.Question, "Add") && !strings.Contains(p.Question, "Subtract") {
			t.Fatalf("expected an explicit operation, got %q", p.Question)
		}
	}
}
func TestFracMixedMultGen(t *testing.T) { fuzzGen(t, &fracMixedMultGen{}) }

func TestFuzz(t *testing.T) {
	reg := generator.NewRegistry()
	Register(reg)
	for _, id := range reg.Concepts() {
		gen, _ := reg.Get(id)
		t.Run(id, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				d := rand.Float64()
				p := gen.Generate(generator.GeneratorContext{Difficulty: d})
				if p.Question == "" || p.Answer == "" || p.Explanation == "" {
					t.Fatalf("empty field at difficulty=%.2f for %s: q=%q a=%q", d, id, p.Question, p.Answer)
				}
			}
		})
	}
}

// The bar was drawn as "X X _ _ _" inside math delimiters, where the spaces
// collapse and the underscores become subscripts, so the learner was asked what
// fraction of a bar they could not see. The question now says the counts.
func TestFractionConceptAsksInWords(t *testing.T) {
	g := &fracConceptGen{}
	for i := 0; i < 200; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64(), Seed: rand.Int63()})
		if strings.Contains(p.Question, "X X") || strings.Contains(p.Question, "\\(") {
			t.Fatalf("question still draws a bar: %q", p.Question)
		}
		if !strings.Contains(p.Question, "equal parts") {
			t.Fatalf("question does not describe the parts: %q", p.Question)
		}
		if p.Answer != p.Facts["filled"]+"/"+p.Facts["parts"] {
			t.Fatalf("answer %q does not match facts %v", p.Answer, p.Facts)
		}
	}
}

// TestRibbonIsNeverCutIntoItsOwnLength pins the degenerate draw. When the piece
// count equalled the length the generator served "a ribbon 1 metres long cut
// into 1 equal pieces" and answered 1/1, which is not a fraction exercise.
func TestRibbonIsNeverCutIntoItsOwnLength(t *testing.T) {
	g := &asDivisionGen{}
	for i := 0; i < 300; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64(), Seed: rand.Int63()})
		length := p.Facts["length"]
		pieces := p.Facts["pieces"]
		if length == pieces {
			t.Fatalf("ribbon %s metres is cut into %s pieces: %q", length, pieces, p.Question)
		}
		if strings.Contains(p.Question, "1 metres") {
			t.Fatalf("singular length reads as plural: %q", p.Question)
		}
	}
}

func TestCrossComparePrintedProductsMatchFacts(t *testing.T) {
	g := &crossCompareGen{}
	re := regexp.MustCompile(`\\\((\d+) \\times (\d+) - (\d+) \\times (\d+)\\\)`)
	for i := 0; i < 500; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64()})
		m := re.FindStringSubmatch(p.Question)
		if m == nil {
			t.Fatalf("question does not show the cross products: %q", p.Question)
		}
		nums := make([]int, 4)
		for j, raw := range m[1:] {
			v, err := strconv.Atoi(raw)
			if err != nil {
				t.Fatalf("parse %q: %v", raw, err)
			}
			nums[j] = v
		}
		if got, want := strconv.Itoa(nums[0]*nums[1]), p.Facts["left"]; got != want {
			t.Fatalf("printed first product is %s, facts say %s: %q", got, want, p.Question)
		}
		if got, want := strconv.Itoa(nums[2]*nums[3]), p.Facts["right"]; got != want {
			t.Fatalf("printed second product is %s, facts say %s: %q", got, want, p.Question)
		}
		if want := p.Facts["difference"] + "/1"; p.Answer != want {
			t.Fatalf("answer %q, want %q", p.Answer, want)
		}
	}
}
