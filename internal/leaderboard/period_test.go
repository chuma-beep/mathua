package leaderboard

import (
	"errors"
	"testing"
	"time"
)

func day(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t.UTC()
}

func TestResolvePeriod_WeekSnapsToMonday(t *testing.T) {
	now := day("2026-11-18")
	period, err := ResolvePeriod(ViewWeekly, "2026-11-12", now)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := period.Start, day("2026-11-09"); !got.Equal(want) {
		t.Fatalf("start = %s, want %s", got, want)
	}
	if got, want := period.End, day("2026-11-16"); !got.Equal(want) {
		t.Fatalf("end = %s, want %s", got, want)
	}
	if !period.Final {
		t.Fatal("a week before now must be final")
	}
}

func TestResolvePeriod_MonthKeepsCalendarBoundary(t *testing.T) {
	now := day("2026-11-18")
	period, err := ResolvePeriod(ViewMonthly, "2026-02-17", now)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := period.Start, day("2026-02-01"); !got.Equal(want) {
		t.Fatalf("start = %s, want %s", got, want)
	}
	if got, want := period.End, day("2026-03-01"); !got.Equal(want) {
		t.Fatalf("end = %s, want %s", got, want)
	}
	if !period.Final {
		t.Fatal("February must be final in November")
	}
}

func TestResolvePeriod_CurrentWeekIsNotFinal(t *testing.T) {
	now := day("2026-11-18").Add(12 * time.Hour)
	period, err := ResolvePeriod(ViewWeekly, "", now)
	if err != nil {
		t.Fatal(err)
	}
	if period.Final {
		t.Fatal("the week containing now is live")
	}
	if got, want := period.Start, day("2026-11-16"); !got.Equal(want) {
		t.Fatalf("start = %s, want %s", got, want)
	}
}

func TestResolvePeriod_RejectsFutureAndUnknown(t *testing.T) {
	now := day("2026-11-18")
	if _, err := ResolvePeriod(ViewWeekly, "2026-11-19", now); !errors.Is(err, ErrInvalidPeriod) {
		t.Fatalf("tomorrow error = %v, want ErrInvalidPeriod", err)
	}
	if _, err := ResolvePeriod("yearly", "", now); !errors.Is(err, ErrInvalidPeriod) {
		t.Fatalf("unknown view error = %v, want ErrInvalidPeriod", err)
	}
	if _, err := ResolvePeriod(ViewWeekly, "18/11/2026", now); !errors.Is(err, ErrInvalidPeriod) {
		t.Fatalf("bad anchor error = %v, want ErrInvalidPeriod", err)
	}
}

func TestResolvePeriod_AllTimeIgnoresAnchor(t *testing.T) {
	period, err := ResolvePeriod(ViewAllTime, "2099-01-01", day("2026-11-18"))
	if err != nil {
		t.Fatal(err)
	}
	if !period.Start.IsZero() || !period.End.IsZero() {
		t.Fatalf("all time has no bounds: %+v", period)
	}
}

func TestPeriod_LeagueWeekUsesTheOpeningFrame(t *testing.T) {
	period, err := ResolvePeriod(ViewMonthly, "2026-11-01", day("2026-11-18"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := period.LeagueWeek(day("2026-11-18")), day("2026-11-30"); !got.Equal(want) {
		t.Fatalf("monthly league week = %s, want the final Monday %s", got, want)
	}
	week, _ := ResolvePeriod(ViewWeekly, "", day("2026-11-18"))
	if got, want := week.LeagueWeek(day("2026-11-18")), day("2026-11-16"); !got.Equal(want) {
		t.Fatalf("weekly league week = %s, want %s", got, want)
	}
}

func TestPeriod_Shift(t *testing.T) {
	period, _ := ResolvePeriod(ViewWeekly, "2026-11-09", day("2026-11-18"))
	if got := period.Shift(-1); got != "2026-11-02" {
		t.Fatalf("previous week = %q", got)
	}
	if got := period.Shift(1); got != "2026-11-16" {
		t.Fatalf("next week = %q", got)
	}
	month, _ := ResolvePeriod(ViewMonthly, "2026-01-01", day("2026-11-18"))
	if got := month.Shift(-1); got != "2025-12-01" {
		t.Fatalf("previous month = %q", got)
	}
}
