package tui

import "github.com/charmbracelet/lipgloss"

const (
	colBg      = "#0b0f1a"
	colBg2     = "#111827"
	colBg3     = "#1a2235"
	colBorder  = "#1e2d45"
	colBorder2 = "#2a3f5f"
	colText    = "#e8e2d5"
	colMuted   = "#9fa8b4"
	colDim     = "#5a6577"
	colGold    = "#c8a96e"
	colGoldDim = "#a88848"
	colTeal    = "#4db8a0"
	colRed     = "#e05c5c"
	colGreen   = "#3fb950"
	colPurple  = "#a8a0f0"
)

var (
	HeaderStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colBg2)).
			Foreground(lipgloss.Color(colGold)).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(colBorder)).
			Padding(0, 2)

	LogoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colGold))

	HeaderDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colDim))

	ContentStyle = lipgloss.NewStyle().
			Padding(1, 3)

	FooterStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colBg2)).
			Foreground(lipgloss.Color(colDim)).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(colBorder)).
			Padding(0, 2)

	DomainStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colDim))

	ConceptIDStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colGold))

	QuestionBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color(colBorder)).
				Background(lipgloss.Color(colBg2)).
				Padding(1, 2).
				MarginBottom(1)

	QuestionLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colDim)).
				MarginBottom(1)

	QuestionTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colText))

	InputLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colGreen))

	InputPromptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colGreen))

	TimerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colMuted))

	TimerWarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colGold))

	TimerCriticalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colRed))

	MasteryBarFilledStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colGreen))

	MasteryBarEmptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colBorder2))

	MasteryLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colDim))

	MasteryValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colGold))

	StreakStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colGold))

	StreakZeroStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colDim))

	FeedbackCorrectBox = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color(colGreen)).
				Background(lipgloss.Color("#0d1f12")).
				Padding(1, 2).
				MarginBottom(1)

	FeedbackCorrectTitle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colGreen))

	FeedbackWrongBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(colRed)).
			Background(lipgloss.Color("#1f0d0d")).
			Padding(1, 2).
			MarginBottom(1)

	FeedbackWrongTitle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colRed))

	FeedbackYourAnswerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colMuted))

	FeedbackCorrectAnswerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colGreen))

	ExplanationStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colMuted)).
			MarginTop(1)

	TimingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colDim))

	TimingValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colText))

	ContinueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colDim)).
			MarginTop(1)

	MasteryAchievedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colGold)).
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color(colGold)).
				Padding(0, 2).
				MarginBottom(1)

	WelcomeTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colGold)).
				MarginBottom(1)

	WelcomeSubtitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colMuted)).
				MarginBottom(2)

	StatLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colDim))

	StatValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colGold))

	StatRowStyle = lipgloss.NewStyle().
			MarginBottom(1)

	DividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colBorder)).
			MarginTop(1).
			MarginBottom(1)

	DomainNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colMuted)).
			Width(18)

	DomainCountStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colDim))

	ProgressTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colGold)).
				MarginBottom(1)

	NextUnlockedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colTeal)).
			MarginTop(1)

	NextReviewStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colGold)).
			MarginTop(1)

	KeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colBorder2))

	KeyDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colDim))

	KeySepStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colBorder))
)

func masteryBar(mastery float64, width int) string {
	filled := int(mastery * float64(width))
	if filled > width {
		filled = width
	}
	bar := MasteryBarFilledStyle.Render(repeat("|", filled))
	bar += MasteryBarEmptyStyle.Render(repeat(".", width-filled))
	return bar
}

func streakDots(current, required int) string {
	result := ""
	for i := 0; i < required; i++ {
		if i < current {
			result += StreakStyle.Render("*")
		} else {
			result += StreakZeroStyle.Render(".")
		}
		result += " "
	}
	return result
}

func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
