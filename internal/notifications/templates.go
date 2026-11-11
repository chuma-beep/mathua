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

// ReviewReminderEmail renders the email body for one delivery job. Every link
// is absolute, the direct action comes first, and the unsubscribe link is the
// final line so it is never hidden behind a preference page.
func ReviewReminderEmail(job storage.NotificationDeliveryJob, reviewURL, settingsURL, unsubscribeURL string) string {
	name := strings.TrimSpace(job.StudentName)
	greeting := "Hello,"
	if name != "" {
		greeting = "Hello " + name + ","
	}
	lines := []string{
		greeting,
		job.Body,
		"Start reviewing: " + reviewURL,
		"Notification settings: " + settingsURL,
		"Unsubscribe from email notifications: " + unsubscribeURL,
	}
	return strings.Join(lines, "\r\n\r\n") + "\r\n"
}
