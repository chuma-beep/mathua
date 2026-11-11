package notifications

import (
	"fmt"
	"strings"

	"github.com/chuma-beep/mathua/internal/storage"
)

// ReviewReminderText is the in-app and email body for a due-review nudge. It is
// plain text on purpose: the in-app feed renders it as text, and the mailer
// turns line breaks into paragraphs.
func ReviewReminderText(due int) (title, body string) {
	if due == 1 {
		return "One review is waiting", "You have one due review. A short session keeps it from piling up."
	}
	return "Reviews are waiting", fmt.Sprintf("You have %d due reviews. A short session keeps them from piling up.", due)
}

// StudyReminderText is the nudge for a learner's chosen study day. It names no
// count because the generator does not know what the learner intends to study,
// and a reminder that guesses wrong is worse than one that does not guess.
func StudyReminderText() (title, body string) {
	return "Today is a study day", "You chose today for study. A short session keeps your streak of work moving."
}

// XPMilestoneText marks an XP threshold. Only the highest crossed threshold is
// passed here, so the copy never lists the ones below it.
func XPMilestoneText(xp int) (title, body string) {
	return fmt.Sprintf("%d XP reached", xp), fmt.Sprintf("You have earned %d XP in total. The next milestone is further along the same path.", xp)
}

// MasteryMilestoneText marks a mastered-concept count. Concepts are the unit
// that matters, so the copy talks about them rather than about points.
func MasteryMilestoneText(mastered int) (title, body string) {
	if mastered == 1 {
		return "First concept mastered", "One concept is now mastered. The next one builds on what you just proved you can do."
	}
	return fmt.Sprintf("%d concepts mastered", mastered), fmt.Sprintf("You have now mastered %d concepts. Keep going at the pace that works for you.", mastered)
}

// InactivityText invites a learner back after an absence. It states the length
// of the absence because a learner who thinks it has been a few days deserves
// to know it has been longer.
func InactivityText(days int) (title, body string) {
	if days == 1 {
		return "Ready for a short session?", "It has been a day since your last session. One review is enough to pick the thread back up."
	}
	return "Ready for a short session?", fmt.Sprintf("It has been %d days since your last session. One review is enough to pick the thread back up.", days)
}

// ReminderEmail renders the email body for one delivery job. Every link is
// absolute, the direct action comes first, and the unsubscribe link is the
// final line so it is never hidden behind a preference page. The action label
// depends on the kind, because "Start reviewing" is wrong on a study reminder.
func ReminderEmail(job storage.NotificationDeliveryJob, actionLabel, actionURL, settingsURL, unsubscribeURL string) string {
	name := strings.TrimSpace(job.StudentName)
	greeting := "Hello,"
	if name != "" {
		greeting = "Hello " + name + ","
	}
	if strings.TrimSpace(actionLabel) == "" {
		actionLabel = "Open Mathua"
	}
	lines := []string{
		greeting,
		job.Body,
		actionLabel + ": " + actionURL,
		"Notification settings: " + settingsURL,
		"Unsubscribe from email notifications: " + unsubscribeURL,
	}
	return strings.Join(lines, "\r\n\r\n") + "\r\n"
}

// DigestEmail renders one message that lists every notification the learner
// has waiting. Each item names its own action so the reader can pick, and the
// shared settings and unsubscribe links stay at the end.
func DigestEmail(studentName string, jobs []storage.NotificationDeliveryJob, frontendURL, settingsURL, unsubscribeURL string) string {
	name := strings.TrimSpace(studentName)
	greeting := "Hello,"
	if name != "" {
		greeting = "Hello " + name + ","
	}
	lines := []string{
		greeting,
		fmt.Sprintf("You have %d notifications waiting.", len(jobs)),
	}
	for _, job := range jobs {
		title := strings.TrimSpace(job.Title)
		if title == "" {
			title = defaultDigestTitle(job.Kind)
		}
		lines = append(lines, title, job.Body,
			ReminderActionLabel(job.Kind)+": "+digestActionURL(job, frontendURL))
	}
	lines = append(lines,
		"Notification settings: "+settingsURL,
		"Unsubscribe from email notifications: "+unsubscribeURL,
	)
	return strings.Join(lines, "\r\n\r\n") + "\r\n"
}

func digestActionURL(job storage.NotificationDeliveryJob, frontendURL string) string {
	href := strings.TrimSpace(job.Href)
	if href == "" {
		href = "/review"
	}
	if frontendURL == "" {
		return href
	}
	return frontendURL + href
}

func defaultDigestTitle(kind string) string {
	switch kind {
	case ReviewKind:
		return "Reviews are due"
	case StudyKind:
		return "A study session is waiting"
	case MilestoneKind:
		return "A milestone reached"
	case InactivityKind:
		return "Ready for a short session?"
	default:
		return "Mathua update"
	}
}

// ReminderActionLabel maps a notification kind to the label shown before the
// action link. An unknown kind falls back to a neutral label rather than to a
// review label that would be wrong.
func ReminderActionLabel(kind string) string {
	switch kind {
	case ReviewKind:
		return "Start reviewing"
	case StudyKind:
		return "Start studying"
	case MilestoneKind:
		return "See your progress"
	case InactivityKind:
		return "Start a short session"
	default:
		return "Open Mathua"
	}
}
