package all

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The variety ratchet.
//
// A generator that returns a hardcoded table of three problems is not a
// generator. A learner who answers four questions from such a concept has seen
// the whole supply, and the fifth question is a repeat. Measured over the whole
// corpus in October 2026, 36 concepts produced at most 4 distinct stems in 200
// draws, 274 produced at most 9, and 527 produced at most 39.
//
// Rewriting several hundred generators in one change is neither reviewable nor
// safe, so this test does the part that keeps working: it records the number of
// distinct stems each concept produces today and fails if any concept falls
// below that number. Converting a concept is therefore a permanent gain rather
// than a gain that the next pull request undoes.
//
// The baseline is measured rather than authored. Regenerate it with:
//
//	go test ./internal/generator/all -run TestConceptVarietyFloor -update-variety
//
// and raise a concept's entry by hand only when its generator has genuinely been
// parametrised. Never lower an entry to make the test pass: a concept that lost
// variety has lost it for learners too.
//
// The counted unit is one mathematical situation: the core stem, the question
// before the "A) …" option block that generator.MultipleChoice appends, together
// with the option texts in sorted order. That key took two corrections to reach,
// and both were found by inspecting a generator rather than by trusting the
// number. Counting the rendered text made option shuffling look like variety: on
// 2026-10-24 the corpus measured 144 concepts up to 150 "stems" whose real
// supply was four problems, and no concept exceeded 99 distinct situations while
// the rendered figures claimed buckets above 200. Stripping the option block
// entirely then hid the opposite defect, because a concept whose question is
// fixed and whose options vary poses one problem per option set:
// complex.topology.domain reported a single problem for fifteen. Sorting the
// options keeps the first fix without the second error. The baseline was
// re-measured from scratch after each correction, by deleting the file and
// re-running -update-variety, because a metric change is the only occasion on
// which a lower number is correct.
//
// A metric aggregated over 1,250 generators will be wrong for some of them, so
// spot-check any figure against the generator it describes before building a
// plan on it. Both corrections above came from doing that and finding the number
// did not match the questions.
const (
	varietyDraws        = 60
	varietyRuns         = 5
	varietyBaselinePath = "variety_baseline.json"
)

var updateVariety = flag.Bool(
	"update-variety",
	false,
	"rewrite the concept variety baseline from a fresh measurement",
)

// corpusConceptIDs lists every concept id in the shipped corpus.
func corpusConceptIDs(t *testing.T) []string {
	t.Helper()
	files, err := filepath.Glob("../../../data/concepts/*.json")
	if err != nil || len(files) == 0 {
		t.Fatalf("no concept files found: %v", err)
	}
	seen := make(map[string]bool)
	var ids []string
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		var list []struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(data, &list); err != nil {
			continue
		}
		for _, c := range list {
			if c.ID != "" && !seen[c.ID] {
				seen[c.ID] = true
				ids = append(ids, c.ID)
			}
		}
	}
	sort.Strings(ids)
	return ids
}

// canonicalProblem reduces a rendered question to the problem it actually
// poses.
//
// generator.MultipleChoice writes its shuffled options into the stem as an
// "A) …" block, so a concept with four options reports up to twenty-four
// distinct question strings for one mathematical situation. Counting those was
// the ratchet's first mistake: measured over the corpus it reported concepts
// above a hundred stems whose real supply was four.
//
// Stripping the options entirely is the second mistake, and it hides the
// opposite defect. complex.topology.domain asks one question, "which of these
// sets is a domain", and varies only which sets are offered, so stripping the
// block reports one problem for a concept that poses four. A permutation is not
// a new problem; a different option set is. The counted unit is therefore the
// core stem plus the option texts in sorted order, which is one key per
// mathematical situation however the letters happen to fall.
// coreStem is the question without its option block.
func coreStem(question string) string {
	if i := strings.Index(question, "\n\nA) "); i >= 0 {
		return question[:i]
	}
	return question
}

func canonicalProblem(question string) string {
	core := coreStem(question)
	options := generator.ChoiceOptions(question)
	if len(options) == 0 {
		return core
	}
	texts := make([]string, 0, len(options))
	for _, text := range options {
		texts = append(texts, text)
	}
	sort.Strings(texts)
	return core + "\n" + strings.Join(texts, "\n")
}

