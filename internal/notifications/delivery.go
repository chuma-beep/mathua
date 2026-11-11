package notifications

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/mailer"
	"github.com/chuma-beep/mathua/internal/storage"
)

const (
	deliveryLease      = 5 * time.Minute
	deliveryBatchSize  = 20
	deliveryMaxAttempt = 5
	deliveryBackoff    = 5 * time.Minute
	deliveryMaxBackoff = 6 * time.Hour
)

// DeliveryOptions are the delivery worker dependencies. UnsubscribeToken is
// injected so this package does not import authentication, and so tests can
// use a deterministic signer.
type DeliveryOptions struct {
	Store            storage.NotificationStore
	Sender           mailer.Sender
	FrontendURL      string
	UnsubscribeToken func(studentID string) (string, error)
	Now              func() time.Time
	Lease            time.Duration
	BatchSize        int
	MaxAttempts      int
}

// Delivery sends due email deliveries. In-app items are already visible in the
// feed and are recorded as sent by the generator, so only email is claimed.
type Delivery struct {
	store            storage.NotificationStore
	sender           mailer.Sender
	frontendURL      string
	unsubscribeToken func(string) (string, error)
	now              func() time.Time
	lease            time.Duration
	batchSize        int
	maxAttempts      int
}

// NewDelivery builds the email worker. The defaults match the Phase 3 plan:
// a twenty-item batch, a five-minute lease and five attempts.
func NewDelivery(opts DeliveryOptions) *Delivery {
	now := opts.Now
	if now == nil {
		now = time.Now
	}
	lease := opts.Lease
	if lease <= 0 {
		lease = deliveryLease
	}
	batchSize := opts.BatchSize
	if batchSize <= 0 || batchSize > 100 {
		batchSize = deliveryBatchSize
	}
	maxAttempts := opts.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = deliveryMaxAttempt
	}
	return &Delivery{
		store:            opts.Store,
		sender:           opts.Sender,
		frontendURL:      strings.TrimSuffix(strings.TrimSpace(opts.FrontendURL), "/"),
		unsubscribeToken: opts.UnsubscribeToken,
		now:              now,
		lease:            lease,
		batchSize:        batchSize,
		maxAttempts:      maxAttempts,
	}
}

// DeliveryResult reports what one claim pass did.
type DeliveryResult struct {
	Claimed    int
	Sent       int
	Retried    int
	Failed     int
	Suppressed int
}

// RunOnce claims one batch and handles each job. A single failed job never
// stops the batch: its state is recorded and the next job is attempted.
func (d *Delivery) RunOnce() (DeliveryResult, error) {
	var result DeliveryResult
	if d.store == nil || d.sender == nil {
		return result, fmt.Errorf("notification delivery is not configured")
	}
	now := d.now().UTC()
	jobs, err := d.store.ClaimNotificationDeliveries(now, d.lease, d.batchSize)
	if err != nil {
		return result, err
	}
	result.Claimed = len(jobs)
	for _, job := range jobs {
		d.deliver(job, now, &result)
	}
	return result, nil
}

