package concepts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Load reads a JSON concept file, validates the graph, topo-sorts, and returns a DAG.
func Load(path string) (*DAG, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read concepts file: %w", err)
	}
	var raw []Concept
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse concepts: %w", err)
	}
	return Build(raw)
}

// LoadDir reads all .json files from a directory, merges, validates, topo-sorts, and returns a DAG.
// enrichment.json is merged as heuristics (Variants/Encompasses/InterferenceGroup/KeyPrerequisites) on top of base concepts.
func LoadDir(dir string) (*DAG, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read concepts dir: %w", err)
	}
	var raw []Concept
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		if e.Name() == "enrichment.json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", e.Name(), err)
		}
		var batch []Concept
		if err := json.Unmarshal(data, &batch); err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		raw = append(raw, batch...)
	}
	// Merge enrichment heuristics (Variants, Encompasses, InterferenceGroup, KeyPrerequisites)
	if err := mergeEnrichment(dir, raw); err != nil {
		return nil, err
	}
	return Build(raw)
}

// enrichmentEntry mirrors the loose schema in data/concepts/enrichment.json.
type enrichmentEntry struct {
	ID               string    `json:"id"`
	Heuristic        string    `json:"heuristic"`
	Note             string    `json:"note"`
	Encompasses      []string  `json:"encompasses"`
	KeyPrerequisites []string  `json:"key_prerequisites"`
	InterferenceGroup string   `json:"interference_group"`
	Members          []string  `json:"members"`
	Variants         []Variant `json:"variants"`
}

func mergeEnrichment(dir string, raw []Concept) error {
	path := filepath.Join(dir, "enrichment.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read enrichment.json: %w", err)
	}
	var entries []enrichmentEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("parse enrichment.json: %w", err)
	}
	// Index concepts by ID for patching
	idx := make(map[string]*Concept, len(raw))
	for i := range raw {
		idx[raw[i].ID] = &raw[i]
	}
	for _, e := range entries {
		if e.ID != "" {
			c, ok := idx[e.ID]
			if !ok {
				return fmt.Errorf("enrichment id %q not found in DAG", e.ID)
			}
			if len(e.Encompasses) > 0 {
				c.Encompasses = append(c.Encompasses, e.Encompasses...)
			}
			if len(e.KeyPrerequisites) > 0 {
				c.KeyPrerequisites = append(c.KeyPrerequisites, e.KeyPrerequisites...)
			}
			if e.InterferenceGroup != "" {
				if c.InterferenceGroup != "" && c.InterferenceGroup != e.InterferenceGroup {
					return fmt.Errorf("concept %q interference_group conflict: %q vs %q", e.ID, c.InterferenceGroup, e.InterferenceGroup)
				}
				c.InterferenceGroup = e.InterferenceGroup
			}
			if len(e.Variants) > 0 {
				c.Variants = append(c.Variants, e.Variants...)
			}
		} else if e.InterferenceGroup != "" && len(e.Members) > 0 {
			for _, mid := range e.Members {
				c, ok := idx[mid]
				if !ok {
					return fmt.Errorf("enrichment member %q (group %q) not found", mid, e.InterferenceGroup)
				}
				if c.InterferenceGroup != "" && c.InterferenceGroup != e.InterferenceGroup {
					return fmt.Errorf("concept %q interference_group conflict: %q vs %q", mid, c.InterferenceGroup, e.InterferenceGroup)
				}
				c.InterferenceGroup = e.InterferenceGroup
			}
		}
	}
	return nil
}

// Build validates a concept slice (no cycles, orphan prereqs, empty/duplicate IDs), topo-sorts, and returns a DAG.
func Build(raw []Concept) (*DAG, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("concepts slice is empty")
	}
	if err := validate(raw); err != nil {
		return nil, err
	}
	return build(raw), nil
}

