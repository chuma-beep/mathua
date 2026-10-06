package concepts

// Mastery units: a *draft* partition of the existing concept DAG, for ADR-046's D3.
//
// What this file is for: making a proposed assessment boundary concrete enough to review.
// It is deliberately not the final design, and two things about it are parameters rather than
// policy:
//
//   - The band width lives in the JSON, not in this package. There is no constant here, because
//     a constant would read as a settled decision, and the width is the open question (how often
//     should a learner be assessed?). The validator uses the width the file declares, so
//     regenerating at 2 or 4 needs no code change and no constant to remember.
//   - Unit prerequisite relationships are **derived from the concept DAG and never stored**.
//     ADR-046 rejected an authored `prerequisiteUnitIds` because two authored hierarchies can
//     contradict each other, and `validate_graph` only knows about concepts. So the file is a
//     partition and nothing else; the ordering falls out of the graph.
//
// The partitioning rule under evaluation is `unit = (domain, floor(depth / bandWidth))`. The
// depth function is per-*concept* (ConceptDepths), chosen over per-subdomain after measuring
// that the subdomain condensation has cycles — see SubdomainDepths, and D3 in
// docs/mastery-assessments.md.

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// MasteryUnitsPath is the draft partition's location, relative to the repo root.
const MasteryUnitsPath = "data/mastery_units.json"

// MasteryUnit is one proposed assessment boundary.
type MasteryUnit struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Domain       string   `json:"domain"`
	Band         int      `json:"band"`
	DepthMin     int      `json:"depth_min"`
	DepthMax     int      `json:"depth_max"`
	Subdomains   []string `json:"subdomains"`
	ConceptIDs   []string `json:"concept_ids"`
	ConceptCount int      `json:"concept_count"`
}

// MasteryUnits is the whole draft file.
//
// BandWidth is the draft parameter and is recorded here rather than in code, so that changing it
// is a regeneration (`go run ./cmd/masteryunits -band N`) and not a code change. The validator
// reads it from here too, which is what lets a regenerated file be checked against the rule it
// was generated from without either side hardcoding a number.
type MasteryUnits struct {
	Draft        bool          `json:"draft"`
	GeneratedBy  string        `json:"generated_by"`
	Rule         string        `json:"rule"`
	BandWidth    int           `json:"band_width"`
	Note         string        `json:"note"`
	Source       string        `json:"source"`
	ConceptCount int           `json:"concept_count"`
	UnitCount    int           `json:"unit_count"`
	Units        []MasteryUnit `json:"units"`
}

// DepthLevel is a concept's or subdomain's position in the prerequisite layering.
type DepthLevel int

// SubdomainDepths returns the topological level of every subdomain, computed on the subdomain
// condensation of the concept DAG: a subdomain is one level deeper than its deepest predecessor
// subdomain. Kahn's algorithm over the condensation, with sorted tie-breaking so the result is
// deterministic.
//
// It returns the levels, the condensation's topological order, and the subdomains it could not
// resolve. Those unresolved subdomains are in a cycle of the condensation.
//
// The condensation is **not** guaranteed to be acyclic even though the concept DAG is: grouping
// concepts by subdomain discards the per-concept ordering, so a concept inside subdomain A can
// be a prerequisite of one inside subdomain B whose other members depend on A. This corpus has
// 9 such cycles covering 26 of 125 subdomains. A caller that ignores the third return value gets
// level 0 for all of them, which is a confidently wrong answer rather than a missing one.
func SubdomainDepths(dag *DAG) (depths map[string]int, order []string, unresolved []string) {
	subs := Subdomains(dag)
	succ := make(map[string]map[string]bool, len(subs))
	indeg := make(map[string]int, len(subs))
	for _, s := range subs {
		succ[s] = map[string]bool{}
	}
	for _, c := range dag.Order() {
		s := c.Subdomain
		for _, p := range dag.PrereqsOf(c.ID) {
			ps := p.Subdomain
			if ps == s || succ[ps][s] {
				continue
			}
			succ[ps][s] = true
			indeg[s]++
		}
	}

	depth := make(map[string]int, len(subs))
	var topo []string
	ready := make([]string, 0, len(subs))
	for _, s := range subs {
		if indeg[s] == 0 {
			ready = append(ready, s)
		}
	}
	sort.Strings(ready)
	for len(ready) > 0 {
		s := ready[0]
		ready = ready[1:]
		topo = append(topo, s)
		for _, t := range sortedKeys(succ[s]) {
			if depth[s]+1 > depth[t] {
				depth[t] = depth[s] + 1
			}
			indeg[t]--
			if indeg[t] == 0 {
				ready = insertSorted(ready, t)
			}
		}
	}
	for _, s := range subs {
		if _, ok := depth[s]; !ok {
			unresolved = append(unresolved, s)
		}
	}
	return depth, topo, unresolved
}

