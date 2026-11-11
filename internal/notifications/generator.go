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
	// ReviewKind reminds a learner that scheduled reviews are due.
	ReviewKind = "review"
	// StudyKind reminds a learner on one of their chosen study days.
	StudyKind = "study"
	// MilestoneKind marks an XP or mastery threshold.
	MilestoneKind = "milestone"
	// InactivityKind nudges a learner who has not studied for a while.
	InactivityKind = "inactivity"

	reviewScanKind     = "review"
	studyScanKind      = "study"
	milestoneScanKind  = "milestone"
	inactivityScanKind = "inactivity"
)

// inactivityDays is both the absence that earns one nudge and the minimum gap
// between nudges. A learner away for a month gets at most four messages rather
// than one per scan.
const inactivityDays = 7

// Milestone thresholds. Only the highest threshold a learner has crossed is
// announced, so a learner who arrives with 12,000 XP receives one message, not
// every threshold below it.
var (
	xpMilestones      = []int{100, 500, 1000, 2500, 5000, 10000, 25000, 50000}
	masteryMilestones = []int{1, 5, 10, 25, 50, 100, 250}
)

// DueCounter is the read-only review count the generator is allowed to use.
// Engine.ReviewsDue satisfies this, but the generator is passed
// ReviewsDueReadOnly so generation can never write learning state.
type DueCounter interface {
	ReviewsDueReadOnly(studentID string) (int, error)
}

// GeneratorOptions are the production dependencies and the test clock.
// EmailAvailable and PushAvailable are false when the deployment has no
// provider, so a delivery row records a suppression instead of retrying a send
// that cannot succeed.
type GeneratorOptions struct {
	Store          storage.NotificationStore
	Due            DueCounter
	EmailAvailable bool
	PushAvailable  bool
	Now            func() time.Time
	PageSize       int
}

// Generator creates at most one notification per kind per learner per local
// day, and only when the learner's preferences say the window is open.
type Generator struct {
	store          storage.NotificationStore
	due            DueCounter
	emailAvailable bool
	pushAvailable  bool
	now            func() time.Time
	pageSize       int
}

// NewGenerator builds the notification generator.
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
		pushAvailable:  opts.PushAvailable,
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
	PushPending      int
	PushSuppressed   int
	AlreadyGenerated int
	OutsideWindow    int
	NoReviews        int
	StudyReminders   int
	Milestones       int
	Inactivity       int
	NoChannel        int
	Failed           int
}

// RunOnce walks every notification preference row in student-id order. The
// unique (student_id, dedupe_key) index is the final concurrency guard; the
// scan watermark only avoids re-counting within the same local day.
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

// channels is the deliverable set for one learner. Email requires a verified
// address and a configured provider; push requires a configured VAPID key pair.
type channels struct {
	inApp bool
	email bool
	push  bool
}

func (c channels) any() bool { return c.inApp || c.email || c.push }

func (g *Generator) channelsFor(r storage.NotificationRecipient) channels {
	prefs := r.Preferences
	return channels{
		inApp: prefs.InAppEnabled,
		email: prefs.EmailEnabled && g.emailAvailable &&
			strings.TrimSpace(r.Email) != "" && r.EmailVerified,
		push: prefs.PushEnabled && g.pushAvailable,
	}
}

func (g *Generator) generateOne(recipient storage.NotificationRecipient, result *GenerationResult) {
	prefs := recipient.Preferences
	now := g.now().UTC()
	if prefs.PausedUntil != nil && prefs.PausedUntil.After(now) {
		result.OutsideWindow++
		return
	}
	ch := g.channelsFor(recipient)
	if !ch.any() {
		result.NoChannel++
		return
	}

	loc := time.UTC
	if loaded, err := time.LoadLocation(prefs.Timezone); err == nil {
		loc = loaded
	}
	local := now.In(loc)
	minute := local.Hour()*60 + local.Minute()

	if prefs.ReviewsEnabled {
		g.reviewReminder(recipient, ch, local, minute, now, result)
	}
	if prefs.StudyEnabled {
		g.studyReminder(recipient, ch, local, minute, now, result)
	}
	if prefs.MilestonesEnabled {
		g.milestoneReminder(recipient, ch, local, minute, now, result)
	}
	if prefs.InactivityEnabled {
		g.inactivityNudge(recipient, ch, local, minute, now, result)
	}
}

