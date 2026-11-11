package notifications

import (
	"encoding/json"
	"errors"
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

// PushSender delivers one encrypted payload to one browser endpoint. The
// delivery worker owns subscription selection and health, so a sender only has
// to make the request and report a permanent rejection as PushGoneError.
type PushSender interface {
	Send(sub storage.NotificationPushSubscription, payload []byte) error
}

// PushGoneError marks an endpoint that will never accept another message, for
// example a browser profile that was cleared. The subscription is deleted
// instead of being retried forever.
type PushGoneError struct {
	Status   int
	Endpoint string
}

func (e *PushGoneError) Error() string {
	return fmt.Sprintf("push endpoint gone (status %d)", e.Status)
}

// PushMessage is the JSON payload the service worker receives. It stays small
// because browsers cap an encrypted payload at 4096 bytes.
type PushMessage struct {
	Title          string `json:"title"`
	Body           string `json:"body"`
	Href           string `json:"href"`
	NotificationID string `json:"notification_id,omitempty"`
}

// DeliveryOptions are the delivery worker dependencies. UnsubscribeToken is
// injected so this package does not import authentication, and so tests can
// use a deterministic signer. Push is optional: a deployment with no VAPID key
// pair still delivers email.
type DeliveryOptions struct {
	Store            storage.NotificationStore
	Sender           mailer.Sender
	Push             PushSender
	FrontendURL      string
	UnsubscribeToken func(studentID string) (string, error)
	Now              func() time.Time
	Lease            time.Duration
	BatchSize        int
	MaxAttempts      int
}

// Delivery sends due email and push deliveries. In-app items are already
// visible in the feed and are recorded as sent by the generator, so only the
// two outbound channels are claimed.
type Delivery struct {
	store            storage.NotificationStore
	sender           mailer.Sender
	push             PushSender
	frontendURL      string
	unsubscribeToken func(string) (string, error)
	now              func() time.Time
	lease            time.Duration
	batchSize        int
	maxAttempts      int
}

// NewDelivery builds the outbound worker. The defaults are a twenty-item
// batch, a five-minute lease and five attempts.
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
		push:             opts.Push,
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
	if d.store == nil {
		return result, fmt.Errorf("notification delivery is not configured")
	}
	if d.sender == nil && d.push == nil {
		return result, fmt.Errorf("notification delivery has no sender")
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
	// A digest settles every item it listed, so a batch claim can still hold a
	// second job from the same learner. That job is no longer sending, and
	// delivering it again would send the digest twice.
	if d.deliverySettled(job) {
		return
	}
	if prefs == nil {
		d.suppress(job, "preferences_missing", now, result)
		return
	}
	if prefs.PausedUntil != nil && prefs.PausedUntil.After(now) {
		// A reminder queued before a pause must not arrive afterwards. The next
		// eligible day will create a fresh notification.
		d.suppress(job, "paused", now, result)
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

	switch job.Channel {
	case storage.NotificationChannelEmail:
		d.deliverEmail(job, prefs, now, result)
	case storage.NotificationChannelPush:
		d.deliverPush(job, prefs, now, result)
	default:
		d.suppress(job, "unknown_channel", now, result)
	}
}

func (d *Delivery) deliverEmail(job storage.NotificationDeliveryJob, prefs *storage.NotificationPreferences, now time.Time, result *DeliveryResult) {
	if !prefs.EmailEnabled {
		d.suppress(job, "email_disabled", now, result)
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
	if d.sender == nil {
		d.suppress(job, "email_not_configured", now, result)
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

	token, err := d.unsubscribeToken(job.StudentID)
	if err != nil {
		d.retryOrFail(job, now, err.Error(), result)
		return
	}
	if jobs := d.digestJobs(job, prefs, now); len(jobs) > 1 {
		body := DigestEmail(job.StudentName, jobs,
			d.frontendURL, d.frontendURL+"/settings#notifications",
			d.frontendURL+"/notifications/unsubscribe?token="+token)
		subject := fmt.Sprintf("Your Mathua digest: %d notifications", len(jobs))
		if err := d.sender.Send(mailer.Message{To: job.Email, Subject: subject, Text: body}); err != nil {
			d.retryOrFail(job, now, err.Error(), result)
			return
		}
		d.completeDigest(jobs, storage.NotificationChannelEmail, "mailer", now, result)
		return
	}
	body := ReminderEmail(job,
		ReminderActionLabel(job.Kind),
		d.actionURL(job),
		d.frontendURL+"/settings#notifications",
		d.frontendURL+"/notifications/unsubscribe?token="+token)
	subject := strings.TrimSpace(job.Title)
	if subject == "" {
		subject = defaultSubject(job.Kind)
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

// deliverPush fans one delivery row out to every device the learner has
// registered. The row is complete when at least one device accepted it; a
// device that reports the endpoint is gone is removed rather than retried.
func (d *Delivery) deliverPush(job storage.NotificationDeliveryJob, prefs *storage.NotificationPreferences, now time.Time, result *DeliveryResult) {
	if !prefs.PushEnabled {
		d.suppress(job, "push_disabled", now, result)
		return
	}
	if d.push == nil {
		d.suppress(job, "push_not_configured", now, result)
		return
	}
	subs, err := d.store.ListNotificationPushSubscriptions(job.StudentID)
	if err != nil {
		d.retryOrFail(job, now, err.Error(), result)
		return
	}
	if len(subs) == 0 {
		d.suppress(job, "push_no_subscription", now, result)
		return
	}
	message := PushMessage{
		Title:          job.Title,
		Body:           job.Body,
		Href:           job.Href,
		NotificationID: job.NotificationID,
	}
	jobs := d.digestJobs(job, prefs, now)
	if len(jobs) > 1 {
		message = PushMessage{
			Title: "Mathua digest",
			Body:  fmt.Sprintf("%d notifications are waiting for you.", len(jobs)),
			Href:  "/notifications",
		}
	}
	payload, err := json.Marshal(message)
	if err != nil {
		d.retryOrFail(job, now, err.Error(), result)
		return
	}
	delivered := 0
	lastErr := ""
	for _, sub := range subs {
		err := d.push.Send(sub, payload)
		if err == nil {
			if markErr := d.store.MarkPushSubscriptionSuccess(sub.Endpoint, now); markErr != nil {
				log.Printf("notifications: mark push success %s: %v", sub.Endpoint, markErr)
			}
			delivered++
			continue
		}
		var gone *PushGoneError
		if errors.As(err, &gone) {
			if delErr := d.store.DeletePushSubscriptionByEndpoint(sub.Endpoint); delErr != nil {
				log.Printf("notifications: delete gone push endpoint: %v", delErr)
			}
			continue
		}
		if markErr := d.store.MarkPushSubscriptionFailure(sub.Endpoint, now); markErr != nil {
			log.Printf("notifications: mark push failure %s: %v", sub.Endpoint, markErr)
		}
		lastErr = err.Error()
	}
	if delivered == 0 {
		if lastErr == "" {
			lastErr = "no reachable push device"
		}
		d.retryOrFail(job, now, lastErr, result)
		return
	}
	if len(jobs) > 1 {
		d.completeDigest(jobs, storage.NotificationChannelPush, "webpush", now, result)
		return
	}
	if err := d.store.CompleteNotificationDelivery(job.NotificationID, job.Channel, "webpush", now); err != nil {
		log.Printf("notifications: complete delivery %s: %v", job.NotificationID, err)
		result.Failed++
		return
	}
	result.Sent++
}

// deliverySettled reports whether another pass already finished this row. A
// claim leaves the row in sending; any other status means the work is done.
func (d *Delivery) deliverySettled(job storage.NotificationDeliveryJob) bool {
	rows, err := d.store.ListNotificationDeliveries(job.NotificationID)
	if err != nil {
		return false
	}
	for _, row := range rows {
		if row.Channel == job.Channel {
			return row.Status != storage.NotificationDeliverySending
		}
	}
	return false
}

// digestJobs returns the waiting items for a learner's digest window. Fewer
// than two items stay on the single-message path: a digest of one is a worse
// version of the reminder it replaces.
func (d *Delivery) digestJobs(job storage.NotificationDeliveryJob, prefs *storage.NotificationPreferences, now time.Time) []storage.NotificationDeliveryJob {
	if prefs == nil || prefs.Digest == storage.NotificationDigestImmediate {
		return nil
	}
	jobs, err := d.store.ListDueNotificationDeliveries(job.StudentID, job.Channel, now)
	if err != nil {
		log.Printf("notifications: list digest items for %s: %v", job.StudentID, err)
		return nil
	}
	return jobs
}

// completeDigest marks every item the one message covered as sent, so the
// worker does not re-send the batch on the next tick.
func (d *Delivery) completeDigest(jobs []storage.NotificationDeliveryJob, channel, providerID string, now time.Time, result *DeliveryResult) {
	for _, item := range jobs {
		if err := d.store.CompleteNotificationDelivery(item.NotificationID, channel, providerID, now); err != nil {
			log.Printf("notifications: complete digest item %s: %v", item.NotificationID, err)
			result.Failed++
			continue
		}
		result.Sent++
	}
}

func (d *Delivery) actionURL(job storage.NotificationDeliveryJob) string {
	href := strings.TrimSpace(job.Href)
	if href == "" {
		href = "/review"
	}
	if d.frontendURL == "" {
		return href
	}
	return d.frontendURL + href
}

func defaultSubject(kind string) string {
	switch kind {
	case ReviewKind:
		return "Mathua review reminder"
	case StudyKind:
		return "Mathua study reminder"
	case MilestoneKind:
		return "Mathua milestone"
	case InactivityKind:
		return "Mathua is still here"
	default:
		return "Mathua notification"
	}
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
