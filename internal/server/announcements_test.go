package server

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/storage"
)

type announcementEnvelope struct {
	Announcement struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Body        string `json:"body"`
		Published   bool   `json:"published"`
		CreatedAt   string `json:"created_at"`
		PublishedAt string `json:"published_at"`
		Read        bool   `json:"read"`
	} `json:"announcement"`
	Announcements []struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Body        string `json:"body"`
		Published   bool   `json:"published"`
		CreatedAt   string `json:"created_at"`
		PublishedAt string `json:"published_at"`
		Read        bool   `json:"read"`
	} `json:"announcements"`
	Unread int `json:"unread"`
}

func decodeAnnouncementEnvelope(t *testing.T, body []byte) announcementEnvelope {
	t.Helper()
	var out announcementEnvelope
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("decode announcements: %v (%s)", err, body)
	}
	return out
}

func TestAnnouncementsLifecycleAndReadState(t *testing.T) {
	e := newAdminEnv(t)

	// A draft is visible to staff but not to learners.
	rec := e.do(t, "POST", "/api/admin/announcements", e.adminTok,
		`{"title":"Mathua beta update","body":"Timed quizzes are live.","published":false}`)
	if rec.Code != 200 {
		t.Fatalf("create draft: %d %s", rec.Code, rec.Body.String())
	}
	draft := decodeAnnouncementEnvelope(t, rec.Body.Bytes()).Announcement
	if draft.ID == "" || draft.Published {
		t.Fatalf("draft = %+v, want an id and unpublished", draft)
	}

	rec = e.do(t, "GET", "/api/announcements", e.studentTok, "")
	if rec.Code != 200 {
		t.Fatalf("learner draft list: %d %s", rec.Code, rec.Body.String())
	}
	if got := decodeAnnouncementEnvelope(t, rec.Body.Bytes()); len(got.Announcements) != 0 || got.Unread != 0 {
		t.Fatalf("learner saw a draft: %+v", got)
	}

	// Publishing makes it visible and unread.
	rec = e.do(t, "PATCH", "/api/admin/announcements/"+draft.ID, e.adminTok, `{"published":true}`)
	if rec.Code != 200 {
		t.Fatalf("publish: %d %s", rec.Code, rec.Body.String())
	}
	published := decodeAnnouncementEnvelope(t, rec.Body.Bytes()).Announcement
	if !published.Published || published.PublishedAt == "" {
		t.Fatalf("published = %+v, want a publication time", published)
	}

	rec = e.do(t, "GET", "/api/announcements", e.studentTok, "")
	visible := decodeAnnouncementEnvelope(t, rec.Body.Bytes())
	if len(visible.Announcements) != 1 || visible.Unread != 1 {
		t.Fatalf("learner list = %+v, want one unread", visible)
	}
	if visible.Announcements[0].Read || visible.Announcements[0].Title != "Mathua beta update" {
		t.Fatalf("announcement = %+v", visible.Announcements[0])
	}

	// Marking read is idempotent.
	for i := 0; i < 2; i++ {
		rec = e.do(t, "POST", "/api/announcements/read", e.studentTok, "")
		if rec.Code != 200 {
			t.Fatalf("mark read %d: %d %s", i, rec.Code, rec.Body.String())
		}
	}
	rec = e.do(t, "GET", "/api/announcements", e.studentTok, "")
	read := decodeAnnouncementEnvelope(t, rec.Body.Bytes())
	if len(read.Announcements) != 1 || read.Unread != 0 || !read.Announcements[0].Read {
		t.Fatalf("after read = %+v", read)
	}

	// Unpublishing hides it from learners but leaves the record for staff.
	rec = e.do(t, "PATCH", "/api/admin/announcements/"+draft.ID, e.adminTok, `{"published":false}`)
	if rec.Code != 200 {
		t.Fatalf("unpublish: %d %s", rec.Code, rec.Body.String())
	}
	rec = e.do(t, "GET", "/api/announcements", e.studentTok, "")
	if got := decodeAnnouncementEnvelope(t, rec.Body.Bytes()); len(got.Announcements) != 0 {
		t.Fatalf("learner saw an unpublished item: %+v", got)
	}
	rec = e.do(t, "GET", "/api/admin/announcements", e.adminTok, "")
	all := decodeAnnouncementEnvelope(t, rec.Body.Bytes())
	if len(all.Announcements) != 1 || all.Announcements[0].Published {
		t.Fatalf("staff list = %+v, want the unpublished row", all)
	}

	// Republishing preserves the first publication time and the read state.
	rec = e.do(t, "PATCH", "/api/admin/announcements/"+draft.ID, e.adminTok, `{"published":true}`)
	if rec.Code != 200 {
		t.Fatalf("republish: %d %s", rec.Code, rec.Body.String())
	}
	again := decodeAnnouncementEnvelope(t, rec.Body.Bytes()).Announcement
	if again.PublishedAt != published.PublishedAt {
		t.Errorf("published_at moved from %q to %q", published.PublishedAt, again.PublishedAt)
	}
	rec = e.do(t, "GET", "/api/announcements", e.studentTok, "")
	if got := decodeAnnouncementEnvelope(t, rec.Body.Bytes()); len(got.Announcements) != 1 || got.Unread != 0 || !got.Announcements[0].Read {
		t.Fatalf("republished read state = %+v", got)
	}
}

