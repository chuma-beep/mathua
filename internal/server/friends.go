package server

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// Friends is a phase-two surface: the storage layer landed behind a flag and
// this file is the only thing that exposes it. Every route answers 404 when
// MATHUA_FRIENDS is not on, so a deployment that has not been cleared for the
// public rollout behaves as though the feature does not exist. That includes
// the API, not just the navigation: a hidden page with a live endpoint is not
// a gate.
//
// The friend graph is deliberately narrow. A friend sees a concept's state and
// whether it is mastered. It never sees weakness, attempts, due dates, elapsed
// time or an answer, because none of those cross this DTO.

type friendAvatarJSON struct {
	URL      string             `json:"avatar_url,omitempty"`
	Dicebear *storage.AvatarRef `json:"avatar_dicebear,omitempty"`
	Preset   *int               `json:"avatar_preset,omitempty"`
}

type friendSummaryJSON struct {
	ProfileID    string           `json:"profile_id"`
	Name         string           `json:"name"`
	Username     string           `json:"username,omitempty"`
	Avatar       friendAvatarJSON `json:"avatar"`
	XPTotal      int              `json:"xp_total"`
	Mastered     int              `json:"mastered"`
	FriendsSince string           `json:"friends_since,omitempty"`
}

type friendRequestJSON struct {
	ProfileID string           `json:"profile_id"`
	Name      string           `json:"name"`
	Username  string           `json:"username,omitempty"`
	Avatar    friendAvatarJSON `json:"avatar"`
	Direction string           `json:"direction"`
	CreatedAt string           `json:"created_at"`
}

type friendIdentityJSON struct {
	ProfileID  string `json:"profile_id"`
	FriendCode string `json:"friend_code"`
}

type friendsListJSON struct {
	Identity friendIdentityJSON  `json:"identity"`
	Friends  []friendSummaryJSON `json:"friends"`
	Incoming []friendRequestJSON `json:"incoming"`
	Outgoing []friendRequestJSON `json:"outgoing"`
}

type friendProfileJSON struct {
	friendSummaryJSON
	FriendsSince string `json:"friends_since,omitempty"`
}

type friendGraphNodeJSON struct {
	ConceptID string `json:"concept_id"`
	Status    string `json:"status"`
}

type friendLeaderboardJSON struct {
	ProfileID string           `json:"profile_id"`
	Name      string           `json:"name"`
	Username  string           `json:"username,omitempty"`
	Avatar    friendAvatarJSON `json:"avatar"`
	XP        int              `json:"xp"`
	Mastered  int              `json:"mastered"`
	IsMe      bool             `json:"is_me"`
}

func friendAvatar(studentID, avatarURL, settings string) friendAvatarJSON {
	url, dicebear, preset := adminAvatar(studentID, avatarURL, settings)
	return friendAvatarJSON{URL: url, Dicebear: dicebear, Preset: preset}
}

func friendSummary(row storage.FriendSummary) friendSummaryJSON {
	return friendSummaryJSON{
		ProfileID:    row.ProfileID,
		Name:         row.Name,
		Username:     row.Username,
		Avatar:       friendAvatar(row.StudentID, row.AvatarURL, row.AvatarSettings),
		XPTotal:      row.XPTotal,
		Mastered:     row.Mastered,
		FriendsSince: friendTime(row.FriendsSince),
	}
}

func friendRequest(row storage.FriendRequest) friendRequestJSON {
	return friendRequestJSON{
		ProfileID: row.ProfileID,
		Name:      row.Name,
		Username:  row.Username,
		Avatar:    friendAvatar(row.StudentID, row.AvatarURL, row.AvatarSettings),
		Direction: row.Direction,
		CreatedAt: friendTime(row.CreatedAt),
	}
}

func friendTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

// friendStore returns the friend surface and writes the failure itself when the
// repository does not implement it. A repository that cannot store friendships
// is a deployment error, not a learner error, so a 500 is honest.
func (s *Server) friendStore(w http.ResponseWriter) (storage.FriendStore, bool) {
	store, ok := s.repo.(storage.FriendStore)
	if !ok {
		writeError(w, "friends are not available", 500)
		return nil, false
	}
	return store, true
}

func (s *Server) friendsAvailable(w http.ResponseWriter) bool {
	if !s.friendsEnabled {
		writeError(w, "not found", 404)
		return false
	}
	return true
}

