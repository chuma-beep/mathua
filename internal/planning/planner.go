package planning

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/mastery"
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

// New returns a Planner over the DAG without a course catalog.
// PrerequisitesOf works catalog-free; course lookups stay empty.
func New(dag *concepts.DAG) *Planner {
	return &Planner{dag: dag, byID: map[string]*Course{}}
}

func (p *Planner) Courses() []*Course {
	return p.courses
}

func (p *Planner) Course(id string) *Course {
	return p.byID[id]
}

func (p *Planner) collectChain(targets []string) (map[string]bool, error) {
	inChain := make(map[string]bool)
	var queue []string
	for _, t := range targets {
		if p.dag.Concept(t) == nil {
			return nil, fmt.Errorf("target concept %q not found in DAG", t)
		}
		queue = append(queue, t)
	}
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
	return inChain, nil
}

func (p *Planner) PrerequisitesOf(conceptIDs []string) (*Path, error) {
	inChain, err := p.collectChain(conceptIDs)
	if err != nil {
		return nil, err
	}
	var chain []*concepts.Concept
	for _, c := range p.dag.Order() {
		if inChain[c.ID] {
			chain = append(chain, c)
		}
	}
	return &Path{
		Course:   &Course{Targets: conceptIDs},
		Concepts: chain,
	}, nil
}

func (p *Planner) PrerequisiteChain(targets []string) (*Path, error) {
	return p.PrerequisitesOf(targets)
}

func (p *Planner) PathForCourse(courseID string) (*Path, error) {
	c := p.byID[courseID]
	if c == nil {
		return nil, fmt.Errorf("course %q not found", courseID)
	}
	path, err := p.PrerequisitesOf(c.Targets)
	if err != nil {
		return nil, err
	}
	path.Course = c
	return path, nil
}

func (p *Planner) Unmastered(path *Path, progress map[string]*storage.ConceptProgress) []*concepts.Concept {
	var out []*concepts.Concept
	for _, c := range path.Concepts {
		if prog, ok := progress[c.ID]; !ok || prog.Status != string(mastery.StatusMastered) {
			out = append(out, c)
		}
	}
	return out
}

func (p *Planner) WeakAreas(path *Path, progress map[string]*storage.ConceptProgress) []*concepts.Concept {
	return p.Unmastered(path, progress)
}

func (p *Planner) Readiness(path *Path, progress map[string]*storage.ConceptProgress) float64 {
	if len(path.Concepts) == 0 {
		return 1.0
	}
	mastered := 0
	for _, c := range path.Concepts {
		if prog, ok := progress[c.ID]; ok && prog.Status == string(mastery.StatusMastered) {
			mastered++
		}
	}
	return float64(mastered) / float64(len(path.Concepts))
}

// CourseProgress is the accreditation-track view of one course.
type CourseProgress struct {
	Total       int      `json:"total"`
	Mastered    int      `json:"mastered"`
	Pct         float64  `json:"pct"`
	Remaining   []string `json:"remaining,omitempty"`
	DaysRemaining int    `json:"days_remaining"` // estimate at daily_xp_goal
}

// ProgressForCourse computes mastery over the course's transitive path.
func (p *Planner) ProgressForCourse(course *Course, progress map[string]*storage.ConceptProgress, dailyXPGoal int) (*CourseProgress, error) {
	path, err := p.PathForCourse(course.ID)
	if err != nil {
		return nil, err
	}
	total := len(path.Concepts)
	mastered := 0
	var remaining []string
	for _, c := range path.Concepts {
		if prog, ok := progress[c.ID]; ok && prog.Status == string(mastery.StatusMastered) {
			mastered++
		} else {
			remaining = append(remaining, c.ID)
		}
	}
	days := 0
	if goal := dailyXPGoal; goal > 0 && len(remaining) > 0 {
		// Estimate: ~5 concepts per 30 XP of focused work (MA ~1 min/XP).
		perDay := goal / 6
		if perDay < 1 {
			perDay = 1
		}
		days = (len(remaining) + perDay - 1) / perDay
	}
	pct := 0.0
	if total > 0 {
		pct = float64(mastered) / float64(total)
	}
	return &CourseProgress{
		Total:         total,
		Mastered:      mastered,
		Pct:           pct,
		Remaining:     remaining,
		DaysRemaining: days,
	}, nil
}

func (p *Planner) SortCourses(courses []*Course) {
	sort.Slice(courses, func(i, j int) bool {
		return courses[i].ID < courses[j].ID
	})
}
