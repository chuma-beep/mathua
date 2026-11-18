package leaderboard

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func testRepo(t *testing.T) *storage.SQLiteStore {
	t.Helper()
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

func seed(t *testing.T, store *storage.SQLiteStore, name, username string, xp int) string {
	t.Helper()
	st, err := store.CreateUser(name, username, "hash")
	if err != nil {
		t.Fatalf("create %q: %v", name, err)
	}
	if xp != 0 {
		if err := store.AwardXP(st.ID, xp, "study", "seed:"+name); err != nil {
			t.Fatalf("award %q: %v", name, err)
		}
	}
	return st.ID
}

func setSettings(t *testing.T, store *storage.SQLiteStore, id, raw string) {
	t.Helper()
	if err := store.UpdateSettings(id, raw); err != nil {
		t.Fatalf("settings for %s: %v", id, err)
	}
}

// Sunday belongs to the week that started six days earlier, and Monday 00:00
// starts the next one. UTC throughout: the old query already used UTC, and a
// boundary that moves with the viewer's clock would give two learners different
// competitions.
func TestWeekStart_IsMondayMidnightUTC(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"2026-10-19T00:00:00Z", "2026-10-19T00:00:00Z"}, // Monday
		{"2026-10-25T23:59:59Z", "2026-10-19T00:00:00Z"}, // Sunday
		{"2026-10-22T12:00:00Z", "2026-10-19T00:00:00Z"}, // Thursday
		{"2026-10-18T23:59:59Z", "2026-10-12T00:00:00Z"}, // previous Sunday
	}
	for _, c := range cases {
		in, _ := time.Parse(time.RFC3339, c.in)
		if got := WeekStart(in).Format(time.RFC3339); got != c.want {
			t.Errorf("WeekStart(%s) = %s, want %s", c.in, got, c.want)
		}
	}
}

// Ranks continue across pages. A page that restarts at 1 would make page two
// claim the twenty-sixth learner is first.
func TestBoard_RanksContinueAcrossPages(t *testing.T) {
	store := testRepo(t)
	for i := 1; i <= 30; i++ {
		seed(t, store, fmt.Sprintf("Learner %02d", i), fmt.Sprintf("l%02d", i), i)
	}
	c := NewComputer(store)
	now := time.Now().UTC()

	first, err := c.Board("", ViewWeekly, "", 1, 25, now)
	if err != nil {
		t.Fatalf("page 1: %v", err)
	}
	if first.Total != 30 || first.TotalPages != 2 || first.PageSize != 25 {
		t.Fatalf("pagination = total %d, pages %d, size %d", first.Total, first.TotalPages, first.PageSize)
	}
	if len(first.Entries) != 25 || first.Entries[0].Rank != 1 || first.Entries[24].Rank != 25 {
		t.Fatalf("page 1 ranks = %d rows", len(first.Entries))
	}
	if want := WeekStart(now).Format(time.RFC3339); first.WeekStart != want {
		t.Errorf("week start = %s, want %s", first.WeekStart, want)
	}
	if want := WeekStart(now).AddDate(0, 0, 7).Format(time.RFC3339); first.WeekEnd != want {
		t.Errorf("week end = %s, want %s", first.WeekEnd, want)
	}

	second, err := c.Board("", ViewWeekly, "", 2, 25, now)
	if err != nil {
		t.Fatalf("page 2: %v", err)
	}
	if len(second.Entries) != 5 || second.Entries[0].Rank != 26 || second.Entries[4].Rank != 30 {
		t.Fatalf("page 2 starts at rank %d with %d rows", second.Entries[0].Rank, len(second.Entries))
	}
	if second.Total != 30 {
		t.Errorf("page 2 total = %d, want 30", second.Total)
	}
}

