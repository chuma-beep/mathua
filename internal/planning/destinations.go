package planning

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Destination is a learner-chosen learning goal: a named bundle of courses.
// The graph stays source of truth — a destination resolves to target
// concepts, whose transitive prerequisite closure is the actual curriculum.
// Shared foundations across member courses are counted once.
type Destination struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Courses     []string `json:"courses"`
}

// LoadDestinations reads destination bundles from a JSON file.
func LoadDestinations(path string) ([]*Destination, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read destinations: %w", err)
	}
	var ds []*Destination
	if err := json.Unmarshal(data, &ds); err != nil {
		return nil, fmt.Errorf("parse destinations: %w", err)
	}
	seen := map[string]bool{}
	for _, d := range ds {
		if d.ID == "" {
			return nil, fmt.Errorf("destination with empty id")
		}
		if seen[d.ID] {
			return nil, fmt.Errorf("duplicate destination id %q", d.ID)
		}
		seen[d.ID] = true
		if len(d.Courses) == 0 {
			return nil, fmt.Errorf("destination %q has no courses", d.ID)
		}
	}
	return ds, nil
}

// SetDestinations registers destination bundles on the planner.
func (p *Planner) SetDestinations(ds []*Destination) {
	if p.destByID == nil {
		p.destByID = map[string]*Destination{}
	}
	for _, d := range ds {
		p.destByID[d.ID] = d
	}
}

// Destinations returns registered destination bundles, sorted by ID.
func (p *Planner) Destinations() []*Destination {
	out := make([]*Destination, 0, len(p.destByID))
	for _, d := range p.destByID {
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Destination returns one bundle by ID, or nil.
func (p *Planner) Destination(id string) *Destination {
	return p.destByID[id]
}

// PathForDestination resolves a destination to its full learning path: the
// union of member courses' target closures. Shared prerequisites appear once.
func (p *Planner) PathForDestination(destID string) (*Path, error) {
	d := p.destByID[destID]
	if d == nil {
		return nil, fmt.Errorf("destination %q not found", destID)
	}
	var targets []string
	for _, cid := range d.Courses {
		c := p.byID[cid]
		if c == nil {
			return nil, fmt.Errorf("destination %q references unknown course %q", destID, cid)
		}
		targets = append(targets, c.Targets...)
	}
	if len(targets) == 0 {
		return nil, fmt.Errorf("destination %q has no target concepts", destID)
	}
	path, err := p.PrerequisitesOf(targets)
	if err != nil {
		return nil, err
	}
	return path, nil
}
