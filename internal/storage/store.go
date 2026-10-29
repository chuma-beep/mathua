package storage

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/ai"
)

// Data types

type Student struct {
	ID                  string
	Name                string
	Username            string
	PasswordHash        string
	CourseID            string
	XPTotal             int
	XPToday             int
	XPTodayDate         string
	DiagnosticCompleted bool
	DailyXPGoal         int
	Settings            string
	CreatedAt           time.Time
	League              string
	LeagueWeek          string
	LeagueMoved         int
	ShareToken          string
	Email               string
	EmailVerified       bool
	GoogleID            string
	AvatarURL           string
}

// Identity links an external OAuth provider account to a student.
// (provider, provider_id) is globally unique across all providers.
type Identity struct {
	Provider      string
	ProviderID    string
	StudentID     string
	Email         string
	EmailVerified bool
	LinkedAt      time.Time
}

type ConceptProgress struct {
	StudentID       string     `json:"student_id,omitempty"`
	ConceptID       string     `json:"concept_id,omitempty"`
	Status          string     `json:"status"` // UNSEEN | LEARNING | PRACTICING | MASTERED
	Streak          int        `json:"streak"`
	BestStreak      int        `json:"best_streak"`
	AvgResponseTime float64    `json:"avg_response_time"`
	Attempts        int        `json:"attempts"`
	LastAttempted   *time.Time `json:"last_attempted,omitempty"`
	LastReviewed    *time.Time `json:"last_reviewed,omitempty"`
	NextReviewDue   *time.Time `json:"next_review_due,omitempty"`
	SM2Repetitions  int        `json:"sm2_repetitions,omitempty"`
	SM2Interval     int        `json:"sm2_interval,omitempty"`
	SM2EFactor      float64    `json:"sm2_efactor,omitempty"`
	MasteredAt      *time.Time `json:"mastered_at,omitempty"`
	WeaknessScore   float64    `json:"weakness_score,omitempty"`
	// PlacementSeeded marks a row written by the diagnostic's placement pass
	// (ApplyGoalResults), not by demonstrated work. It grants reachability only:
	// a placement-seeded concept satisfies its successors' prerequisites so the
	// learner can start at their frontier, but it is never MASTERED and never
	// counts toward mastery, XP, the quiz gate or review. Placement is not
	// evidence (ADR-008, ADR-050).
	PlacementSeeded bool `json:"placement_seeded,omitempty"`
	// IsProvisional marks conditional completion from diagnostic placement.
	// When a learner struggles on a dependent concept, provisional prerequisites
	// demote to LEARNING and schedule for remedial review.
	IsProvisional bool `json:"is_provisional,omitempty"`
	// FSRS-6 state fields
	FSRSStability  *float64   `json:"fsrs_stability,omitempty"`
	FSRSDifficulty *float64   `json:"fsrs_difficulty,omitempty"`
	FSRSLastReview *time.Time `json:"fsrs_last_review,omitempty"`
}

type Session struct {
	ID        string
	StudentID string
	StartedAt time.Time
}

type ActiveSession struct {
	SessionID string
	StudentID string
	ConceptID string
	// Difficulty the currently-served question was generated at. Zero means the question
	// was not generated at this difficulty, and the attempt is written with a NULL.
	Difficulty     *float64
	ConceptName    string
	ExpectedAnswer string
	AttemptID      string
	Question       string
	Explanation    string
	Diagram        string
	IsReview       bool
	Answered       bool
	LastConceptID  string
	SessionReview  int
	SessionNew     int
	KPOrdinal      int
	UpdatedAt      time.Time
}

// AttemptRow is the minimum a replay needs from the attempts table. The XP
// amount was never stored there, but a correct attempt's award depends only on
// the concept's time threshold, which the engine can read from the curriculum.
type AttemptRow struct {
	ID        int64
	StudentID string
	ConceptID string
	Correct   bool
	Source    string
}

type AttemptEntry struct {
	// ID is the persistent row id (set on read; ignored on insert).
	ID             int64     `json:"id,omitempty"`
	SessionID      string    `json:"session_id"`
	StudentID      string    `json:"student_id"`
	ConceptID      string    `json:"concept_id"`
	Answer         string    `json:"answer"`
	Expected       string    `json:"expected"`
	Correct        bool      `json:"correct"`
	ElapsedSeconds float64   `json:"elapsed_seconds"`
	Timestamp      time.Time `json:"timestamp"`
	// Question is the served question text (for mistakes review).
	// Source is diagnostic|quiz|practice|review. Explanation is the
	// teaching content shown with feedback. Pre-migration rows carry "".
	Question    string `json:"question"`
	Source      string `json:"source"`
	Explanation string `json:"explanation"`
	// Difficulty the generator served this question at (0.3-1.0). Nil means unknown:
	// a pre-migration row, or a question that was not generated. Never fabricate it,
	// "not recorded" and "recorded as trivial" are different facts, and an evidence model
	// that cannot tell them apart is worse than one that knows less.
	Difficulty *float64 `json:"difficulty,omitempty"`
}

// StudySubmission is the durable record of one accepted study answer. The
// submission id is chosen by the client once per logical answer, so a retry
// carries the same id and can replay the stored response without writing
// another attempt or another XP event.
type StudySubmission struct {
	StudentID    string
	SubmissionID string
	Operation    string
	RequestHash  string
	ResponseJSON string
	AttemptID    int64
	CreatedAt    time.Time
}

// StudySubmissionCommit is one atomic study-answer write. The caller grades
// first, then hands the finished progress row and response here; the store
// never grades.
type StudyEffectInput struct {
	Effect      string
	PayloadJSON string
}

type StudyAnswerEffect struct {
	ID            int64
	StudentID     string
	SubmissionID  string
	Effect        string
	PayloadJSON   string
	Attempts      int
	CreatedAt     time.Time
	NextAttemptAt time.Time
}

type StudySubmissionCommit struct {
	StudentID    string
	SubmissionID string
	Operation    string
	RequestHash  string
	ResponseJSON string
	Attempt      AttemptEntry
	Progress     *ConceptProgress
	TopicSpeed   *TopicSpeed
	XPAmount     int
	XPReason     string
	Effects      []StudyEffectInput
}

// StudySubmissionOutcome reports whether the commit wrote a new answer or
// replayed the response stored by an earlier identical submission.
type StudySubmissionOutcome struct {
	Replayed     bool
	ResponseJSON string
	AttemptID    int64
}