// The learner's own row is a position, not an extra entry. Someone ranked 30th
// looking at page one must still be told where they stand.
func TestBoard_MeOutsideThePage(t *testing.T) {
	store := testRepo(t)
	var me string
	for i := 1; i <= 30; i++ {
		sid := seed(t, store, fmt.Sprintf("L%02d", i), fmt.Sprintf("m%02d", i), i)
		if i == 1 {
			me = sid
		}
	}
	c := NewComputer(store)
	now := time.Now().UTC()

	page, err := c.Board(me, ViewWeekly, "", 1, 25, now)
	if err != nil {
		t.Fatalf("board: %v", err)
	}
	if page.Me == nil || !page.Me.Eligible || page.Me.Rank != 30 || page.Me.XP != 1 {
		t.Fatalf("me = %+v, want rank 30 with 1 XP", page.Me)
	}
	if page.Me.InPage {
		t.Error("rank 30 reported as in page 1")
	}
	if page.Me.Guest {
		t.Error("account reported as guest")
	}
	for _, e := range page.Entries {
		if e.IsMe {
			t.Fatal("the caller was injected into the entries")
		}
	}

	last, err := c.Board(me, ViewWeekly, "", 2, 25, now)
	if err != nil {
		t.Fatalf("page 2: %v", err)
	}
	if last.Me == nil || !last.Me.InPage || last.Me.Rank != 30 {
		t.Fatalf("me on page 2 = %+v, want in page at rank 30", last.Me)
	}
	if !last.Entries[4].IsMe {
		t.Error("the caller's entry is not marked on their own page")
	}
}

// Guests are not ranked, but a guest is still holding the phone. Telling them
// they are a guest is the honest answer; a fabricated rank is not.
func TestBoard_GuestCallerIsNotRanked(t *testing.T) {
	store := testRepo(t)
	seed(t, store, "Real", "real", 12)
	c := NewComputer(store)
	now := time.Now().UTC()

	board, err := c.Board("guest_abc123", ViewWeekly, "", 1, 25, now)
	if err != nil {
		t.Fatalf("board: %v", err)
	}
	if board.Me == nil || !board.Me.Guest || board.Me.Eligible {
		t.Fatalf("guest me = %+v, want guest and ineligible", board.Me)
	}
}

// A zero-XP account is eligible in principle but has nothing to rank. It must
// not be given rank 0 or a row of its own.
func TestBoard_ZeroXPCallerHasNoRank(t *testing.T) {
	store := testRepo(t)
	seed(t, store, "Active", "active", 5)
	idle := seed(t, store, "Idle", "idle", 0)
	c := NewComputer(store)
	now := time.Now().UTC()

	board, err := c.Board(idle, ViewWeekly, "", 1, 25, now)
	if err != nil {
		t.Fatalf("board: %v", err)
	}
	if board.Me == nil {
		t.Fatal("no position object")
	}
	if board.Me.Eligible || board.Me.Rank != 0 {
		t.Fatalf("zero-XP me = %+v, want ineligible", board.Me)
	}
}

// The payload is public. It must not carry the account id, and it must not
// carry a guest credential.
func TestBoard_JSONHasNoStudentID(t *testing.T) {
	store := testRepo(t)
	named := seed(t, store, "Named", "named", 9)
	unnamed := seed(t, store, "Unnamed", "", 8)
	c := NewComputer(store)
	now := time.Now().UTC()

	board, err := c.Board(named, ViewWeekly, "", 1, 25, now)
	if err != nil {
		t.Fatalf("board: %v", err)
	}
	raw, err := json.Marshal(board)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(raw), "student_id") {
		t.Errorf("payload contains student_id: %s", raw)
	}
	if strings.Contains(string(raw), named) || strings.Contains(string(raw), unnamed) {
		t.Errorf("payload contains an account id: %s", raw)
	}
	if !strings.Contains(string(raw), `"id":"named"`) {
		t.Errorf("named learner's public id missing: %s", raw)
	}
	if !strings.Contains(string(raw), `"id":"learner-`) {
		t.Errorf("username-less learner has no fallback id: %s", raw)
	}
}

