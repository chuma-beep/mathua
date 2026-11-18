package leaderboard

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/levels"
	"github.com/chuma-beep/mathua/internal/storage"
)

// Public views. "weekly" is the current competition, measured in XP earned
// inside this week's UTC window; "all_time" is lifetime XP. They are separate
// queries on purpose, the old board ranked lifetime mastery by this week's
// count, so the two questions had one wrong answer between them.
const (
	ViewWeekly  = "weekly"
	ViewMonthly = "monthly"
	ViewAllTime = "all_time"

	DefaultPageSize = 25
	MaxPageSize     = 100
)

// Entry is one line of a public ranking. There is no student id here and there
// must not be: a guest id is that client's credential, and an account id is an
// internal key. PublicID is a username, or a one-way derivation of the id for
// accounts that never set one.
type Entry struct {
	Rank     int     `json:"rank"`
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Avatar   *Avatar `json:"avatar,omitempty"`
	XP       int     `json:"xp"`
	XPTotal  int     `json:"xp_total"`
	Mastered int     `json:"mastered"`
	Level    string  `json:"level"`
	// League carries the earned frame. It is empty when the row has no seat in
	// the week being ranked, so a stale tier never decorates an avatar.
	League   string `json:"league,omitempty"`
	LeagueXP int    `json:"league_xp,omitempty"`
	IsMe     bool   `json:"is_me,omitempty"`
}

// Avatar is the small amount the board needs to draw a face. Generated avatars
// are seeded from the public id, never from the account id.
type Avatar struct {
	Kind    string `json:"kind"` // "custom" or "generated"
	URL     string `json:"url,omitempty"`
	Style   string `json:"style,omitempty"`
	Version int    `json:"version,omitempty"`
}

// Position is the learner's own place, kept out of Entries so the board never
// contains a duplicate identity to make room for it.
type Position struct {
	Rank     int  `json:"rank"`
	XP       int  `json:"xp"`
	Eligible bool `json:"eligible"`
	Guest    bool `json:"guest,omitempty"`
	InPage   bool `json:"in_page"`
}

// Board is one page of one view.
type Board struct {
	View           string    `json:"view"`
	PeriodKind     string    `json:"period_kind"`
	PeriodStart    string    `json:"period_start,omitempty"`
	PeriodEnd      string    `json:"period_end,omitempty"`
	Final          bool      `json:"final"`
	EarliestAnchor string    `json:"earliest_anchor,omitempty"`
	Page           int       `json:"page"`
	PageSize       int       `json:"page_size"`
	Total          int       `json:"total"`
	TotalPages     int       `json:"total_pages"`
	Entries        []Entry   `json:"entries"`
	Me             *Position `json:"me,omitempty"`
	// WeekStart and WeekEnd stay for one release so a cached bundle can still
	// read a weekly response while the new period fields roll out.
	WeekStart string `json:"week_start,omitempty"`
	WeekEnd   string `json:"week_end,omitempty"`
}

type Computer struct {
	repo storage.Repository
}

func NewComputer(repo storage.Repository) *Computer {
	return &Computer{repo: repo}
}

// Board answers one page of one view, plus where the caller stands.
//
// Total and rank come from the store's own predicates. Nothing is cached: a
// page query is one indexed read, and the previous 30-second whole-board cache
// meant a learner could answer a question and watch their rank stay stale.
func (c *Computer) Board(studentID, view, anchor string, page, pageSize int, now time.Time) (*Board, error) {
	period, err := ResolvePeriod(view, anchor, now)
	if err != nil {
		return nil, err
	}
	if pageSize <= 0 || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	rows, total, err := c.repo.LeaderboardPage(view, period.Start, period.End, period.LeagueWeek(now), pageSize, offset)
	if err != nil {
		return nil, err
	}
	board := &Board{
		View:       view,
		PeriodKind: view,
		Final:      period.Final,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: (total + pageSize - 1) / pageSize,
		Entries:    make([]Entry, 0, len(rows)),
	}
	if !period.Start.IsZero() {
		board.PeriodStart = period.Start.Format(time.RFC3339)
		board.PeriodEnd = period.End.Format(time.RFC3339)
	}
	if view == ViewWeekly {
		board.WeekStart = period.Start.Format(time.RFC3339)
		board.WeekEnd = period.End.Format(time.RFC3339)
	}
	if earliest, ok, err := c.repo.EarliestLeaderboardEvent(); err == nil && ok {
		board.EarliestAnchor = earliest.Format("2006-01-02")
	}
	for i, r := range rows {
		board.Entries = append(board.Entries, entryFromRow(r, offset+i+1, studentID))
	}
	if studentID != "" {
		board.Me = c.position(studentID, view, period.Start, period.End, offset, pageSize)
	}
	return board, nil
}

// position never fails the page: a learner who cannot be located is simply not
// shown a rank, which is better than a 500 on a board that did load.
func (c *Computer) position(studentID, view string, start, end time.Time, offset, pageSize int) *Position {
	me := &Position{}
	if storage.IsGuestID(studentID) {
		me.Guest = true
		return me
	}
	rank, xp, found, err := c.repo.LeaderboardPosition(studentID, view, start, end)
	if err != nil || !found {
		return me
	}
	me.Rank = rank
	me.XP = xp
	me.Eligible = true
	me.InPage = rank > offset && rank <= offset+pageSize
	return me
}

