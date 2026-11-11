package notifications

import (
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// Retention windows. A notification older than RetentionDays is history: the
// feed is a reminder surface, not an archive. A push endpoint that has failed
// PushFailureLimit times is dead, and keeping it would make every send pay for
// a device that will never answer again.
const (
	RetentionDays            = 180
	PushFailureRetentionDays = 30
	PushFailureLimit         = 10
)

// RetentionResult reports one pass of the daily cleanup.
type RetentionResult struct {
	Notifications     int
	Deliveries        int
	PushSubscriptions int
}

// RunRetention deletes history older than the policy allows. It is safe to run
// more than once a day: every statement is bounded by a timestamp, so a second
// pass simply finds nothing left to remove.
func RunRetention(store storage.NotificationStore, now time.Time) (RetentionResult, error) {
	var result RetentionResult
	pruned, err := store.PruneNotificationHistory(now.UTC().AddDate(0, 0, -RetentionDays))
	if err != nil {
		return result, err
	}
	result.Notifications = pruned.Notifications
	result.Deliveries = pruned.Deliveries
	subs, err := store.PrunePushSubscriptions(
		now.UTC().AddDate(0, 0, -PushFailureRetentionDays), PushFailureLimit)
	if err != nil {
		return result, err
	}
	result.PushSubscriptions = subs
	return result, nil
}