// distinctStems draws from one concept and counts the distinct problems it
// poses, as defined by canonicalProblem.
//
// Generators reach for the global rand source rather than the seed in the
// context, so a single measurement is a sample. That is why the baseline stores
// the minimum of several runs, which is a lower bound on the normal outcome.
func distinctStems(reg *generator.Registry, id string, draws int) int {
	seen := make(map[string]bool)
	for i := 0; i < draws; i++ {
		p, err := reg.GenerateContext(id, generator.GeneratorContext{
			Difficulty: 0.2 + 0.8*rand.Float64(),
		})
		if err != nil {
			continue
		}
		seen[canonicalProblem(p.Question)] = true
	}
	return len(seen)
}

func TestConceptVarietyFloor(t *testing.T) {
	ids := corpusConceptIDs(t)
	reg := Registry()

	if *updateVariety {
		previous := make(map[string]int)
		if raw, err := os.ReadFile(varietyBaselinePath); err == nil {
			_ = json.Unmarshal(raw, &previous)
		}
		baseline := make(map[string]int, len(ids))
		kept := 0
		for _, id := range ids {
			lowest := varietyDraws + 1
			for run := 0; run < varietyRuns; run++ {
				if n := distinctStems(reg, id, varietyDraws); n < lowest {
					lowest = n
				}
			}
			// A ratchet only ratchets up. A re-measure below the recorded figure
			// on an untouched generator is sampling noise, and accepting it would
			// quietly lower the floor this file exists to hold.
			if old, ok := previous[id]; ok && old > lowest {
				lowest = old
				kept++
			}
			baseline[id] = lowest
		}
		out, err := json.MarshalIndent(baseline, "", "  ")
		if err != nil {
			t.Fatalf("marshal baseline: %v", err)
		}
		if err := os.WriteFile(varietyBaselinePath, append(out, '\n'), 0o644); err != nil {
			t.Fatalf("write baseline: %v", err)
		}
		t.Logf("wrote %s for %d concepts (%d kept above the re-measured value)", varietyBaselinePath, len(baseline), kept)
		return
	}

	raw, err := os.ReadFile(varietyBaselinePath)
	if err != nil {
		t.Fatalf("read %s: %v", varietyBaselinePath, err)
	}
	baseline := make(map[string]int)
	if err := json.Unmarshal(raw, &baseline); err != nil {
		t.Fatalf("parse %s: %v", varietyBaselinePath, err)
	}

	present := make(map[string]bool, len(ids))
	for _, id := range ids {
		present[id] = true
	}

	var offenders []string
	for _, id := range ids {
		want, ok := baseline[id]
		if !ok {
			offenders = append(offenders, id+": no baseline entry, run -update-variety once and commit the result")
			continue
		}
		// The comparison is deliberately loose. Generators draw from the global
		// random source, so a concept whose support is at least as large as the
		// draw count saturates at the draw count and varies by a handful from
		// run to run. The number this test defends against is a collapse, a
		// parametrised generator replaced by a table of three or four entries,
		// so it asks for three fifths of the recorded floor rather than all of
		// it. A tighter bound would fail on sampling alone and would be deleted
		// by whoever it failed for.
		required := want * 3 / 5
		if required < 2 {
			required = 2
		}
		// A concept recorded at one distinct situation cannot be asked for two
		// without failing on the measurement that produced its own baseline.
		// The two concepts at that floor are real defects, tracked for authored
		// shapes, and the ratchet cannot repair them by demanding more.
		if required > want {
			required = want
		}
		got := distinctStems(reg, id, varietyDraws)
		if got < required {
			offenders = append(offenders, fmt.Sprintf("%s: %d distinct stems in %d draws, baseline %d", id, got, varietyDraws, want))
		}
	}
	for id := range baseline {
		if !present[id] {
			offenders = append(offenders, id+": baseline entry for a concept that is no longer in the corpus")
		}
	}

	if len(offenders) > 0 {
		sort.Strings(offenders)
		t.Fatalf("%d concepts fell below their recorded variety floor:\n%s",
			len(offenders), strings.Join(offenders, "\n"))
	}

	low := 0
	for _, id := range ids {
		if baseline[id] <= 4 {
			low++
		}
	}
	t.Logf("%d concepts hold their variety floor, %d of them still produce at most 4 distinct stems",
		len(ids), low)
}
