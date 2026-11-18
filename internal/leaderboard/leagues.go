package leaderboard

import (
	"fmt"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// MyLeague is one learner's weekly league: the seat, the standings in it, and
// their place among them. It is deliberately not a list of every tier, the
// old page rendered the whole pyramid, which told a bronze learner nothing
// about the thirty people they are actually competing with.
//
// Promotion and relegation counts travel with the payload so the interface can
// state the rule without hardcoding it, and so the rule cannot drift between
// the transition and the sentence describing it.
type MyLeague struct {
	Week         string             `json:"week"`
	WeekStart    time.Time          `json:"week_start"`
	WeekEnd      time.Time          `json:"week_end"`
	Final        bool               `json:"final"`
	EarliestWeek string             `json:"earliest_week,omitempty"`
	Movement     string             `json:"movement,omitempty"`
	FromTier     string             `json:"from_tier,omitempty"`
	ToTier       string             `json:"to_tier,omitempty"`
	Tier         string             `json:"tier"`
	LeagueNo     int                `json:"league_no"`
	Size         int                `json:"size"`
	Promote      int                `json:"promote"`
	Demote       int                `json:"demote"`
	Members      []LeagueMemberLine `json:"members"`
	Me           *LeaguePlace       `json:"me,omitempty"`
}

// LeagueMemberLine is one seat. As on the public board, there is no account id
// here: the public id is a username or a one-way derivation of the id, and the
// avatar is resolved through it.
type LeagueMemberLine struct {
	Rank   int     `json:"rank"`
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Avatar *Avatar `json:"avatar,omitempty"`
	XP     int     `json:"xp"`
	IsMe   bool    `json:"is_me,omitempty"`
}

// LeaguePlace is the caller's own line, also present inside Members when they
// are seated. Kept separate so a future interface can pin it without scanning.
type LeaguePlace struct {
	Rank int `json:"rank"`
	XP   int `json:"xp"`
}

// MyLeague returns the caller's league for the week containing now.
func (c *Computer) MyLeague(studentID string, now time.Time) (*MyLeague, error) {
	return c.MyLeagueAt(studentID, "", now)
}

// MyLeagueAt returns the caller's league for the week containing anchor. An
// empty anchor is the current week. A closed week is read from the stored
// snapshot, so later ledger corrections cannot rewrite a result a learner
// already saw.
//
// The transition runs here rather than on a schedule because there is no
// scheduler that survives a quiet period: if nobody reads, nobody needs seats,
// and the first reader of a new week is exactly the moment the previous week's
// results must be settled. ResolveLeagueWeek is idempotent, so every reader can
// run it.
//
// A guest or a learner with no XP that week gets an empty value rather than an
// invented league: the interface explains how to enter, which is truer than
// seating them alone in a league of one.
func (c *Computer) MyLeagueAt(studentID, anchor string, now time.Time) (*MyLeague, error) {
	now = now.UTC()
	start := WeekStart(now)
	if anchor != "" {
		at, err := time.Parse("2006-01-02", anchor)
		if err != nil {
			return nil, fmt.Errorf("%w: week must be YYYY-MM-DD", ErrInvalidPeriod)
		}
		if at.After(now) {
			return nil, fmt.Errorf("%w: week is in the future", ErrInvalidPeriod)
		}
		start = WeekStart(at)
	}
	week := start.Format("2006-01-02")
	end := start.AddDate(0, 0, 7)
	current := WeekStart(now)
	final := start.Before(current)

	out := &MyLeague{
		Week:      week,
		WeekStart: start,
		WeekEnd:   end,
		Final:     final,
		Members:   []LeagueMemberLine{},
	}
	if studentID == "" {
		return out, nil
	}
	if err := c.repo.ResolveLeagueWeek(current.Format("2006-01-02"), current.AddDate(0, 0, -7), current); err != nil {
		return nil, err
	}
	if earliest, ok, err := c.repo.EarliestLeagueWeek(); err == nil && ok {
		out.EarliestWeek = earliest
	}
	if !final {
		if err := c.repo.EnsureLeagueMembership(studentID, week); err != nil {
			return nil, err
		}
	}
	tier, leagueNo, found, err := c.repo.LeagueMembershipOf(studentID, week)
	if err != nil {
		return nil, err
	}
	if !found {
		return out, nil
	}
	var rows []storage.LeagueRow
	if final {
		rows, err = c.repo.LeagueRosterFinal(week, tier, leagueNo)
	} else {
		rows, err = c.repo.LeagueRoster(week, tier, leagueNo, start, end)
	}
	if err != nil {
		return nil, err
	}
	out.Tier = tier
	out.LeagueNo = leagueNo
	out.Size = len(rows)
	out.Promote = storage.LeagueMovers(len(rows))
	out.Demote = out.Promote
	if final {
		movement, fromTier, toTier, found, err := c.repo.LeagueOutcome(studentID, end.Format("2006-01-02"))
		if err == nil && found {
			out.Movement, out.FromTier, out.ToTier = movement, fromTier, toTier
		}
	}
	for _, r := range rows {
		publicID := PublicID(r.Username, r.StudentID)
		line := LeagueMemberLine{
			Rank: r.Rank,
			ID:   publicID,
			Name: displayName(r.Name),
			XP:   r.XP,
		}
		if av := avatarFor(storage.LeaderboardRow{
			StudentID:      r.StudentID,
			Username:       r.Username,
			AvatarURL:      r.AvatarURL,
			AvatarSettings: r.Settings,
		}, publicID); av != nil {
			line.Avatar = av
		}
		// Same rule as the public board: an opted-out learner competes under a
		// rank-derived placeholder, and the caller still recognises their own
		// seat through IsMe.
		if IsAnonymous(r.Settings) {
			line.ID, line.Name, line.Avatar = anonymousIdentity(r.StudentID)
		}
		if r.StudentID == studentID {
			line.IsMe = true
			out.Me = &LeaguePlace{Rank: r.Rank, XP: r.XP}
		}
		out.Members = append(out.Members, line)
	}
	return out, nil
}
