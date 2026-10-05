package all

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// How many distinct questions can a generator actually supply?
//
// This is the measurement that was missing when a learner hit a dead end four questions into
// `geo.basic.points_lines`. That generator has exactly four distinct question texts, so once a
// session had seen all four, `/api/lessons/{id}/practice` returned `questions: []` with a
// 200 — and the Learn feed read an empty list as success, stopped appending, and left the last
// verdict as the end of the page with no next step and no message.
//
// Nothing in the suite could see it. `TestEasyPoolsDiscriminate` checks that the *answers*
// discriminate, which `geo.basic.points_lines` passes with four distinct answers; question
// *texts* were never counted anywhere.
//
// **The counts below are lower bounds, not exact banks.** Generators draw from the global
// `math/rand`, which Go 1.24+ seeds automatically and `rand.Seed` can no longer override, so
// repeated runs differ in the high digits — seeding was attempted and did not help. That is
// fine for a guard, and the reason is directional: sampling can only *miss* variants, never
// invent them, so a sample never over-reports. A concept measured at N has at least N
// distinct questions. The property being gated ("cannot sustain a section") is therefore
// decided the same way every run; only the large banks wobble, and they are not gated.
var difficultiesSwept = []float64{0.3, 0.5, 0.8, 1.0}

const (
	// sampleCount is how many times each generator is called per difficulty. Small banks are
	// exhaustive well before this; large ones are merely lower-bounded.
	sampleCount = 600
	// difficultiesSwept covers the range the Learn feed moves through, since the feed raises
	// difficulty on a correct answer and lowers it on a miss. A generator whose variety
	// depends on difficulty can be tiny at one setting and large at another, and the learner
	// meets the small one.
	difficultiesSweptStr = "0.3, 0.5, 0.8, 1.0"

	// minDistinctQuestions is the floor below which a concept cannot even fill one practice
	// request (the feed asks for 3 at a time). Below this, a learner sees the same question
	// twice within the first two cards, which is a correctness-adjacent problem rather than
	// a variety one.
	//
	// Deliberately not higher. The measured distribution is what it is — see the log below —
	// and a threshold set above the corpus median would be a gate 378 of 657 concepts cannot
	// pass, which is a content backlog wearing a test's clothes. The variety gap is reported
	// instead; the dead end is fixed independently, in the practice endpoint and the Learn
	// feed, so nothing here is load-bearing for correctness.
	minDistinctQuestions = 3

	// varietyReference is the count reported as a quality line rather than gated: below it, a
	// session long enough to exhaust `exclude` starts repeating. Reported, not enforced.
	varietyReference = 12
)

// lowBankAllowlist is the measured set of concepts that cannot even fill one practice request,
// with the reason each is accepted.
//
// Every entry is a content gap, not a code gap: the generator is correct, there is simply not
// enough written. They are here so the set is visible and cannot grow silently, and so
// retiring one is a one-line deletion rather than an archaeology exercise.
var lowBankAllowlist = map[string]string{
	// All three scale an integer range by difficulty: `scale := int(1 + ctx.Difficulty*5)`,
	// so at difficulty 0.3 the range is 2 wide and `rand.Intn(scale)+1` yields two distinct
	// questions. The feed opens a concept at difficulty 0.4, where the range is 3, so a
	// learner meets two variants before any repeats — degrades to a repeat, not a dead end,
	// now that the practice endpoint falls back rather than returning empty.
	"calc.integral.power_rule":   "difficulty-scaled integer range is 2 wide at difficulty 0.3; widen the base range",
	"calc.integral.substitution": "difficulty-scaled integer range is 2 wide at difficulty 0.3; widen the base range",
	"calc.integral.ftc":          "difficulty-scaled integer range is 2 wide at difficulty 0.3; widen the base range",
}

// distinctQuestionLowerBound samples one generator and returns the smallest distinct-question
// count seen across the difficulty sweep.
//
// The minimum, not the union: the feed generates at one difficulty at a time, so the pool the
// learner draws from is the one for the difficulty they are currently at.
func distinctQuestionLowerBound(reg *generator.Registry, conceptID string) int {
	gen, err := reg.Get(conceptID)
	if err != nil {
		return 0
	}
	best := -1
	for _, difficulty := range difficultiesSwept {
		seen := map[string]bool{}
		for i := 0; i < sampleCount; i++ {
			seed := int64(i)*2654435761 + int64(difficulty*1000)
			// Mirrors Registry.GenerateContext. Not required to be effective (see the file
			// comment on Go 1.24+), but it keeps this sampler honest if that ever changes.
			rand.Seed(seed)
			p := gen.Generate(generator.GeneratorContext{Difficulty: difficulty, Seed: seed})
			seen[p.Question] = true
		}
		if n := len(seen); best < 0 || n < best {
			best = n
		}
	}
	return best
}

func TestGeneratorsCanFillAPracticeRequest(t *testing.T) {
	reg := Registry()
	ids := reg.Concepts()
	sort.Strings(ids)

	histogram := map[int]int{}
	var belowVariety, failing []string

	for _, id := range ids {
		n := distinctQuestionLowerBound(reg, id)
		histogram[n]++
		if n < varietyReference {
			belowVariety = append(belowVariety, id)
		}
		if n < minDistinctQuestions {
			failing = append(failing, fmt.Sprintf("%s (%d distinct)", id, n))
		}
	}

	// The shape is the finding. Print it compactly — the raw bucket list is a hundred long.
	var buckets []string
	for _, n := range sortedKeys(histogram) {
		buckets = append(buckets, fmt.Sprintf("%d:%d", n, histogram[n]))
	}
	t.Logf("%d generators; sampled distinct question texts (lower bound, %d samples x %v)",
		len(ids), sampleCount, difficultiesSweptStr)
	t.Logf("  distribution n:concepts — %s", strings.Join(buckets, " "))
	t.Logf("  %d of %d concepts (%.0f%%) fall under the %d-question variety reference",
		len(belowVariety), len(ids), float64(len(belowVariety))/float64(len(ids))*100, varietyReference)

	if len(failing) == 0 {
		return
	}
	for _, f := range failing {
		if _, ok := lowBankAllowlist[strings.Split(f, " ")[0]]; ok {
			t.Logf("  known low bank: %s", f)
			continue
		}
		t.Errorf("generator cannot fill one %d-question practice request: %s\n"+
			"Author more variants, or add the concept to lowBankAllowlist with a reason — "+
			"which is a decision, not a silencing.", minDistinctQuestions, f)
	}
}

func sortedKeys(m map[int]int) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Ints(out)
	return out
}
