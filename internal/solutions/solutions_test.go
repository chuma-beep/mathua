package solutions

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func load(t *testing.T, files map[string]string) *Loader {
	t.Helper()
	dir := t.TempDir()
	for name, body := range files {
		write(t, dir, name, body)
	}
	l, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return l
}

func facts() map[string]string {
	return map[string]string{"a_num": "3", "a_den": "4", "b_num": "1", "b_den": "2", "lcm": "4", "b_scaled": "2", "num": "5", "answer": "5/4"}
}

const fracSchema = `{
  "concept": "frac.add.diff",
  "setup": "You are adding $[[a_num]]/[[a_den]]$ and $[[b_num]]/[[b_den]]$.",
  "steps": [
    {"fact": "lcm", "say": "Put both over LCM([[a_den]], [[b_den]]) = [[lcm]]."},
    {"fact": "b_scaled", "say": "Then $[[b_num]]/[[b_den]]$ becomes $[[b_scaled]]/[[lcm]]$."}
  ],
  "answer": "Adding the tops gives [[num]]/[[lcm]]."
}`

func TestAssemble_InterpolatesEveryFact(t *testing.T) {
	l := load(t, map[string]string{"frac.add.diff.json": fracSchema})
	got := l.Assemble("frac.add.diff", facts(), "FALLBACK")
	if got == "FALLBACK" {
		t.Fatal("schema was skipped; it should have resolved")
	}
	for _, want := range []string{
		"You are adding $3/4$ and $1/2$.",
		"Put both over LCM(4, 2) = 4.",
		"Then $1/2$ becomes $2/4$.",
		"Adding the tops gives 5/4.",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	// No placeholder may survive into what a learner reads.
	if strings.Contains(got, "[[") {
		t.Errorf("unsubstituted placeholder left in output:\n%s", got)
	}
}

func TestAssemble_SkipsRatherThanGuessing(t *testing.T) {
	cases := map[string]struct {
		files   map[string]string
		concept string
		f       map[string]string
		why     string
	}{
		"no schema for the concept": {
			files:   map[string]string{"other.json": strings.Replace(fracSchema, `"concept": "frac.add.diff"`, `"concept": "frac.sub.diff"`, 1)},
			concept: "frac.add.diff",
			f:       facts(),
			why:     "a concept without a schema must keep the generator's explanation",
		},
		"placeholder has no fact": {
			files:   map[string]string{"frac.add.diff.json": fracSchema},
			concept: "frac.add.diff",
			f:       map[string]string{"a_num": "3", "a_den": "4", "b_num": "1", "b_den": "2", "lcm": "4", "b_scaled": "2"},
			why:     "a missing num must not render an empty result",
		},
		"step fact was renamed in the generator": {
			files:   map[string]string{"frac.add.diff.json": fracSchema},
			concept: "frac.add.diff",
			f:       map[string]string{"a_num": "3", "a_den": "4", "b_num": "1", "b_den": "2", "num": "5", "answer": "5/4"},
			why:     "a step whose fact is gone has lost its footing",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			l := load(t, tc.files)
			if got := l.Assemble(tc.concept, tc.f, "FALLBACK"); got != "FALLBACK" {
				t.Errorf("expected the fallback, got:\n%s\n(%s)", got, tc.why)
			}
		})
	}
}

// One generator often asks two different questions — find the complement of an
// angle, or classify it — and a single schema covers both. Each question's
// steps must apply only to the question it came from, without the other
// question's steps leaking in or the schema being dropped.
func TestRender_SkipsStepsWhoseFactIsAbsent(t *testing.T) {
	s := &Schema{
		Setup: "An angle of [[given]] degrees.",
		Steps: []Step{
			{Fact: "missing", Say: "A right angle is 90, so the other is 90 - [[given]] = [[missing]]."},
			{Fact: "cls", Say: "[[given]] degrees falls in the range for [[cls]]."},
		},
		Answer: "So the answer is [[answer]].",
	}
	complement, missing := Render(s, map[string]string{"given": "40", "missing": "50", "answer": "50"})
	if len(missing) > 0 {
		t.Fatalf("unexpected missing: %v", missing)
	}
	if strings.Contains(complement, "falls in the range") {
		t.Errorf("the classification step leaked into the complement question:\n%s", complement)
	}
	classify, missing := Render(s, map[string]string{"given": "120", "cls": "obtuse", "answer": "obtuse"})
	if len(missing) > 0 {
		t.Fatalf("unexpected missing: %v", missing)
	}
	if strings.Contains(classify, "A right angle is 90") {
		t.Errorf("the complement step leaked into the classification question:\n%s", classify)
	}
	// The drift is still reported to the test even though the runtime copes.
	got := Missing(s, map[string]string{"given": "40", "answer": "50"})
	if len(got) != 2 || got[0] != "missing" || got[1] != "cls" {
		t.Errorf("Missing should report both absent step facts, got %v", got)
	}
}