var (
	// ErrStudySubmissionPayloadMismatch means the same idempotency key was
	// reused with a different concept, answer, elapsed time or question.
	ErrStudySubmissionPayloadMismatch = errors.New("study submission payload does not match the stored submission")
	// ErrStudySubmissionInProgress means a concurrent transaction owns the key
	// and has not stored its response yet. The caller should retry.
	ErrStudySubmissionInProgress = errors.New("study submission is still in progress")
)

// StudySubmissionStore is implemented by the SQL stores. It is separate from
// Repository so test doubles and the quiz path do not need to implement it.
type StudySubmissionStore interface {
	CommitStudySubmission(c StudySubmissionCommit) (StudySubmissionOutcome, error)
	GetStudySubmission(studentID, operation, submissionID string) (*StudySubmission, error)
}

// StudyEffectStore is the durable outbox for effects that must run after the
// answer transaction. Claim is atomic, so two workers cannot process one row.
type StudyEffectStore interface {
	ClaimStudyEffects(studentID string, now time.Time, limit int) ([]StudyAnswerEffect, error)
	CompleteStudyEffect(id int64) error
	FailStudyEffect(id int64, reason string, retryAt time.Time) error
	PruneStudyEffects(olderThan time.Time, limit int) (int64, error)
}

// parseAttemptTimestamp parses attempt timestamps from either store.
// SQLite writes RFC3339; Postgres writes now()::text ("2006-01-02
// 15:04:05.999999-07"), which RFC3339 parsing rejects.
func parseAttemptTimestamp(ts string) time.Time {
	if ts == "" {
		return time.Time{}
	}
	if t, err := time.Parse(time.RFC3339, ts); err == nil {
		return t
	}
	for _, layout := range []string{
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999999-07",
		"2006-01-02 15:04:05-07:00",
		"2006-01-02 15:04:05-07",
		"2006-01-02T15:04:05.999999999-07:00",
	} {
		if t, err := time.Parse(layout, ts); err == nil {
			return t
		}
	}
	return time.Time{}
}

// parseQuizCompletedAt parses the completion string each store writes.
// SQLite defaults to datetime('now') ("2006-01-02 15:04:05", no zone, UTC),
// Postgres writes now()::text, and older rows may carry RFC3339.
func parseQuizCompletedAt(ts string) (time.Time, bool) {
	if ts == "" {
		return time.Time{}, false
	}
	if t := parseAttemptTimestamp(ts); !t.IsZero() {
		return t, true
	}
	for _, layout := range []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	} {
		if t, err := time.Parse(layout, ts); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

type Question struct {
	ID          int
	ConceptID   string
	Question    string
	Answer      string
	Explanation string
	Source      string
	Difficulty  float64
}

type LeaderboardRow struct {
	StudentID      string
	Name           string
	Username       string
	AvatarURL      string
	AvatarSettings string
	TotalMastered  int
	WeeklyMastered int
	// XP is the score for the requested view: the week's ledger sum for
	// "weekly", lifetime total for "all_time". XPTotal is always the lifetime
	// total, so a weekly board can break a tie on earned history.
	XP        int
	XPTotal   int
	CreatedAt string
}

// AvatarRef is the public dicebear pick surfaced on leaderboard entries.
type AvatarRef struct {
	Style string `json:"style"`
	Seed  string `json:"seed"`
}

// avatarSettingsJSON mirrors the opaque settings keys the frontend owns
// (web/next-app UserSettings); only leaderboard-relevant bits are decoded.
type avatarSettingsJSON struct {
	AvatarCustom   bool       `json:"avatar_custom"`
	AvatarDicebear *AvatarRef `json:"avatar_dicebear"`
	AvatarVersion  int        `json:"avatar_version"`
}

// AvatarBits extracts leaderboard-relevant avatar data from a student's
// opaque settings blob: custom-upload flag, dicebear pick, photo version.
// Corrupt/empty blobs yield zero values (initial-avatar fallback downstream).
func AvatarBits(raw string) (custom bool, pick *AvatarRef, version int) {
	if strings.TrimSpace(raw) == "" {
		return false, nil, 0
	}
	var s avatarSettingsJSON
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		return false, nil, 0
	}
	return s.AvatarCustom, s.AvatarDicebear, s.AvatarVersion
}

// AvatarPreset extracts the legacy preset number, which predates the DiceBear
// pick and is still honoured for accounts created before those were offered.
// It is separate from AvatarBits so existing leaderboard callers do not pay for
// a field they no longer render.
func AvatarPreset(raw string) *int {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var s struct {
		AvatarPreset *int `json:"avatar_preset"`
	}
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		return nil
	}
	return s.AvatarPreset
}

type DailyActivity struct {
	Date      string   `json:"date"`
	Questions int      `json:"questions"`
	Correct   int      `json:"correct"`
	Concepts  []string `json:"concepts"`
}

// LearningEvent is one substantiated item in a learner's history. It is built
// from rows that exist, never from a current aggregate: an answer, a concept
// that reached mastery, a quiz completion or a diagnostic completion. The Ref
// field is stable within its Kind, so a page can key on Kind+Ref without
// inventing an id.
type LearningEvent struct {
	Kind      string    `json:"kind"`
	Ref       string    `json:"ref"`
	At        time.Time `json:"at"`
	ConceptID string    `json:"concept_id,omitempty"`
	Source    string    `json:"source,omitempty"`
	Correct   bool      `json:"correct,omitempty"`
	XP        int       `json:"xp,omitempty"`
}

type TopicSpeed struct {
	StudentID     string  `json:"student_id"`
	ConceptID     string  `json:"concept_id"`
	EFactor       float64 `json:"efactor"`
	Interval      int     `json:"interval"`
	Repetitions   int     `json:"repetitions"`
	LearningSpeed float64 `json:"learning_speed"`
}

// QuizCompletion records one finished mastery-check quiz (Batch 1 gate).
type QuizCompletion struct {
	StudentID   string `json:"student_id"`
	CompletedAt string `json:"completed_at"`
	XPTotal     int    `json:"xp_total"`
}

// TimedQuizBaseline is the XP floor from which a learner's timed-quiz gate is
// measured. Existing learners are seeded at rollout, so historical XP never
// creates an immediate timed quiz. A new account starts at zero.
type TimedQuizBaseline struct {
	StudentID  string
	BaselineXP int
	BaselineAt time.Time
}

// TimedQuizRun is one durable timed quiz. The server owns the deadline and the
// question set; a client can resume it after a reload or a server restart.
type TimedQuizRun struct {
	ID              string
	StudentID       string
	StartedAt       time.Time
	DeadlineAt      time.Time
	CompletedAt     *time.Time
	Status          string
	QuestionCount   int
	CorrectCount    int
	UnansweredCount int
	XPEarned        int
	CreatedAt       time.Time
}