// GET /api/friends
func (s *Server) handleFriends(w http.ResponseWriter, r *http.Request) {
	if !s.friendsAvailable(w) {
		return
	}
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	store, ok := s.friendStore(w)
	if !ok {
		return
	}
	identity, err := store.GetFriendIdentity(studentID)
	if err != nil {
		writeFriendError(w, err)
		return
	}
	friends, err := store.ListFriends(studentID)
	if err != nil {
		writeFriendError(w, err)
		return
	}
	incoming, outgoing, err := store.ListFriendRequests(studentID)
	if err != nil {
		writeFriendError(w, err)
		return
	}
	out := friendsListJSON{
		Identity: friendIdentityJSON{ProfileID: identity.ProfileID, FriendCode: identity.FriendCode},
		Friends:  make([]friendSummaryJSON, 0, len(friends)),
		Incoming: make([]friendRequestJSON, 0, len(incoming)),
		Outgoing: make([]friendRequestJSON, 0, len(outgoing)),
	}
	for _, f := range friends {
		out.Friends = append(out.Friends, friendSummary(f))
	}
	for _, f := range incoming {
		out.Incoming = append(out.Incoming, friendRequest(f))
	}
	for _, f := range outgoing {
		out.Outgoing = append(out.Outgoing, friendRequest(f))
	}
	writeJSON(w, out)
}

// /api/friends/<action>, all authenticated and write-limited where mutating.
func (s *Server) handleFriendRoute(w http.ResponseWriter, r *http.Request) {
	if !s.friendsAvailable(w) {
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	store, ok := s.friendStore(w)
	if !ok {
		return
	}
	action := strings.TrimPrefix(r.URL.Path, "/api/friends/")
	switch action {
	case "requests":
		s.handleFriendSend(w, r, store, studentID)
	case "requests/accept":
		s.handleFriendDecision(w, r, store, studentID, "accept")
	case "requests/decline":
		s.handleFriendDecision(w, r, store, studentID, "decline")
	case "requests/cancel":
		s.handleFriendDecision(w, r, store, studentID, "cancel")
	case "remove":
		s.handleFriendDecision(w, r, store, studentID, "remove")
	case "block":
		s.handleFriendDecision(w, r, store, studentID, "block")
	case "unblock":
		s.handleFriendDecision(w, r, store, studentID, "unblock")
	case "code/rotate":
		s.handleFriendCodeRotate(w, r, store, studentID)
	case "profile":
		s.handleFriendProfile(w, r, store, studentID)
	case "graph":
		s.handleFriendGraph(w, r, store, studentID)
	case "leaderboard":
		s.handleFriendLeaderboard(w, r, store, studentID)
	default:
		writeError(w, "not found", 404)
	}
}

type friendBody struct {
	Query     string `json:"query"`
	ProfileID string `json:"profile_id"`
}

func decodeFriendBody(w http.ResponseWriter, r *http.Request) (friendBody, bool) {
	var body friendBody
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, "invalid request", 400)
		return friendBody{}, false
	}
	body.Query = strings.TrimSpace(body.Query)
	body.ProfileID = strings.TrimSpace(body.ProfileID)
	return body, true
}

func (s *Server) handleFriendSend(w http.ResponseWriter, r *http.Request, store storage.FriendStore, studentID string) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", 405)
		return
	}
	body, ok := decodeFriendBody(w, r)
	if !ok {
		return
	}
	if body.Query == "" {
		writeError(w, "enter a username or friend code", 400)
		return
	}
	profileID, err := store.ResolveFriendTarget(body.Query)
	if err != nil {
		writeFriendError(w, err)
		return
	}
	friendship, err := store.SendFriendRequest(studentID, profileID, time.Now())
	if err != nil {
		writeFriendError(w, err)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "status": friendship.Status})
}

