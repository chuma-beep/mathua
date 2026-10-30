package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

type friendAccess struct {
	db     *sql.DB
	rebind func(string) string
}

func (s *SQLiteStore) friends() *friendAccess {
	return &friendAccess{db: s.db, rebind: rebindQuestion}
}

func (s *PostgresStore) friends() *friendAccess {
	return &friendAccess{db: s.db, rebind: rebindDollars}
}

type friendSettings struct {
	FriendsVisibility  string `json:"friends_visibility"`
	LeaderboardPrivacy string `json:"leaderboard_privacy"`
}

func parseFriendSettings(raw string) friendSettings {
	var s friendSettings
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &s)
	}
	if s.FriendsVisibility == "" {
		s.FriendsVisibility = "friends"
	}
	if s.LeaderboardPrivacy == "" {
		s.LeaderboardPrivacy = "named"
	}
	return s
}

func friendVisibility(raw string) string {
	if parseFriendSettings(raw).FriendsVisibility == "nobody" {
		return "nobody"
	}
	return "friends"
}

func friendLeaderboardHidden(raw string) bool {
	return parseFriendSettings(raw).LeaderboardPrivacy == "hidden"
}

func canonicalPair(a, b string) (string, string) {
	if a < b {
		return a, b
	}
	return b, a
}

func friendTime(ns int64) time.Time {
	if ns <= 0 {
		return time.Time{}
	}
	return time.Unix(0, ns).UTC()
}

func friendUnix(t time.Time) int64 {
	if t.IsZero() {
		return time.Now().UTC().UnixNano()
	}
	return t.UTC().UnixNano()
}

type friendTarget struct {
	ID            string
	ProfileID     string
	Name          string
	Username      string
	AvatarURL     string
	Settings      string
	EmailVerified bool
	XPTotal       int
}

func (a *friendAccess) resolveFriendTarget(query string) (string, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return "", ErrFriendNotFound
	}
	var profileID string
	err := a.db.QueryRow(
		a.rebind("SELECT profile_id FROM students WHERE profile_id = ? OR friend_code = ? OR lower(username) = lower(?) LIMIT 1"),
		q, strings.ToUpper(q), q,
	).Scan(&profileID)
	if errors.Is(err, sql.ErrNoRows) || profileID == "" {
		return "", ErrFriendNotFound
	}
	if err != nil {
		return "", fmt.Errorf("resolve friend target: %w", err)
	}
	return profileID, nil
}

func (a *friendAccess) loadStudent(id string) (*friendTarget, error) {
	return a.loadStudentWhere("id = ?", id)
}

func (a *friendAccess) loadStudentByProfile(profileID string) (*friendTarget, error) {
	if strings.TrimSpace(profileID) == "" {
		return nil, ErrFriendNotFound
	}
	return a.loadStudentWhere("profile_id = ?", profileID)
}

func (a *friendAccess) loadStudentWhere(where, value string) (*friendTarget, error) {
	var t friendTarget
	var verified int
	err := a.db.QueryRow(
		a.rebind("SELECT id, profile_id, name, COALESCE(username, ''), avatar_url, settings, email_verified, xp_total FROM students WHERE "+where),
		value,
	).Scan(&t.ID, &t.ProfileID, &t.Name, &t.Username, &t.AvatarURL, &t.Settings, &verified, &t.XPTotal)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrFriendNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("load friend student: %w", err)
	}
	t.EmailVerified = verified != 0
	return &t, nil
}

func (a *friendAccess) blockedEither(x, y string) (bool, error) {
	var one int
	err := a.db.QueryRow(
		a.rebind("SELECT 1 FROM friend_blocks WHERE (blocker_id = ? AND blocked_id = ?) OR (blocker_id = ? AND blocked_id = ?) LIMIT 1"),
		x, y, y, x,
	).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check friend block: %w", err)
	}
	return true, nil
}