// An opt-out must remove every handle the public payload carried: the display
// name, the username, the hashed fallback id and the learner's own avatar pick.
// Both spellings of the setting do the same thing.
func TestBoard_AnonymisesOptedOutLearners(t *testing.T) {
	cases := []struct {
		mode     string
		settings string
	}{
		{"anonymous", `{"leaderboard_privacy":"anonymous"}`},
		// The hidden case also carries avatar bits, because a custom upload or a
		// dicebear style is itself a way to name somebody.
		{"hidden", `{"leaderboard_privacy":"hidden","avatar_custom":true,"avatar_dicebear":{"style":"bottts","seed":"ada"},"avatar_version":4}`},
	}
	for _, tc := range cases {
		t.Run(tc.mode, func(t *testing.T) {
			store := testRepo(t)
			hidden := seed(t, store, "Ada Lovelace", "ada-user", 30)
			visible := seed(t, store, "Grace Hopper", "grace", 20)
			seed(t, store, "Alan Turing", "alan", 10)
			if err := store.SetAvatarURL(hidden, "https://photos.example/ada.jpg"); err != nil {
				t.Fatalf("avatar url: %v", err)
			}
			setSettings(t, store, hidden, tc.settings)

			now := time.Now().UTC()
			board, err := NewComputer(store).Board("", ViewWeekly, "", 1, 25, now)
			if err != nil {
				t.Fatalf("board: %v", err)
			}
			if len(board.Entries) != 3 {
				t.Fatalf("entries = %d, want 3", len(board.Entries))
			}

			got := board.Entries[0]
			if got.Name == "Anonymous" {
				t.Errorf("name = %q, want a generated alias", got.Name)
			}
			if parts := strings.Fields(got.Name); len(parts) != 3 || len(parts[2]) != 4 {
				t.Errorf("name = %q, want adjective noun 4-digit suffix", got.Name)
			}
			if !strings.HasPrefix(got.ID, "anon-") {
				t.Errorf("id = %q, want a stable anon- alias", got.ID)
			}
			// The alias is stable across reads, not rerolled per request.
			again, err := NewComputer(store).Board("", ViewWeekly, "", 1, 25, now)
			if err != nil {
				t.Fatalf("board again: %v", err)
			}
			if again.Entries[0].Name != got.Name || again.Entries[0].ID != got.ID {
				t.Errorf("alias changed between reads: %+v then %+v", got, again.Entries[0])
			}
			if got.Avatar == nil || got.Avatar.Kind != "generated" || got.Avatar.Style != "notionists" {
				t.Errorf("avatar = %+v, want a generic generated face", got.Avatar)
			}
			if got.Avatar != nil && (got.Avatar.URL != "" || got.Avatar.Version != 0) {
				t.Errorf("avatar still carries learner data: %+v", got.Avatar)
			}

			// The score is public, only the identity is not.
			if got.Rank != 1 || got.XP != 30 {
				t.Errorf("ranking changed: %+v", got)
			}

			raw, err := json.Marshal(board)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			// Check whole handles rather than short substrings. The stable
			// anonymous id is hexadecimal, so a short handle such as "ada" can
			// appear inside it by chance and turn the assertion into a flake.
			for _, leak := range []string{"Ada Lovelace", "ada-user", "photos.example", "bottts", hidden, visible} {
				if leak == "" {
					continue
				}
				if strings.Contains(string(raw), leak) {
					t.Errorf("payload leaks %q: %s", leak, raw)
				}
			}

			// The rest of the page is unaffected.
			other := board.Entries[1]
			if other.Name != "Grace Hopper" || other.ID != "grace" {
				t.Errorf("second row = %+v, want the named learner unchanged", other)
			}
		})
	}
}

