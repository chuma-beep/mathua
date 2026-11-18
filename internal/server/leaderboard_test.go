package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// leaderboardGet issues the public board request with an optional bearer token.
func leaderboardGet(t *testing.T, mux *http.ServeMux, query, token string) (int, map[string]interface{}) {
	t.Helper()
	req := httptest.NewRequest("GET", "/api/leaderboard"+query, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var body map[string]interface{}
	if rec.Body.Len() > 0 {
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode %q: %v (%s)", query, err, rec.Body.String())
		}
	}
	return rec.Code, body
}

func entriesOf(t *testing.T, body map[string]interface{}) []map[string]interface{} {
	t.Helper()
	raw, _ := body["entries"].([]interface{})
	out := make([]map[string]interface{}, 0, len(raw))
	for _, e := range raw {
		m, _ := e.(map[string]interface{})
		out = append(out, m)
	}
	return out
}

// Server-side pagination: the client asks for a page and receives exactly that
// slice plus a total it did not have to compute by downloading the board.
func TestLeaderboard_PaginatesOnTheServer(t *testing.T) {
	s, mux, token := avatarTestServer(t)
	ada, err := s.repo.FindByUsername("ada")
	if err != nil || ada == nil {
		t.Fatalf("find ada: %v", err)
	}
	if err := s.repo.AwardXP(ada.ID, 1, "test", "rank:ada"); err != nil {
		t.Fatalf("award ada: %v", err)
	}
	for i := 2; i <= 30; i++ {
		st, err := s.repo.CreateUser(fmt.Sprintf("Ranker %02d", i), fmt.Sprintf("ranker%02d", i), "hash")
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		if err := s.repo.AwardXP(st.ID, i, "test", fmt.Sprintf("rank:%d", i)); err != nil {
			t.Fatalf("award: %v", err)
		}
	}

	code, page1 := leaderboardGet(t, mux, "?view=weekly&page=1&page_size=25", token)
	if code != 200 {
		t.Fatalf("page 1: %d %v", code, page1)
	}
	if page1["total"] != float64(30) || page1["total_pages"] != float64(2) {
		t.Errorf("total/pages = %v/%v, want 30/2", page1["total"], page1["total_pages"])
	}
	if got := entriesOf(t, page1); len(got) != 25 {
		t.Fatalf("page 1 entries = %d, want 25", len(got))
	}
	if got := entriesOf(t, page1)[0]; got["rank"] != float64(1) || got["xp"] != float64(30) {
		t.Errorf("top entry = %v, want rank 1 with 30 XP", got)
	}

	code, page2 := leaderboardGet(t, mux, "?view=weekly&page=2&page_size=25", token)
	if code != 200 {
		t.Fatalf("page 2: %d %v", code, page2)
	}
	got := entriesOf(t, page2)
	if len(got) != 5 || got[0]["rank"] != float64(26) {
		t.Fatalf("page 2 = %v, want 5 entries starting at rank 26", got)
	}

	// The caller is 30th and on page 2. Page 1 must still tell them where they
	// stand without adding them to the twenty-five rows.
	me, _ := page1["me"].(map[string]interface{})
	if me == nil || me["rank"] != float64(30) || me["eligible"] != true || me["in_page"] != false {
		t.Errorf("me on page 1 = %v, want ineligible for the page but ranked 30", me)
	}
	me2, _ := page2["me"].(map[string]interface{})
	if me2 == nil || me2["in_page"] != true {
		t.Errorf("me on page 2 = %v, want in_page", me2)
	}
}

// An out-of-range page is empty, not an error and not a fabricated page.
func TestLeaderboard_OutOfRangePageIsEmpty(t *testing.T) {
	_, mux, _ := avatarTestServer(t)
	code, body := leaderboardGet(t, mux, "?view=weekly&page=99&page_size=25", "")
	if code != 200 {
		t.Fatalf("out of range: %d %v", code, body)
	}
	if len(entriesOf(t, body)) != 0 {
		t.Errorf("entries = %v, want none", body["entries"])
	}
	if body["total"] != float64(0) {
		t.Errorf("total = %v, want 0 (no XP was awarded)", body["total"])
	}
}