// ConceptDepths returns each concept's topological level by a single pass over the DAG's
// existing topological order: a concept is one deeper than its deepest prerequisite.
//
// This is included because it is the *other* candidate depth function and the difference
// between the two explains every backward band edge. Concept depth is a proper ranking — by
// construction every prerequisite is strictly shallower — so banding on it produces zero
// backward edges by construction, which is exactly why the 22 that subdomain banding produces
// are a measurement artefact rather than a property of the corpus.
func ConceptDepths(dag *DAG) map[string]int {
	depth := make(map[string]int, dag.Count())
	for _, c := range dag.Order() {
		d := 0
		for _, p := range dag.PrereqsOf(c.ID) {
			if depth[p.ID]+1 > d {
				d = depth[p.ID] + 1
			}
		}
		depth[c.ID] = d
	}
	return depth
}

// Subdomains returns the distinct subdomains present in the corpus, sorted.
func Subdomains(dag *DAG) []string {
	seen := map[string]bool{}
	for _, c := range dag.Concepts() {
		seen[c.Subdomain] = true
	}
	out := make([]string, 0, len(seen))
	for s := range seen {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// BuildMasteryUnitDraft partitions the corpus by (domain, floor(conceptDepth / bandWidth)).
//
// The depth function is **concept** depth, not subdomain depth, and that changed during D3.
// Banding on subdomain depth is not merely lossy here — it is *undefined* for 26 of the 125
// subdomains, because the subdomain condensation contains cycles. Concept depth is a proper
// topological ranking on an acyclic graph, so it resolves all 657 concepts and produces zero
// backward band edges by construction.
//
// The band width is a parameter with no default in this package on purpose: a caller that
// forgets it should not silently get 3. Everything else about a unit — id, title, depth range,
// subdomains, membership — is derived, so the file is reproducible from the corpus alone.
func BuildMasteryUnitDraft(dag *DAG, bandWidth int) (*MasteryUnits, error) {
	if bandWidth < 1 {
		return nil, fmt.Errorf("band width must be >= 1, got %d", bandWidth)
	}
	depths := ConceptDepths(dag)

	type key struct {
		domain string
		band   int
	}
	byKey := map[key]*MasteryUnit{}
	for _, c := range dag.Order() {
		d := depths[c.ID]
		k := key{domain: c.Domain, band: d / bandWidth}
		u, ok := byKey[k]
		if !ok {
			u = &MasteryUnit{
				ID:         UnitID(k.domain, k.band),
				Domain:     k.domain,
				Band:       k.band,
				DepthMin:   d,
				DepthMax:   d,
				Subdomains: []string{},
				ConceptIDs: []string{},
			}
			byKey[k] = u
		}
		u.ConceptIDs = append(u.ConceptIDs, c.ID)
		u.Title = UnitTitle(c.Domain, k.band)
		if d < u.DepthMin {
			u.DepthMin = d
		}
		if d > u.DepthMax {
			u.DepthMax = d
		}
		if !contains(u.Subdomains, c.Subdomain) {
			u.Subdomains = append(u.Subdomains, c.Subdomain)
		}
	}

	keys := make([]key, 0, len(byKey))
	for k := range byKey {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].domain != keys[j].domain {
			return keys[i].domain < keys[j].domain
		}
		return keys[i].band < keys[j].band
	})

	units := make([]MasteryUnit, 0, len(keys))
	for _, k := range keys {
		u := byKey[k]
		sort.Strings(u.Subdomains)
		sort.Strings(u.ConceptIDs)
		u.ConceptCount = len(u.ConceptIDs)
		units = append(units, *u)
	}

	return &MasteryUnits{
		Draft:       true,
		GeneratedBy: "go run ./cmd/masteryunits",
		Rule:        "unit = (domain, floor(conceptDepth / bandWidth)); conceptDepth = the concept's level in the concept DAG's own topological order",
		BandWidth:   bandWidth,
		Note: "DRAFT — ADR-046 D3. Band width is the open decision (how often a learner is " +
			"assessed), not product policy. Unit prerequisite order is DERIVED from the concept " +
			"DAG and deliberately not stored: authoring it would be a second hierarchy. " +
			"Do not hand-edit concept_ids; regenerate with `go run ./cmd/masteryunits -band N`.",
		Source:       "data/concepts",
		ConceptCount: dag.Count(),
		UnitCount:    len(units),
		Units:        units,
	}, nil
}

