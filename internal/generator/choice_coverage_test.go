package generator_test

import (
	"regexp"
	"strings"
	"testing"
	"unicode"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/all"
)

// A multi-word answer cannot be typed into the math field. The field strips
// spaces, so "line segment" arrives as "linesegment" or worse, and the grader
// compares that to the authored phrase. Five reports on 23 October 2026 were
// exactly this: three asked a yes/no question, one asked for a geometric term
// and one asked for a sentence.
//
// So: a concept graded multiple_choice whose answer contains a space must serve
// an option list, and the answer must be one of the options. Numeric answers
// are exempt, because they stay typeable and most of the corpus relies on that
// fallback today.
func TestMultipleChoiceProseAnswersServeOptions(t *testing.T) {
	dag, err := concepts.LoadDir("../../data/concepts")
	if err != nil {
		t.Fatalf("load concepts: %v", err)
	}
	reg := generator.NewRegistry()
	all.Register(reg)

	checked, offenders := 0, 0
	for id, c := range dag.Concepts() {
		if c == nil || c.GradingType != "multiple_choice" || !learnerDomains[c.Domain] {
			continue
		}
		gen, err := reg.Get(id)
		if err != nil {
			continue
		}
		prose := false
		bad := ""
		for i := 0; i < 60; i++ {
			p := gen.Generate(generator.GeneratorContext{Difficulty: float64(i%5) / 4, Seed: int64(i)})
			if !looksLikeWords(p.Answer) {
				continue
			}
			prose = true
			options := generator.ChoiceOptions(p.Question)
			if len(options) < 2 {
				bad = "answer is prose and the question carries no options"
				break
			}
			found := false
			for _, text := range options {
				if text == p.Answer {
					found = true
					break
				}
			}
			if !found {
				bad = "the answer is not one of the options"
				break
			}
		}
		if !prose {
			continue
		}
		checked++
		if bad != "" {
			offenders++
			sample := gen.Generate(generator.GeneratorContext{Difficulty: 0.5, Seed: 1})
			t.Errorf("%s: %s\n  question: %s\n  answer: %s", id, bad, sample.Question, sample.Answer)
		}
	}
	if checked == 0 {
		t.Skip("no multiple_choice concept remains to serve a prose answer")
	}
	t.Logf("checked %d prose-answer concepts, %d offenders", checked, offenders)
}

// mathWords are words that belong inside an expression rather than in the
// answer. "pi/2 - 1" is typeable, so it is not prose; "two real" is.
var latexCommand = regexp.MustCompile(`\\[A-Za-z]+`)

var mathWords = map[string]bool{
	"sin": true, "cos": true, "tan": true, "cot": true, "sec": true, "csc": true,
	"arc": true, "log": true, "ln": true, "exp": true, "sqrt": true, "inf": true,
}

// looksLikeWords reports whether the answer contains a word a learner could not
// type into a math field. Splitting on non-letters keeps expressions out.
func looksLikeWords(answer string) bool {
	answer = latexCommand.ReplaceAllString(answer, " ")
	for _, field := range strings.FieldsFunc(answer, func(r rune) bool { return !unicode.IsLetter(r) }) {
		w := strings.ToLower(field)
		if len(w) >= 3 && !mathWords[w] {
			return true
		}
	}
	return false
}

// A yes/no answer is a word, and a maths field cannot take a word: it strips
// the spaces and reads the letters as variables. One hundred and eighteen
// concepts across topology, calculus, complex analysis, abstract algebra and
// machine learning answered yes/no with no options, and a learner who typed
// "y" was marked wrong (report 118). The registry now wraps every boolean
// answer in options, so this sweep must never find one without them. The
// production conversions have since replaced many boolean surfaces with
// counts, so the floor is a sanity check on the sweep, not a target.
func TestBooleanAnswersAlwaysServeOptions(t *testing.T) {
	reg := generator.NewRegistry()
	all.Register(reg)

	checked, offenders := 0, 0
	for _, id := range reg.Concepts() {
		for i := 0; i < 8; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: float64(i) / 8,
				Seed:       int64(i*7919 + len(id)),
			})
			if err != nil {
				continue
			}
			switch strings.ToLower(strings.TrimSpace(p.Answer)) {
			case "yes", "no", "true", "false":
			default:
				continue
			}
			checked++
			if generator.ChoiceOptions(p.Question) == nil {
				offenders++
				t.Errorf("%s serves the answer %q with no options: %q", id, p.Answer, p.Question)
			}
			if strings.Contains(strings.ToLower(p.Question), "(yes/no)") {
				t.Errorf("%s still asks for a typed yes or no: %q", id, p.Question)
			}
			break
		}
	}
	if checked < 40 {
		t.Fatalf("the sweep checked only %d boolean answers; it is not reading the corpus", checked)
	}
	if offenders == 0 {
		t.Logf("checked %d boolean answers, all serve options", checked)
	}
}

// A multiple-choice concept whose answer is written with a symbol a phone
// keyboard cannot produce needs options for the same reason a word answer does:
// there is nothing to type that the choice grader will match. Report 28 was
// exactly this, tan(60) answered with a radical and no options to tap.
//
// latexOnlyOptions is the backlog this guard found: multiple-choice concepts
// that answer in LaTeX with no option list, so a learner has nothing to tap.
// It is empty now. Add an id here only while a newly found concept is being
// fixed, and remove it when the concept serves options, because a permanent
// entry turns the guard off for the concept it exists to protect.
var latexOnlyOptions = map[string]bool{}

func TestMultipleChoiceSymbolAnswersServeOptions(t *testing.T) {
	dag, err := concepts.LoadDir("../../data/concepts")
	if err != nil {
		t.Fatalf("load concepts: %v", err)
	}
	reg := generator.NewRegistry()
	all.Register(reg)

	checked := 0
	for id, c := range dag.Concepts() {
		if c == nil || c.GradingType != "multiple_choice" || !learnerDomains[c.Domain] {
			continue
		}
		if latexOnlyOptions[id] {
			continue
		}
		gen, err := reg.Get(id)
		if err != nil {
			continue
		}
		for i := 0; i < 60; i++ {
			p := gen.Generate(generator.GeneratorContext{Difficulty: float64(i%5) / 4, Seed: int64(i)})
			if !strings.ContainsAny(p.Answer, "\u221a\u00b0") {
				continue
			}
			checked++
			options := generator.ChoiceOptions(p.Question)
			if len(options) < 2 {
				t.Fatalf("%s: answer %q carries a symbol and the question carries no options:\n%s", id, p.Answer, p.Question)
			}
			found := false
			for _, text := range options {
				if text == p.Answer {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("%s: answer %q is not one of the options", id, p.Answer)
			}
			break
		}
	}
	if checked == 0 {
		t.Skip("no multiple-choice concept serves a symbolic answer now")
	}
}
