package leaderboard

import (
	"errors"
	"fmt"
	"time"
)

// ErrInvalidPeriod is returned for an unknown view, an unparseable anchor or a
// period that starts in the future. The HTTP layer turns it into a 400 rather
// than a 500, because the caller can fix every one of those.
var ErrInvalidPeriod = errors.New("invalid period")

// Period is one half-open window [Start, End). All time has zero bounds. Final
// means the window has closed, so the board is an archive rather than a live
// race.
type Period struct {
	Kind  string
	Start time.Time
	End   time.Time
	Final bool
}

// ResolvePeriod maps a view and an optional YYYY-MM-DD anchor to the period
// containing it. An empty anchor means the period containing now.
//
// A week is Monday 00:00 UTC to the next Monday, a month is the first 00:00 UTC
// to the first of the next month. The boundaries are fixed by the product, and
// a client that sends a Wednesday gets that week's Monday rather than a
// different answer per browser timezone.
func ResolvePeriod(view, anchor string, now time.Time) (Period, error) {
	now = now.UTC()
	period := Period{Kind: view}
	switch view {
	case ViewWeekly, ViewMonthly:
		start := WeekStart(now)
		if view == ViewMonthly {
			start = MonthStart(now)
		}
		if anchor != "" {
			at, err := time.Parse("2006-01-02", anchor)
			if err != nil {
				return Period{}, fmt.Errorf("%w: anchor must be YYYY-MM-DD", ErrInvalidPeriod)
			}
			at = at.UTC()
			if at.After(now) {
				return Period{}, fmt.Errorf("%w: anchor is in the future", ErrInvalidPeriod)
			}
			if view == ViewMonthly {
				start = MonthStart(at)
			} else {
				start = WeekStart(at)
			}
		}
		period.Start = start
		if view == ViewMonthly {
			period.End = start.AddDate(0, 1, 0)
		} else {
			period.End = start.AddDate(0, 0, 7)
		}
		period.Final = !now.Before(period.End)
		return period, nil
	case ViewAllTime:
		return period, nil
	default:
		return Period{}, fmt.Errorf("%w: view must be weekly, monthly or all_time", ErrInvalidPeriod)
	}
}

// MonthStart returns the first moment of t's calendar month in UTC.
func MonthStart(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// LeagueWeek is the Monday whose tier frame describes the period. A week uses
// its own Monday, a month uses the Monday of its final week, and all time uses
// the current Monday. The frame is the tier a learner held when the period
// opened, not a tier earned after it closed.
func (p Period) LeagueWeek(now time.Time) time.Time {
	switch p.Kind {
	case ViewWeekly:
		return p.Start
	case ViewMonthly:
		return WeekStart(p.End.Add(-time.Nanosecond))
	default:
		return WeekStart(now)
	}
}

// Shift returns the canonical anchor for the period delta steps away. It is
// used by the Previous and Next controls and clamps at the current period.
func (p Period) Shift(delta int) string {
	if p.Kind == ViewAllTime || delta == 0 {
		return ""
	}
	start := p.Start
	if p.Kind == ViewMonthly {
		start = start.AddDate(0, delta, 0)
	} else {
		start = start.AddDate(0, 0, 7*delta)
	}
	return start.Format("2006-01-02")
}