// UnitID is a unit's stable identifier. Derived from the partition key so it is stable across
// regenerations at the same band width and changes when the band does.
func UnitID(domain string, band int) string {
	return fmt.Sprintf("%s.b%d", domain, band)
}

// UnitTitle is a human-readable label for a band of a domain.
//
// Derived from the domain slug rather than from a table, because the corpus has no display
// labels in Go and a second copy of one would be another thing to keep in sync. The learner's
// display labels live in the web app (lib/graphDomains.ts).
func UnitTitle(domain string, band int) string {
	words := strings.Split(strings.ReplaceAll(domain, "_", " "), " ")
	for i, w := range words {
		if w == "" {
			continue
		}
		words[i] = strings.ToUpper(w[:1]) + w[1:]
	}
	return fmt.Sprintf("%s · Stage %d", strings.Join(words, " "), band+1)
}

// BackwardBandEdge is a prerequisite edge that crosses the proposed unit ordering backwards.
type BackwardBandEdge struct {
	ConceptID       string `json:"concept_id"`
	ConceptLabel    string `json:"concept_label"`
	Domain          string `json:"domain"`
	Subdomain       string `json:"subdomain"`
	UnitID          string `json:"unit_id"`
	Band            int    `json:"band"`
	BandDepth       int    `json:"band_depth"`
	PrereqID        string `json:"prereq_id"`
	PrereqLabel     string `json:"prereq_label"`
	PrereqSubdomain string `json:"prereq_subdomain"`
	PrereqUnitID    string `json:"prereq_unit_id"`
	PrereqBand      int    `json:"prereq_band"`
	PrereqBandDepth int    `json:"prereq_band_depth"`
	SameSubdomain   bool   `json:"same_subdomain"`
}

// BandIndex is the partition's key function, exposed so the validator, the generator and the
// backward-edge scan cannot disagree about which unit a concept belongs to.
func BandIndex(mu *MasteryUnits, dag *DAG) func(conceptID string) (unitID string, band int, bandDepth int, ok bool) {
	depths := ConceptDepths(dag)
	unitOf := make(map[string]string, dag.Count())
	for _, u := range mu.Units {
		for _, id := range u.ConceptIDs {
			unitOf[id] = u.ID
		}
	}
	return func(conceptID string) (string, int, int, bool) {
		c := dag.Concept(conceptID)
		if c == nil {
			return "", 0, 0, false
		}
		d := depths[conceptID]
		band := d / mu.BandWidth
		return unitOf[conceptID], band, d, true
	}
}