func (a *friendAccess) friendship(x, y string) (*Friendship, error) {
	low, high := canonicalPair(x, y)
	var f Friendship
	var status string
	var created, responded int64
	err := a.db.QueryRow(
		a.rebind("SELECT status, requested_by, created_at, responded_at FROM friendships WHERE student_a = ? AND student_b = ?"),
		low, high,
	).Scan(&status, &f.RequestedBy, &created, &responded)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrFriendNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("load friendship: %w", err)
	}
	f.StudentA = low
	f.StudentB = high
	f.Status = FriendshipStatus(status)
	f.CreatedAt = friendTime(created)
	f.RespondedAt = friendTime(responded)
	return &f, nil
}

func (a *friendAccess) loadPair(studentID, otherProfileID string) (*friendTarget, *friendTarget, error) {
	caller, err := a.loadStudent(studentID)
	if err != nil {
		return nil, nil, err
	}
	target, err := a.loadStudentByProfile(otherProfileID)
	if err != nil {
		return nil, nil, err
	}
	if caller.ID == target.ID {
		return nil, nil, ErrFriendSelf
	}
	blocked, err := a.blockedEither(caller.ID, target.ID)
	if err != nil {
		return nil, nil, err
	}
	if blocked {
		return nil, nil, ErrFriendBlocked
	}
	return caller, target, nil
}

