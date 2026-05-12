package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/scoring"
	"github.com/chuma-beep/mathua/internal/storage"
)

type screen int

const (
	screenWelcome screen = iota
	screenPractice
	screenStats
)

type statsMsg struct {
	scores   *scoring.Scores
	progress map[string]*storage.ConceptProgress
}

type Model struct {
	eng  *engine.Engine
	repo storage.Repository

	current screen
	w, h    int

	nameInput textarea.Model

	student       *storage.Student
	session       *storage.Session
	question      *engine.Question
	answerInput   textarea.Model
	feedback      string
	fbStyle       string
	lessonVP      viewport.Model
	showingLesson bool

	statsMsg *statsMsg
	err      error
}

func New(eng *engine.Engine, repo storage.Repository) *Model {
	ni := textarea.New()
	ni.Placeholder = "Enter your name..."
	ni.Focus()
	ni.SetHeight(1)

	ai := textarea.New()
	ai.Placeholder = "Type your answer and press Enter..."
	ai.SetHeight(1)

	vp := viewport.New(60, 20)

	return &Model{
		eng:         eng,
		repo:        repo,
		current:     screenWelcome,
		nameInput:   ni,
		answerInput: ai,
		lessonVP:    vp,
	}
}

func (m *Model) Init() tea.Cmd { return textarea.Blink }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
		m.lessonVP.Width = msg.Width - 8
		m.lessonVP.Height = msg.Height - 10

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.current == screenStats || m.showingLesson {
			if msg.String() == "esc" || msg.String() == "q" {
				m.current = screenPractice
				m.showingLesson = false
				m.answerInput.Focus()
				return m, nil
			}
			if m.showingLesson {
				var cmd tea.Cmd
				m.lessonVP, cmd = m.lessonVP.Update(msg)
				return m, cmd
			}
			return m, nil
		}
		return m, m.handleKey(msg)

	case *engine.Question:
		m.question = msg
		m.feedback = ""
		m.fbStyle = ""
		m.current = screenPractice
		m.answerInput.Reset()
		m.answerInput.Focus()
		if msg == nil {
			m.feedback = "All concepts mastered! Come back tomorrow."
			m.fbStyle = "info"
		}

	case *engine.AnswerResult:
		if msg.Correct {
			m.feedback = fmt.Sprintf("Correct! Streak %d/%d (%.1fs) → %s",
				msg.Streak, msg.RequiredStreak, 0.0, msg.NewStatus)
			m.fbStyle = "ok"
		} else {
			m.feedback = fmt.Sprintf("Wrong — %s", msg.Explanation)
			m.fbStyle = "err"
		}

	case statsMsg:
		m.statsMsg = &msg

	case error:
		m.feedback = msg.Error()
		m.fbStyle = "err"
	}
	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) tea.Cmd {
	switch m.current {
	case screenWelcome:
		return m.welcomeKey(msg)
	case screenPractice:
		return m.practiceKey(msg)
	}
	return nil
}

func (m *Model) welcomeKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter":
		name := strings.TrimSpace(m.nameInput.Value())
		if name == "" {
			return nil
		}
		st, err := m.eng.CreateStudent(name)
		if err != nil {
			m.feedback = fmt.Sprintf("Error: %v", err)
			m.fbStyle = "err"
			return nil
		}
		m.student = st
		sess, err := m.repo.CreateSession(st.ID)
		if err != nil {
			m.feedback = fmt.Sprintf("Error: %v", err)
			m.fbStyle = "err"
			return nil
		}
		m.session = sess
		m.current = screenPractice
		return loadQuestion(m.eng, m.session.ID, m.student.ID)
	default:
		var cmd tea.Cmd
		m.nameInput, cmd = m.nameInput.Update(msg)
		return cmd
	}
}

func (m *Model) practiceKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter":
		answer := strings.TrimSpace(m.answerInput.Value())
		if answer == "" {
			return nil
		}
		m.answerInput.Reset()
		return submitAnswer(m.eng, m.session.ID, m.student.ID, answer)
	case "tab":
		m.current = screenStats
		return loadStats(m.eng, m.student.ID)
	case "l":
		if m.question != nil && m.question.Lesson != nil {
			m.showingLesson = true
			m.lessonVP.SetContent(m.question.Lesson.Body)
			m.lessonVP.GotoTop()
		}
	case "n":
		m.current = screenPractice
		m.feedback = ""
		return loadQuestion(m.eng, m.session.ID, m.student.ID)
	default:
		var cmd tea.Cmd
		m.answerInput, cmd = m.answerInput.Update(msg)
		return cmd
	}
	return nil
}

