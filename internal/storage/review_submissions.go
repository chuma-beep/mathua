package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// commitReviewSubmissionTx writes one spaced-review answer as a single unit:
// the review event, the attempt, the progress transition (including FSRS
// state), the topic speed and the XP award. If any write fails the transaction
// rolls back, so a valid FSRS transition can never be committed without its
// review event.
func commitReviewSubmissionTx(db *sql.DB, rebind func(string) string, insertAttempt attemptInserter, c ReviewSubmissionCommit) (ReviewSubmissionOutcome, error) {
	if c.StudentID == "" || c.SubmissionID == "" {
		return ReviewSubmissionOutcome{}, fmt.Errorf("review submission: missing scope")
	}
	if c.Event.AttemptID == "" {
		return ReviewSubmissionOutcome{}, fmt.Errorf("review submission: missing attempt key")
	}
	if c.Event.StudentID == "" {
		c.Event.StudentID = c.StudentID
	}
	if c.Event.StudentID != c.StudentID {
		return ReviewSubmissionOutcome{}, fmt.Errorf("review submission: event student mismatch")
	}
	if c.XPReason == "" {
		c.XPReason = "study"
	}

	tx, err := db.Begin()
	if err != nil {
		return ReviewSubmissionOutcome{}, fmt.Errorf("review submission: begin: %w", err)
	}
	defer tx.Rollback()

	now := time.Now().UTC().UnixNano()
	res, err := tx.Exec(rebind(`
		INSERT INTO review_events
			(student_id, concept_id, attempt_id, rating, correct, elapsed_seconds, reviewed_at,
			 stability_before, difficulty_before, due_before, stability_after, difficulty_after,
			 due_after, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(attempt_id) DO NOTHING
	`), c.Event.StudentID, c.Event.ConceptID, c.Event.AttemptID, c.Event.Rating,
		boolToInt(c.Event.Correct), c.Event.ElapsedSeconds, c.Event.ReviewedAt.UTC().UnixNano(),
		nullableFloat(c.Event.StabilityBefore), nullableFloat(c.Event.DifficultyBefore),
		nullableTime(c.Event.DueBefore), nullableFloat(c.Event.StabilityAfter),
		nullableFloat(c.Event.DifficultyAfter), nullableTime(c.Event.DueAfter), now)
	if err != nil {
		return ReviewSubmissionOutcome{}, fmt.Errorf("review submission: event: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		// The event is the idempotency key. A replay must not write a second
		// attempt, a second progress transition or a second XP award.
		return ReviewSubmissionOutcome{Replayed: true}, nil
	}

	attemptID, err := insertAttempt(tx, c.Attempt)
	if err != nil {
		return ReviewSubmissionOutcome{}, fmt.Errorf("review submission: attempt: %w", err)
	}
	if c.Progress != nil {
		if _, err := tx.Exec(rebind(upsertProgressSQL),
			c.Progress.StudentID, c.Progress.ConceptID, c.Progress.Status, c.Progress.Streak, c.Progress.BestStreak,
			c.Progress.AvgResponseTime, c.Progress.Attempts,
			nullTime(c.Progress.LastAttempted), nullTime(c.Progress.LastReviewed),
			nullTime(c.Progress.NextReviewDue), c.Progress.SM2Repetitions, c.Progress.SM2Interval,
			c.Progress.SM2EFactor, nullTime(c.Progress.MasteredAt), c.Progress.WeaknessScore,
			boolToInt(c.Progress.PlacementSeeded), boolToInt(c.Progress.IsProvisional),
			nullFloat(c.Progress.FSRSStability), nullFloat(c.Progress.FSRSDifficulty), nullTime(c.Progress.FSRSLastReview),
		); err != nil {
			return ReviewSubmissionOutcome{}, fmt.Errorf("review submission: progress: %w", err)
		}
	}
	if c.TopicSpeed != nil {
		ts := *c.TopicSpeed
		if ts.LearningSpeed < 0.5 {
			ts.LearningSpeed = 0.5
		}
		if ts.LearningSpeed > 2.0 {
			ts.LearningSpeed = 2.0
		}
		if ts.EFactor == 0 {
			ts.EFactor = 2.5
		}
		if _, err := tx.Exec(rebind(upsertTopicSpeedSQL),
			ts.StudentID, ts.ConceptID, ts.EFactor, ts.Interval, ts.Repetitions, ts.LearningSpeed,
			time.Now().UTC().Format("2006-01-02 15:04:05"),
		); err != nil {
			return ReviewSubmissionOutcome{}, fmt.Errorf("review submission: topic speed: %w", err)
		}
	}
	if c.XPAmount != 0 {
		eventKey := fmt.Sprintf("attempt:%d", attemptID)
		xpRes, err := tx.Exec(rebind(`
			INSERT INTO xp_events (student_id, amount, reason, event_key, created_at)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(student_id, event_key) DO NOTHING
		`), c.StudentID, c.XPAmount, c.XPReason, eventKey, now)
		if err != nil {
			return ReviewSubmissionOutcome{}, fmt.Errorf("review submission: xp event: %w", err)
		}
		if n, _ := xpRes.RowsAffected(); n > 0 {
			today := time.Now().UTC().Format("2006-01-02")
			if _, err := tx.Exec(rebind(`
				UPDATE students
				SET xp_total = xp_total + ?,
				    xp_today = CASE WHEN xp_date = ? THEN xp_today + ? ELSE ? END,
				    xp_date  = ?
				WHERE id = ?
			`), c.XPAmount, today, c.XPAmount, c.XPAmount, today, c.StudentID); err != nil {
				return ReviewSubmissionOutcome{}, fmt.Errorf("review submission: xp counters: %w", err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return ReviewSubmissionOutcome{}, fmt.Errorf("review submission: commit: %w", err)
	}
	return ReviewSubmissionOutcome{AttemptID: attemptID}, nil
}