// emit creates the notification and one delivery row per channel the learner
// enabled. A suppressed row is deliberate: it records why a channel the learner
// asked for produced nothing, without blocking the channels that worked.
func (g *Generator) emit(r storage.NotificationRecipient, ch channels, kind, dedupeKey, title, body, href string, now time.Time, result *GenerationResult) bool {
	prefs := r.Preferences
	id, err := newNotificationID()
	if err != nil {
		log.Printf("notifications: id for %s: %v", r.StudentID, err)
		result.Failed++
		return false
	}
	created, err := g.store.CreateNotification(&storage.Notification{
		ID:        id,
		StudentID: r.StudentID,
		Kind:      kind,
		DedupeKey: dedupeKey,
		Title:     title,
		Body:      body,
		Href:      href,
		CreatedAt: now,
	})
	if err != nil {
		log.Printf("notifications: create %s for %s: %v", kind, r.StudentID, err)
		result.Failed++
		return false
	}
	if !created {
		result.AlreadyGenerated++
		return false
	}
	result.Notifications++

	nextAttempt := now
	if prefs.Digest != storage.NotificationDigestImmediate {
		loc := time.UTC
		if loaded, err := time.LoadLocation(prefs.Timezone); err == nil {
			loc = loaded
		}
		nextAttempt = nextDigestTime(now.In(loc), prefs.Digest, prefs.PreferredTimeMin).UTC()
	}

	if ch.inApp {
		if _, err := g.store.CreateNotificationDelivery(&storage.NotificationDelivery{
			NotificationID: id,
			Channel:        storage.NotificationChannelInApp,
			Status:         storage.NotificationDeliverySent,
			SentAt:         now,
			ProviderID:     "app",
		}); err != nil {
			log.Printf("notifications: record in-app delivery for %s: %v", r.StudentID, err)
			result.Failed++
		}
	}
	if prefs.EmailEnabled {
		delivery := &storage.NotificationDelivery{
			NotificationID: id,
			Channel:        storage.NotificationChannelEmail,
			Status:         storage.NotificationDeliveryPending,
			NextAttemptAt:  nextAttempt,
		}
		switch {
		case !g.emailAvailable:
			delivery.Status = storage.NotificationDeliverySuppressed
			delivery.SuppressedReason = "email_not_configured"
		case strings.TrimSpace(r.Email) == "":
			delivery.Status = storage.NotificationDeliverySuppressed
			delivery.SuppressedReason = "email_missing"
		case !r.EmailVerified:
			delivery.Status = storage.NotificationDeliverySuppressed
			delivery.SuppressedReason = "email_unverified"
		}
		if _, err := g.store.CreateNotificationDelivery(delivery); err != nil {
			log.Printf("notifications: record email delivery for %s: %v", r.StudentID, err)
			result.Failed++
		} else if delivery.Status == storage.NotificationDeliveryPending {
			result.EmailPending++
		} else {
			result.EmailSuppressed++
		}
	}
	if prefs.PushEnabled {
		delivery := &storage.NotificationDelivery{
			NotificationID: id,
			Channel:        storage.NotificationChannelPush,
			Status:         storage.NotificationDeliveryPending,
			NextAttemptAt:  nextAttempt,
		}
		if !g.pushAvailable {
			delivery.Status = storage.NotificationDeliverySuppressed
			delivery.SuppressedReason = "push_not_configured"
		}
		if _, err := g.store.CreateNotificationDelivery(delivery); err != nil {
			log.Printf("notifications: record push delivery for %s: %v", r.StudentID, err)
			result.Failed++
		} else if delivery.Status == storage.NotificationDeliveryPending {
			result.PushPending++
		} else {
			result.PushSuppressed++
		}
	}
	return true
}