// BackwardBandEdges lists every prerequisite edge that crosses the unit ordering backwards.
//
// Reported, never corrected. These are the 22 edges the D3 analysis surfaced, and the point of
// the function is to keep them visible: a partition whose derived ordering is contradicted by the
// graph is a finding about the partition, not something to make the validator quiet about.
func BackwardBandEdges(mu *MasteryUnits, dag *DAG) []BackwardBandEdge {
	locate := BandIndex(mu, dag)
	var out []BackwardBandEdge
	for _, c := range dag.Order() {
		uid, band, bd, ok := locate(c.ID)
		if !ok {
			continue
		}
		for _, p := range dag.PrereqsOf(c.ID) {
			puid, pband, pbd, pok := locate(p.ID)
			if !pok || pband <= band {
				continue
			}
			out = append(out, BackwardBandEdge{
				ConceptID: c.ID, ConceptLabel: c.Label, Domain: c.Domain, Subdomain: c.Subdomain,
				UnitID: uid, Band: band, BandDepth: bd,
				PrereqID: p.ID, PrereqLabel: p.Label, PrereqSubdomain: p.Subdomain,
				PrereqUnitID: puid, PrereqBand: pband, PrereqBandDepth: pbd,
				SameSubdomain: p.Subdomain == c.Subdomain,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ConceptID < out[j].ConceptID })
	return out
}

// DerivedUnitPrereq returns unit → unit edges derived from concept prerequisites, and the
// concept-level edge that justifies each. Nothing about this is stored: it is recomputed from the
// DAG every time, which is the property ADR-046 asks for.
func DerivedUnitPrereqs(mu *MasteryUnits, dag *DAG) (edges []UnitEdge, err error) {
	locate := BandIndex(mu, dag)
	seen := map[string]bool{}
	type acc struct {
		unit, prereq string
		via          string
	}
	var accs []acc
	for _, c := range dag.Order() {
		uid, _, _, ok := locate(c.ID)
		if !ok {
			continue
		}
		for _, p := range dag.PrereqsOf(c.ID) {
			puid, _, _, pok := locate(p.ID)
			if !pok || puid == uid {
				continue
			}
			k := puid + "\x00" + uid
			if seen[k] {
				continue
			}
			seen[k] = true
			accs = append(accs, acc{unit: uid, prereq: puid, via: p.ID + " → " + c.ID})
		}
	}
	sort.Slice(accs, func(i, j int) bool {
		if accs[i].prereq != accs[j].prereq {
			return accs[i].prereq < accs[j].prereq
		}
		return accs[i].unit < accs[j].unit
	})
	for _, a := range accs {
		edges = append(edges, UnitEdge{Prereq: a.prereq, Unit: a.unit, Via: a.via})
	}
	return edges, nil
}

// UnitEdge is one derived unit-to-unit prerequisite relationship.
type UnitEdge struct {
	Prereq string `json:"prereq"`
	Unit   string `json:"unit"`
	Via    string `json:"via"`
}

// Severity distinguishes a malformed file from a finding about the data.
type Severity string

const (
	// SeverityError means the file is wrong: coverage, duplicate membership, a dangling
	// reference, or membership that has drifted from the rule it declares.
	SeverityError Severity = "error"
	// SeverityAnomaly means the file is internally consistent but the corpus or the rule
	// produces something worth a human's attention. Backward band edges are anomalies: the
	// partition contradicts the graph it is supposed to be derived from, which is a finding,
	// not a parse error.
	SeverityAnomaly Severity = "anomaly"
)

// Finding is one validation result.
type Finding struct {
	Severity Severity `json:"severity"`
	Kind     string   `json:"kind"`
	Detail   string   `json:"detail"`
}

// LoadMasteryUnits reads and parses the draft.
func LoadMasteryUnits(path string) (*MasteryUnits, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var mu MasteryUnits
	if err := json.Unmarshal(raw, &mu); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &mu, nil
}

// ValidateMasteryUnits checks a draft against the corpus and against the rule the draft itself
// declares.
//
// The drift check is the important one: membership is recomputed from `mu.BandWidth` and
// compared, so hand-editing concept_ids cannot quietly produce a boundary nobody chose. Using
// the file's declared width (rather than a constant here) is also what lets the same validator
// check a width-2 or width-4 regeneration without a code change.
func ValidateMasteryUnits(dag *DAG, mu *MasteryUnits) []Finding {
	var out []Finding
	add := func(sev Severity, kind, format string, args ...any) {
		out = append(out, Finding{Severity: sev, Kind: kind, Detail: fmt.Sprintf(format, args...)})
	}

	if mu.BandWidth < 1 {
		add(SeverityError, "band_width", "band_width is %d; must be >= 1", mu.BandWidth)
		return out
	}

	// Referential integrity and duplicate identity.
	unitIDs := map[string]bool{}
	for _, u := range mu.Units {
		if unitIDs[u.ID] {
			add(SeverityError, "duplicate_unit_id", "unit id %q appears more than once", u.ID)
		}
		unitIDs[u.ID] = true
	}
	owner := map[string]string{}
	for _, u := range mu.Units {
		for _, id := range u.ConceptIDs {
			if dag.Concept(id) == nil {
				add(SeverityError, "dangling_concept", "unit %q references unknown concept %q", u.ID, id)
				continue
			}
			if prev, dup := owner[id]; dup {
				add(SeverityError, "duplicate_membership",
					"concept %q is in both %q and %q; a partition must cover each concept once", id, prev, u.ID)
				continue
			}
			owner[id] = u.ID
		}
		if u.ConceptCount != len(u.ConceptIDs) {
			add(SeverityError, "stale_count", "unit %q declares concept_count %d but lists %d",
				u.ID, u.ConceptCount, len(u.ConceptIDs))
		}
	}
	for _, c := range dag.Order() {
		if _, ok := owner[c.ID]; !ok {
			add(SeverityError, "uncovered_concept",
				"concept %q (%s/%s) belongs to no mastery unit", c.ID, c.Domain, c.Subdomain)
		}
	}

	// Boundary consistency: membership must equal what the declared rule produces.
	want, err := BuildMasteryUnitDraft(dag, mu.BandWidth)
	if err != nil {
		add(SeverityError, "rule", "cannot rebuild draft at band width %d: %v", mu.BandWidth, err)
		return out
	}
	wantOwner := map[string]string{}
	for _, u := range want.Units {
		for _, id := range u.ConceptIDs {
			wantOwner[id] = u.ID
		}
	}
	for _, u := range mu.Units {
		for _, id := range u.ConceptIDs {
			if w, ok := wantOwner[id]; ok && w != u.ID {
				add(SeverityError, "boundary_drift",
					"concept %q is in %q but the declared rule (domain, floor(depth/%d)) puts it in %q; regenerate rather than hand-edit",
					id, u.ID, mu.BandWidth, w)
			}
		}
	}

	// Anomaly: the derived unit ordering is contradicted by the concept DAG.
	if back := BackwardBandEdges(mu, dag); len(back) > 0 {
		add(SeverityAnomaly, "backward_band_edges",
			"%d prerequisite edge(s) cross the unit ordering backwards, so unit order is not a "+
				"topological order of the graph. Surfaced, not corrected — see docs/mastery-assessments.md D3.",
			len(back))
	}

	// Anomaly: the partition must actually partition. A unit of one concept is permitted here
	// and reported, because whether it is an error is a content decision (D3), not a schema rule.
	singletons := 0
	for _, u := range mu.Units {
		if len(u.ConceptIDs) == 1 {
			singletons++
		}
	}
	if singletons > 0 {
		add(SeverityAnomaly, "singleton_units",
			"%d unit(s) hold a single concept; the spec forbids a quiz per concept, so these need merging",
			singletons)
	}

	return out
}

// Errors returns only the findings that make the file invalid.
func Errors(findings []Finding) []Finding {
	var out []Finding
	for _, f := range findings {
		if f.Severity == SeverityError {
			out = append(out, f)
		}
	}
	return out
}

// Anomalies returns the findings that are about the data rather than the file.
func Anomalies(findings []Finding) []Finding {
	var out []Finding
	for _, f := range findings {
		if f.Severity == SeverityAnomaly {
			out = append(out, f)
		}
	}
	return out
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func insertSorted(s []string, v string) []string {
	i := sort.SearchStrings(s, v)
	s = append(s, "")
	copy(s[i+1:], s[i:])
	s[i] = v
	return s
}

func contains(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}