// TimedQuizItem is one served question. ExpectedAnswer is deliberately stored
// server-side and is never sent with the question.
type TimedQuizItem struct {
	ID             string
	RunID          string
	ConceptID      string
	Position       int
	Question       string
	ExpectedAnswer string
	Explanation    string
	GradingType    string
	Difficulty     *float64
	ServedAt       time.Time
	AnsweredAt     *time.Time
	ElapsedSeconds float64
	Answer         string
	Correct        bool
	Status         string
	SubmissionID   string
	XPEarned       int
}

// ReviewEvent is the durable record of one spaced-review attempt and the FSRS
// transition it caused. AttemptID is unique, so a replayed submission cannot
// create a second event or a second scheduler update.
type ReviewEvent struct {
	StudentID        string
	ConceptID        string
	AttemptID        string
	Rating           int
	Correct          bool
	ElapsedSeconds   float64
	ReviewedAt       time.Time
	StabilityBefore  *float64
	DifficultyBefore *float64
	DueBefore        *time.Time
	StabilityAfter   *float64
	DifficultyAfter  *float64
	DueAfter         *time.Time
}

// TimedQuizStore is the durable timed-quiz and review-event surface. It is
// separate from Repository so existing test doubles do not need every method.
type TimedQuizStore interface {
	EnsureTimedQuizBaseline(studentID string) (TimedQuizBaseline, error)
	GetTimedQuizBaseline(studentID string) (*TimedQuizBaseline, error)
	TimedQuizEligibleXP(studentID string, since time.Time) (int, error)
	CreateTimedQuizRun(run TimedQuizRun, items []TimedQuizItem) error
	GetActiveTimedQuizRun(studentID string) (*TimedQuizRun, error)
	GetTimedQuizRun(id, studentID string) (*TimedQuizRun, error)
	LatestTimedQuizRun(studentID string) (*TimedQuizRun, error)
	GetTimedQuizItems(runID string) ([]TimedQuizItem, error)
	ClaimTimedQuizItem(itemID, runID, submissionID string) (bool, error)
	UpdateTimedQuizItem(itemID, runID, answer, submissionID string, correct bool, elapsed float64, xp int, status string) (bool, error)
	FinalizeTimedQuizRun(runID, studentID string, now time.Time) (*TimedQuizRun, error)
	RecordReviewEvent(event ReviewEvent) error
}

// ReviewEventStore is the narrow review-event writer used by the engine.
type ReviewEventStore interface {
	RecordReviewEvent(event ReviewEvent) error
}

// ReviewSubmissionCommit is one spaced-review answer. The review event, the
// graded attempt, the progress transition and the XP award are committed
// together so an FSRS change can never exist without its audit row.
type ReviewSubmissionCommit struct {
	StudentID    string
	SubmissionID string
	Event        ReviewEvent
	Attempt      AttemptEntry
	Progress     *ConceptProgress
	TopicSpeed   *TopicSpeed
	XPAmount     int
	XPReason     string
}

// ReviewSubmissionOutcome reports the attempt created by a review commit, or
// that an event with the same idempotency key already exists.
type ReviewSubmissionOutcome struct {
	AttemptID int64
	Replayed  bool
}

// ReviewSubmissionStore is the atomic review-answer writer. It is separate
// from Repository so existing test doubles do not need the method.
type ReviewSubmissionStore interface {
	CommitReviewSubmission(ReviewSubmissionCommit) (ReviewSubmissionOutcome, error)
}

// Announcement is a product update shown in the learner notification center.
// Body is Markdown rendered with the same sanitizer as lesson content.
type Announcement struct {
	ID          string
	Title       string
	Body        string
	Published   bool
	CreatedAt   time.Time
	PublishedAt time.Time
	CreatedBy   string
	Read        bool
}

// AnnouncementRead records that a learner has seen a published announcement.
type AnnouncementRead struct {
	StudentID      string
	AnnouncementID string
	ReadAt         time.Time
}

// AnnouncementStore is the announcement surface. It is separate from
// Repository so existing test doubles do not need every method.
type AnnouncementStore interface {
	ListPublishedAnnouncements(studentID string) ([]Announcement, error)
	MarkAnnouncementsRead(studentID string, now time.Time) error
	ClearAnnouncements(studentID string, clearedAt time.Time) error
	ListAnnouncements() ([]Announcement, error)
	GetAnnouncement(id string) (*Announcement, error)
	CreateAnnouncement(a Announcement) (*Announcement, error)
	UpdateAnnouncement(id, title, body string, published bool, now time.Time) (*Announcement, error)
}

// MaxRemedialQueue caps the per-student remedial queue.
const MaxRemedialQueue = 10

// RemedialQueueEntry records one concept queued for targeted remediation.
type RemedialQueueEntry struct {
	StudentID string `json:"student_id"`
	ConceptID string `json:"concept_id"`
	Priority  int    `json:"priority"`
	CreatedAt int64  `json:"created_at"`
}

// QuestionReport is a user complaint about a question, explanation,
// lesson body, worked example, or diagram.
//
// A report is a moderation record, not a queue item: resolving it sets Status and the outcome
// fields and leaves the row in place. Nothing in the system deletes one, because "what was
// reported and what we decided" is the evidence a moderation log exists to hold.
type QuestionReport struct {
	ID               int64  `json:"id"`
	ReporterID       string `json:"reporter_id"`
	ReporterUsername string `json:"reporter_username,omitempty"`
	ReporterName     string `json:"reporter_name,omitempty"`
	ConceptID        string `json:"concept_id"`
	Kind             string `json:"kind"` // question | explanation | lesson_body | worked_example | diagram | message
	Question         string `json:"question"`
	Expected         string `json:"expected"`
	Explanation      string `json:"explanation"`
	LessonID         string `json:"lesson_id"`
	Source           string `json:"source"`
	SessionID        string `json:"session_id"`
	AttemptID        string `json:"attempt_id"`
	// The answer the learner actually typed, and the format the field offered
	// when they typed it. A moderator could previously see the question and the
	// expected answer but not what was on screen, so a report about an
	// untypeable answer could not be reproduced: the input mode and the hint
	// line are what the learner faced, and the answer given is what the grader
	// judged. Empty on reports filed before these were captured.
	AnswerGiven     string `json:"answer_given"`
	GradingType     string `json:"grading_type"`
	AnswerHint      string `json:"answer_hint"`
	AnswerInputMode string `json:"answer_input_mode"`
	Reason          string `json:"reason"` // wrong_answer | bad_explanation | unclear | formatting | other
	Detail          string `json:"detail"`
	// Status is the moderation state: open | reviewing | resolved | dismissed.
	Status string `json:"status"`
	// Resolution is the moderator's reason for the decision, written when the report leaves
	// open/reviewing. Kept separate from Detail, which is the reporter's own words.
	Resolution string `json:"resolution"`
	// ResolvedBy and ResolvedAt name who decided and when. Text, not a foreign key, for the same
	// reason the audit actor is: the record must outlive the account it names.
	ResolvedBy string    `json:"resolved_by"`
	ResolvedAt string    `json:"resolved_at"`
	CreatedAt  time.Time `json:"created_at"`
	// OriginReportID is the id this report carried in the database it was exported from, and
	// "" for a report filed here. The row id is not that number: SQLite does not reuse an id
	// from a rolled-back insert, so the two sequences drift as soon as anything fails, and a
	// local "#94" can name a different report from the production #94 a colleague is looking
	// at. Anything quoted outside the local database uses this field.
	OriginReportID string `json:"origin_report_id,omitempty"`
	// HasImage says whether a report_images row exists. The bytes are served
	// separately behind the moderation permission, so the queue payload stays
	// small and an image cannot be fetched from a plain report list.
	HasImage bool `json:"has_image"`
	// Unread is set by the learner-facing history endpoint when a resolved or
	// dismissed report has changed since that account last read the list. It is
	// not stored on the report: one row can be unread for its reporter and read
	// for a moderator, and the report itself has no reader.
	Unread bool `json:"unread,omitempty"`
}