func (g *Generator) reviewReminder(r storage.NotificationRecipient, ch channels, local time.Time, minute int, now time.Time, result *GenerationResult) {
	prefs := r.Preferences
	if minute < prefs.PreferredTimeMin || inQuietHours(minute, prefs.QuietStartMin, prefs.QuietEndMin) {
		result.OutsideWindow++
		return
	}
	dayKey := localDateKey(local)
	watermark, err := g.store.GetNotificationScanState(r.StudentID, reviewScanKind)
	if err != nil {
		log.Printf("notifications: scan state for %s: %v", r.StudentID, err)
		result.Failed++
		return
	}
	if watermark >= dayKey {
		result.AlreadyGenerated++
		return
	}
	due, err := g.due.ReviewsDueReadOnly(r.StudentID)
	if err != nil {
		log.Printf("notifications: reviews due for %s: %v", r.StudentID, err)
		result.Failed++
		return
	}
	// The watermark records that this local day was evaluated even when the
	// count is zero, so a quiet learner is not re-counted every worker tick.
	if due <= 0 {
		g.setScan(r.StudentID, reviewScanKind, dayKey, now, result)
		result.NoReviews++
		return
	}
	title, body := ReviewReminderText(due)
	g.emit(r, ch, ReviewKind, ReviewDedupeKey(local), title, body, "/review", now, result)
	g.setScan(r.StudentID, reviewScanKind, dayKey, now, result)
}

func (g *Generator) studyReminder(r storage.NotificationRecipient, ch channels, local time.Time, minute int, now time.Time, result *GenerationResult) {
	prefs := r.Preferences
	if !studyDayEnabled(prefs.StudyDays, local.Weekday()) {
		return
	}
	if minute < prefs.StudyTimeMin || inQuietHours(minute, prefs.QuietStartMin, prefs.QuietEndMin) {
		result.OutsideWindow++
		return
	}
	dayKey := localDateKey(local)
	watermark, err := g.store.GetNotificationScanState(r.StudentID, studyScanKind)
	if err != nil {
		log.Printf("notifications: study scan state for %s: %v", r.StudentID, err)
		result.Failed++
		return
	}
	if watermark >= dayKey {
		result.AlreadyGenerated++
		return
	}
	title, body := StudyReminderText()
	if g.emit(r, ch, StudyKind, "study:"+local.Format("2006-01-02"), title, body, "/study", now, result) {
		result.StudyReminders++
	}
	g.setScan(r.StudentID, studyScanKind, dayKey, now, result)
}

func (g *Generator) milestoneReminder(r storage.NotificationRecipient, ch channels, local time.Time, minute int, now time.Time, result *GenerationResult) {
	prefs := r.Preferences
	if minute < prefs.PreferredTimeMin || inQuietHours(minute, prefs.QuietStartMin, prefs.QuietEndMin) {
		return
	}
	xp := highestReached(xpMilestones, r.XPTotal)
	mastered := highestReached(masteryMilestones, r.MasteredCount)
	// Nothing has been reached yet, so there is nothing to remember and no
	// reason to record a watermark that would hide a crossing later today.
	if xp == 0 && mastered == 0 {
		return
	}
	dayKey := localDateKey(local)
	watermark, err := g.store.GetNotificationScanState(r.StudentID, milestoneScanKind)
	if err != nil {
		log.Printf("notifications: milestone scan state for %s: %v", r.StudentID, err)
		result.Failed++
		return
	}
	if watermark >= dayKey {
		return
	}
	// The dedupe key is the real guard: the highest crossed threshold is
	// announced once, and a learner who arrives with a large total gets one
	// message rather than every threshold below it.
	if xp > 0 {
		title, body := XPMilestoneText(xp)
		if g.emit(r, ch, MilestoneKind, fmt.Sprintf("milestone:xp:%d", xp), title, body, "/profile", now, result) {
			result.Milestones++
		}
	}
	if mastered > 0 {
		title, body := MasteryMilestoneText(mastered)
		if g.emit(r, ch, MilestoneKind, fmt.Sprintf("milestone:mastered:%d", mastered), title, body, "/profile", now, result) {
			result.Milestones++
		}
	}
	g.setScan(r.StudentID, milestoneScanKind, dayKey, now, result)
}

