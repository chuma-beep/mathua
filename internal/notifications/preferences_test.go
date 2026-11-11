package notifications

import (
	"strings"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func validPreferences() storage.NotificationPreferences {
	return storage.DefaultNotificationPreferences("s1", time.Unix(0, 0).UTC())
}

func TestNormalizePreferencesDefaultsAndCanonicalDays(t *testing.T) {
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	p := validPreferences()
	p.StudyEnabled = true
	p.StudyDays = "fri, mon,wed,mon"
	p.Timezone = ""
	got, err := NormalizePreferencesAt(p, now)
	if err != nil {
		t.Fatalf("normalize: %v", err)
	}
	if got.Timezone != "UTC" {
		t.Errorf("timezone = %q, want UTC", got.Timezone)
	}
	if got.StudyDays != "mon,wed,fri" {
		t.Errorf("study days = %q, want mon,wed,fri", got.StudyDays)
	}
}

func TestNormalizePreferencesRejectsInvalidTimezones(t *testing.T) {
	p := validPreferences()
	p.Timezone = "Mars/Olympus_Mons"
	if _, err := NormalizePreferencesAt(p, time.Now().UTC()); err == nil {
		t.Fatal("expected invalid timezone to fail")
	}
}

func TestNormalizePreferencesRejectsMinuteOutsideDay(t *testing.T) {
	for _, field := range []string{"preferred", "study", "quiet_start", "quiet_end"} {
		for _, value := range []int{-1, 1440} {
			t.Run(field, func(t *testing.T) {
				p := validPreferences()
				switch field {
				case "preferred":
					p.PreferredTimeMin = value
				case "study":
					p.StudyTimeMin = value
				case "quiet_start":
					p.QuietStartMin = value
				case "quiet_end":
					p.QuietEndMin = value
				}
				if _, err := NormalizePreferencesAt(p, time.Now().UTC()); err == nil {
					t.Fatalf("%s=%d should fail", field, value)
				}
			})
		}
	}
}

func TestNormalizePreferencesStudyDays(t *testing.T) {
	now := time.Now().UTC()
	p := validPreferences()
	p.StudyEnabled = true
	p.StudyDays = ""
	if _, err := NormalizePreferencesAt(p, now); err == nil {
		t.Fatal("study enabled with no days should fail")
	}

	for _, raw := range []string{"monday", "mon,xyz", "mon,,tue"} {
		p = validPreferences()
		p.StudyDays = raw
		if _, err := NormalizePreferencesAt(p, now); err == nil {
			t.Fatalf("study days %q should fail", raw)
		}
	}
}

func TestNormalizePreferencesPauseWindow(t *testing.T) {
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	p := validPreferences()
	past := now.Add(-time.Minute)
	p.PausedUntil = &past
	if _, err := NormalizePreferencesAt(p, now); err == nil || !strings.Contains(err.Error(), "future") {
		t.Fatalf("past pause error = %v, want future error", err)
	}

	p = validPreferences()
	tooFar := now.Add(366 * 24 * time.Hour)
	p.PausedUntil = &tooFar
	if _, err := NormalizePreferencesAt(p, now); err == nil || !strings.Contains(err.Error(), "365") {
		t.Fatalf("long pause error = %v, want 365-day error", err)
	}

	p = validPreferences()
	future := now.Add(24 * time.Hour)
	p.PausedUntil = &future
	got, err := NormalizePreferencesAt(p, now)
	if err != nil {
		t.Fatalf("valid pause: %v", err)
	}
	if got.PausedUntil == nil || !got.PausedUntil.Equal(future) {
		t.Errorf("paused_until = %v, want %v", got.PausedUntil, future)
	}
}

func TestValidatePreferencesUsesCurrentClock(t *testing.T) {
	if err := ValidatePreferences(validPreferences()); err != nil {
		t.Fatalf("defaults should validate: %v", err)
	}
}
