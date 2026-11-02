package all

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The shipping contract for a concept is three files and one generator:
//
//   - a lesson mapping whose source file exists,
//   - exactly three knowledge-point shards,
//   - a generator that returns a non-empty problem for every ordinal 1..3.
//
// Existing concepts are grandfathered by section_alignment_baseline.json. A
// concept not in that file must also produce a different question at each
// ordinal under the same seed. This is the guard for the class of defect that
// let prealg.eq.one_step_add ask an integer question on the "Equations with
// Fractions" section: the files existed, but nothing compared the section to
// the question.
func TestConceptShippingContract(t *testing.T) {
	concepts := loadConceptIDs(t)
	lessons := loadLessonSources(t)
	kps := loadKPShards(t)
	baseline := loadAlignmentBaseline(t)

	reg := generator.NewRegistry()
	Register(reg)

	var missingLesson, missingKP, wrongKPCount, missingGenerator, emptyProblem, ordinalMismatch []string
	var unaligned []string
	for _, cid := range concepts {
		if len(lessons[cid]) == 0 {
			missingLesson = append(missingLesson, cid)
		}
		shards, ok := kps[cid]
		if !ok {
			missingKP = append(missingKP, cid)
		} else if len(shards) != 3 {
			wrongKPCount = append(wrongKPCount, fmt.Sprintf("%s (%d)", cid, len(shards)))
		}
		if !reg.Has(cid) {
			missingGenerator = append(missingGenerator, cid)
			continue
		}
		first := make([]string, 0, 3)
		for ordinal := 1; ordinal <= 3; ordinal++ {
			problems, err := reg.BatchGenerateContext(cid, 3, generator.GeneratorContext{
				Difficulty: 0.5,
				KPOrdinal:  ordinal,
			})
			if err != nil || len(problems) == 0 {
				emptyProblem = append(emptyProblem, fmt.Sprintf("%s kp%d: %v", cid, ordinal, err))
				continue
			}
			for _, p := range problems {
				if p.Question == "" || p.Answer == "" {
					emptyProblem = append(emptyProblem, fmt.Sprintf("%s kp%d: blank question or answer", cid, ordinal))
					break
				}
				if p.KPOrdinal != ordinal {
					ordinalMismatch = append(ordinalMismatch, fmt.Sprintf("%s: got kp%d for kp%d", cid, p.KPOrdinal, ordinal))
					break
				}
			}
			one, err := reg.GenerateContext(cid, generator.GeneratorContext{
				Difficulty: 0.5,
				Seed:       4242,
				KPOrdinal:  ordinal,
			})
			if err == nil {
				first = append(first, one.Question)
			}
		}
		if baseline[cid] {
			continue
		}
		if len(first) == 3 {
			if first[0] == first[1] || first[1] == first[2] || first[0] == first[2] {
				unaligned = append(unaligned, fmt.Sprintf("%s: same seeded question at two ordinals", cid))
			}
		}
		for ordinal := 1; ordinal <= 3; ordinal++ {
			problems, _ := reg.BatchGenerateContext(cid, 12, generator.GeneratorContext{
				Difficulty: 0.5,
				Seed:       int64(5000 + ordinal),
				KPOrdinal:  ordinal,
			})
			seen := map[string]bool{}
			for _, p := range problems {
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				unaligned = append(unaligned, fmt.Sprintf("%s kp%d: only %d distinct questions in 12 draws", cid, ordinal, len(seen)))
			}
		}
	}

	sort.Strings(missingLesson)
	sort.Strings(missingKP)
	sort.Strings(wrongKPCount)
	sort.Strings(missingGenerator)
	sort.Strings(emptyProblem)
	sort.Strings(ordinalMismatch)
	sort.Strings(unaligned)
	for _, group := range []struct {
		name  string
		items []string
	}{
		{"missing lesson mapping", missingLesson},
		{"missing KP shard", missingKP},
		{"wrong KP shard count", wrongKPCount},
		{"missing generator", missingGenerator},
		{"empty generator output", emptyProblem},
		{"ordinal not recorded", ordinalMismatch},
		{"new concept is not section-aligned", unaligned},
	} {
		for _, item := range group.items {
			t.Errorf("%s: %s", group.name, item)
		}
	}
	if len(concepts) < 1000 {
		t.Fatalf("only %d concepts loaded; the contract would pass vacuously", len(concepts))
	}
}

func loadConceptIDs(t *testing.T) []string {
	t.Helper()
	files, err := filepath.Glob("../../../data/concepts/*.json")
	if err != nil || len(files) == 0 {
		t.Fatalf("no concept files found: %v", err)
	}
	var ids []string
	for _, f := range files {
		if filepath.Base(f) == "enrichment.json" {
			continue
		}
		var batch []struct {
			ID string `json:"id"`
		}
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		if err := json.Unmarshal(b, &batch); err != nil {
			t.Fatalf("parse %s: %v", f, err)
		}
		for _, c := range batch {
			if c.ID != "" {
				ids = append(ids, c.ID)
			}
		}
	}
	sort.Strings(ids)
	return ids
}

func loadLessonSources(t *testing.T) map[string][]string {
	t.Helper()
	b, err := os.ReadFile("../../../data/lessons/lessons.json")
	if err != nil {
		t.Fatalf("read lessons.json: %v", err)
	}
	var rows []struct {
		ConceptID string `json:"concept_id"`
		Source    string `json:"source"`
	}
	if err := json.Unmarshal(b, &rows); err != nil {
		t.Fatalf("parse lessons.json: %v", err)
	}
	out := make(map[string][]string)
	for _, row := range rows {
		path := filepath.Join("../../../data/lessons", row.Source)
		if _, err := os.Stat(path); err != nil {
			continue
		}
		out[row.ConceptID] = append(out[row.ConceptID], row.Source)
	}
	return out
}

func loadKPShards(t *testing.T) map[string][]json.RawMessage {
	t.Helper()
	files, err := filepath.Glob("../../../data/lessons/kp/*.json")
	if err != nil {
		t.Fatalf("glob kp shards: %v", err)
	}
	out := make(map[string][]json.RawMessage, len(files))
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		var shards []json.RawMessage
		if err := json.Unmarshal(b, &shards); err != nil {
			t.Fatalf("parse %s: %v", f, err)
		}
		out[filepath.Base(f[:len(f)-len(".json")])] = shards
	}
	return out
}

func loadAlignmentBaseline(t *testing.T) map[string]bool {
	t.Helper()
	b, err := os.ReadFile("section_alignment_baseline.json")
	if err != nil {
		t.Fatalf("read section alignment baseline: %v", err)
	}
	var ids []string
	if err := json.Unmarshal(b, &ids); err != nil {
		t.Fatalf("parse section alignment baseline: %v", err)
	}
	out := make(map[string]bool, len(ids))
	for _, id := range ids {
		out[id] = true
	}
	return out
}
