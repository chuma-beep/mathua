package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Screen int

const (
	ScreenWelcome Screen = iota
	ScreenProblem
	ScreenFeedback
	ScreenProgress
)

type TickMsg time.Time

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
	startTime  time.Time
	elapsed    float64
	roundCount int

	result   SubmitResultMsg
	progress ProgressMsg

	OnSubmit   func(conceptID, answer string, elapsed float64) SubmitResultMsg
	OnNext     func() SessionMsg
	OnProgress func() ProgressMsg
	OnWelcome  func() WelcomeStatsMsg
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
		m.startTime = time.Now()
		m.elapsed = 0
		m.input.Reset()
		m.input.Focus()
		return m, tickCmd()

	case TickMsg:
		if m.screen == ScreenProblem {
			m.elapsed = time.Since(m.startTime).Seconds()
			return m, tickCmd()
		}
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
		case "p":
			return m, loadProgress(m.OnProgress)
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
			elapsed := time.Since(m.startTime).Seconds()
			return m, submitAnswer(m.OnSubmit, m.session.ConceptID, answer, elapsed)
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

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/10, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
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
