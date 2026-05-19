package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Screen int

const (
	ScreenWelcome Screen = iota
	ScreenProblem
	ScreenFeedback
	ScreenProgress
	ScreenStudy
	ScreenDiagnostic
	ScreenDiagFeedback
	ScreenBrowse
)

type ConceptNode struct {
	ID         string
	Label      string
	Unlocked   bool
	MasteryPct float64
	Streak     int
	Status     string
}

type SubdomainNode struct {
	Name     string
	Concepts []ConceptNode
}

type DomainNode struct {
	Name       string
	Subdomains []SubdomainNode
}

type ConceptTreeMsg struct {
	Domains []DomainNode
}

type StudyLesson struct {
	Title     string
	Body      string
	ConceptID string
	Domain    string
}

type DiagQuestionMsg struct {
	ConceptID   string
	ConceptName string
	Question    string
	Count       int
}

type DiagFeedbackMsg struct {
	Correct       bool
	UserAnswer    string
	CorrectAnswer string
	Explanation   string
	Done          bool
}

type StudyListMsg struct {
	Lessons []StudyLesson
}

type SessionMsg struct {
	ConceptID    string
	Domain       string
	Subdomain    string
	Label        string
	Question     string
	Expected     string
	Difficulty   float64
	TimeLimit    float64
	MasteryPct   float64
	Streak       int
	StreakNeeded int
	Status       string
	LessonTitle  string
}

type SubmitResultMsg struct {
	Correct          bool
	UserAnswer       string
	CorrectAnswer    string
	Explanation      string
	ElapsedSecs      float64
	NewMastery       float64
	NewStreak        int
	NewStatus        string
	MasteryAchieved  bool
	UnlockedConcepts []string
}

type ProgressMsg struct {
	Domains       []DomainProgress
	NextReviewDue string
	NewlyUnlocked []string
}

type DomainProgress struct {
	Domain     string
	Mastered   int
	Total      int
	Subdomains []SubdomainProgress
}

type SubdomainProgress struct {
	Name     string
	Mastered int
	Total    int
}

type WelcomeStatsMsg struct {
	TotalMastered  int
	TotalConcepts  int
	DayStreak      int
	Level          string
	LevelName      string
	WeeklyScore    int
	DomainProgress []DomainProgress
}

