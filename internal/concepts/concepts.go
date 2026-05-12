package concepts

type MasteryThreshold struct {
	Streak         int     `json:"streak"`
	AvgTimeSeconds float64 `json:"avg_time_seconds"`
}

type Concept struct {
	ID               string           `json:"id"`
	Label            string           `json:"label"`
	Domain           string           `json:"domain"`
	Subdomain        string           `json:"subdomain"`
	GradingType      string           `json:"grading_type"`
	Prerequisites    []string         `json:"prerequisites"`
	MasteryThreshold MasteryThreshold `json:"mastery_threshold"`
}

type DAG struct {
	concepts     map[string]*Concept
	order        []*Concept
	domains      []string
	prereqsOf    map[string][]*Concept
	dependentsOf map[string][]*Concept
}

func (d *DAG) Concepts() map[string]*Concept {
	c := make(map[string]*Concept, len(d.concepts))
	for k, v := range d.concepts {
		c[k] = v
	}
	return c
}

func (d *DAG) Concept(id string) *Concept {
	return d.concepts[id]
}

func (d *DAG) Order() []*Concept {
	o := make([]*Concept, len(d.order))
	copy(o, d.order)
	return o
}

func (d *DAG) Count() int {
	return len(d.concepts)
}

func (d *DAG) Domains() []string {
	dom := make([]string, len(d.domains))
	copy(dom, d.domains)
	return dom
}

func (d *DAG) PrereqsOf(id string) []*Concept {
	if ps, ok := d.prereqsOf[id]; ok {
		out := make([]*Concept, len(ps))
		copy(out, ps)
		return out
	}
	return nil
}

func (d *DAG) DependentsOf(id string) []*Concept {
	if deps, ok := d.dependentsOf[id]; ok {
		out := make([]*Concept, len(deps))
		copy(out, deps)
		return out
	}
	return nil
}

// Available returns concepts whose prerequisites are all present in the
// mastered set. Only includes concepts that are NOT themselves in mastered.
func (d *DAG) Available(mastered map[string]bool) []*Concept {
	var avail []*Concept
	for _, c := range d.order {
		if mastered[c.ID] {
			continue
		}
		if d.allPrereqsMastered(c, mastered) {
			avail = append(avail, c)
		}
	}
	return avail
}

func (d *DAG) allPrereqsMastered(c *Concept, mastered map[string]bool) bool {
	for _, pid := range c.Prerequisites {
		if !mastered[pid] {
			return false
		}
	}
	return true
}
