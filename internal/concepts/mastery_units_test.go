package concepts

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testDAG(t *testing.T) *DAG {
	t.Helper()
	dag, err := LoadDir(filepath.Join("..", "..", "data", "concepts"))
	if err != nil {
		t.Fatalf("load corpus: %v", err)
	}
	return dag
}

// A partition is a partition: every concept in exactly one unit. This is the property the whole
// feature rests on, and the only one that is cheap to break by hand-editing the JSON.
func TestMasteryUnitDraft_CoversEveryConceptExactlyOnce(t *testing.T) {
	dag := testDAG(t)
	mu, err := BuildMasteryUnitDraft(dag, 3)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if mu.ConceptCount != dag.Count() {
		t.Errorf("ConceptCount = %d, want %d", mu.ConceptCount, dag.Count())
	}
	owner := map[string]string{}
	for _, u := range mu.Units {
		for _, id := range u.ConceptIDs {
			if prev, dup := owner[id]; dup {
				t.Errorf("concept %q in both %q and %q", id, prev, u.ID)
			}
			owner[id] = u.ID
			if dag.Concept(id) == nil {
				t.Errorf("unit %q references unknown concept %q", u.ID, id)
			}
		}
	}
	for _, c := range dag.Order() {
		if _, ok := owner[c.ID]; !ok {
			t.Errorf("concept %q belongs to no unit", c.ID)
		}
	}
	if len(owner) != dag.Count() {
		t.Errorf("partition has %d concepts, corpus has %d", len(owner), dag.Count())
	}
}

// The file on disk must be the partition the rule produces, at the width the file itself
// declares. Regeneration at a different width therefore needs no code change, and a hand-edit
// that quietly moves a concept is caught.
func TestMasteryUnitFile_MatchesItsOwnDeclaredRule(t *testing.T) {
	dag := testDAG(t)
	path := filepath.Join("..", "..", MasteryUnitsPath)
	mu, err := LoadMasteryUnits(path)
	if err != nil {
		t.Fatalf("load %s: %v (run `go run ./cmd/masteryunits`)", path, err)
	}
	if !mu.Draft {
		t.Error("the file must declare itself a draft")
	}
	for _, f := range Errors(ValidateMasteryUnits(dag, mu)) {
		t.Errorf("%s: %s", f.Kind, f.Detail)
	}
}

// Membership must be reproducible from the corpus, which is what makes the band width
// reversible: regenerate, diff, done.
func TestMasteryUnitDraft_IsReproducible(t *testing.T) {
	dag := testDAG(t)
	a, err := BuildMasteryUnitDraft(dag, 3)
	if err != nil {
		t.Fatal(err)
	}
	b, err := BuildMasteryUnitDraft(dag, 3)
	if err != nil {
		t.Fatal(err)
	}
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	if string(ja) != string(jb) {
		t.Error("two builds at the same width differ; the draft is not reproducible")
	}
}

// Different widths must produce different partitions, or "the width is a parameter" is not
// true and the regeneration workflow is an illusion.
func TestMasteryUnitDraft_BandWidthChangesThePartition(t *testing.T) {
	dag := testDAG(t)
	narrow, err := BuildMasteryUnitDraft(dag, 3)
	if err != nil {
		t.Fatal(err)
	}
	wide, err := BuildMasteryUnitDraft(dag, 12)
	if err != nil {
		t.Fatal(err)
	}
	if narrow.UnitCount == wide.UnitCount {
		t.Errorf("width 3 and 12 both give %d units; the width is not a parameter", narrow.UnitCount)
	}
	if wide.BandWidth != 12 {
		t.Errorf("BandWidth = %d, want 12", wide.BandWidth)
	}
}

// Unit prerequisite order is derived from the concept DAG, never stored. If it were stored, the
// file would be a second hierarchy that could contradict the graph — which is what ADR-046
// rejected. This asserts the absence of the field, which is easy to reintroduce by accident.
func TestMasteryUnitFile_StoresNoPrerequisiteGraph(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", MasteryUnitsPath))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	for _, banned := range []string{"prerequisiteUnitIds", "prereq_unit", "prerequisite_units"} {
		if strings.Contains(string(raw), banned) {
			t.Errorf("the draft file contains %q; unit order must be derived, not authored", banned)
		}
	}
}