func (m *Model) View() string {
	switch m.current {
	case screenWelcome:
		return m.viewWelcome()
	case screenPractice:
		return m.viewPractice()
	case screenStats:
		return m.viewStats()
	}
	return ""
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)
	okStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	infoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	hintStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
)

func (m *Model) viewWelcome() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("=== Mathua ==="))
	b.WriteString("\nWelcome! What's your name?\n\n")
	b.WriteString(m.nameInput.View())
	b.WriteString("\n\n" + hintStyle.Render("Press Enter to start"))
	return center(m.w, m.h, b.String())
}

func (m *Model) viewPractice() string {
	if m.showingLesson {
		return m.viewLesson()
	}
	var b strings.Builder
	b.WriteString(titleStyle.Render("Mathua"))

	if m.question != nil {
		if m.question.Lesson != nil {
			b.WriteString(fmt.Sprintf(" [L]esson: %s", m.question.Lesson.Title))
		}
		b.WriteString(fmt.Sprintf("\n\n%s: %s", m.question.ConceptName, m.question.Question))
	}

	b.WriteString("\n\n" + m.answerInput.View())

	if m.feedback != "" {
		var style lipgloss.Style
		switch m.fbStyle {
		case "ok":
			style = okStyle
		case "err":
			style = errStyle
		default:
			style = infoStyle
		}
		b.WriteString("\n\n" + style.Render(m.feedback))
		b.WriteString("\n\n" + hintStyle.Render("Press 'n' for next question"))
	}

	b.WriteString("\n\n" + hintStyle.Render("Enter = Submit | Tab = Stats | Esc = Quit"))
	return b.String()
}

func (m *Model) viewLesson() string {
	var b strings.Builder
	if m.question != nil && m.question.Lesson != nil {
		b.WriteString(titleStyle.Render(m.question.Lesson.Title))
	}
	b.WriteString("\n" + m.lessonVP.View())
	b.WriteString("\n" + hintStyle.Render("Space/Enter = Back to practice"))
	return b.String()
}

func (m *Model) viewStats() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Stats"))
	if m.statsMsg != nil {
		s := m.statsMsg.scores
		b.WriteString(fmt.Sprintf("\nLevel: %s\n", s.Level))
		b.WriteString(fmt.Sprintf("Mastered: %d  Streak: %d  Score: %d\n",
			s.ConceptsMastered, s.CurrentStreak, s.WeeklyScore))
		if m.statsMsg.progress != nil {
			learning := 0
			for _, p := range m.statsMsg.progress {
				if p.Status == "LEARNING" || p.Status == "PRACTICING" {
					learning++
				}
			}
			b.WriteString(fmt.Sprintf("In Progress: %d  Unseen: %d\n",
				learning, 284-s.ConceptsMastered-learning))
		}
	}
	b.WriteString("\n" + hintStyle.Render("Tab/esc = Back to practice"))
	return center(m.w, m.h, b.String())
}

func center(w, h int, s string) string {
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, s)
}

func loadQuestion(eng *engine.Engine, sessionID, studentID string) tea.Cmd {
	return func() tea.Msg {
		q, err := eng.NextQuestion(sessionID, studentID)
		if err != nil {
			return err
		}
		return q
	}
}

func submitAnswer(eng *engine.Engine, sessionID, studentID, answer string) tea.Cmd {
	return func() tea.Msg {
		result, err := eng.SubmitAnswer(sessionID, studentID, answer, 0)
		if err != nil {
			return err
		}
		return result
	}
}

func loadStats(eng *engine.Engine, studentID string) tea.Cmd {
	return func() tea.Msg {
		scores, err := eng.GetScores(studentID)
		if err != nil {
			return err
		}
		prog, _ := eng.GetProgress(studentID)
		return statsMsg{scores, prog}
	}
}