func (g *Generator) inactivityNudge(r storage.NotificationRecipient, ch channels, local time.Time, minute int, now time.Time, result *GenerationResult) {
	prefs := r.Preferences
	if minute < prefs.PreferredTimeMin || inQuietHours(minute, prefs.QuietStartMin, prefs.QuietEndMin) {
		return
	}
	last := lastActivity(r, local)
	if last.IsZero() {
		return
	}
	today := localDayNumber(local)
	away := today - localDayNumber(last)
	if away < inactivityDays {
		return
	}
	watermark, err := g.store.GetNotificationScanState(r.StudentID, inactivityScanKind)
	if err != nil {
		log.Printf("notifications: inactivity scan state for %s: %v", r.StudentID, err)
		result.Failed++
		return
	}
	if watermark > 0 && today-watermark < inactivityDays {
		return
	}
	title, body := InactivityText(int(away))
	if g.emit(r, ch, InactivityKind, "inactivity:"+local.Format("2006-01-02"), title, body, "/study", now, result) {
		result.Inactivity++
	}
	g.setScan(r.StudentID, inactivityScanKind, today, now, result)
}

func (g *Generator) setScan(studentID, kind string, watermark int64, now time.Time, result *GenerationResult) {
	if err := g.store.SetNotificationScanState(studentID, kind, watermark, now); err != nil {
		log.Printf("notifications: set %s scan state for %s: %v", kind, studentID, err)
		result.Failed++
	}
}

// lastActivity prefers the last day XP was awarded, then the account creation
// date, then the preference row's own creation. A learner who signed up
// yesterday is never nudged for inactivity.
func lastActivity(r storage.NotificationRecipient, local time.Time) time.Time {
	if t, err := time.Parse("2006-01-02", strings.TrimSpace(r.XPDate)); err == nil {
		return t
	}
	if !r.AccountCreated.IsZero() {
		return r.AccountCreated.In(local.Location())
	}
	if !r.Preferences.CreatedAt.IsZero() {
		return r.Preferences.CreatedAt.In(local.Location())
	}
	return time.Time{}
}

// localDayNumber is the local calendar day as a day count, which makes the
// "at least seven days apart" rule simple subtraction across month and year
// boundaries.
func localDayNumber(t time.Time) int64 {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix() / 86400
}

func highestReached(thresholds []int, value int) int {
	highest := 0
	for _, t := range thresholds {
		if value >= t {
			highest = t
		}
	}
	return highest
}

// nextDigestTime is the first preferred-time slot the digest may send in. Daily
// means the next occurrence of the preferred minute, weekly means the next
// Monday at that minute, so a learner who writes three notes on Tuesday gets
// one message rather than three.
func nextDigestTime(local time.Time, digest string, preferredMin int) time.Time {
	at := time.Date(local.Year(), local.Month(), local.Day(),
		preferredMin/60, preferredMin%60, 0, 0, local.Location())
	if digest == storage.NotificationDigestWeekly {
		at = at.AddDate(0, 0, (8-int(local.Weekday()))%7)
	}
	if !at.After(local) {
		if digest == storage.NotificationDigestWeekly {
			at = at.AddDate(0, 0, 7)
		} else {
			at = at.AddDate(0, 0, 1)
		}
	}
	return at
}

// studyDayEnabled matches the canonical comma-separated day list against a
// local weekday.
func studyDayEnabled(days string, weekday time.Weekday) bool {
	name := studyDayOrder[(int(weekday)+6)%7]
	for _, day := range strings.Split(days, ",") {
		if strings.TrimSpace(day) == name {
			return true
		}
	}
	return false
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
