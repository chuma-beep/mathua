package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	var body string
	switch m.screen {
	case ScreenWelcome:
		body = m.viewWelcome()
	case ScreenProblem:
		body = m.viewProblem()
	case ScreenFeedback:
		body = m.viewFeedback()
	case ScreenProgress:
		body = m.viewProgress()
	}

	header := m.viewHeader()
	footer := m.viewFooter()

	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		body,
		footer,
	)
}

func (m Model) viewHeader() string {
	logo := LogoStyle.Render("Mathua")
	sep := HeaderDimStyle.Render("  |  ")
	stats := HeaderDimStyle.Render(
		fmt.Sprintf("%d/%d mastered", m.welcomeStats.TotalMastered, m.welcomeStats.TotalConcepts),
	)
	streak := ""
	if m.welcomeStats.DayStreak > 0 {
		streak = HeaderDimStyle.Render(fmt.Sprintf("  |  %dd streak", m.welcomeStats.DayStreak))
	}
	level := HeaderDimStyle.Render(fmt.Sprintf("  |  %s", m.welcomeStats.LevelName))

	return HeaderStyle.Width(m.width).Render(logo + sep + stats + streak + level)
}

func (m Model) viewFooter() string {
	var items [][2]string
	switch m.screen {
	case ScreenWelcome:
		items = [][2]string{{"enter", "start"}, {"p", "progress"}, {"q", "quit"}}
	case ScreenProblem:
		items = [][2]string{{"enter", "submit"}, {"ctrl+p", "progress"}, {"ctrl+c", "quit"}}
	case ScreenFeedback:
		items = [][2]string{{"enter", "next"}, {"p", "progress"}, {"q", "quit"}}
	case ScreenProgress:
		items = [][2]string{{"enter", "continue"}, {"w", "welcome"}, {"q", "quit"}}
	}
	var parts []string
	for i, kv := range items {
		p := KeyStyle.Render(kv[0]) + " " + KeyDescStyle.Render(kv[1])
		if i > 0 {
			p = KeySepStyle.Render("  |  ") + p
		}
		parts = append(parts, p)
	}
	return FooterStyle.Width(m.width).Render(strings.Join(parts, ""))
}

// Welcome screen

func (m Model) viewWelcome() string {
	s := m.welcomeStats
	w := m.width - 6

	var b strings.Builder
	b.WriteString(WelcomeTitleStyle.Render("Mathua"))
	b.WriteString("\n")
	b.WriteString(WelcomeSubtitleStyle.Render("Master the foundation. Earn the abstraction."))
	b.WriteString("\n\n")

	b.WriteString(DividerStyle.Render(strings.Repeat("-", w)))
	b.WriteString("\n\n")

	b.WriteString(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			StatRowStyle.Render(
				StatLabelStyle.Render("mastered  ")+
					StatValueStyle.Render(fmt.Sprintf("%d", s.TotalMastered))+
					StatLabelStyle.Render(fmt.Sprintf("/%d", s.TotalConcepts)),
			),
			StatRowStyle.MarginLeft(4).Render(
				StatLabelStyle.Render("streak  ")+
					StatValueStyle.Render(fmt.Sprintf("%dd", s.DayStreak)),
			),
			StatRowStyle.MarginLeft(4).Render(
				StatLabelStyle.Render("level  ")+
					StatValueStyle.Render(fmt.Sprintf("%s", s.LevelName)),
			),
			StatRowStyle.MarginLeft(4).Render(
				StatLabelStyle.Render("score  ")+
					StatValueStyle.Render(fmt.Sprintf("%d", s.WeeklyScore)),
			),
		),
	)
	b.WriteString("\n\n")

	b.WriteString(DividerStyle.Render(strings.Repeat("-", w)))
	b.WriteString("\n")
	b.WriteString(StatLabelStyle.Render("domain progress"))
	b.WriteString("\n\n")

	barWidth := 24
	for _, d := range s.DomainProgress {
		pct := 0.0
		if d.Total > 0 {
			pct = float64(d.Mastered) / float64(d.Total)
		}
		name := DomainNameStyle.Render(d.Domain)
		bar := masteryBar(pct, barWidth)
		count := DomainCountStyle.Render(
			fmt.Sprintf("  %d/%d", d.Mastered, d.Total),
		)
		b.WriteString(name + bar + count + "\n")
	}

	b.WriteString("\n")
	b.WriteString(DividerStyle.Render(strings.Repeat("-", w)))
	b.WriteString("\n\n")
	b.WriteString(ContinueStyle.Render("press enter to begin"))

	return ContentStyle.Render(b.String())
}

// Problem screen

func (m Model) viewProblem() string {
	s := m.session
	w := m.width - 10

	var b strings.Builder

	b.WriteString(
		DomainStyle.Render(s.Domain) +
			DomainStyle.Render("  |  ") +
			DomainStyle.Render(s.Subdomain) +
			DomainStyle.Render("  |  ") +
			ConceptIDStyle.Render(s.ConceptID),
	)
	b.WriteString("\n\n")

	bar := masteryBar(s.MasteryPct, 20)
	pctStr := MasteryValueStyle.Render(fmt.Sprintf("%.0f%%", s.MasteryPct*100))
	dots := streakDots(s.Streak, s.StreakNeeded)
	b.WriteString(
		MasteryLabelStyle.Render("mastery  ") +
			bar + "  " + pctStr +
			MasteryLabelStyle.Render("   streak  ") +
			dots,
	)
	b.WriteString("\n\n")

	questionContent := QuestionLabelStyle.Render(
		fmt.Sprintf("question %d", m.roundCount+1),
	) + "\n\n" +
		QuestionTextStyle.Render(s.Question)

	box := QuestionBoxStyle.Width(w).Render(questionContent)
	b.WriteString(box)
	b.WriteString("\n")

	timerStr := m.viewTimer(s.TimeLimit)
	b.WriteString(timerStr)
	b.WriteString("\n\n")

	b.WriteString(InputLabelStyle.Render("answer"))
	b.WriteString("\n")
	b.WriteString(m.input.View())

	return ContentStyle.Render(b.String())
}