// ReportImage is one optional screenshot attached to a report. Bytes live in
// the same database as the report so an upload cannot be separated from the
// moderation record it belongs to.
type ReportImage struct {
	ContentType string
	Bytes       []byte
}

// ReportImageStore is separate from Repository because only the report
// endpoints need it; the moderation row and its image still share one store.
type ReportImageStore interface {
	SetReportImage(reportID int64, contentType string, data []byte) error
	GetReportImage(reportID int64) (ReportImage, bool, error)
}

// AdminInvitation is a pending offer of a staff role to an account.
//
// The raw token is never stored. Only its SHA-256 hash is, exactly as with password resets: a
// database read must not yield a usable credential, so a leaked backup cannot be replayed as
// access. The row records who offered it and what it grants, so accepting it is a decision the
// audit log can attribute to a person.
type AdminInvitation struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Role      Role   `json:"role"`
	InvitedBy string `json:"invited_by"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
	// AcceptedAt/RevokedAt are "" while the invitation is pending. A row is never deleted: a
	// revoked or expired invitation is history.
	AcceptedAt string `json:"accepted_at"`
	AcceptedBy string `json:"accepted_by"`
	RevokedAt  string `json:"revoked_at"`
	RevokedBy  string `json:"revoked_by"`
}

// ErrRoleUnchanged is returned when an account already holds the requested role.
//
// Its own sentinel rather than a success, because the caller has to answer two different
// questions: did the role change (no), and is there an event for it (also no). Reporting
// success with a changed flag would work, but the audit write is now inside the same
// transaction as the mutation, so "nothing happened" has to be a real outcome rather than a
// return value the caller can forget to check.
var ErrRoleUnchanged = errors.New("account already has that role")

// RoleChange is what one audited role change did.
type RoleChange struct {
	// Changed is false only when the account already held the requested role, which is also the
	// case that produces no audit event.
	Changed bool
	// Action names the audit event this change wrote, or "" when it wrote none. Carried on the
	// result so a caller can confirm which event to expect without re-reading the log and so
	// the CLI can print it rather than reconstructing the name from the requested role.
	Action string
	Before AdminUser
	After  AdminUser
}

// ErrLastAdmin is returned when a demotion would leave the system with no administrator.
//
// It is a sentinel rather than a bool so a caller cannot confuse "refused" with "no such
// user", the two are reported differently, and the first one is a rule while the second is
// a lookup miss.
var ErrLastAdmin = errors.New("cannot demote the last administrator")

// ErrLastOwner is returned when a change would leave the system with no owner.
//
// The invariant moved up the hierarchy in Admin V2. V1 guarded the last `admin`; the role that
// can hand out every other role is now `owner`, so that is the one the system may never run out
// of. An admin or moderator count of zero is recoverable, an owner can always create one, but
// an owner count of zero is not, because nothing but an owner can mint the first one back.
var ErrLastOwner = errors.New("cannot remove the last owner")

// ErrReportUnchanged is returned when a moderation decision would not change the report's
// status. Like ErrRoleUnchanged it means "no change, no event": reopening a report that is
// already open is not a moderation action, and recording it would put a fact in the trail that
// did not happen.
var ErrReportUnchanged = errors.New("report is already in that status")

// Role is what an account may do beyond learning.
//
// It lives on the account and nowhere else. It is not derivable from XP, mastery, level,
// curriculum progress, email domain or anything the client sends, because every one of those
// is either something a learner controls or something they can change without meaning to. The
// only thing that grants it is an explicit write by an administrator or the operator.
//
// Four values, ordered least to most privileged: student, moderator, admin, owner. The ordering
// matters, every guard that stops someone acting above their authority is a comparison of
// ranks, and having one definition of "above" is what keeps the guards from disagreeing.
type Role string

const (
	RoleStudent   Role = "student"
	RoleModerator Role = "moderator"
	RoleAdmin     Role = "admin"
	RoleOwner     Role = "owner"
)

// roleRank orders the roles by privilege. An unknown role ranks below student, so anything that
// slips through NormalizeRole is treated as no privilege rather than as some privilege nobody
// has audited.
func roleRank(r Role) int {
	switch r {
	case RoleOwner:
		return 4
	case RoleAdmin:
		return 3
	case RoleModerator:
		return 2
	case RoleStudent:
		return 1
	}
	return 0
}

// AtLeast reports whether r is at least as privileged as other.
func (r Role) AtLeast(other Role) bool {
	return roleRank(r) >= roleRank(other)
}

// IsStaff reports whether the role grants any administrative surface at all. It is the check
// `/api/auth/me` uses to decide whether to offer the Admin entry, and the one the middleware
// uses to decide whether a caller is staff before it consults a specific permission.
func (r Role) IsStaff() bool {
	return roleRank(r) >= roleRank(RoleModerator)
}

// NormalizeRole maps any stored or supplied string to a role, defaulting to student.
//
// Anything unrecognised becomes `student`, which is the only safe direction: a future role
// name, a typo, a hand-edited database cell and an empty column all read as "no privilege"
// rather than as "some privilege nobody has audited yet". Callers that need to reject bad
// input use ParseRole.
func NormalizeRole(s string) Role {
	switch Role(s) {
	case RoleOwner:
		return RoleOwner
	case RoleAdmin:
		return RoleAdmin
	case RoleModerator:
		return RoleModerator
	default:
		return RoleStudent
	}
}

// ParseRole is NormalizeRole for input that must be one of the values exactly.
//
// Role assignment goes through this, so a PATCH body carrying "Admin", "ADMIN" or "root" is
// rejected rather than silently becoming something. Silently accepting it would mean the
// response said one thing and the row said another.
func ParseRole(s string) (Role, bool) {
	switch Role(s) {
	case RoleStudent:
		return RoleStudent, true
	case RoleModerator:
		return RoleModerator, true
	case RoleAdmin:
		return RoleAdmin, true
	case RoleOwner:
		return RoleOwner, true
	}
	return RoleStudent, false
}

// AdminUser is the non-sensitive account projection administration reads.
//
// Deliberately not `Student`. That struct carries PasswordHash, ShareToken and the learner's
// XP and streak, and an admin listing is not a place any of those should be reachable, not
// because an admin has no business seeing them, but because a projection that cannot express a
// field cannot leak it by a later edit. Everything here is something an administrator would
// be shown anyway, and nothing here is derived from learning state.
type AdminUser struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      Role   `json:"role"`
	CreatedAt string `json:"created_at"`
	// HasPassword distinguishes "signed up with a password" from "arrived through an OAuth
	// provider". It is a boolean about a fact, never the hash.
	HasPassword bool `json:"has_password"`
}

// AdminAuditRecord is one privileged mutation, written once and never updated.
//
// ActorID is deliberately not a foreign key. The point of an audit trail is that it survives
// the account being deleted, and a foreign key to `students` would delete the evidence along
// with the subject. The actor's id is kept as text so the record still names who did it.
type AdminAuditRecord struct {
	ID         int64  `json:"id"`
	ActorID    string `json:"actor_id"`
	Action     string `json:"action"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	// BeforeJSON and AfterJSON are the entity's non-sensitive state either side of the change,
	// or "" when there was none. They are built from AdminUser, never from Student, so a
	// password hash cannot reach them by accident.
	BeforeJSON string `json:"before_json"`
	AfterJSON  string `json:"after_json"`
	CreatedAt  string `json:"created_at"`
}

