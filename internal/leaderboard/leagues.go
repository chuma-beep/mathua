package leaderboard

import (
	"sort"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// Tier ladder, low → high. Students are grouped by tier; each week the top
// performers promote and the bottom demote (improve.md:43 leagues weekly
// promotion/demotion, fresh start).
var Tiers = []string{"bronze", "silver", "gold", "platinum", "diamond"}

// LeagueSize is the number of students grouped per league table.
const LeagueSize = 10

// PromoteCount / DemoteCount are how many move each weekly reset.
const (
	PromoteCount = 2
	DemoteCount  = 2
)

// League is one tier's current standings after any pending reset.
type League struct {
	Tier    string          `json:"tier"`
	Members []storage.LeagueMember `json:"members"`
}

// LeagueBoard is the full response for GET /api/leagues.
type LeagueBoard struct {
	Week  string   `json:"week"`
	Leagues []League `json:"leagues"`
}

// weekKey is the ISO week the standings are tracked against.
func weekKey(t time.Time) string {
	y, w := t.ISOWeek()
	return time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, (w-1)*7).Format("2006") + "-W" + two(w)
}

func two(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}

// ResolveWeeklyReset applies promotion/demotion for the previous week once
// per student when the tracked league_week differs from the current week.
// It reads all standings, groups by tier, sorts by weekly mastered, and
// moves PromoteCount up / DemoteCount down within each tier.
func ResolveWeeklyReset(repo storage.Repository, now time.Time) (string, error) {
	key := weekKey(now)
	members, err := repo.GetLeagueStandings()
	if err != nil {
		return "", err
	}
	// Group by tier, skipping students whose league_week is already current.
	byTier := make(map[string][]storage.LeagueMember)
	for _, m := range members {
		st, err := repo.GetStudent(m.StudentID)
		if err != nil || st == nil {
			continue
		}
		if st.LeagueWeek == key {
			continue
		}
		byTier[m.Tier] = append(byTier[m.Tier], m)
	}
	for tier, list := range byTier {
		sort.Slice(list, func(i, j int) bool {
			if list[i].WeeklyMastered != list[j].WeeklyMastered {
				return list[i].WeeklyMastered > list[j].WeeklyMastered
			}
			return list[i].Name < list[j].Name
		})
		tierIdx := indexOf(Tiers, tier)
		if tierIdx < 0 {
			tierIdx = 0
		}
		for i, m := range list {
			moved := 0
			newTier := tier
			if len(list) >= PromoteCount*2 && i < PromoteCount && tierIdx < len(Tiers)-1 {
				moved = 1
				newTier = Tiers[tierIdx+1]
			} else if len(list) > LeagueSize/2 && i >= len(list)-DemoteCount && tierIdx > 0 {
				moved = -1
				newTier = Tiers[tierIdx-1]
			}
			if err := repo.SetLeague(m.StudentID, newTier, key, moved); err != nil {
				return "", err
			}
		}
	}
	return key, nil
}

// Standings returns the full league board (all tiers), resolving the weekly
// reset first.
func Standings(repo storage.Repository, now time.Time) (*LeagueBoard, error) {
	key, err := ResolveWeeklyReset(repo, now)
	if err != nil {
		return nil, err
	}
	members, err := repo.GetLeagueStandings()
	if err != nil {
		return nil, err
	}
	board := &LeagueBoard{Week: key}
	for _, tier := range Tiers {
		var list []storage.LeagueMember
		for _, m := range members {
			if m.Tier == tier {
				list = append(list, m)
			}
		}
		sort.Slice(list, func(i, j int) bool {
			if list[i].WeeklyMastered != list[j].WeeklyMastered {
				return list[i].WeeklyMastered > list[j].WeeklyMastered
			}
			return list[i].Name < list[j].Name
		})
		if len(list) == 0 {
			continue
		}
		board.Leagues = append(board.Leagues, League{Tier: tier, Members: list})
	}
	return board, nil
}

func indexOf(s []string, v string) int {
	for i, x := range s {
		if x == v {
			return i
		}
	}
	return -1
}
