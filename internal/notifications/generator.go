package notifications

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

const (
	// ReviewKind is the only notification kind Phase 3 generates.
	ReviewKind = "review"
	// reviewScanKind is the notification_scan_state row that records the last
	// local day evaluated for review reminders.
	reviewScanKind = "review"
)

// DueCounter is the read-only review count the generator is allowed to use.
// Engine.ReviewsDue satisfies this, but the generator is passed
// ReviewsDueReadOnly so generation can never write learning state.
type DueCounter interface {
	ReviewsDueReadOnly(studentID string) (int, error)
}

// GeneratorOptions are the production dependencies and the test clock.
type GeneratorOptions struct {
	Store          storage.NotificationStore
	Due            DueCounter
	EmailAvailable bool
	Now            func() time.Time
	PageSize       int
}

// Generator creates one review reminder per learner per local day when the
// learner's preferences say the reminder window is open.
type Generator struct {
	store          storage.NotificationStore
	due            DueCounter
	emailAvailable bool
	now            func() time.Time
	pageSize       int
}

// NewGenerator builds a review-reminder generator. emailAvailable is false when
// no mail provider is configured, so the delivery row records a suppression
// instead of retrying a send that cannot succeed.
func NewGenerator(opts GeneratorOptions) *Generator {
	now := opts.Now
	if now == nil {
		now = time.Now
	}
	pageSize := opts.PageSize
	if pageSize <= 0 || pageSize > 500 {
		pageSize = 100
	}
	return &Generator{
		store:          opts.Store,
		due:            opts.Due,
		emailAvailable: opts.EmailAvailable,
		now:            now,
		pageSize:       pageSize,
	}
}

// GenerationResult reports what one scan changed.
type GenerationResult struct {
	Evaluated        int
	Notifications    int
	EmailPending     int
	EmailSuppressed  int
	AlreadyGenerated int
	OutsideWindow    int
	NoReviews        int
	NoChannel        int
	Failed           int
}

// RunOnce walks every notification preference row in student-id order. The
// unique (student_id, dedupe_key) index is the final concurrency guard; the
// scan watermark only avoids re-counting reviews within the same local day.
func (g *Generator) RunOnce() (GenerationResult, error) {
	var result GenerationResult
	if g.store == nil || g.due == nil {
		return result, fmt.Errorf("notification generator is not configured")
	}
	after := ""
	for {
		recipients, err := g.store.ListNotificationRecipients(after, g.pageSize)
		if err != nil {
			return result, err
		}
		for _, recipient := range recipients {
			result.Evaluated++
			g.generateOne(recipient, &result)
		}
		if len(recipients) < g.pageSize {
			break
		}
		after = recipients[len(recipients)-1].StudentID
	}
	return result, nil
}

