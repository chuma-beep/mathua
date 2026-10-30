package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/storage"
)

type friendsEnv struct {
	mux       *http.ServeMux
	store     *storage.SQLiteStore
	firstTok  string
	firstID   string
	secondTok string
	secondID  string
}

func newFriendsEnv(t *testing.T, enabled bool) *friendsEnv {
	t.Helper()
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	authSvc := auth.New(store)
	s := New(engine.New(store, d, reg, nil, nil), store, authSvc, WithFriends(enabled))
	mux := http.NewServeMux()
	s.Register(mux)

	env := &friendsEnv{mux: mux, store: store}
	signup := func(name, username string) (string, string) {
		t.Helper()
		body := fmt.Sprintf(`{"name":%q,"username":%q,"password":"Engine!n1","email":%q}`, name, username, username+"@example.com")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/auth/signup", bytes.NewReader([]byte(body))))
		if rec.Code != 200 {
			t.Fatalf("signup %s: %d %s", username, rec.Code, rec.Body.String())
		}
		var out struct {
			Token string `json:"token"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
			t.Fatalf("decode signup: %v", err)
		}
		me := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+out.Token)
		mux.ServeHTTP(me, req)
		var who struct {
			StudentID string `json:"student_id"`
		}
		json.Unmarshal(me.Body.Bytes(), &who)
		if who.StudentID == "" {
			t.Fatalf("signup %s returned no student id: %s", username, me.Body.String())
		}
		if err := store.SetEmailVerified(who.StudentID, true); err != nil {
			t.Fatalf("verify %s: %v", username, err)
		}
		return out.Token, who.StudentID
	}
	env.firstTok, env.firstID = signup("Ada", "ada")
	env.secondTok, env.secondID = signup("Grace", "grace")
	return env
}

func (e *friendsEnv) do(t *testing.T, method, path, token, body string) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body == "" {
		r = httptest.NewRequest(method, path, nil)
	} else {
		r = httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	}
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	e.mux.ServeHTTP(rec, r)
	return rec
}

func TestFriendsRoutesHonourTheFlag(t *testing.T) {
	off := newFriendsEnv(t, false)
	rec := off.do(t, "GET", "/api/friends", off.firstTok, "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("disabled GET /api/friends = %d, want 404: %s", rec.Code, rec.Body.String())
	}

	on := newFriendsEnv(t, true)
	rec = on.do(t, "GET", "/api/friends", on.firstTok, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("enabled GET /api/friends = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	var feed struct {
		Identity struct {
			ProfileID  string `json:"profile_id"`
			FriendCode string `json:"friend_code"`
		} `json:"identity"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &feed); err != nil {
		t.Fatalf("decode feed: %v", err)
	}
	if feed.Identity.ProfileID == "" || feed.Identity.FriendCode == "" {
		t.Fatalf("feed has no identity: %s", rec.Body.String())
	}
}

func TestFriendsRequestLifecycle(t *testing.T) {
	env := newFriendsEnv(t, true)

	secondIdentity, err := env.store.GetFriendIdentity(env.secondID)
	if err != nil {
		t.Fatalf("identity: %v", err)
	}

	rec := env.do(t, "POST", "/api/friends/requests", env.firstTok, fmt.Sprintf(`{"query":%q}`, secondIdentity.FriendCode))
	if rec.Code != http.StatusOK {
		t.Fatalf("send request = %d: %s", rec.Code, rec.Body.String())
	}

	rec = env.do(t, "GET", "/api/friends", env.secondTok, "")
	var feed struct {
		Incoming []struct {
			ProfileID string `json:"profile_id"`
		} `json:"incoming"`
		Friends []struct {
			ProfileID string `json:"profile_id"`
		} `json:"friends"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &feed); err != nil {
		t.Fatalf("decode feed: %v", err)
	}
	if len(feed.Incoming) != 1 {
		t.Fatalf("incoming = %d, want 1: %s", len(feed.Incoming), rec.Body.String())
	}

	firstIdentity, err := env.store.GetFriendIdentity(env.firstID)
	if err != nil {
		t.Fatalf("identity: %v", err)
	}
	rec = env.do(t, "POST", "/api/friends/requests/accept", env.secondTok, fmt.Sprintf(`{"profile_id":%q}`, firstIdentity.ProfileID))
	if rec.Code != http.StatusOK {
		t.Fatalf("accept = %d: %s", rec.Code, rec.Body.String())
	}

	rec = env.do(t, "GET", "/api/friends", env.firstTok, "")
	if err := json.Unmarshal(rec.Body.Bytes(), &feed); err != nil {
		t.Fatalf("decode feed: %v", err)
	}
	if len(feed.Friends) != 1 || feed.Friends[0].ProfileID != secondIdentity.ProfileID {
		t.Fatalf("friends after accept = %s, want %s", rec.Body.String(), secondIdentity.ProfileID)
	}

	rec = env.do(t, "GET", "/api/friends/profile?profile_id="+secondIdentity.ProfileID, env.firstTok, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("friend profile = %d: %s", rec.Code, rec.Body.String())
	}
	rec = env.do(t, "GET", "/api/friends/graph?profile_id="+secondIdentity.ProfileID, env.firstTok, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("friend graph = %d: %s", rec.Code, rec.Body.String())
	}
}

func TestFriendsRejectSelfAndMissingAuth(t *testing.T) {
	env := newFriendsEnv(t, true)
	identity, err := env.store.GetFriendIdentity(env.firstID)
	if err != nil {
		t.Fatalf("identity: %v", err)
	}
	rec := env.do(t, "POST", "/api/friends/requests", env.firstTok, fmt.Sprintf(`{"query":%q}`, identity.FriendCode))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("self request = %d, want 400: %s", rec.Code, rec.Body.String())
	}
	rec = env.do(t, "GET", "/api/friends", "", "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous GET = %d, want 401: %s", rec.Code, rec.Body.String())
	}
}