func entryFromRow(r storage.LeaderboardRow, rank int, meID string) Entry {
	id := PublicID(r.Username, r.StudentID)
	e := Entry{
		Rank:     rank,
		ID:       id,
		Name:     displayName(r.Name),
		XP:       r.XP,
		XPTotal:  r.XPTotal,
		Mastered: r.TotalMastered,
		Level:    levels.Compute(r.TotalMastered),
		League:   r.League,
		LeagueXP: r.LeagueXP,
		IsMe:     meID != "" && r.StudentID == meID,
		Avatar:   avatarFor(r, id),
	}
	if IsAnonymous(r.AvatarSettings) {
		e.ID, e.Name, e.Avatar = anonymousIdentity(r.StudentID)
	}
	return e
}

// Leaderboard privacy modes, read from the opaque settings blob the frontend
// owns. Anything else, a missing key, an unknown value or a corrupt blob, means
// "named": the default must stay the behaviour that existed before the setting.
const (
	PrivacyNamed     = "named"
	PrivacyAnonymous = "anonymous"
	PrivacyHidden    = "hidden"
)

// The alias word lists are deliberately short and neutral. The suffix comes
// from the same hash, so two learners with the same two words still get
// different names in practice.
var (
	anonymousAdjectives = []string{
		"Amber", "Bold", "Calm", "Clever", "Cosmic", "Crisp", "Eager", "Golden",
		"Honest", "Jolly", "Kind", "Lucid", "Nimble", "Quiet", "Sunny", "Swift",
	}
	anonymousNouns = []string{
		"Otter", "Falcon", "Maple", "Comet", "River", "Panda", "Cedar", "Lark",
		"Dolphin", "Sparrow", "Willow", "Badger", "Heron", "Fox", "Koala", "Wren",
	}
)

// leaderboardSettings mirrors only the key this package reads. The blob is
// owned by the frontend, so unknown keys are ignored by encoding/json.
type leaderboardSettings struct {
	LeaderboardPrivacy string `json:"leaderboard_privacy"`
}

// Privacy normalises the leaderboard privacy mode held in a student's settings
// blob. It is tolerant on purpose: a learner must not become visible again, or
// disappear, because a settings blob was written by an older client.
func Privacy(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return PrivacyNamed
	}
	var s leaderboardSettings
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		return PrivacyNamed
	}
	switch mode := strings.ToLower(strings.TrimSpace(s.LeaderboardPrivacy)); mode {
	case PrivacyAnonymous, PrivacyHidden:
		return mode
	default:
		return PrivacyNamed
	}
}

// IsAnonymous reports whether the settings blob asks the board to withhold the
// learner's name and public id.
func IsAnonymous(raw string) bool {
	return Privacy(raw) != PrivacyNamed
}

// anonymousIdentity is the identity a board or league line publishes in place
// of a learner who opted out. The alias is stable for the account and carries
// nothing derived from the username or display name. The hash is namespaced so
// changing the alias vocabulary later does not accidentally make the same
// learner look like a different person mid-week.
func anonymousIdentity(studentID string) (string, string, *Avatar) {
	seed := studentID
	if seed == "" {
		seed = "anonymous"
	}
	sum := sha256.Sum256([]byte("mathua-anon-v1:" + seed))
	adjective := anonymousAdjectives[binary.BigEndian.Uint32(sum[0:4])%uint32(len(anonymousAdjectives))]
	noun := anonymousNouns[binary.BigEndian.Uint32(sum[4:8])%uint32(len(anonymousNouns))]
	number := binary.BigEndian.Uint32(sum[8:12]) % 10000
	id := "anon-" + hex.EncodeToString(sum[12:18])
	return id, fmt.Sprintf("%s %s %04d", adjective, noun, number), anonymousAvatar()
}

// anonymousAvatar is the generic generated face. It deliberately drops the
// learner's own style, seed and version, a style plus a version is already a
// usable fingerprint.
func anonymousAvatar() *Avatar {
	return &Avatar{Kind: "generated", Style: "notionists"}
}

func avatarFor(r storage.LeaderboardRow, publicID string) *Avatar {
	custom, pick, version := storage.AvatarBits(r.AvatarSettings)
	u := r.AvatarURL
	// A custom upload lives behind a public route keyed by username: the board
	// must not need the account id, and the browser must not receive it.
	if custom && r.Username != "" {
		u = "/api/avatar/u/" + url.PathEscape(r.Username)
	}
	if u != "" {
		return &Avatar{Kind: "custom", URL: u, Version: version}
	}
	if pick != nil && pick.Style != "" {
		return &Avatar{Kind: "generated", Style: pick.Style, Version: version}
	}
	return &Avatar{Kind: "generated", Style: "notionists"}
}

// PublicID is what the outside world may call a learner. A username is already
// public and unique; the hash is a stable stand-in for accounts created before
// usernames were required, not a secret and not reversible in practice.
func PublicID(username, studentID string) string {
	if username != "" {
		return username
	}
	sum := sha256.Sum256([]byte(studentID))
	return "learner-" + hex.EncodeToString(sum[:4])
}

func displayName(name string) string {
	if name == "" {
		return "Learner"
	}
	return name
}

// WeekStart returns the Monday 00:00 UTC that opens the competition containing
// t. UTC is the product rule (the existing weekly query used it), and Monday is
// the boundary the UI has always counted down to.
func WeekStart(t time.Time) time.Time {
	t = t.UTC()
	day := t.Weekday()
	if day == time.Sunday {
		day = 7
	}
	monday := t.AddDate(0, 0, -int(day-time.Monday))
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)
}