// The caller still recognises their own anonymised row.
func TestBoard_AnonymisedCallerIsStillMe(t *testing.T) {
	store := testRepo(t)
	hidden := seed(t, store, "Ada Lovelace", "ada", 30)
	seed(t, store, "Grace Hopper", "grace", 20)
	setSettings(t, store, hidden, `{"leaderboard_privacy":"anonymous"}`)

	now := time.Now().UTC()
	board, err := NewComputer(store).Board(hidden, ViewWeekly, "", 1, 25, now)
	if err != nil {
		t.Fatalf("board: %v", err)
	}
	if len(board.Entries) != 2 || !board.Entries[0].IsMe {
		t.Fatalf("entries = %+v, want the caller marked first", board.Entries)
	}
	if board.Entries[0].Name == "Anonymous" || !strings.HasPrefix(board.Entries[0].ID, "anon-") {
		t.Errorf("caller's own row = %+v, want a generated alias", board.Entries[0])
	}
	if board.Me == nil || board.Me.Rank != 1 {
		t.Errorf("me = %+v, want rank 1 with no name field", board.Me)
	}
	if board.Entries[1].IsMe {
		t.Error("another learner's row marked as the caller")
	}
}

// The alias must be stable for one account and different for another, so two
// anonymous learners do not collapse into the same visible competitor.
func TestAnonymousIdentityIsStableAndDistinct(t *testing.T) {
	idA, nameA, avatarA := anonymousIdentity("student-a")
	idAgain, nameAgain, _ := anonymousIdentity("student-a")
	idB, nameB, _ := anonymousIdentity("student-b")

	if nameA != nameAgain || idA != idAgain {
		t.Fatalf("alias changed for one account: %q/%q then %q/%q", nameA, idA, nameAgain, idAgain)
	}
	if nameA == nameB || idA == idB {
		t.Fatalf("two accounts share an alias: %q/%q and %q/%q", nameA, idA, nameB, idB)
	}
	if !strings.HasPrefix(idA, "anon-") || !strings.HasPrefix(idB, "anon-") {
		t.Errorf("aliases = %q and %q, want anon- ids", idA, idB)
	}
	if avatarA == nil || avatarA.Kind != "generated" || avatarA.Style != "notionists" {
		t.Errorf("alias avatar = %+v, want the generic generated face", avatarA)
	}
}

// Anything that is not a recognised opt-out keeps the learner named, including
// a missing key, an unknown value and a corrupt blob. Defaulting the other way
// would silently hide learners who never asked for it.
func TestBoard_KeepsNamedLearners(t *testing.T) {
	cases := []struct {
		name     string
		settings string
	}{
		{"explicit named", `{"leaderboard_privacy":"named"}`},
		{"missing key", `{"avatar_dicebear":{"style":"bottts"}}`},
		{"unknown value", `{"leaderboard_privacy":"friends-only"}`},
		{"corrupt json", `not-json`},
		{"empty blob", ``},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := testRepo(t)
			id := seed(t, store, "Grace Hopper", "grace", 20)
			setSettings(t, store, id, tc.settings)

			now := time.Now().UTC()
			board, err := NewComputer(store).Board("", ViewWeekly, "", 1, 25, now)
			if err != nil {
				t.Fatalf("board: %v", err)
			}
			if len(board.Entries) != 1 {
				t.Fatalf("entries = %d, want 1", len(board.Entries))
			}
			entry := board.Entries[0]
			if entry.Name != "Grace Hopper" || entry.ID != "grace" {
				t.Errorf("row = %+v, want the real name and username", entry)
			}
		})
	}
}

func TestPrivacy_NormalisesSettingsBlob(t *testing.T) {
	cases := []struct {
		raw  string
		want string
	}{
		{``, PrivacyNamed},
		{`{}`, PrivacyNamed},
		{`{"leaderboard_privacy":"named"}`, PrivacyNamed},
		{`{"leaderboard_privacy":"anonymous"}`, PrivacyAnonymous},
		{`{"leaderboard_privacy":"hidden"}`, PrivacyHidden},
		{`{"leaderboard_privacy":"  Anonymous  "}`, PrivacyAnonymous},
		{`{"leaderboard_privacy":42}`, PrivacyNamed},
		{`not-json`, PrivacyNamed},
	}
	for _, tc := range cases {
		if got := Privacy(tc.raw); got != tc.want {
			t.Errorf("Privacy(%q) = %q, want %q", tc.raw, got, tc.want)
		}
	}
	if !IsAnonymous(`{"leaderboard_privacy":"hidden"}`) {
		t.Error("hidden not reported as anonymous")
	}
	if IsAnonymous(`{"leaderboard_privacy":"named"}`) || IsAnonymous(`{}`) {
		t.Error("named or missing reported as anonymous")
	}
}