func (a *friendAccess) sendFriendRequest(studentID, targetProfileID string, now time.Time) (Friendship, error) {
	caller, target, err := a.loadPair(studentID, targetProfileID)
	if err != nil {
		return Friendship{}, err
	}
	if friendVisibility(caller.Settings) == "nobody" || friendVisibility(target.Settings) == "nobody" {
		return Friendship{}, ErrFriendVisibilityOff
	}
	if !caller.EmailVerified && !a.hasVerifiedIdentity(caller.ID) {
		return Friendship{}, ErrFriendNotVerified
	}

	low, high := canonicalPair(caller.ID, target.ID)
	tx, err := a.db.Begin()
	if err != nil {
		return Friendship{}, fmt.Errorf("send friend request: %w", err)
	}
	at := friendUnix(now)
	_, err = tx.Exec(
		a.rebind(`INSERT INTO friendships (student_a, student_b, status, requested_by, created_at, responded_at)
			VALUES (?, ?, 'pending', ?, ?, 0)
			ON CONFLICT (student_a, student_b) DO UPDATE SET
				status = CASE
					WHEN friendships.status = 'pending' AND friendships.requested_by <> excluded.requested_by THEN 'accepted'
					ELSE friendships.status
				END,
				responded_at = CASE
					WHEN friendships.status = 'pending' AND friendships.requested_by <> excluded.requested_by THEN excluded.created_at
					ELSE friendships.responded_at
				END`),
		low, high, caller.ID, at,
	)
	if err != nil {
		tx.Rollback()
		return Friendship{}, fmt.Errorf("send friend request: %w", err)
	}
	var f Friendship
	var status string
	var created, responded int64
	err = tx.QueryRow(
		a.rebind("SELECT status, requested_by, created_at, responded_at FROM friendships WHERE student_a = ? AND student_b = ?"),
		low, high,
	).Scan(&status, &f.RequestedBy, &created, &responded)
	if err != nil {
		tx.Rollback()
		return Friendship{}, fmt.Errorf("send friend request: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return Friendship{}, fmt.Errorf("send friend request: %w", err)
	}
	f.StudentA = low
	f.StudentB = high
	f.Status = FriendshipStatus(status)
	f.CreatedAt = friendTime(created)
	f.RespondedAt = friendTime(responded)
	return f, nil
}

func (a *friendAccess) acceptFriendRequest(studentID, otherProfileID string, now time.Time) (Friendship, error) {
	caller, target, err := a.loadPair(studentID, otherProfileID)
	if err != nil {
		return Friendship{}, err
	}
	if friendVisibility(caller.Settings) == "nobody" || friendVisibility(target.Settings) == "nobody" {
		return Friendship{}, ErrFriendVisibilityOff
	}
	low, high := canonicalPair(caller.ID, target.ID)
	res, err := a.db.Exec(
		a.rebind("UPDATE friendships SET status = 'accepted', responded_at = ? WHERE student_a = ? AND student_b = ? AND status = 'pending' AND requested_by = ?"),
		friendUnix(now), low, high, target.ID,
	)
	if err != nil {
		return Friendship{}, fmt.Errorf("accept friend request: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		existing, err := a.friendship(low, high)
		if err != nil || existing.Status != FriendshipAccepted {
			return Friendship{}, ErrFriendNotFound
		}
		return *existing, nil
	}
	f, err := a.friendship(low, high)
	if err != nil {
		return Friendship{}, err
	}
	return *f, nil
}

func (a *friendAccess) declineFriendRequest(studentID, otherProfileID string, now time.Time) error {
	caller, target, err := a.loadPair(studentID, otherProfileID)
	if err != nil {
		return err
	}
	_ = now
	low, high := canonicalPair(caller.ID, target.ID)
	res, err := a.db.Exec(
		a.rebind("DELETE FROM friendships WHERE student_a = ? AND student_b = ? AND status = 'pending' AND requested_by = ?"),
		low, high, target.ID,
	)
	if err != nil {
		return fmt.Errorf("decline friend request: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrFriendNotFound
	}
	return nil
}

func (a *friendAccess) cancelFriendRequest(studentID, otherProfileID string, now time.Time) error {
	caller, target, err := a.loadPair(studentID, otherProfileID)
	if err != nil {
		return err
	}
	_ = now
	low, high := canonicalPair(caller.ID, target.ID)
	res, err := a.db.Exec(
		a.rebind("DELETE FROM friendships WHERE student_a = ? AND student_b = ? AND status = 'pending' AND requested_by = ?"),
		low, high, caller.ID,
	)
	if err != nil {
		return fmt.Errorf("cancel friend request: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrFriendNotFound
	}
	return nil
}

func (a *friendAccess) removeFriend(studentID, otherProfileID string, now time.Time) error {
	caller, target, err := a.loadPair(studentID, otherProfileID)
	if err != nil {
		return err
	}
	_ = now
	low, high := canonicalPair(caller.ID, target.ID)
	res, err := a.db.Exec(
		a.rebind("DELETE FROM friendships WHERE student_a = ? AND student_b = ? AND status = 'accepted'"),
		low, high,
	)
	if err != nil {
		return fmt.Errorf("remove friend: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrFriendNotFound
	}
	return nil
}

func (a *friendAccess) blockUser(studentID, otherProfileID string, now time.Time) error {
	caller, err := a.loadStudent(studentID)
	if err != nil {
		return err
	}
	target, err := a.loadStudentByProfile(otherProfileID)
	if err != nil {
		return err
	}
	if caller.ID == target.ID {
		return ErrFriendSelf
	}
	low, high := canonicalPair(caller.ID, target.ID)
	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("block user: %w", err)
	}
	if _, err := tx.Exec(a.rebind("DELETE FROM friendships WHERE student_a = ? AND student_b = ?"), low, high); err != nil {
		tx.Rollback()
		return fmt.Errorf("block user: %w", err)
	}
	if _, err := tx.Exec(
		a.rebind("INSERT INTO friend_blocks (blocker_id, blocked_id, created_at) VALUES (?, ?, ?) ON CONFLICT DO NOTHING"),
		caller.ID, target.ID, friendUnix(now),
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("block user: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("block user: %w", err)
	}
	return nil
}

func (a *friendAccess) unblockUser(studentID, otherProfileID string) error {
	caller, err := a.loadStudent(studentID)
	if err != nil {
		return err
	}
	target, err := a.loadStudentByProfile(otherProfileID)
	if err != nil {
		return err
	}
	if caller.ID == target.ID {
		return ErrFriendSelf
	}
	if _, err := a.db.Exec(
		a.rebind("DELETE FROM friend_blocks WHERE blocker_id = ? AND blocked_id = ?"),
		caller.ID, target.ID,
	); err != nil {
		return fmt.Errorf("unblock user: %w", err)
	}
	return nil
}

func (a *friendAccess) listFriends(studentID string) ([]FriendSummary, error) {
	caller, err := a.loadStudent(studentID)
	if err != nil {
		return nil, err
	}
	if friendVisibility(caller.Settings) == "nobody" {
		return []FriendSummary{}, nil
	}
	rows, err := a.db.Query(
		a.rebind(`SELECT s.id, s.profile_id, s.name, COALESCE(s.username, ''), s.avatar_url, s.settings,
			s.xp_total, f.created_at, f.responded_at,
			(SELECT COUNT(*) FROM concept_progress cp WHERE cp.student_id = s.id AND cp.status = 'MASTERED')
			FROM friendships f
			JOIN students s ON s.id = CASE WHEN f.student_a = ? THEN f.student_b ELSE f.student_a END
			WHERE f.status = 'accepted' AND (f.student_a = ? OR f.student_b = ?)
			ORDER BY s.name, s.id`),
		studentID, studentID, studentID,
	)
	if err != nil {
		return nil, fmt.Errorf("list friends: %w", err)
	}
	defer rows.Close()
	out := []FriendSummary{}
	for rows.Next() {
		var f FriendSummary
		var created, responded int64
		if err := rows.Scan(&f.StudentID, &f.ProfileID, &f.Name, &f.Username, &f.AvatarURL, &f.AvatarSettings, &f.XPTotal, &created, &responded, &f.Mastered); err != nil {
			return nil, fmt.Errorf("list friends: %w", err)
		}
		if friendVisibility(f.AvatarSettings) == "nobody" {
			continue
		}
		if responded > 0 {
			f.FriendsSince = friendTime(responded)
		} else {
			f.FriendsSince = friendTime(created)
		}
		out = append(out, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list friends: %w", err)
	}
	return out, nil
}

func (a *friendAccess) listFriendRequests(studentID string) ([]FriendRequest, []FriendRequest, error) {
	caller, err := a.loadStudent(studentID)
	if err != nil {
		return nil, nil, err
	}
	incoming, outgoing := []FriendRequest{}, []FriendRequest{}
	if friendVisibility(caller.Settings) == "nobody" {
		return incoming, outgoing, nil
	}
	rows, err := a.db.Query(
		a.rebind(`SELECT s.id, s.profile_id, s.name, COALESCE(s.username, ''), s.avatar_url, s.settings, f.requested_by, f.created_at
			FROM friendships f
			JOIN students s ON s.id = CASE WHEN f.student_a = ? THEN f.student_b ELSE f.student_a END
			WHERE f.status = 'pending' AND (f.student_a = ? OR f.student_b = ?)
			ORDER BY f.created_at DESC`),
		studentID, studentID, studentID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("list friend requests: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var r FriendRequest
		var requestedBy string
		var created int64
		if err := rows.Scan(&r.StudentID, &r.ProfileID, &r.Name, &r.Username, &r.AvatarURL, &r.AvatarSettings, &requestedBy, &created); err != nil {
			return nil, nil, fmt.Errorf("list friend requests: %w", err)
		}
		if friendVisibility(r.AvatarSettings) == "nobody" {
			continue
		}
		r.CreatedAt = friendTime(created)
		if requestedBy == studentID {
			r.Direction = FriendDirectionOutgoing
			outgoing = append(outgoing, r)
		} else {
			r.Direction = FriendDirectionIncoming
			incoming = append(incoming, r)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("list friend requests: %w", err)
	}
	return incoming, outgoing, nil
}

func (a *friendAccess) friendProfile(studentID, otherProfileID string) (*FriendProfile, error) {
	caller, err := a.loadStudent(studentID)
	if err != nil {
		return nil, err
	}
	target, err := a.loadStudentByProfile(otherProfileID)
	if err != nil {
		return nil, ErrFriendNotFound
	}
	if caller.ID == target.ID {
		return nil, ErrFriendNotFound
	}
	blocked, err := a.blockedEither(caller.ID, target.ID)
	if err != nil {
		return nil, err
	}
	if blocked || friendVisibility(caller.Settings) == "nobody" || friendVisibility(target.Settings) == "nobody" {
		return nil, ErrFriendNotFound
	}
	f, err := a.friendship(caller.ID, target.ID)
	if err != nil || f.Status != FriendshipAccepted {
		return nil, ErrFriendNotFound
	}
	mastered, err := a.masteredCount(target.ID)
	if err != nil {
		return nil, err
	}
	since := f.CreatedAt
	if !f.RespondedAt.IsZero() {
		since = f.RespondedAt
	}
	return &FriendProfile{
		StudentID:      target.ID,
		ProfileID:      target.ProfileID,
		Name:           target.Name,
		Username:       target.Username,
		AvatarURL:      target.AvatarURL,
		AvatarSettings: target.Settings,
		XPTotal:        target.XPTotal,
		Mastered:       mastered,
		FriendsSince:   since,
	}, nil
}

func (a *friendAccess) friendGraph(studentID, otherProfileID string) ([]FriendGraphNode, error) {
	if _, err := a.friendProfile(studentID, otherProfileID); err != nil {
		return nil, err
	}
	target, err := a.loadStudentByProfile(otherProfileID)
	if err != nil {
		return nil, ErrFriendNotFound
	}
	rows, err := a.db.Query(
		a.rebind("SELECT concept_id, COALESCE(status, 'UNSEEN'), mastered_at FROM concept_progress WHERE student_id = ? AND status <> 'UNSEEN' ORDER BY concept_id"),
		target.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("friend graph: %w", err)
	}
	defer rows.Close()
	out := []FriendGraphNode{}
	for rows.Next() {
		var n FriendGraphNode
		var mastered sql.NullString
		if err := rows.Scan(&n.ConceptID, &n.Status, &mastered); err != nil {
			return nil, fmt.Errorf("friend graph: %w", err)
		}
		if mastered.Valid && mastered.String != "" {
			if at, err := time.Parse(time.RFC3339, mastered.String); err == nil {
				n.MasteredAt = &at
			}
		}
		out = append(out, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("friend graph: %w", err)
	}
	return out, nil
}

func (a *friendAccess) friendsLeaderboard(studentID, view string, week time.Time) ([]FriendLeaderboardRow, error) {
	if view != "weekly" && view != "all_time" {
		return nil, fmt.Errorf("unknown friends leaderboard view %q", view)
	}
	caller, err := a.loadStudent(studentID)
	if err != nil {
		return nil, err
	}
	if friendVisibility(caller.Settings) == "nobody" {
		return []FriendLeaderboardRow{}, nil
	}
	friends, err := a.listFriends(studentID)
	if err != nil {
		return nil, err
	}
	type entry struct {
		studentID, profileID, name, username, avatarURL, settings string
		xpTotal                                                   int
		mastered                                                  int
		isMe                                                      bool
	}
	entries := make([]entry, 0, len(friends)+1)
	entries = append(entries, entry{
		studentID: caller.ID, profileID: caller.ProfileID, name: caller.Name,
		username: caller.Username, avatarURL: caller.AvatarURL, settings: caller.Settings,
		xpTotal: caller.XPTotal, isMe: true,
	})
	for _, f := range friends {
		entries = append(entries, entry{
			studentID: f.StudentID, profileID: f.ProfileID, name: f.Name,
			username: f.Username, avatarURL: f.AvatarURL, settings: f.AvatarSettings,
			xpTotal: f.XPTotal, mastered: f.Mastered,
		})
	}
	since := week
	if since.IsZero() {
		since = time.Now().UTC().AddDate(0, 0, -7)
	}
	out := make([]FriendLeaderboardRow, 0, len(entries))
	for _, e := range entries {
		if friendLeaderboardHidden(e.settings) {
			continue
		}
		xp := e.xpTotal
		if view == "weekly" {
			xp, err = a.weeklyXP(e.studentID, since)
			if err != nil {
				return nil, err
			}
		}
		mastered := e.mastered
		if e.isMe {
			mastered, err = a.masteredCount(e.studentID)
			if err != nil {
				return nil, err
			}
		}
		out = append(out, FriendLeaderboardRow{
			StudentID: e.studentID, ProfileID: e.profileID, Name: e.name,
			Username: e.username, AvatarURL: e.avatarURL, AvatarSettings: e.settings,
			XP: xp, Mastered: mastered, IsMe: e.isMe,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].XP != out[j].XP {
			return out[i].XP > out[j].XP
		}
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].StudentID < out[j].StudentID
	})
	return out, nil
}

func (a *friendAccess) masteredCount(studentID string) (int, error) {
	var n int
	if err := a.db.QueryRow(
		a.rebind("SELECT COUNT(*) FROM concept_progress WHERE student_id = ? AND status = 'MASTERED'"),
		studentID,
	).Scan(&n); err != nil {
		return 0, fmt.Errorf("count mastered concepts: %w", err)
	}
	return n, nil
}

func (a *friendAccess) weeklyXP(studentID string, since time.Time) (int, error) {
	var n int
	if err := a.db.QueryRow(
		a.rebind("SELECT COALESCE(SUM(amount), 0) FROM xp_events WHERE student_id = ? AND created_at >= ?"),
		studentID, since.UTC().UnixNano(),
	).Scan(&n); err != nil {
		return 0, fmt.Errorf("sum weekly XP: %w", err)
	}
	return n, nil
}

func (a *friendAccess) rotateFriendCode(studentID string, now time.Time) (string, error) {
	_ = now
	for attempt := 0; attempt < 5; attempt++ {
		code, err := newFriendCode()
		if err != nil {
			return "", err
		}
		res, err := a.db.Exec(a.rebind("UPDATE students SET friend_code = ? WHERE id = ?"), code, studentID)
		if err != nil {
			if IsUniqueViolation(err) {
				continue
			}
			return "", fmt.Errorf("rotate friend code: %w", err)
		}
		n, _ := res.RowsAffected()
		if n == 0 {
			return "", ErrFriendNotFound
		}
		return code, nil
	}
	return "", fmt.Errorf("rotate friend code: could not allocate a unique code")
}

func (a *friendAccess) getFriendIdentity(studentID string) (FriendIdentity, error) {
	var id FriendIdentity
	err := a.db.QueryRow(
		a.rebind("SELECT profile_id, friend_code FROM students WHERE id = ?"),
		studentID,
	).Scan(&id.ProfileID, &id.FriendCode)
	if errors.Is(err, sql.ErrNoRows) {
		return FriendIdentity{}, ErrFriendNotFound
	}
	if err != nil {
		return FriendIdentity{}, fmt.Errorf("get friend identity: %w", err)
	}
	return id, nil
}

// FriendStore methods on the concrete stores. They keep the interface stable
// while the implementation stays shared between SQLite and PostgreSQL.

func (s *SQLiteStore) ResolveFriendTarget(query string) (string, error) {
	return s.friends().resolveFriendTarget(query)
}

func (s *SQLiteStore) GetFriendIdentity(studentID string) (FriendIdentity, error) {
	return s.friends().getFriendIdentity(studentID)
}

func (s *SQLiteStore) SendFriendRequest(studentID, targetProfileID string, now time.Time) (Friendship, error) {
	return s.friends().sendFriendRequest(studentID, targetProfileID, now)
}

func (s *SQLiteStore) AcceptFriendRequest(studentID, otherProfileID string, now time.Time) (Friendship, error) {
	return s.friends().acceptFriendRequest(studentID, otherProfileID, now)
}

func (s *SQLiteStore) DeclineFriendRequest(studentID, otherProfileID string, now time.Time) error {
	return s.friends().declineFriendRequest(studentID, otherProfileID, now)
}

func (s *SQLiteStore) CancelFriendRequest(studentID, otherProfileID string, now time.Time) error {
	return s.friends().cancelFriendRequest(studentID, otherProfileID, now)
}

func (s *SQLiteStore) RemoveFriend(studentID, otherProfileID string, now time.Time) error {
	return s.friends().removeFriend(studentID, otherProfileID, now)
}

func (s *SQLiteStore) BlockUser(studentID, otherProfileID string, now time.Time) error {
	return s.friends().blockUser(studentID, otherProfileID, now)
}

func (s *SQLiteStore) UnblockUser(studentID, otherProfileID string) error {
	return s.friends().unblockUser(studentID, otherProfileID)
}

func (s *SQLiteStore) ListFriends(studentID string) ([]FriendSummary, error) {
	return s.friends().listFriends(studentID)
}

func (s *SQLiteStore) ListFriendRequests(studentID string) ([]FriendRequest, []FriendRequest, error) {
	return s.friends().listFriendRequests(studentID)
}

func (s *SQLiteStore) FriendProfile(studentID, otherProfileID string) (*FriendProfile, error) {
	return s.friends().friendProfile(studentID, otherProfileID)
}

func (s *SQLiteStore) FriendGraph(studentID, otherProfileID string) ([]FriendGraphNode, error) {
	return s.friends().friendGraph(studentID, otherProfileID)
}

func (s *SQLiteStore) FriendsLeaderboard(studentID, view string, week time.Time) ([]FriendLeaderboardRow, error) {
	return s.friends().friendsLeaderboard(studentID, view, week)
}

func (s *SQLiteStore) RotateFriendCode(studentID string, now time.Time) (string, error) {
	return s.friends().rotateFriendCode(studentID, now)
}

func (s *PostgresStore) ResolveFriendTarget(query string) (string, error) {
	return s.friends().resolveFriendTarget(query)
}

func (s *PostgresStore) GetFriendIdentity(studentID string) (FriendIdentity, error) {
	return s.friends().getFriendIdentity(studentID)
}

func (s *PostgresStore) SendFriendRequest(studentID, targetProfileID string, now time.Time) (Friendship, error) {
	return s.friends().sendFriendRequest(studentID, targetProfileID, now)
}

func (s *PostgresStore) AcceptFriendRequest(studentID, otherProfileID string, now time.Time) (Friendship, error) {
	return s.friends().acceptFriendRequest(studentID, otherProfileID, now)
}

func (s *PostgresStore) DeclineFriendRequest(studentID, otherProfileID string, now time.Time) error {
	return s.friends().declineFriendRequest(studentID, otherProfileID, now)
}

func (s *PostgresStore) CancelFriendRequest(studentID, otherProfileID string, now time.Time) error {
	return s.friends().cancelFriendRequest(studentID, otherProfileID, now)
}

func (s *PostgresStore) RemoveFriend(studentID, otherProfileID string, now time.Time) error {
	return s.friends().removeFriend(studentID, otherProfileID, now)
}

func (s *PostgresStore) BlockUser(studentID, otherProfileID string, now time.Time) error {
	return s.friends().blockUser(studentID, otherProfileID, now)
}

func (s *PostgresStore) UnblockUser(studentID, otherProfileID string) error {
	return s.friends().unblockUser(studentID, otherProfileID)
}

func (s *PostgresStore) ListFriends(studentID string) ([]FriendSummary, error) {
	return s.friends().listFriends(studentID)
}

func (s *PostgresStore) ListFriendRequests(studentID string) ([]FriendRequest, []FriendRequest, error) {
	return s.friends().listFriendRequests(studentID)
}

func (s *PostgresStore) FriendProfile(studentID, otherProfileID string) (*FriendProfile, error) {
	return s.friends().friendProfile(studentID, otherProfileID)
}

func (s *PostgresStore) FriendGraph(studentID, otherProfileID string) ([]FriendGraphNode, error) {
	return s.friends().friendGraph(studentID, otherProfileID)
}

func (s *PostgresStore) FriendsLeaderboard(studentID, view string, week time.Time) ([]FriendLeaderboardRow, error) {
	return s.friends().friendsLeaderboard(studentID, view, week)
}

func (s *PostgresStore) RotateFriendCode(studentID string, now time.Time) (string, error) {
	return s.friends().rotateFriendCode(studentID, now)
}

// hasVerifiedIdentity reports whether any linked provider attests the mailbox.
// Google sign-in sets this on the identity row, so an account whose
// students.email_verified flag was never updated can still send friend
// requests without being asked to verify an address the provider already
// verified.
func (a *friendAccess) hasVerifiedIdentity(studentID string) bool {
	var one int
	err := a.db.QueryRow(
		a.rebind("SELECT 1 FROM identities WHERE student_id = ? AND email_verified = 1 LIMIT 1"),
		studentID,
	).Scan(&one)
	return err == nil
}