// Admin action names. A closed set on purpose: the audit log is a vocabulary, and a
// free-text action column turns into prose nobody can filter on.
const (
	ActionAdminPromoteUser = "ADMIN_PROMOTE_USER"
	ActionAdminDemoteUser  = "ADMIN_DEMOTE_USER"
	// Admin V2 additions.
	ActionAdminRoleChange     = "ADMIN_ROLE_CHANGE" // any role change that is not a plain promote/demote
	ActionAdminInviteCreate   = "ADMIN_INVITE_CREATE"
	ActionAdminInviteRevoke   = "ADMIN_INVITE_REVOKE"
	ActionAdminInviteAccept   = "ADMIN_INVITE_ACCEPT"
	ActionReportStatusChange  = "REPORT_STATUS_CHANGE"
	ActionContentStatusChange = "CONTENT_STATUS_CHANGE"
	ActionAnnouncementCreate  = "ANNOUNCEMENT_CREATE"
	ActionAnnouncementUpdate  = "ANNOUNCEMENT_UPDATE"
)

// Repository interface

type Repository interface {
	CreateStudent(name string) (*Student, error)
	// ClaimGuestStudent adopts a client-generated guest ID into a real
	// students row (progress/attempt rows keyed by that ID predate the row,
	// there are no FKs). Idempotent: existing rows are returned as-is.
	ClaimGuestStudent(id, name string) (*Student, error)
	GetStudent(id string) (*Student, error)
	FindByUsername(username string) (*Student, error)
	FindByEmail(email string) (*Student, error)
	CreateUser(name, username, passwordHash string) (*Student, error)
	// AdoptGuestAsUser turns a credential-less guest row into a named account
	// without changing its ID, so every progress, attempt and XP row keyed by
	// that ID stays attached. Callers must verify the row is claimable first.
	AdoptGuestAsUser(id, name, username, passwordHash string) (*Student, error)
	SetUsername(studentID, username string) error
	ListStudentsMissingUsernames() ([]string, error)
	CreateGoogleUser(name, email, googleID, avatarURL string) (*Student, error)
	CreateOAuthUser(provider, providerID, name, email string, emailVerified bool, avatarURL string) (*Student, error)
	SetAvatarURL(studentID, avatarURL string) error
	SetCourseID(studentID, courseID string) error
	UpdateStudentName(studentID, name string) error
	SetEmail(studentID, email string) error
	SetEmailVerified(studentID string, verified bool) error
	SetPasswordHash(studentID, hash string) error
	CreatePasswordReset(tokenHash, studentID string, expiresAt time.Time) error
	ConsumePasswordReset(tokenHash string) (studentID string, ok bool, err error)

	GetProgress(studentID, conceptID string) (*ConceptProgress, error)
	GetAllProgress(studentID string) (map[string]*ConceptProgress, error)
	UpsertProgress(p *ConceptProgress) error
	// UpsertProgressBatch persists many rows in one transaction, the
	// N+1 antidote for PropagateWeakness and ApplyGoalResults.
	UpsertProgressBatch(ps []*ConceptProgress) error
	// SeedFSRSFromSM2 backfills FSRS state from existing SM-2 intervals for unseeded rows.
	SeedFSRSFromSM2(now time.Time) (int, error)

	CreateSession(studentID string) (*Session, error)
	GetSession(id string) (*Session, error)
	// EnsureSession inserts a sessions row idempotently (ephemeral
	// diagnostic/quiz UUIDs are not created through CreateSession but
	// attempts.session_id references sessions(id)).
	EnsureSession(id, studentID string) error
	// ServerSessions is a durable KV for restart-proof server state
	// (study anti-cheat expected answers, admin logins). Values carry an
	// RFC3339 expires_at; readers treat expired rows as missing.
	UpsertServerSession(kind, key, value, expiresAt string) error
	GetServerSession(kind, key string) (value, expiresAt string, found bool, err error)
	DeleteServerSession(kind, key string) error
	SweepServerSessions() error
	// ReserveAIUsage consumes one request from every daily limit in a single
	// transaction, and commits only when every counter is within its ceiling.
	//
	// When any ceiling is reached the whole transaction rolls back, so a refused
	// request spends no provider call and no budget. The check is on the
	// post-increment value inside the transaction, never on a count read before
	// it, so two concurrent requests can never both take the last free slot,
	// which is the whole reason this is one call and not a read followed by a
	// write. The limits are (identity, application-wide): the identity is the
	// server's choice, a student id when authenticated and an IP for a guest,
	// and the reserved global key bounds the shared free tier.
	ReserveAIUsage(ctx context.Context, day string, limits []ai.AIUsageLimit) (allowed bool, err error)
	// PruneAIUsage deletes counters for days before beforeDay (a YYYY-MM-DD
	// string) and reports how many rows went. Counters are keyed by day, and a
	// day that has passed can never be reserved again, so the rows are dead
	// weight: without this, one row per identity per day accumulates forever.
	// The current day is never a valid argument, callers pass a cutoff weeks
	// back, because pruning it would refill every learner's budget.
	PruneAIUsage(beforeDay string) (int64, error)
	// AIUsage reads one counter without incrementing. Diagnostics only; it must
	// never decide a reservation.
	AIUsage(ctx context.Context, usageKey, day string) (int, error)
	GetActiveSession(sessionID string) (*ActiveSession, error)
	UpsertActiveSession(a *ActiveSession) error
	DeleteActiveSession(sessionID string) error
	// RecordAttempt returns the id of the row it wrote. The id is the event key
	// for the XP award: it is unique per answer, so the award and the attempt it
	// rewards cannot drift apart under retries.
	RecordAttempt(entry AttemptEntry) (int64, error)
	// UpdateAttemptCorrect flips a recorded attempt's correctness. Used by the
	// grading data-repair pass to fix false negatives.
	UpdateAttemptCorrect(id int64, correct bool) error
	// ResetProgress wipes one student's learning record in a single
	// transaction: attempts, sessions, progress, topic speeds, quiz
	// completions, study-scoped server sessions, plus XP, diagnostic flag,
	// and league fields on the student row. Account, credentials,
	// settings, avatar, course, and plan prefs survive. The AI usage counter
	// also survives: a reset that refilled the daily AI budget would make the
	// budget meaningless.
	ResetProgress(studentID string) error
	// DeleteAccount removes one student and every owned row in a single
	// transaction: credential-adjacent tables, learning evidence, sessions,
	// study plans and prefs, filed reports, the AI usage counter, and the
	// student row itself. Irreversible; the caller confirms explicitly.
	DeleteAccount(studentID string) error
	GetSessionAttempts(studentID, sessionID string) ([]AttemptEntry, error)
	GetAttemptsForStudent(studentID string) ([]AttemptEntry, error)
	// GetRecentAttemptsForConcept returns up to `limit` of the most recent attempts on one
	// concept, newest last.
	//
	// This exists because the evidence model needs a bounded window per concept on every
	// graded answer, and GetAttemptsForStudent returns the learner's entire history,
	// 15,774 rows for someone who has worked through the corpus. Reading that per answer
	// would be a different kind of mistake from the one it fixes.
	//
	// Backed by idx_attempts_cover (student_id, concept_id, timestamp) in both stores.
	GetRecentAttemptsForConcept(studentID, conceptID string, limit int) ([]AttemptEntry, error)
	GetAllAttempts() ([]AttemptEntry, error)

	GetQuestions(conceptID string, count int) ([]Question, error)
	GetQuestionCount(conceptID string) (int, error)
	ImportQuestions(qs []Question) error
	PurgeGeneratedQuestions(conceptIDs map[string]bool) (int64, error)

	GetWeeklyLeaderboard() ([]LeaderboardRow, error)
	// LeaderboardPage returns one page and the count of eligible rows. The
	// count comes from the same predicate as the rows, so "page 3 of 9" is a
	// fact rather than a guess.
	LeaderboardPage(view string, since time.Time, limit, offset int) ([]LeaderboardRow, int, error)
	// LeaderboardPosition locates one learner without loading the board.
	LeaderboardPosition(studentID, view string, since time.Time) (rank, xp int, found bool, err error)
	// ResolveLeagueWeek materialises week's league seats from the previous
	// week's final standings, in one transaction. It is safe to call on every
	// read: a completed week is recorded in league_weeks and a second call
	// does nothing.
	ResolveLeagueWeek(week string, prevFrom, prevTo time.Time) error
	// EnsureLeagueMembership seats a learner who is not in this week's rows,
	// used for accounts that first earn XP after the transition ran.
	EnsureLeagueMembership(studentID, week string) error
	LeagueMembershipOf(studentID, week string) (tier string, leagueNo int, found bool, err error)
	LeagueRoster(week, tier string, leagueNo int, from, to time.Time) ([]LeagueRow, error)
	// LeagueCohorts lists every league that exists in a week, in tier order and
	// then league number. It reads seats only, so an administrator can see the
	// whole pyramid without loading one roster per league.
	LeagueCohorts(week string) ([]LeagueCohort, error)
	SetShareToken(studentID, token string) error
	GetStudentByShareToken(token string) (*Student, error)
	GetDailyActivity(studentID string, days int) ([]DailyActivity, error)
	// AssessmentAttempts returns the learner's recorded answers whose source is
	// quiz or diagnostic, newest first, bounded by the window and by limit.
	// Practice and review answers are learning activity, not assessments.
	AssessmentAttempts(studentID string, since, until time.Time, limit int) ([]AttemptEntry, error)
	// XPInRange sums the ledger between since (inclusive) and until (exclusive).
	// The ledger is the only authoritative XP source; students.xp_total is a
	// running counter that has already been moved by the same transactions.
	XPInRange(studentID string, since, until time.Time) (int, error)
	// LearningEventsPage returns one page of substantiated history, newest
	// first. It unions answers, mastery transitions, quiz completions and the
	// diagnostic completion, so pagination is over the real event stream rather
	// than over one table that happens to be the largest.
	LearningEventsPage(studentID string, limit, offset int) ([]LearningEvent, int, error)
	// QuizCompletions returns mastery-check completions at or after since,
	// newest first. The table stores the lifetime XP snapshot at completion, not
	// a score, which is why the Reports page must not show a quiz percentage.
	QuizCompletions(studentID string, since time.Time, limit int) ([]QuizCompletion, error)

	// AwardXP appends to the XP ledger and moves the counters in one
	// transaction. event_key identifies the learning event that justifies the
	// award; a second award for the same key is a no-op, so a retried request
	// cannot pay twice. The ledger is what makes "XP earned this week" a query
	// rather than a guess, and it is append-only: a negative rush penalty is a
	// row with a negative amount, not an edit of an earlier one.
	AwardXP(studentID string, amount int, reason, eventKey string) error
	// WeeklyXPByStudent sums the ledger from since (inclusive) to now.
	WeeklyXPByStudent(since time.Time) (map[string]int, error)
	// RecordXPEvent appends to the ledger without moving the lifetime or daily
	// counters. It exists for replaying awards made before the ledger existed:
	// the counters already contain them, so this must not add again.
	RecordXPEvent(studentID string, amount int, reason, eventKey string) (bool, error)
	// AttemptsSince returns attempts at or after since, oldest first.
	AttemptsSince(since time.Time) ([]AttemptRow, error)
	// RecordDiagnosticAttempt stores the verdict for one answered diagnostic
	// question so a retry can replay it instead of grading again.
	RecordDiagnosticAttempt(attemptID, response string) error
	DiagnosticAttemptResponse(attemptID string) (string, bool, error)
	GetXP(studentID string) (total int, today int, err error)
	SetDiagnosticCompleted(studentID string) error
	// DiagnosticCompletedAt reports when placement completed. A separate call
	// rather than a Student field for the same reason GetStudentRole is: widening
	// Student means editing every fixed-column SELECT in both stores, and a
	// missing timestamp is better read once where it is displayed than carried on
	// every student load.
	DiagnosticCompletedAt(studentID string) (time.Time, bool, error)
	SetDailyXPGoal(studentID string, goal int) error
	GetSettings(studentID string) (string, error)
	UpdateSettings(studentID string, settings string) error
	SetAvatarImage(studentID, contentType string, data []byte) error
	GetAvatarImage(studentID string) (contentType string, data []byte, found bool, err error)
	ClearAvatarImage(studentID string) error

	GetTopicSpeed(studentID, conceptID string) (*TopicSpeed, error)
	GetAllTopicSpeeds(studentID string) (map[string]*TopicSpeed, error)
	UpsertTopicSpeed(ts *TopicSpeed) error

	// Batch 1: Quiz 150 XP gate completions.
	RecordQuizCompletion(studentID string, xpTotal int) error
	LastQuizCompletion(studentID string) (*QuizCompletion, error)

	CreateReport(r QuestionReport) (int64, error)
	ListReports(status string, limit, offset int) ([]QuestionReport, error)
	// ListReportsByReporter returns one account's reports, newest first, whatever their
	// status. It answers "what has this person sent us", which is the question the admin
	// user page asks, and it is deliberately not the moderation queue: a resolved report
	// still explains the learner's experience.
	ListReportsByReporter(reporterID string, limit int) ([]QuestionReport, error)
	// CountReportsByReporter counts that account's reports across every status. The list is
	// capped so the page stays quick; the count is not, so a repeat reporter is visible even
	// when the list is truncated.
	CountReportsByReporter(reporterID string) (int, error)
	// GetReportByOrigin finds an imported report by the id it had in the database it came
	// from, which is the number the admin page and a colleague both quote.
	GetReportByOrigin(originID string) (*QuestionReport, error)

	// GetReportSeenAt returns the RFC3339 time the account last opened its
	// report history. An empty string means it has never opened it.
	GetReportSeenAt(studentID string) (string, error)
	// SetReportSeenAt records that the account has now seen its report history.
	SetReportSeenAt(studentID, seenAt string) error
	// GetReportClearedAt returns the RFC3339 time the account last cleared its
	// report history. Older reports stay in the moderation log but leave this
	// learner's list.
	GetReportClearedAt(studentID string) (string, error)
	// ClearReports hides the account's current report history and marks it seen.
	ClearReports(studentID, clearedAt string) error

	// ImportReport copies a report that came from another installation into this one.
	//
	// originID is that report's id where it was filed, and the unique index on it is what
	// makes re-running an export a no-op: a partially applied import can be retried without
	// duplicating the queue. The returned bool is false when the report was already present.
	//
	// There is deliberately no write-back. The import exists so the person fixing a reported
	// question can reproduce it against a local database; the fix still travels as a commit.
	ImportReport(originID string, r QuestionReport) (bool, error)

	// Diagnostic sessions outlive the process that started them.
	//
	// They used to live only in the server's memory, so every deploy erased the
	// run behind every open tab and the learner got "Session expired, the
	// server restarted and forgot this run" on their next answer. Eight deploys
	// in a day made that a routine event rather than an edge case. The encoded
	// snapshot is opaque here: the diagnostic package owns its format, and a
	// store that parsed it would have to change whenever the adaptive model did.
	SaveDiagnosticSession(id, studentID string, data []byte) error
	LoadDiagnosticSession(id string) (data []byte, found bool, err error)
	DeleteDiagnosticSession(id string) error
	DeleteDiagnosticSessionsForStudent(studentID string) error
	SweepDiagnosticSessions(before time.Time) (int, error)

	// Remedial queue persistence: 2-miss halt remedial items survive server restarts.
	EnqueueRemedial(studentID string, conceptIDs ...string) error
	GetRemedialQueue(studentID string) ([]string, error)
	RemoveRemedial(studentID, conceptID string) error
	ClearRemedialQueue(studentID string) error

	// ListUnsettledReportIDsForConcept returns the ids of a concept's reports
	// that are still open or under review, oldest first. Bulk moderation closes
	// a fixed bug's reports as a set, and selecting that set by concept is what
	// makes the operation about the bug rather than about a hand-picked list.
	ListUnsettledReportIDsForConcept(conceptID string) ([]int64, error)
	UpdateReportStatus(id int64, status string) error

	// ── Administration ────────────────────────────────────────────────────
	//
	// Every method here sits behind the server's admin boundary; the repository is not the
	// security boundary and does no authorization of its own.

	// GetStudentRole is the only way to learn a role.
	//
	// It is a separate call rather than a field on Student on purpose. Every existing read of
	// Student selects a fixed column list in two stores, so adding a column there meant editing
	// roughly twenty query strings to keep them in sync and a missed one is a scan error at
	// runtime, in whichever endpoint nobody was testing. Worse, a Role field that is zero on
	// most Student values reads as `student` when it has not been loaded at all, which is the
	// worst possible default for a privilege field: forgetting to populate it looks exactly
	// like correctly concluding that someone is a learner.
	GetStudentRole(studentID string) (Role, error)

	// SetStudentRole writes a role directly, with no last-admin protection. Callers that can
	// be reached from the network must use SetRoleAudited instead; this exists for the
	// operator's bootstrap, which is trusted and which must be able to make the very first
	// promotion, and for tests arranging a starting state.
	SetStudentRole(studentID string, role Role) error

	// SetRoleAudited changes one account's role and records the event in a single transaction.
	//
	// This is the network-facing path, and the transaction is the point. Applying the mutation
	// and then writing the audit row leaves a window in which a change exists with no event; an
	// administrator whose own promotion failed to be recorded has an action nobody can account
	// for. Here the insert is part of the same unit of work, so either both land or neither
	// does, a failing audit insert rolls the role back rather than orphaning it.
	//
	// Returns ErrLastAdmin if the demotion would leave no administrator, ErrRoleUnchanged if the
	// account already holds the requested role (and writes nothing), sql.ErrNoRows if there is
	// no such account.
	SetRoleAudited(actorID, studentID string, role Role) (RoleChange, error)

	// DemoteAdminSafely demotes one admin to student, and refuses if that would leave the
	// system with no administrator.
	//
	// The check and the write are one statement. Reading a count and then writing is a race:
	// two administrators demoting each other at the same moment both read two, and the system
	// reaches zero. Making the condition part of the UPDATE means the second one affects no
	// rows. Returns ErrLastAdmin when it refused.
	DemoteAdminSafely(studentID string) error

	// CountAdmins is how many accounts may reach the admin surface right now.
	CountAdmins() (int, error)

	// SearchAdminUsers finds accounts by email, id, username or name. An empty query lists
	// newest-first. Password hashes, share tokens and learning state are not selected, so
	// they cannot be returned even by a bug.
	SearchAdminUsers(query string, limit, offset int) ([]AdminUser, error)

	// GetAdminUser reads one account's non-sensitive projection, or nil when there is no such
	// row. It is the read a role change does before deciding whether it is allowed, so it has to
	// distinguish "no such account" from "an account that is a student".
	GetAdminUser(studentID string) (*AdminUser, error)

	// ListStaff returns accounts holding moderator or above, newest first. It is the
	// contributors page's read, and it is a query rather than a filter over SearchAdminUsers so
	// the page does not page through every learner to find the few staff.
	ListStaff(limit, offset int) ([]AdminUser, error)

	// CountAdminUsers is the total number of accounts, for the overview.
	CountAdminUsers() (int, error)

	// ActiveLearnerCounts answers "is anyone using this right now" without a
	// third-party analytics service. activeNow counts distinct learners whose
	// active_sessions row was touched at or after since; activeToday counts
	// distinct learners with an attempt at or after dayStart, which measures
	// people who answered something rather than people with a tab open. Guests
	// are reported separately because a guest-heavy number looks healthy while
	// converting to nothing.
	ActiveLearnerCounts(since, dayStart time.Time) (activeNow, activeNowGuests, activeToday, activeTodayGuests int, err error)

	// CountStoredQuestions is how many question rows exist, for the overview.
	CountStoredQuestions() (int, error)

	// RecordAdminAudit appends one event. There is no update and no delete, on either store or
	// through any route: the log is append-only because an audit trail you can edit is not one.
	RecordAdminAudit(rec AdminAuditRecord) (int64, error)

	// ListAdminAudit returns the most recent events first.
	ListAdminAudit(limit, offset int) ([]AdminAuditRecord, error)

	// ListAdminAuditFor returns the most recent events about one entity, newest first. It is how
	// a person's role history and a report's moderation history are reconstructed: "who granted
	// this, and when" is a query over the trail, not a column anyone maintains by hand.
	ListAdminAuditFor(entityType, entityID string, limit int) ([]AdminAuditRecord, error)

	// ── Admin V2: roles, moderation and invitations ───────────────────────
	//
	// Additive to V1. The role guards and the audit trail are the same machinery; V2 adds the
	// roles above `admin`, the moderation state on reports, and invitations.

	// CountStaff is how many accounts hold moderator or above. It is the number the overview
	// shows; CountAdmins (role = 'admin' exactly) no longer describes "who can administer".
	CountStaff() (int, error)

	// CountOwners is how many accounts hold the owner role. The invariant is that this is never
	// zero, so it is also the number the last-owner guard protects.
	CountOwners() (int, error)

	// CountReportsByStatus returns one count per moderation status. Statuses with no rows are
	// absent rather than zero, so the caller decides what "none" reads as.
	CountReportsByStatus() (map[string]int, error)

	// GetReport reads one report, for the detail view a moderator decides from.
	GetReport(id int64) (*QuestionReport, error)

	// UpdateReportStatusAudited applies a moderation decision and records it in one transaction.
	// The report row is never deleted; Status, Resolution, ResolvedBy and ResolvedAt change and
	// an audit event is appended. Returns sql.ErrNoRows if there is no such report.
	UpdateReportStatusAudited(actorID string, id int64, status, resolution string) (QuestionReport, error)

	// CreateAdminInvitation stores a pending invitation. Only the token hash is passed; the raw
	// token exists in memory for the length of the response and nowhere else.
	CreateAdminInvitation(inv AdminInvitation, tokenHash string) (int64, error)

	// ListAdminInvitations returns invitations newest first, pending or not, so the contributors
	// page can show history rather than only what is outstanding.
	ListAdminInvitations(limit, offset int) ([]AdminInvitation, error)

	// RevokeAdminInvitation marks a pending invitation revoked. Returns the updated row, or
	// sql.ErrNoRows if it does not exist or is already settled.
	RevokeAdminInvitation(actorID string, id int64) (AdminInvitation, error)

	// FindAdminInvitationByTokenHash reads a pending-or-settled invitation by its token hash,
	// or nil. It is read-only and is used to check the invited address before burning the token
	//, the burn itself is ConsumeAdminInvitation, which is the atomic step.
	FindAdminInvitationByTokenHash(tokenHash string) (*AdminInvitation, error)

	// ConsumeAdminInvitation burns a raw-token hash and returns the invitation it belonged to.
	// ok is false for unknown, expired, revoked or already-accepted tokens, which are
	// indistinguishable to the caller on purpose: a token prober learns nothing about which
	// near-miss they are close to.
	ConsumeAdminInvitation(tokenHash, studentID string) (*AdminInvitation, bool, error)

	CreateIdentity(provider, providerID, studentID, email string, emailVerified bool) error
	FindStudentByIdentity(provider, providerID string) (*Student, error)
	ListIdentities(studentID string) ([]Identity, error)
	DeleteIdentity(provider, studentID string) error
	CreateEmailVerification(tokenHash, studentID string, expiresAt time.Time) error
	ConsumeEmailVerification(tokenHash string) (studentID string, ok bool, err error)
	CreateLinkToken(tokenHash, studentID string, expiresAt time.Time) error
	ConsumeLinkToken(tokenHash string) (studentID string, ok bool, err error)

	Migrate() error
	Close() error
}

// UUID

func newUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("uuid-%x", time.Now().UnixNano())
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// IsUniqueViolation reports UNIQUE-constraint failures across stores
// (SQLite "UNIQUE constraint failed", Postgres "duplicate key").
func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed") ||
		strings.Contains(msg, "duplicate key")
}