func TestPublicID_IsStableAndNonReversible(t *testing.T) {
	a := PublicID("", "0b6f3a4c-1111-2222-3333-444455556666")
	b := PublicID("", "0b6f3a4c-1111-2222-3333-444455556666")
	if a != b {
		t.Fatalf("fallback id not stable: %s then %s", a, b)
	}
	if !strings.HasPrefix(a, "learner-") || len(a) != len("learner-")+8 {
		t.Errorf("fallback id = %q, want learner-<8 hex>", a)
	}
	if strings.Contains(a, "0b6f3a4c") {
		t.Errorf("fallback id leaks the account id: %s", a)
	}
	if got := PublicID("ada", "ignored"); got != "ada" {
		t.Errorf("username not preferred: %s", got)
	}
}

func TestBoard_ViewAndPageBounds(t *testing.T) {
	repo := testRepo(t)
	comp := NewComputer(repo)
	now := time.Now().UTC()

	// Unknown view returns error
	if _, err := comp.Board("", "yearly", "", 1, 10, now); err == nil {
		t.Fatal("expected error for unknown view")
	}

	// Page size <= 0 resets to DefaultPageSize (50)
	b, err := comp.Board("", ViewAllTime, "", 0, 0, now)
	if err != nil {
		t.Fatalf("Board: %v", err)
	}
	if b.PageSize != DefaultPageSize {
		t.Errorf("PageSize = %d, want %d", b.PageSize, DefaultPageSize)
	}
	if b.Page != 1 {
		t.Errorf("Page = %d, want 1", b.Page)
	}

	// Page size > MaxPageSize resets to DefaultPageSize
	b, err = comp.Board("", ViewAllTime, "", 1, MaxPageSize+100, now)
	if err != nil {
		t.Fatalf("Board: %v", err)
	}
	if b.PageSize != DefaultPageSize {
		t.Errorf("PageSize = %d, want %d", b.PageSize, DefaultPageSize)
	}

	// Anonymous call leaves Me as nil
	if b.Me != nil {
		t.Errorf("expected Me to be nil for anonymous call, got %+v", b.Me)
	}
}

func TestDisplayNameAndAvatarFor(t *testing.T) {
	if got := displayName(""); got != "Learner" {
		t.Errorf("displayName('') = %q, want 'Learner'", got)
	}
	if got := displayName("Alice"); got != "Alice" {
		t.Errorf("displayName('Alice') = %q, want 'Alice'", got)
	}

	// Custom upload without username falls back to raw AvatarURL
	av1 := avatarFor(storage.LeaderboardRow{
		AvatarURL:      "https://example.com/avatar.png",
		AvatarSettings: `{"custom":true}`,
	}, "user1")
	if av1.Kind != "custom" || av1.URL != "https://example.com/avatar.png" {
		t.Errorf("av1 = %+v", av1)
	}

	// Generated avatar with custom style
	av2 := avatarFor(storage.LeaderboardRow{
		AvatarSettings: `{"avatar_dicebear":{"style":"bottts"}}`,
	}, "user2")
	if av2.Kind != "generated" || av2.Style != "bottts" {
		t.Errorf("av2 = %+v", av2)
	}

	// Default fallback
	av3 := avatarFor(storage.LeaderboardRow{}, "user3")
	if av3.Kind != "generated" || av3.Style != "notionists" {
		t.Errorf("av3 = %+v", av3)
	}
}
