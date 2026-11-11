// Package notifications owns the preference policy for learner notifications.
// Phase 1 validates and stores preferences only; generation and delivery are
// deliberately absent.
package notifications

import (
	"fmt"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

var studyDayOrder = []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}

var studyDayIndex = func() map[string]int {
	out := make(map[string]int, len(studyDayOrder))
	for i, day := range studyDayOrder {
		out[day] = i
	}
	return out
}()

// ValidatePreferences validates a preference row against the current time.
func ValidatePreferences(p storage.NotificationPreferences) error {
	_, err := NormalizePreferencesAt(p, time.Now().UTC())
	return err
}

// NormalizePreferencesAt validates and canonicalizes preferences. It is the
// boundary used by storage and by tests that need a deterministic clock.
func NormalizePreferencesAt(p storage.NotificationPreferences, now time.Time) (storage.NotificationPreferences, error) {
	p.Timezone = strings.TrimSpace(p.Timezone)
	if p.Timezone == "" {
		p.Timezone = storage.NotificationDefaultTimezone
	}
	loc, err := time.LoadLocation(p.Timezone)
	if err != nil {
		return p, fmt.Errorf("invalid timezone %q", p.Timezone)
	}
	p.Timezone = loc.String()

	if err := validateMinute("preferred_time_min", p.PreferredTimeMin); err != nil {
		return p, err
	}
	if err := validateMinute("study_time_min", p.StudyTimeMin); err != nil {
		return p, err
	}
	if err := validateMinute("quiet_start_min", p.QuietStartMin); err != nil {
		return p, err
	}
	if err := validateMinute("quiet_end_min", p.QuietEndMin); err != nil {
		return p, err
	}

	days, err := normalizeStudyDays(p.StudyDays)
	if err != nil {
		return p, err
	}
	p.StudyDays = days
	if p.StudyEnabled && days == "" {
		return p, fmt.Errorf("study_days is required when study reminders are enabled")
	}

	if p.PausedUntil != nil {
		paused := p.PausedUntil.UTC()
		if !paused.After(now.UTC()) {
			return p, fmt.Errorf("paused_until must be in the future")
		}
		if paused.After(now.UTC().Add(365 * 24 * time.Hour)) {
			return p, fmt.Errorf("paused_until must be within 365 days")
		}
		p.PausedUntil = &paused
	}
	return p, nil
}

func validateMinute(name string, value int) error {
	if value < 0 || value > 1439 {
		return fmt.Errorf("%s must be between 0 and 1439", name)
	}
	return nil
}

func normalizeStudyDays(raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", nil
	}
	seen := make(map[string]bool, len(studyDayOrder))
	for _, part := range strings.Split(raw, ",") {
		day := strings.ToLower(strings.TrimSpace(part))
		if day == "" {
			return "", fmt.Errorf("invalid study day %q", part)
		}
		if _, ok := studyDayIndex[day]; !ok {
			return "", fmt.Errorf("invalid study day %q", day)
		}
		seen[day] = true
	}
	canonical := make([]string, 0, len(seen))
	for _, day := range studyDayOrder {
		if seen[day] {
			canonical = append(canonical, day)
		}
	}
	return strings.Join(canonical, ","), nil
}