func (s *Server) handleFriendDecision(w http.ResponseWriter, r *http.Request, store storage.FriendStore, studentID, action string) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", 405)
		return
	}
	body, ok := decodeFriendBody(w, r)
	if !ok {
		return
	}
	if body.ProfileID == "" {
		writeError(w, "missing profile", 400)
		return
	}
	now := time.Now()
	var err error
	switch action {
	case "accept":
		_, err = store.AcceptFriendRequest(studentID, body.ProfileID, now)
	case "decline":
		err = store.DeclineFriendRequest(studentID, body.ProfileID, now)
	case "cancel":
		err = store.CancelFriendRequest(studentID, body.ProfileID, now)
	case "remove":
		err = store.RemoveFriend(studentID, body.ProfileID, now)
	case "block":
		err = store.BlockUser(studentID, body.ProfileID, now)
	case "unblock":
		err = store.UnblockUser(studentID, body.ProfileID)
	}
	if err != nil {
		writeFriendError(w, err)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

func (s *Server) handleFriendCodeRotate(w http.ResponseWriter, r *http.Request, store storage.FriendStore, studentID string) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", 405)
		return
	}
	code, err := store.RotateFriendCode(studentID, time.Now())
	if err != nil {
		writeFriendError(w, err)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "friend_code": code})
}

func (s *Server) handleFriendProfile(w http.ResponseWriter, r *http.Request, store storage.FriendStore, studentID string) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	profileID := strings.TrimSpace(r.URL.Query().Get("profile_id"))
	if profileID == "" {
		writeError(w, "missing profile", 400)
		return
	}
	profile, err := store.FriendProfile(studentID, profileID)
	if err != nil {
		writeFriendError(w, err)
		return
	}
	out := friendProfileJSON{friendSummaryJSON: friendSummary(storage.FriendSummary{
		StudentID:      profile.StudentID,
		ProfileID:      profile.ProfileID,
		Name:           profile.Name,
		Username:       profile.Username,
		AvatarURL:      profile.AvatarURL,
		AvatarSettings: profile.AvatarSettings,
		XPTotal:        profile.XPTotal,
		Mastered:       profile.Mastered,
		FriendsSince:   profile.FriendsSince,
	})}
	out.FriendsSince = friendTime(profile.FriendsSince)
	writeJSON(w, out)
}

func (s *Server) handleFriendGraph(w http.ResponseWriter, r *http.Request, store storage.FriendStore, studentID string) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	profileID := strings.TrimSpace(r.URL.Query().Get("profile_id"))
	if profileID == "" {
		writeError(w, "missing profile", 400)
		return
	}
	nodes, err := store.FriendGraph(studentID, profileID)
	if err != nil {
		writeFriendError(w, err)
		return
	}
	out := make([]friendGraphNodeJSON, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, friendGraphNodeJSON{ConceptID: n.ConceptID, Status: n.Status})
	}
	writeJSON(w, map[string]interface{}{"nodes": out})
}

func (s *Server) handleFriendLeaderboard(w http.ResponseWriter, r *http.Request, store storage.FriendStore, studentID string) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	view := strings.TrimSpace(r.URL.Query().Get("view"))
	if view == "" {
		view = "weekly"
	}
	if view != "weekly" && view != "all_time" {
		writeError(w, "invalid view", 400)
		return
	}
	rows, err := store.FriendsLeaderboard(studentID, view, time.Now())
	if err != nil {
		writeFriendError(w, err)
		return
	}
	out := make([]friendLeaderboardJSON, 0, len(rows))
	for _, row := range rows {
		out = append(out, friendLeaderboardJSON{
			ProfileID: row.ProfileID,
			Name:      row.Name,
			Username:  row.Username,
			Avatar:    friendAvatar(row.StudentID, row.AvatarURL, row.AvatarSettings),
			XP:        row.XP,
			Mastered:  row.Mastered,
			IsMe:      row.IsMe,
		})
	}
	writeJSON(w, map[string]interface{}{"view": view, "rows": out})
}

// writeFriendError maps the storage errors to the two answers a caller is
// allowed to distinguish. Unknown and blocked targets share 404 so a probe
// cannot learn whether an account exists or has blocked the caller.
func writeFriendError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, storage.ErrFriendNotFound), errors.Is(err, storage.ErrFriendBlocked):
		writeError(w, "not found", 404)
	case errors.Is(err, storage.ErrFriendSelf):
		writeError(w, "you cannot friend yourself", 400)
	case errors.Is(err, storage.ErrFriendNotVerified):
		writeError(w, "verify your email address before sending friend requests", 403)
	case errors.Is(err, storage.ErrFriendVisibilityOff):
		writeError(w, "this learner has turned off friends", 403)
	default:
		log.Printf("friends: %v", err)
		writeError(w, "friends request failed", 500)
	}
}