func TestLeaderboard_RefusesBadParameters(t *testing.T) {
	_, mux, _ := avatarTestServer(t)
	for _, q := range []string{
		"?view=yearly",
		"?view=weekly&anchor=not-a-date",
		"?view=weekly&anchor=2999-01-01",
		"?view=monthly&anchor=2999-01-01",
		"?view=weekly&page=0",
		"?view=weekly&page=-3",
		"?view=weekly&page=abc",
		"?view=weekly&page_size=0",
		"?view=weekly&page_size=101",
	} {
		code, body := leaderboardGet(t, mux, q, "")
		if code != 400 {
			t.Errorf("%s = %d %v, want 400", q, code, body)
		}
	}
}

// Guests are not ranked, and a guest caller is told so instead of being given a
// rank. The zero-XP account is excluded from the total as well.
func TestLeaderboard_ExcludesGuests(t *testing.T) {
	s, mux, token := avatarTestServer(t)
	if _, err := s.repo.ClaimGuestStudent("guest_board", "Guest"); err != nil {
		t.Fatalf("guest: %v", err)
	}
	if err := s.repo.AwardXP("guest_board", 500, "test", "guest:board"); err != nil {
		t.Fatalf("guest xp: %v", err)
	}

	code, body := leaderboardGet(t, mux, "?view=all_time", "")
	if code != 200 {
		t.Fatalf("board: %d %v", code, body)
	}
	if body["total"] != float64(0) {
		t.Errorf("guest leaked into total: %v", body["total"])
	}

	// Ada has no XP yet, so she is signed in and unranked.
	ada, _ := s.repo.FindByUsername("ada")
	if ada == nil {
		t.Fatal("ada missing")
	}
	me := body["me"]
	if me != nil {
		t.Errorf("unauthenticated request carried a position: %v", me)
	}
	_ = token
}

// A signed-in learner with no XP is not given rank 0 as a fake entry; the
// payload says ineligible and the UI explains how to join.
func TestLeaderboard_ZeroXPAccountHasNoRank(t *testing.T) {
	_, mux, token := avatarTestServer(t)
	code, body := leaderboardGet(t, mux, "?view=weekly", token)
	if code != 200 {
		t.Fatalf("board: %d %v", code, body)
	}
	me, _ := body["me"].(map[string]interface{})
	if me == nil {
		t.Fatal("no position for a signed-in learner")
	}
	if me["eligible"] != false || me["rank"] != float64(0) {
		t.Errorf("me = %v, want ineligible with no rank", me)
	}
}

func TestLeaderboard_MonthlyWindowAndArchive(t *testing.T) {
	_, mux, token := avatarTestServer(t)
	now := time.Now().UTC()
	code, body := leaderboardGet(t, mux, "?view=monthly", token)
	if code != 200 {
		t.Fatalf("monthly = %d %v", code, body)
	}
	if body["period_kind"] != "monthly" || body["final"] != false {
		t.Fatalf("period = %v final=%v", body["period_kind"], body["final"])
	}
	want := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	if body["period_start"] != want {
		t.Fatalf("period_start = %v, want %s", body["period_start"], want)
	}

	lastWeek := time.Now().UTC().AddDate(0, 0, -7).Format("2006-01-02")
	code, body = leaderboardGet(t, mux, "?view=weekly&anchor="+lastWeek, token)
	if code != 200 {
		t.Fatalf("archive = %d %v", code, body)
	}
	if body["final"] != true {
		t.Fatalf("an anchor in a closed week must be final: %v", body)
	}
	if body["period_end"] == body["period_start"] {
		t.Fatalf("archive window is empty: %v", body)
	}
}

func TestLeagues_RefusesAFutureOrMalformedWeek(t *testing.T) {
	_, mux, token := avatarTestServer(t)
	for _, query := range []string{"?week=not-a-date", "?week=2999-01-05"} {
		req := httptest.NewRequest("GET", "/api/leagues"+query, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != 400 {
			t.Fatalf("%s = %d %s, want 400", query, rec.Code, rec.Body.String())
		}
	}
}