func TestAnnouncementsClearHidesTheLearnerFeed(t *testing.T) {
	e := newAdminEnv(t)
	rec := e.do(t, "POST", "/api/admin/announcements", e.adminTok,
		`{"title":"Clear me","body":"This leaves the feed.","published":true}`)
	if rec.Code != 200 {
		t.Fatalf("create: %d %s", rec.Code, rec.Body.String())
	}
	rec = e.do(t, "GET", "/api/announcements", e.studentTok, "")
	if got := decodeAnnouncementEnvelope(t, rec.Body.Bytes()); len(got.Announcements) != 1 || got.Unread != 1 {
		t.Fatalf("before clear = %+v", got)
	}
	rec = e.do(t, "POST", "/api/announcements/clear", e.studentTok, "")
	if rec.Code != 200 {
		t.Fatalf("clear: %d %s", rec.Code, rec.Body.String())
	}
	rec = e.do(t, "GET", "/api/announcements", e.studentTok, "")
	if got := decodeAnnouncementEnvelope(t, rec.Body.Bytes()); len(got.Announcements) != 0 || got.Unread != 0 {
		t.Fatalf("after clear = %+v, want empty", got)
	}
}

func TestAnnouncementsRequireAuthentication(t *testing.T) {
	e := newAdminEnv(t)
	rec := e.do(t, "GET", "/api/announcements", "", "")
	if rec.Code != 401 {
		t.Errorf("GET without a token: got %d, want 401", rec.Code)
	}
	rec = e.do(t, "POST", "/api/announcements/read", "", "")
	if rec.Code != 401 {
		t.Errorf("POST read without a token: got %d, want 401", rec.Code)
	}
	rec = e.do(t, "POST", "/api/announcements/clear", "", "")
	if rec.Code != 401 {
		t.Errorf("POST clear without a token: got %d, want 401", rec.Code)
	}
}

func TestAnnouncementsModeratorCanReadButNotWrite(t *testing.T) {
	e := newAdminEnv(t)
	if err := e.store.SetStudentRole(e.studentID, storage.RoleModerator); err != nil {
		t.Fatalf("make moderator: %v", err)
	}
	rec := e.do(t, "GET", "/api/admin/announcements", e.studentTok, "")
	if rec.Code != 200 {
		t.Fatalf("moderator read: %d %s", rec.Code, rec.Body.String())
	}
	rec = e.do(t, "POST", "/api/admin/announcements", e.studentTok,
		`{"title":"Nope","body":"Nope","published":false}`)
	if rec.Code != 403 {
		t.Errorf("moderator create: got %d, want 403", rec.Code)
	}
	rec = e.do(t, "POST", "/api/admin/announcements", e.adminTok,
		`{"title":"Existing","body":"Existing","published":false}`)
	if rec.Code != 200 {
		t.Fatalf("create for patch test: %d %s", rec.Code, rec.Body.String())
	}
	id := decodeAnnouncementEnvelope(t, rec.Body.Bytes()).Announcement.ID
	rec = e.do(t, "PATCH", "/api/admin/announcements/"+id, e.studentTok, `{"published":true}`)
	if rec.Code != 403 {
		t.Errorf("moderator update: got %d, want 403", rec.Code)
	}
}

func TestAnnouncementValidation(t *testing.T) {
	e := newAdminEnv(t)
	cases := []struct {
		name string
		body string
	}{
		{"empty title", `{"title":"","body":"Body","published":false}`},
		{"empty body", `{"title":"Title","body":"","published":false}`},
		{"title too long", fmt.Sprintf(`{"title":%q,"body":"Body","published":false}`, strings.Repeat("t", announcementTitleLimit+1))},
		{"body too long", fmt.Sprintf(`{"title":"Title","body":%q,"published":false}`, strings.Repeat("b", announcementBodyLimit+1))},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := e.do(t, "POST", "/api/admin/announcements", e.adminTok, tc.body)
			if rec.Code != 400 {
				t.Errorf("got %d, want 400 (%s)", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestAnnouncementAuditTrail(t *testing.T) {
	e := newAdminEnv(t)
	rec := e.do(t, "POST", "/api/admin/announcements", e.adminTok,
		`{"title":"Audited","body":"Body","published":false}`)
	if rec.Code != 200 {
		t.Fatalf("create: %d %s", rec.Code, rec.Body.String())
	}
	id := decodeAnnouncementEnvelope(t, rec.Body.Bytes()).Announcement.ID
	rec = e.do(t, "PATCH", "/api/admin/announcements/"+id, e.adminTok, `{"published":true}`)
	if rec.Code != 200 {
		t.Fatalf("publish: %d %s", rec.Code, rec.Body.String())
	}

	audits, err := e.store.ListAdminAuditFor("announcement", id, 10)
	if err != nil {
		t.Fatalf("audit: %v", err)
	}
	seen := map[string]bool{}
	for _, audit := range audits {
		seen[audit.Action] = true
		if audit.ActorID != e.adminID {
			t.Errorf("actor = %q, want %q", audit.ActorID, e.adminID)
		}
	}
	if !seen[storage.ActionAnnouncementCreate] || !seen[storage.ActionAnnouncementUpdate] {
		t.Errorf("audit actions = %+v", seen)
	}
}
