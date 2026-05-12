package concepts

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Load reads concept definitions from a JSON file, validates the graph,
// topologically sorts, and returns the populated DAG.
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

// Build validates a slice of concepts (no cycles, no orphan prerequisites,
// no empty or duplicate IDs), topologically sorts, and returns a DAG.
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

	// Topological sort via Kahn's algorithm
	inDegree := make(map[string]int, len(raw))
	for _, c := range raw {
		inDegree[c.ID] = len(c.Prerequisites)
	}
	var queue []string
	for _, c := range raw {
		if inDegree[c.ID] == 0 {
			queue = append(queue, c.ID)
		}
	}
	var order []*Concept
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		order = append(order, cmap[id])
		for _, dep := range depMap[id] {
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
		concepts:     cmap,
		order:        order,
		domains:      domains,
		prereqsOf:    prereqMap,
		dependentsOf: depMap,
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
	}

	// Detect cycles via topological sort (Kahn's)
	cmap := make(map[string]*Concept, len(raw))
	depMap := make(map[string][]string)
	inDegree := make(map[string]int, len(raw))
	for i := range raw {
		c := &raw[i]
		cmap[c.ID] = c
		inDegree[c.ID] = len(c.Prerequisites)
	}
	for _, c := range raw {
		for _, pid := range c.Prerequisites {
			depMap[pid] = append(depMap[pid], c.ID)
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
		return fmt.Errorf("concept graph contains a cycle")
	}

	return nil
}
