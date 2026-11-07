package storage

import "time"

// XPSinceLastQuiz returns the positive XP earned since the last completed quiz.
// It is the one gate calculation both the engine and the scoring dashboard use.
//
// A raw balance delta is not enough once lifetime XP can be negative. Summing
// positive ledger events is the definition of "earned XP": a deduction, even
// one that takes the balance below the previous snapshot, cannot reduce it.
func XPSinceLastQuiz(repo Repository, studentID string) (int, error) {
	total, _, err := repo.GetXP(studentID)
	if err != nil {
		return 0, err
	}
	last, err := repo.LastQuizCompletion(studentID)
	if err != nil {
		return 0, err
	}
	earnedStore, _ := repo.(EffortStore)
	if earnedStore == nil {
		if last == nil {
			return total, nil
		}
		since := total - last.XPTotal
		if since < 0 {
			since = 0
		}
		return since, nil
	}
	if last == nil {
		return earnedStore.SumXPEarnedSince(studentID, time.Time{})
	}
	if last.CompletedAt == "" {
		return 0, nil
	}
	completedAt, ok := completionCutoff(last.CompletedAt)
	if !ok {
		// An unparseable completion timestamp is not a licence to use a
		// negative balance delta. Treat the gate as unearned until fixed.
		return 0, nil
	}
	return earnedStore.SumXPEarnedSince(studentID, completedAt)
}

// completionCutoff parses every completion timestamp format Mathua has
// written. Legacy SQLite rows have second precision, so their window starts
// one second later to avoid counting a pre-completion event in the same second.
func completionCutoff(raw string) (time.Time, bool) {
	layouts := []struct {
		layout string
		shift  time.Duration
	}{
		{time.RFC3339Nano, 0},
		{time.RFC3339, 0},
		{"2006-01-02 15:04:05.999999-07", 0},
		{"2006-01-02 15:04:05-07", 0},
		{"2006-01-02 15:04:05", time.Second},
	}
	for _, candidate := range layouts {
		completedAt, err := time.Parse(candidate.layout, raw)
		if err == nil {
			return completedAt.Add(candidate.shift), true
		}
	}
	return time.Time{}, false
}