// The derived unit edges must each be backed by a concept-DAG edge, which is the whole claim.
func TestDerivedUnitPrereqs_AreBackedByConceptEdges(t *testing.T) {
	dag := testDAG(t)
	mu, err := BuildMasteryUnitDraft(dag, 3)
	if err != nil {
		t.Fatal(err)
	}
	edges, err := DerivedUnitPrereqs(mu, dag)
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) == 0 {
		t.Fatal("no unit edges derived; the corpus certainly has some")
	}
	locate := BandIndex(mu, dag)
	via := map[string]bool{}
	for _, c := range dag.Order() {
		cu, _, _, _ := locate(c.ID)
		for _, p := range dag.PrereqsOf(c.ID) {
			pu, _, _, _ := locate(p.ID)
			via[pu+"\x00"+cu] = true
		}
	}
	for _, e := range edges {
		if !via[e.Prereq+"\x00"+e.Unit] {
			t.Errorf("derived edge %s -> %s has no concept-DAG edge behind it", e.Prereq, e.Unit)
		}
	}
}

// Concept depth is a proper topological ranking, so the derived unit order cannot be
// contradicted by the graph. This is the property the first candidate rule lacked, and it is
// worth pinning at every width rather than only the drafted one.
func TestBackwardBandEdges_NoneAtAnyWidth(t *testing.T) {
	dag := testDAG(t)
	for _, b := range []int{1, 2, 3, 4, 6, 12, 45} {
		mu, err := BuildMasteryUnitDraft(dag, b)
		if err != nil {
			t.Fatalf("width %d: %v", b, err)
		}
		if back := BackwardBandEdges(mu, dag); len(back) > 0 {
			t.Errorf("width %d: %d backward band edge(s), e.g. %q needs %q",
				b, len(back), back[0].ConceptID, back[0].PrereqID)
		}
	}
}

// The reason depth is measured per concept. The subdomain condensation of an acyclic concept
// DAG is NOT necessarily acyclic, so a subdomain depth function does not exist for every
// subdomain — and code that assumes it does gets level 0 rather than an error, which is the
// failure mode this test exists to document.
//
// This corpus has 9 such cycles covering 26 of 125 subdomains.
func TestSubdomainDepths_ReportsCyclesRatherThanDefaultingToZero(t *testing.T) {
	dag := testDAG(t)
	depths, _, unresolved := SubdomainDepths(dag)
	if len(unresolved) == 0 {
		t.Skip("no cyclic subdomains; the corpus changed, so re-check whether concept depth is still required")
	}
	for _, s := range unresolved {
		if _, ok := depths[s]; ok {
			t.Errorf("subdomain %q is reported unresolved but has a depth", s)
		}
	}
	// Sanity: the resolved ones are real levels, not placeholders.
	if len(depths) == 0 || len(depths) == len(Subdomains(dag)) {
		t.Errorf("resolved %d of %d subdomains; expected a partial resolution", len(depths), len(Subdomains(dag)))
	}
}

// Concept depth must resolve every concept, which is the reason it was chosen.
func TestConceptDepths_ResolvesEveryConceptAndOrdersPrerequisites(t *testing.T) {
	dag := testDAG(t)
	depth := ConceptDepths(dag)
	if len(depth) != dag.Count() {
		t.Fatalf("resolved %d of %d concepts", len(depth), dag.Count())
	}
	for _, c := range dag.Order() {
		for _, p := range dag.PrereqsOf(c.ID) {
			if depth[p.ID] >= depth[c.ID] {
				t.Errorf("%q (depth %d) has prerequisite %q at depth %d; a prerequisite must be shallower",
					c.ID, depth[c.ID], p.ID, depth[p.ID])
			}
		}
	}
}

// The report is how D3's finding gets read, so the numbers that decide the width must be in it
// and must match reality. A report that quietly describes a different partition than the JSON
// is worse than no report.
func TestMasteryUnitReport_DescribesTheFileBesideIt(t *testing.T) {
	root := filepath.Join("..", "..")
	mu, err := LoadMasteryUnits(filepath.Join(root, MasteryUnitsPath))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(root, "data", "mastery_units.report.md"))
	if err != nil {
		t.Fatalf("read report: %v", err)
	}
	rep := string(raw)
	for _, want := range []string{
		"| Total units | " + itoa(mu.UnitCount) + " |",
		"| Total concepts | " + itoa(mu.ConceptCount) + " |",
		"| Backward band edges | 0 |",
		"Band width: **" + itoa(mu.BandWidth) + "**",
	} {
		if !strings.Contains(rep, want) {
			t.Errorf("report is missing %q; it must describe the file beside it", want)
		}
	}
	// Every unit must be listed, or the report is a summary of a partition nobody can see.
	for _, u := range mu.Units {
		if !strings.Contains(rep, "`"+u.ID+"`") {
			t.Errorf("report does not list unit %q", u.ID)
		}
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		return "-" + string(b)
	}
	return string(b)
}