func build(raw []Concept) *DAG {
	cmap := make(map[string]*Concept, len(raw))
	for i := range raw {
		c := raw[i]
		cmap[c.ID] = &c
	}

	// Compute prereqsOf and dependentsOf
	prereqMap := make(map[string][]*Concept, len(raw))
	depMap := make(map[string][]*Concept, len(raw))
	for _, c := range raw {
		for _, pid := range c.Prerequisites {
			if p, ok := cmap[pid]; ok {
				prereqMap[c.ID] = append(prereqMap[c.ID], p)
				depMap[pid] = append(depMap[pid], cmap[c.ID])
			}
		}
	}

	// Compute encompassesOf and encompassedBy
	encompassedBy := make(map[string][]*Concept, len(raw))
	encompassesOf := make(map[string][]*Concept, len(raw))
	for _, c := range raw {
		for _, eid := range c.Encompasses {
			if target, ok := cmap[eid]; ok {
				encompassesOf[c.ID] = append(encompassesOf[c.ID], target)
				encompassedBy[eid] = append(encompassedBy[eid], cmap[c.ID])
			}
		}
	}

	// Compute interferersOf by InterferenceGroup
	groupMembers := make(map[string][]*Concept)
	for _, c := range raw {
		if c.InterferenceGroup != "" {
			groupMembers[c.InterferenceGroup] = append(groupMembers[c.InterferenceGroup], cmap[c.ID])
		}
	}
	interferersOf := make(map[string][]*Concept, len(raw))
	for _, members := range groupMembers {
		for _, c := range members {
			for _, other := range members {
				if other.ID != c.ID {
					interferersOf[c.ID] = append(interferersOf[c.ID], other)
				}
			}
		}
	}

	// Topological sort via Kahn's algorithm — includes prerequisites + encompasses as edges
	inDegree := make(map[string]int, len(raw))
	for _, c := range raw {
		inDegree[c.ID] = len(c.Prerequisites) + len(c.Encompasses)
	}
	var queue []string
	for _, c := range raw {
		if inDegree[c.ID] == 0 {
			queue = append(queue, c.ID)
		}
	}
	// depMap for Kahn includes both prerequisites and encompasses dependents
	combinedDepMap := make(map[string][]*Concept, len(raw))
	for k, v := range depMap {
		combinedDepMap[k] = append([]*Concept(nil), v...)
	}
	for eid, encompassors := range encompassedBy {
		for _, enc := range encompassors {
			combinedDepMap[eid] = append(combinedDepMap[eid], enc)
		}
	}
	var order []*Concept
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		order = append(order, cmap[id])
		for _, dep := range combinedDepMap[id] {
			inDegree[dep.ID]--
			if inDegree[dep.ID] == 0 {
				queue = append(queue, dep.ID)
			}
		}
	}

	// Collect unique domains, preserving order of first appearance
	seen := make(map[string]bool)
	var domains []string
	for _, c := range order {
		if !seen[c.Domain] {
			seen[c.Domain] = true
			domains = append(domains, c.Domain)
		}
	}
	sort.Strings(domains)

	return &DAG{
		concepts:      cmap,
		order:         order,
		domains:       domains,
		prereqsOf:     prereqMap,
		dependentsOf:  depMap,
		encompassedBy: encompassedBy,
		encompassesOf: encompassesOf,
		interferersOf: interferersOf,
	}
}