func (d *Delivery) deliver(job storage.NotificationDeliveryJob, now time.Time, result *DeliveryResult) {
	prefs, err := d.store.GetNotificationPreferences(job.StudentID)
	if err != nil {
		log.Printf("notifications: delivery preferences for %s: %v", job.StudentID, err)
		d.retryOrFail(job, now, err.Error(), result)
		return
	}
	if prefs == nil || !prefs.EmailEnabled {
		d.suppress(job, "email_disabled", now, result)
		return
	}
	if prefs.PausedUntil != nil && prefs.PausedUntil.After(now) {
		// A reminder queued before a pause must not arrive afterwards. The next
		// eligible day will create a fresh notification.
		d.suppress(job, "paused", now, result)
		return
	}
	if strings.TrimSpace(job.Email) == "" {
		d.suppress(job, "email_missing", now, result)
		return
	}
	if !job.EmailVerified {
		d.suppress(job, "email_unverified", now, result)
		return
	}
	if d.frontendURL == "" {
		d.suppress(job, "frontend_not_configured", now, result)
		return
	}
	if d.unsubscribeToken == nil {
		d.suppress(job, "unsubscribe_not_configured", now, result)
		return
	}

	loc := time.UTC
	if loaded, err := time.LoadLocation(prefs.Timezone); err == nil {
		loc = loaded
	}
	local := now.In(loc)
	minute := local.Hour()*60 + local.Minute()
	if minute < prefs.PreferredTimeMin || inQuietHours(minute, prefs.QuietStartMin, prefs.QuietEndMin) {
		next := nextDeliveryWindow(local, prefs.PreferredTimeMin, prefs.QuietStartMin, prefs.QuietEndMin)
		if err := d.store.RetryNotificationDelivery(job.NotificationID, job.Channel, "outside delivery window", next.UTC()); err != nil {
			log.Printf("notifications: defer delivery %s: %v", job.NotificationID, err)
			result.Failed++
			return
		}
		result.Retried++
		return
	}

	token, err := d.unsubscribeToken(job.StudentID)
	if err != nil {
		d.retryOrFail(job, now, err.Error(), result)
		return
	}
	body := ReviewReminderEmail(job,
		d.frontendURL+"/review",
		d.frontendURL+"/settings#notifications",
		d.frontendURL+"/notifications/unsubscribe?token="+token)
	subject := strings.TrimSpace(job.Title)
	if subject == "" {
		subject = "Mathua review reminder"
	}
	if err := d.sender.Send(mailer.Message{To: job.Email, Subject: subject, Text: body}); err != nil {
		d.retryOrFail(job, now, err.Error(), result)
		return
	}
	if err := d.store.CompleteNotificationDelivery(job.NotificationID, job.Channel, "mailer", now); err != nil {
		log.Printf("notifications: complete delivery %s: %v", job.NotificationID, err)
		result.Failed++
		return
	}
	result.Sent++
}

func (d *Delivery) retryOrFail(job storage.NotificationDeliveryJob, now time.Time, reason string, result *DeliveryResult) {
	if job.Attempts >= d.maxAttempts {
		if err := d.store.FailNotificationDelivery(job.NotificationID, job.Channel, reason, now); err != nil {
			log.Printf("notifications: fail delivery %s: %v", job.NotificationID, err)
		}
		result.Failed++
		return
	}
	next := now.Add(retryDelay(job.Attempts))
	if err := d.store.RetryNotificationDelivery(job.NotificationID, job.Channel, reason, next); err != nil {
		log.Printf("notifications: retry delivery %s: %v", job.NotificationID, err)
		result.Failed++
		return
	}
	result.Retried++
}

func (d *Delivery) suppress(job storage.NotificationDeliveryJob, reason string, now time.Time, result *DeliveryResult) {
	if err := d.store.SuppressNotificationDelivery(job.NotificationID, job.Channel, reason, now); err != nil {
		log.Printf("notifications: suppress delivery %s: %v", job.NotificationID, err)
		result.Failed++
		return
	}
	result.Suppressed++
}

// retryDelay grows 5m, 10m, 20m, 40m up to six hours.
func retryDelay(attempts int) time.Duration {
	if attempts < 1 {
		attempts = 1
	}
	delay := deliveryBackoff
	for i := 1; i < attempts; i++ {
		delay *= 2
		if delay >= deliveryMaxBackoff {
			return deliveryMaxBackoff
		}
	}
	return delay
}

// nextDeliveryWindow returns the next instant at or after the preferred time
// and outside quiet hours. Equal quiet start and end means no quiet period.
func nextDeliveryWindow(local time.Time, preferredMin, quietStart, quietEnd int) time.Time {
	candidate := local
	at := func(day time.Time, minute int) time.Time {
		return time.Date(day.Year(), day.Month(), day.Day(), minute/60, minute%60, 0, 0, day.Location())
	}
	if local.Hour()*60+local.Minute() < preferredMin {
		candidate = at(local, preferredMin)
	}
	if quietStart == quietEnd {
		return candidate
	}
	for i := 0; i < 2; i++ {
		minute := candidate.Hour()*60 + candidate.Minute()
		if !inQuietHours(minute, quietStart, quietEnd) {
			break
		}
		// Quiet end is today when the window is overnight and the candidate is
		// before midnight-to-end, otherwise it is tomorrow's end.
		end := at(candidate, quietEnd)
		if quietStart > quietEnd && minute >= quietStart {
			end = end.Add(24 * time.Hour)
		}
		if !end.After(candidate) {
			end = end.Add(24 * time.Hour)
		}
		candidate = end
	}
	return candidate
}