func (g *Generator) generateOne(recipient storage.NotificationRecipient, result *GenerationResult) {
	prefs := recipient.Preferences
	if !prefs.ReviewsEnabled {
		return
	}
	emailDeliverable := prefs.EmailEnabled && g.emailAvailable &&
		strings.TrimSpace(recipient.Email) != "" && recipient.EmailVerified
	if !prefs.InAppEnabled && !emailDeliverable {
		result.NoChannel++
		return
	}
	now := g.now().UTC()
	if prefs.PausedUntil != nil && prefs.PausedUntil.After(now) {
		result.OutsideWindow++
		return
	}

	loc := time.UTC
	if loaded, err := time.LoadLocation(prefs.Timezone); err == nil {
		loc = loaded
	}
	local := now.In(loc)
	minute := local.Hour()*60 + local.Minute()
	if minute < prefs.PreferredTimeMin {
		result.OutsideWindow++
		return
	}
	if inQuietHours(minute, prefs.QuietStartMin, prefs.QuietEndMin) {
		result.OutsideWindow++
		return
	}

	dayKey := localDateKey(local)
	watermark, err := g.store.GetNotificationScanState(recipient.StudentID, reviewScanKind)
	if err != nil {
		log.Printf("notifications: scan state for %s: %v", recipient.StudentID, err)
		result.Failed++
		return
	}
	if watermark >= dayKey {
		result.AlreadyGenerated++
		return
	}

	due, err := g.due.ReviewsDueReadOnly(recipient.StudentID)
	if err != nil {
		log.Printf("notifications: reviews due for %s: %v", recipient.StudentID, err)
		result.Failed++
		return
	}
	// The watermark records that this local day was evaluated even when the
	// count is zero, so a quiet learner is not re-counted every worker tick.
	if due <= 0 {
		if err := g.store.SetNotificationScanState(recipient.StudentID, reviewScanKind, dayKey, now); err != nil {
			log.Printf("notifications: set scan state for %s: %v", recipient.StudentID, err)
			result.Failed++
			return
		}
		result.NoReviews++
		return
	}

	title, body := ReviewReminderText(due)
	id, err := newNotificationID()
	if err != nil {
		log.Printf("notifications: id for %s: %v", recipient.StudentID, err)
		result.Failed++
		return
	}
	created, err := g.store.CreateNotification(&storage.Notification{
		ID:        id,
		StudentID: recipient.StudentID,
		Kind:      ReviewKind,
		DedupeKey: ReviewDedupeKey(local),
		Title:     title,
		Body:      body,
		Href:      "/review",
		CreatedAt: now,
	})
	if err != nil {
		log.Printf("notifications: create review reminder for %s: %v", recipient.StudentID, err)
		result.Failed++
		return
	}
	if !created {
		if err := g.store.SetNotificationScanState(recipient.StudentID, reviewScanKind, dayKey, now); err != nil {
			log.Printf("notifications: set scan state for %s: %v", recipient.StudentID, err)
		}
		result.AlreadyGenerated++
		return
	}
	result.Notifications++

	if prefs.InAppEnabled {
		if _, err := g.store.CreateNotificationDelivery(&storage.NotificationDelivery{
			NotificationID: id,
			Channel:        storage.NotificationChannelInApp,
			Status:         storage.NotificationDeliverySent,
			SentAt:         now,
			ProviderID:     "app",
		}); err != nil {
			log.Printf("notifications: record in-app delivery for %s: %v", recipient.StudentID, err)
			result.Failed++
		}
	}
	if prefs.EmailEnabled {
		delivery := &storage.NotificationDelivery{
			NotificationID: id,
			Channel:        storage.NotificationChannelEmail,
			Status:         storage.NotificationDeliveryPending,
			NextAttemptAt:  now,
		}
		switch {
		case !g.emailAvailable:
			delivery.Status = storage.NotificationDeliverySuppressed
			delivery.SuppressedReason = "email_not_configured"
		case recipient.Email == "":
			delivery.Status = storage.NotificationDeliverySuppressed
			delivery.SuppressedReason = "email_missing"
		case !recipient.EmailVerified:
			delivery.Status = storage.NotificationDeliverySuppressed
			delivery.SuppressedReason = "email_unverified"
		}
		if _, err := g.store.CreateNotificationDelivery(delivery); err != nil {
			log.Printf("notifications: record email delivery for %s: %v", recipient.StudentID, err)
			result.Failed++
		} else if delivery.Status == storage.NotificationDeliveryPending {
			result.EmailPending++
		} else {
			result.EmailSuppressed++
		}
	}

	if err := g.store.SetNotificationScanState(recipient.StudentID, reviewScanKind, dayKey, now); err != nil {
		log.Printf("notifications: set scan state for %s: %v", recipient.StudentID, err)
		result.Failed++
	}
}

// ReviewDedupeKey is the per-learner, per-local-day guard. Its local date is
// deliberate: a learner who travels still gets at most one reminder per day in
// the timezone they chose.
func ReviewDedupeKey(local time.Time) string {
	return "review:" + local.Format("2006-01-02")
}

func localDateKey(local time.Time) int64 {
	return int64(local.Year())*10000 + int64(local.Month())*100 + int64(local.Day())
}

// inQuietHours handles both an overnight window (22:00 to 07:00) and a
// same-day window. Equal start and end means no quiet period.
func inQuietHours(minute, start, end int) bool {
	if start == end {
		return false
	}
	if start < end {
		return minute >= start && minute < end
	}
	return minute >= start || minute < end
}

func newNotificationID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "ntf_" + hex.EncodeToString(b), nil
}