func validate(raw []Concept) error {
	ids := make(map[string]bool, len(raw))
	for _, c := range raw {
		if c.ID == "" {
			return fmt.Errorf("concept with empty ID")
		}
		if ids[c.ID] {
			return fmt.Errorf("duplicate concept ID %q", c.ID)
		}
		ids[c.ID] = true
	}

	// Detect orphan prerequisites
	for _, c := range raw {
		for _, pid := range c.Prerequisites {
			if !ids[pid] {
				return fmt.Errorf("concept %q requires %q which does not exist", c.ID, pid)
			}
		}
		for _, eid := range c.Encompasses {
			if !ids[eid] {
				return fmt.Errorf("concept %q encompasses %q which does not exist", c.ID, eid)
			}
			if eid == c.ID {
				return fmt.Errorf("concept %q cannot encompass itself", c.ID)
			}
		}
		for _, kpid := range c.KeyPrerequisites {
			if !ids[kpid] {
				return fmt.Errorf("concept %q key_prerequisite %q does not exist", c.ID, kpid)
			}
		}
		for _, kpid := range c.KeyPrerequisites {
			found := false
			for _, pid := range c.Prerequisites {
				if pid == kpid {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("concept %q key_prerequisite %q must be a prerequisite", c.ID, kpid)
			}
		}
		for _, v := range c.Variants {
			if v.Difficulty < 0.1 || v.Difficulty > 1.0 {
				return fmt.Errorf("concept %q variant %q has invalid difficulty %v (must be 0.1-1.0)", c.ID, v.Label, v.Difficulty)
			}
			if v.TimeThresh <= 0 {
				return fmt.Errorf("concept %q variant %q has invalid time_thresh %v", c.ID, v.Label, v.TimeThresh)
			}
			if v.Label == "" {
				return fmt.Errorf("concept %q variant missing label", c.ID)
			}
		}
		// Grading type enum check (allow empty for test fixtures)
		allowedGrading := map[string]bool{
			"multiple_choice": true, "numeric": true, "symbolic": true, "complex": true,
			"expression": true, "ordering": true, "tuple": true, "polynomial": true, "comparison": true,
		}
		if c.GradingType != "" && !allowedGrading[c.GradingType] {
			return fmt.Errorf("concept %q has invalid grading_type %q", c.ID, c.GradingType)
		}
		if c.Domain == "" {
			return fmt.Errorf("concept %q has empty domain", c.ID)
		}
		if c.Label == "" {
			return fmt.Errorf("concept %q has empty label", c.ID)
		}
		if c.MasteryThreshold.Streak <= 0 {
			return fmt.Errorf("concept %q has invalid streak %d", c.ID, c.MasteryThreshold.Streak)
		}
		if c.MasteryThreshold.AvgTimeSeconds <= 0 {
			return fmt.Errorf("concept %q has invalid avg_time_seconds %v", c.ID, c.MasteryThreshold.AvgTimeSeconds)
		}
		if c.MasteryThreshold.AvgTimeSeconds > 300 {
			return fmt.Errorf("concept %q avg_time_seconds %v exceeds sanity cap 300s", c.ID, c.MasteryThreshold.AvgTimeSeconds)
		}
	}

	// Detect non-trivial interference groups (singletons are errors)
	groupCounts := make(map[string]int)
	for _, c := range raw {
		if c.InterferenceGroup != "" {
			groupCounts[c.InterferenceGroup]++
		}
	}
	for g, cnt := range groupCounts {
		if cnt < 2 {
			return fmt.Errorf("interference group %q is non-trivial: only %d member", g, cnt)
		}
	}

	// Detect cycles via topological sort (Kahn's) — includes prerequisites + encompasses
	cmap := make(map[string]*Concept, len(raw))
	depMap := make(map[string][]string)
	inDegree := make(map[string]int, len(raw))
	for i := range raw {
		c := &raw[i]
		cmap[c.ID] = c
		inDegree[c.ID] = len(c.Prerequisites) + len(c.Encompasses)
	}
	for _, c := range raw {
		for _, pid := range c.Prerequisites {
			depMap[pid] = append(depMap[pid], c.ID)
		}
		for _, eid := range c.Encompasses {
			depMap[eid] = append(depMap[eid], c.ID)
		}
	}
	var queue []string
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}
	processed := 0
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		processed++
		for _, dep := range depMap[id] {
			inDegree[dep]--
			if inDegree[dep] == 0 {
				queue = append(queue, dep)
			}
		}
	}
	if processed != len(raw) {
		return fmt.Errorf("concept graph contains a cycle (including encompasses)")
	}

	return nil
}