func TestAssemble_NilLoaderIsTheFallback(t *testing.T) {
	var l *Loader
	if got := l.Assemble("anything", facts(), "FALLBACK"); got != "FALLBACK" {
		t.Errorf("a nil loader must fall back, got %q", got)
	}
}

// LaTeX braces are everywhere in this corpus, so they cannot be the
// placeholder delimiter. A schema full of \frac must survive untouched.
func TestRender_LeavesLatexBracesAlone(t *testing.T) {
	s := &Schema{
		Steps:  []Step{{Fact: "m", Say: `The matrix $\begin{pmatrix} 1 \\ 0 \end{pmatrix}$ has [[m]] rows.`}},
		Answer: "So [[answer]].",
	}
	got, missing := Render(s, map[string]string{"m": "2", "answer": "done"})
	if len(missing) > 0 {
		t.Fatalf("unexpected missing: %v", missing)
	}
	if !strings.Contains(got, `\begin{pmatrix} 1 \\ 0 \end{pmatrix}`) {
		t.Errorf("LaTeX was mangled:\n%s", got)
	}
}

func TestLoad_RejectsUnusableSchemas(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "good.json", fracSchema)
	write(t, dir, "no-answer.json", `{"concept":"x","setup":"hi"}`)
	write(t, dir, "broken.json", `{"concept":`)
	l, err := Load(dir)
	if err == nil {
		t.Fatal("expected an error naming the unusable schemas")
	}
	for _, want := range []string{"no-answer.json", "broken.json"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error should name %s: %v", want, err)
		}
	}
	// A bad file must not take the good ones down with it.
	if _, ok := l.Schema("frac.add.diff"); !ok {
		t.Error("a usable schema should still load alongside a broken one")
	}
}

func TestLoad_MissingDirectoryIsEmpty(t *testing.T) {
	l, err := Load(filepath.Join(t.TempDir(), "nope"))
	if err != nil {
		t.Fatalf("a missing directory must not be an error: %v", err)
	}
	if len(l.Concepts()) != 0 {
		t.Errorf("expected no concepts, got %v", l.Concepts())
	}
}

// The row-echelon form a systems question is graded against is written
// [[a,b,c],[0,d,e]]. An unrestricted [[...]] placeholder would eat that, so
// contents must look like a fact name to be treated as one.
func TestRender_LeavesMatrixNotationAlone(t *testing.T) {
	s := &Schema{
		Steps:  []Step{{Fact: "ans", Say: "The echelon form is [[2,1,7],[0,3,5]]."}},
		Answer: "So the answer is [[answer]].",
	}
	got, missing := Render(s, map[string]string{"ans": "ok", "answer": "[[4,3,6],[0,1,1]]"})
	if len(missing) > 0 {
		t.Fatalf("unexpected missing: %v", missing)
	}
	if !strings.Contains(got, "[[2,1,7],[0,3,5]]") {
		t.Errorf("matrix notation in a step was eaten:\n%s", got)
	}
	if !strings.Contains(got, "[[4,3,6],[0,1,1]]") {
		t.Errorf("matrix notation in the answer was eaten:\n%s", got)
	}
	// A capitalised name is not a fact name either, so it survives verbatim.
	// The leak test is what catches that typo, not the assembler — reporting it
	// here would mean reporting every matrix span as missing.
	got, missing = Render(&Schema{Answer: "[[Num]]"}, map[string]string{"Num": "7"})
	if len(missing) > 0 {
		t.Errorf("unexpected missing: %v", missing)
	}
	if !strings.Contains(got, "[[Num]]") {
		t.Errorf("a capitalised name should be left literal, got %q", got)
	}
}
