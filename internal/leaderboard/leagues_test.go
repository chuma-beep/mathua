package leaderboard

import (
	"fmt"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func leagueRepo(t *testing.T) *storage.SQLiteStore {
	t.Helper()
	repo, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { repo.Close() })
	return repo
}

// TestMyLeague_SeatsAndRanksTheCallersCohort covers the whole read: the first
// reader of a week bootstraps the leagues, the caller is seated, and the
// standings are ordered by this week's XP with the caller's place returned
// separately rather than injected into the list.
func TestMyLeague_SeatsAndRanksTheCallersCohort(t *testing.T) {
	repo := leagueRepo(t)
	var ids []string
	for i := 0; i < 8; i++ {
		st, err := repo.CreateStudent(fmt.Sprintf("learner%d", i))
		if err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
		if err := repo.SetUsername(st.ID, fmt.Sprintf("learner%d", i)); err != nil {
			t.Fatalf("username %d: %v", i, err)
		}
		if err := repo.AwardXP(st.ID, 10*(i+1), "test", fmt.Sprintf("week:%d", i)); err != nil {
			t.Fatalf("award %d: %v", i, err)
		}
		ids = append(ids, st.ID)
	}

	// AwardXP stamps the live week. Pinning a past Monday made this test pass
	// until the calendar moved and then fail with an empty board, so the week
	// is read from the same clock the awards used.
	now := time.Now().UTC()
	board, err := NewComputer(repo).MyLeague(ids[0], now)
	if err != nil {
		t.Fatalf("my league: %v", err)
	}
	if board.Tier != "bronze" || board.LeagueNo != 1 {
		t.Errorf("seat = %s/%d, want bronze/1", board.Tier, board.LeagueNo)
	}
	if board.Size != 8 || len(board.Members) != 8 {
		t.Fatalf("size = %d members = %d, want 8/8", board.Size, len(board.Members))
	}
	if want := WeekStart(now); !board.WeekStart.Equal(want) {
		t.Errorf("week start = %s, want %s", board.WeekStart, want)
	}
	if board.Me == nil || board.Me.Rank != 8 || board.Me.XP != 10 {
		t.Errorf("me = %+v, want rank 8 with 10 XP", board.Me)
	}
	// Highest XP first, and each row carries the public username instead of an
	// account id.
	if board.Members[0].XP != 80 || board.Members[0].ID != "learner7" {
		t.Errorf("first row = %+v, want learner7 at 80 XP", board.Members[0])
	}
	if !board.Members[7].IsMe {
		t.Error("caller's line is not marked")
	}
	// The rules travel with the payload so the interface cannot describe a
	// different threshold from the one applied.
	if board.Promote != storage.LeagueMovers(8) || board.Demote != board.Promote {
		t.Errorf("rules = %d/%d, want %d", board.Promote, board.Demote, storage.LeagueMovers(8))
	}
}

func TestMyLeague_GuestAndZeroXPGetNoSeat(t *testing.T) {
	repo := leagueRepo(t)
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	computer := NewComputer(repo)

	guest, err := computer.MyLeague("guest_abc", now)
	if err != nil {
		t.Fatalf("guest league: %v", err)
	}
	if guest.Tier != "" || len(guest.Members) != 0 || guest.Me != nil {
		t.Errorf("guest got a league: %+v", guest)
	}

	quiet, err := repo.CreateStudent("quiet")
	if err != nil {
		t.Fatalf("create quiet: %v", err)
	}
	empty, err := computer.MyLeague(quiet.ID, now)
	if err != nil {
		t.Fatalf("zero-xp league: %v", err)
	}
	if empty.Tier != "" || empty.Me != nil {
		t.Errorf("zero-XP learner got a league: %+v", empty)
	}
}

// TestMyLeague_SecondReadDoesNotReseatAnyone pins the retry case: two reads of
// the same week return the same seats and the same standings.
func TestMyLeague_SecondReadDoesNotReseatAnyone(t *testing.T) {
	repo := leagueRepo(t)
	var ids []string
	for i := 0; i < 12; i++ {
		st, err := repo.CreateStudent(fmt.Sprintf("learner%d", i))
		if err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
		if err := repo.AwardXP(st.ID, i+1, "test", fmt.Sprintf("twice:%d", i)); err != nil {
			t.Fatalf("award %d: %v", i, err)
		}
		ids = append(ids, st.ID)
	}
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	computer := NewComputer(repo)

	first, err := computer.MyLeague(ids[3], now)
	if err != nil {
		t.Fatalf("first: %v", err)
	}
	second, err := computer.MyLeague(ids[3], now)
	if err != nil {
		t.Fatalf("second: %v", err)
	}
	if first.Tier != second.Tier || first.LeagueNo != second.LeagueNo || first.Size != second.Size {
		t.Errorf("seat changed between reads: %+v then %+v", first, second)
	}
	if first.Me == nil || second.Me == nil || first.Me.Rank != second.Me.Rank {
		t.Errorf("rank changed between reads: %+v then %+v", first.Me, second.Me)
	}
}

func TestMyLeague_EmptyOrGuestStudentID(t *testing.T) {
	repo := leagueRepo(t)
	comp := NewComputer(repo)
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)

	// Empty student ID
	emptyRes, err := comp.MyLeague("", now)
	if err != nil {
		t.Fatalf("MyLeague empty student: %v", err)
	}
	if len(emptyRes.Members) != 0 || emptyRes.Me != nil {
		t.Errorf("expected empty members for empty student, got %+v", emptyRes)
	}

	// Guest student ID
	guestRes, err := comp.MyLeague("guest-abc-123", now)
	if err != nil {
		t.Fatalf("MyLeague guest student: %v", err)
	}
	if len(guestRes.Members) != 0 || guestRes.Me != nil {
		t.Errorf("expected empty members for guest student, got %+v", guestRes)
	}
}