func (m Model) viewTimer(timeLimit float64) string {
	elapsed := m.elapsed
	remaining := timeLimit - elapsed
	if remaining < 0 {
		remaining = 0
	}

	ratio := elapsed / timeLimit
	var style lipgloss.Style
	switch {
	case ratio > 0.9:
		style = TimerCriticalStyle
	case ratio > 0.7:
		style = TimerWarningStyle
	default:
		style = TimerStyle
	}

	return style.Render(fmt.Sprintf("%.1fs / %.0fs", elapsed, timeLimit))
}

// Feedback screen

func (m Model) viewFeedback() string {
	r := m.result
	s := m.session
	w := m.width - 10

	var b strings.Builder

	b.WriteString(
		ConceptIDStyle.Render(s.ConceptID) +
			DomainStyle.Render("  |  "+s.Domain),
	)
	b.WriteString("\n\n")

	var box string
	if r.Correct {
		content := FeedbackCorrectTitle.Render("correct") + "\n\n" +
			FeedbackYourAnswerStyle.Render("your answer  ") +
			FeedbackCorrectAnswerStyle.Render(r.UserAnswer)
		box = FeedbackCorrectBox.Width(w).Render(content)
	} else {
		content := FeedbackWrongTitle.Render("incorrect") + "\n\n" +
			FeedbackYourAnswerStyle.Render("your answer  ") +
			lipgloss.NewStyle().Foreground(lipgloss.Color(colRed)).Render(r.UserAnswer) + "\n" +
			FeedbackYourAnswerStyle.Render("correct      ") +
			FeedbackCorrectAnswerStyle.Render(r.CorrectAnswer)
		box = FeedbackWrongBox.Width(w).Render(content)
	}
	b.WriteString(box)
	b.WriteString("\n")

	b.WriteString(
		TimingStyle.Render("time  ") +
			TimingValueStyle.Render(fmt.Sprintf("%.1fs", r.ElapsedSecs)) +
			TimingStyle.Render("  |  new mastery  ") +
			MasteryValueStyle.Render(fmt.Sprintf("%.0f%%", r.NewMastery*100)),
	)
	b.WriteString("\n")

	if r.Correct {
		b.WriteString(
			TimingStyle.Render("streak  ") +
				streakDots(r.NewStreak, s.StreakNeeded),
		)
	} else {
		b.WriteString(
			TimingStyle.Render("streak reset  ") +
				StreakZeroStyle.Render("this concept will resurface soon"),
		)
	}
	b.WriteString("\n\n")

	if r.Explanation != "" {
		b.WriteString(ExplanationStyle.Render(r.Explanation))
		b.WriteString("\n\n")
	}

	if r.MasteryAchieved {
		msg := "concept mastered"
		if len(r.UnlockedConcepts) > 0 {
			msg += fmt.Sprintf("  |  unlocked: %s", strings.Join(r.UnlockedConcepts, ", "))
		}
		b.WriteString(MasteryAchievedStyle.Render(msg))
		b.WriteString("\n\n")
	}

	b.WriteString(ContinueStyle.Render("press enter for next  |  p for progress"))

	return ContentStyle.Render(b.String())
}

// Progress screen

func (m Model) viewProgress() string {
	p := m.progress
	w := m.width - 6

	var b strings.Builder

	b.WriteString(ProgressTitleStyle.Render("session checkpoint"))
	b.WriteString("\n")
	b.WriteString(DividerStyle.Render(strings.Repeat("-", w)))
	b.WriteString("\n\n")

	barWidth := 20
	domains := p.Domains
	sort.Slice(domains, func(i, j int) bool {
		return domains[i].Mastered > domains[j].Mastered
	})

	for _, d := range domains {
		pct := 0.0
		if d.Total > 0 {
			pct = float64(d.Mastered) / float64(d.Total)
		}

		name := DomainNameStyle.Render(d.Domain)
		bar := masteryBar(pct, barWidth)
		count := DomainCountStyle.Render(fmt.Sprintf("  %d/%d", d.Mastered, d.Total))
		b.WriteString(name + bar + count + "\n")

		for _, sd := range d.Subdomains {
			sdPct := 0.0
			if sd.Total > 0 {
				sdPct = float64(sd.Mastered) / float64(sd.Total)
			}
			sdName := lipgloss.NewStyle().Foreground(lipgloss.Color(colDim)).
				Width(20).Render("  " + sd.Name)
			sdBar := masteryBar(sdPct, 14)
			sdCount := DomainCountStyle.Render(fmt.Sprintf("  %d/%d", sd.Mastered, sd.Total))
			b.WriteString(sdName + sdBar + sdCount + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(DividerStyle.Render(strings.Repeat("-", w)))
	b.WriteString("\n\n")

	if p.NextReviewDue != "" {
		b.WriteString(
			StatLabelStyle.Render("next review  ") +
				NextReviewStyle.Render(p.NextReviewDue),
		)
		b.WriteString("\n")
	}

	if len(p.NewlyUnlocked) > 0 {
		b.WriteString(
			StatLabelStyle.Render("newly unlocked  ") +
				NextUnlockedStyle.Render(strings.Join(p.NewlyUnlocked, "  |  ")),
		)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(ContinueStyle.Render("press enter to continue  |  w for welcome"))

	return ContentStyle.Render(b.String())
}