type Model struct {
	screen Screen
	width  int
	height int

	welcomeStats WelcomeStatsMsg

	session    SessionMsg
	input      textinput.Model
	roundCount int

	result   SubmitResultMsg
	progress ProgressMsg

	studyLessons []StudyLesson
	studyContent *struct{ Title, Body string }

	diagQuestion DiagQuestionMsg
	diagFeedback DiagFeedbackMsg

	browseData     ConceptTreeMsg
	browseDomain   int
	browseSub      int
	browseLevel    int // 0=domains, 1=subdomains, 2=concepts

	OnSubmit   func(conceptID, answer string, elapsed float64) SubmitResultMsg
	OnNext     func() SessionMsg
	OnProgress func() ProgressMsg
	OnWelcome  func() WelcomeStatsMsg
	OnStudy    func() []StudyLesson
	OnDiagStart  func() DiagQuestionMsg
	OnDiagSubmit func(conceptID, answer string, elapsed float64) DiagFeedbackMsg
	OnBrowse   func() ConceptTreeMsg
	OnSelectConcept func(conceptID string) SessionMsg
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "your answer"
	ti.CharLimit = 128
	ti.Width = 40
	ti.Prompt = "> "
	ti.PromptStyle = InputPromptStyle
	return Model{
		screen: ScreenWelcome,
		input:  ti,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		loadWelcomeStats(m.OnWelcome),
		textinput.Blink,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case WelcomeStatsMsg:
		m.welcomeStats = msg
		return m, nil

	case SessionMsg:
		m.session = msg
		m.screen = ScreenProblem
		m.input.Reset()
		m.input.Focus()
		return m, nil

	case SubmitResultMsg:
		m.result = msg
		m.screen = ScreenFeedback
		m.input.Blur()
		return m, nil

	case ProgressMsg:
		m.progress = msg
		m.screen = ScreenProgress
		return m, nil

	case DiagQuestionMsg:
		m.diagQuestion = msg
		m.diagFeedback = DiagFeedbackMsg{}
		m.input.Reset()
		m.input.Focus()
		m.screen = ScreenDiagnostic
		return m, nil

	case DiagFeedbackMsg:
		m.diagFeedback = msg
		m.input.Blur()
		m.screen = ScreenDiagFeedback
		return m, nil

	case StudyListMsg:
		m.studyLessons = msg.Lessons
		m.studyContent = nil
		m.screen = ScreenStudy
		return m, nil

	case ConceptTreeMsg:
		m.browseData = msg
		m.browseLevel = 0
		m.browseDomain = 0
		m.browseSub = 0
		m.screen = ScreenBrowse
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	if m.screen == ScreenProblem {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	}

	switch m.screen {

	case ScreenWelcome:
		switch msg.String() {
		case "enter", " ":
			return m, loadNextProblem(m.OnNext)
		case "b":
			return m, loadConceptTree(m.OnBrowse)
		case "d":
			return m, startDiagnostic(m.OnDiagStart)
		case "p":
			return m, loadProgress(m.OnProgress)
		case "s":
			return m, loadStudyLessons(m.OnStudy)
		case "q":
			return m, tea.Quit
		}

	case ScreenProblem:
		switch msg.String() {
		case "ctrl+p":
			return m, loadProgress(m.OnProgress)
		case "enter":
			answer := m.input.Value()
			if answer == "" {
				return m, nil
			}
			return m, submitAnswer(m.OnSubmit, m.session.ConceptID, answer, 0)
		default:
			if len(msg.String()) > 3 {
				return m, nil
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

	case ScreenFeedback:
		switch msg.String() {
		case "enter", " ", "n":
			m.roundCount++
			if m.roundCount%10 == 0 {
				return m, loadProgress(m.OnProgress)
			}
			return m, loadNextProblem(m.OnNext)
		case "p":
			return m, loadProgress(m.OnProgress)
		case "q":
			return m, tea.Quit
		}

	case ScreenStudy:
		switch {
		case msg.String() == "q":
			return m, tea.Quit
		case msg.String() == "esc" || msg.String() == "w":
			if m.studyContent != nil {
				m.studyContent = nil
				return m, nil
			}
			m.screen = ScreenWelcome
			return m, loadWelcomeStats(m.OnWelcome)
		case msg.String() == "enter" && m.studyContent == nil:
			return m, nil // No-op, handled by numbered selection
		default:
			if m.studyContent == nil {
				// Number key selects a lesson
				if len(msg.String()) == 1 && msg.String()[0] >= '1' && msg.String()[0] <= '9' {
					idx := int(msg.String()[0] - '1')
					if idx >= 0 && idx < len(m.studyLessons) {
						l := m.studyLessons[idx]
						m.studyContent = &struct{ Title, Body string }{l.Title, l.Body}
						return m, nil
					}
				}
			}
			return m, nil
		}

	case ScreenBrowse:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "esc", "w":
			if m.browseLevel == 0 {
				m.screen = ScreenWelcome
				return m, loadWelcomeStats(m.OnWelcome)
			}
			m.browseLevel--
			return m, nil
		case "enter":
			if m.browseLevel == 2 {
				domains := m.browseData.Domains
				if m.browseDomain < len(domains) {
					subs := domains[m.browseDomain].Subdomains
					if m.browseSub < len(subs) {
						concepts := subs[m.browseSub].Concepts
						if len(concepts) > 0 {
							c := concepts[0]
							if c.Unlocked {
								return m, loadConceptPractice(m.OnSelectConcept, c.ID)
							}
						}
					}
				}
			}
			return m, nil
		default:
			if len(msg.String()) == 1 && msg.String()[0] >= '0' && msg.String()[0] <= '9' {
				idx := int(msg.String()[0] - '0')
				switch m.browseLevel {
				case 0:
					if idx >= 1 && idx <= len(m.browseData.Domains) {
						m.browseDomain = idx - 1
						m.browseSub = 0
						m.browseLevel = 1
					}
				case 1:
					domains := m.browseData.Domains
					if m.browseDomain < len(domains) {
						subs := domains[m.browseDomain].Subdomains
						if idx >= 1 && idx <= len(subs) {
							m.browseSub = idx - 1
							m.browseLevel = 2
						}
					}
				case 2:
					domains := m.browseData.Domains
					if m.browseDomain < len(domains) {
						subs := domains[m.browseDomain].Subdomains
						if m.browseSub < len(subs) {
							concepts := subs[m.browseSub].Concepts
							if idx >= 1 && idx <= len(concepts) {
								c := concepts[idx-1]
								if c.Unlocked {
									return m, loadConceptPractice(m.OnSelectConcept, c.ID)
								}
							}
						}
					}
				}
			}
			return m, nil
		}

	case ScreenDiagnostic:
		switch msg.String() {
		case "enter":
			answer := m.input.Value()
			if answer == "" {
				return m, nil
			}
			return m, submitDiagAnswer(m.OnDiagSubmit, m.diagQuestion.ConceptID, answer, 0)
		case "q":
			return m, tea.Quit
		default:
			if len(msg.String()) > 3 {
				return m, nil
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

	case ScreenDiagFeedback:
		switch msg.String() {
		case "enter", " ", "n":
			if m.diagFeedback.Done {
				return m, loadNextProblem(m.OnNext)
			}
			return m, startDiagnostic(m.OnDiagStart)
		case "q":
			return m, tea.Quit
		}

	case ScreenProgress:
		switch msg.String() {
		case "enter", " ", "n":
			return m, loadNextProblem(m.OnNext)
		case "w":
			m.screen = ScreenWelcome
			return m, loadWelcomeStats(m.OnWelcome)
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func loadNextProblem(fn func() SessionMsg) tea.Cmd {
	return func() tea.Msg {
		return fn()
	}
}

func submitAnswer(fn func(conceptID, answer string, elapsed float64) SubmitResultMsg, conceptID, answer string, elapsed float64) tea.Cmd {
	return func() tea.Msg {
		return fn(conceptID, answer, elapsed)
	}
}

func loadProgress(fn func() ProgressMsg) tea.Cmd {
	return func() tea.Msg {
		return fn()
	}
}

func loadWelcomeStats(fn func() WelcomeStatsMsg) tea.Cmd {
	return func() tea.Msg {
		return fn()
	}
}

func loadStudyLessons(fn func() []StudyLesson) tea.Cmd {
	return func() tea.Msg {
		return StudyListMsg{Lessons: fn()}
	}
}

func loadConceptTree(fn func() ConceptTreeMsg) tea.Cmd {
	return func() tea.Msg {
		return fn()
	}
}

func loadConceptPractice(fn func(conceptID string) SessionMsg, conceptID string) tea.Cmd {
	return func() tea.Msg {
		return fn(conceptID)
	}
}

func startDiagnostic(fn func() DiagQuestionMsg) tea.Cmd {
	return func() tea.Msg {
		return fn()
	}
}

func submitDiagAnswer(fn func(conceptID, answer string, elapsed float64) DiagFeedbackMsg, conceptID, answer string, elapsed float64) tea.Cmd {
	return func() tea.Msg {
		return fn(conceptID, answer, elapsed)
	}
}
