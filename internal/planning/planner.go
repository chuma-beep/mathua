package planning

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/storage"
)

type Course struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Grade       string   `json:"grade"`
	Description string   `json:"description"`
	Targets     []string `json:"targets"`
}

type Path struct {
	Course   *Course
	Concepts []*concepts.Concept
}

type Planner struct {
	dag     *concepts.DAG
	courses []*Course
	byID    map[string]*Course
}

func Load(coursesPath string, dag *concepts.DAG) (*Planner, error) {
	data, err := os.ReadFile(coursesPath)
	if err != nil {
		return nil, fmt.Errorf("read courses: %w", err)
	}
	var courses []*Course
	if err := json.Unmarshal(data, &courses); err != nil {
		return nil, fmt.Errorf("parse courses: %w", err)
	}
	byID := make(map[string]*Course, len(courses))
	for _, c := range courses {
		byID[c.ID] = c
	}
	return &Planner{dag: dag, courses: courses, byID: byID}, nil
}

func (p *Planner) Courses() []*Course {
	return p.courses
}

func (p *Planner) Course(id string) *Course {
	return p.byID[id]
}

func (p *Planner) PrerequisiteChain(targets []string) (*Path, error) {
	inChain := make(map[string]bool)
	var queue []string

	for _, t := range targets {
		if p.dag.Concept(t) == nil {
			return nil, fmt.Errorf("target concept %q not found in DAG", t)
		}
		queue = append(queue, t)
	}
	// BFS backward through prerequisites
	for len(queue) > 0 {
		cid := queue[0]
		queue = queue[1:]
		if inChain[cid] {
			continue
		}
		inChain[cid] = true
		for _, prereq := range p.dag.PrereqsOf(cid) {
			if !inChain[prereq.ID] {
				queue = append(queue, prereq.ID)
			}
		}
	}
	// Collect and topo-sort (respecting existing DAG order)
	var chain []*concepts.Concept
	for _, c := range p.dag.Order() {
		if inChain[c.ID] {
			chain = append(chain, c)
		}
	}
	return &Path{
		Course:   &Course{Targets: targets},
		Concepts: chain,
	}, nil
}

func (p *Planner) Unmastered(path *Path, progress map[string]*storage.ConceptProgress) []*concepts.Concept {
	var out []*concepts.Concept
	for _, c := range path.Concepts {
		if prog, ok := progress[c.ID]; !ok || prog.Status != "MASTERED" {
			out = append(out, c)
		}
	}
	return out
}

func (p *Planner) PathForCourse(courseID string) (*Path, error) {
	c := p.byID[courseID]
	if c == nil {
		return nil, fmt.Errorf("course %q not found", courseID)
	}
	path, err := p.PrerequisiteChain(c.Targets)
	if err != nil {
		return nil, err
	}
	path.Course = c
	return path, nil
}

func (p *Planner) SortCourses(courses []*Course) {
	sort.Slice(courses, func(i, j int) bool {
		return courses[i].ID < courses[j].ID
	})
}
