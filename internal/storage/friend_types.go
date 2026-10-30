package storage

import (
	"errors"
	"time"
)

type FriendshipStatus string

const (
	FriendshipPending  FriendshipStatus = "pending"
	FriendshipAccepted FriendshipStatus = "accepted"
)

const (
	FriendDirectionIncoming = "incoming"
	FriendDirectionOutgoing = "outgoing"
)

var (
	ErrFriendNotFound      = errors.New("friend target not found")
	ErrFriendBlocked       = errors.New("friend target blocked")
	ErrFriendNotVerified   = errors.New("email verification required")
	ErrFriendSelf          = errors.New("cannot friend yourself")
	ErrFriendVisibilityOff = errors.New("friends visibility disabled")
)

type Friendship struct {
	StudentA    string
	StudentB    string
	Status      FriendshipStatus
	RequestedBy string
	CreatedAt   time.Time
	RespondedAt time.Time
}

type FriendRequest struct {
	StudentID      string
	ProfileID      string
	Name           string
	Username       string
	AvatarURL      string
	AvatarSettings string
	Direction      string
	CreatedAt      time.Time
}

type FriendSummary struct {
	StudentID      string
	ProfileID      string
	Name           string
	Username       string
	AvatarURL      string
	AvatarSettings string
	XPTotal        int
	Mastered       int
	FriendsSince   time.Time
}

type FriendProfile struct {
	StudentID      string
	ProfileID      string
	Name           string
	Username       string
	AvatarURL      string
	AvatarSettings string
	XPTotal        int
	Mastered       int
	FriendsSince   time.Time
}

type FriendGraphNode struct {
	ConceptID  string
	Status     string
	MasteredAt *time.Time
}

type FriendLeaderboardRow struct {
	StudentID      string
	ProfileID      string
	Name           string
	Username       string
	AvatarURL      string
	AvatarSettings string
	XP             int
	Mastered       int
	IsMe           bool
}

type FriendIdentity struct {
	ProfileID  string
	FriendCode string
}

type FriendStore interface {
	ResolveFriendTarget(query string) (profileID string, err error)
	GetFriendIdentity(studentID string) (FriendIdentity, error)
	SendFriendRequest(studentID, targetProfileID string, now time.Time) (Friendship, error)
	AcceptFriendRequest(studentID, otherProfileID string, now time.Time) (Friendship, error)
	DeclineFriendRequest(studentID, otherProfileID string, now time.Time) error
	CancelFriendRequest(studentID, otherProfileID string, now time.Time) error
	RemoveFriend(studentID, otherProfileID string, now time.Time) error
	BlockUser(studentID, otherProfileID string, now time.Time) error
	UnblockUser(studentID, otherProfileID string) error
	ListFriends(studentID string) ([]FriendSummary, error)
	ListFriendRequests(studentID string) (incoming, outgoing []FriendRequest, err error)
	FriendProfile(studentID, otherProfileID string) (*FriendProfile, error)
	FriendGraph(studentID, otherProfileID string) ([]FriendGraphNode, error)
	FriendsLeaderboard(studentID, view string, week time.Time) ([]FriendLeaderboardRow, error)
	RotateFriendCode(studentID string, now time.Time) (string, error)
}

var _ FriendStore = (*SQLiteStore)(nil)
var _ FriendStore = (*PostgresStore)(nil)
