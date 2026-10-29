package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/admin"
	"github.com/chuma-beep/mathua/internal/storage"
)

const (
	announcementTitleLimit = 120
	announcementBodyLimit  = 5000
)

// announcementJSON is the learner-facing shape. CreatedAt and PublishedAt are
// RFC3339 so the notification centre can render them without knowing the
// storage representation.
type announcementJSON struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Published   bool   `json:"published"`
	CreatedAt   string `json:"created_at"`
	PublishedAt string `json:"published_at"`
	Read        bool   `json:"read"`
}

func announcementDTO(a storage.Announcement) announcementJSON {
	out := announcementJSON{
		ID:        a.ID,
		Title:     a.Title,
		Body:      a.Body,
		Published: a.Published,
		Read:      a.Read,
	}
	if !a.CreatedAt.IsZero() {
		out.CreatedAt = a.CreatedAt.UTC().Format(time.RFC3339)
	}
	if !a.PublishedAt.IsZero() {
		out.PublishedAt = a.PublishedAt.UTC().Format(time.RFC3339)
	}
	return out
}

func announcementStore(w http.ResponseWriter, repo storage.Repository) (storage.AnnouncementStore, bool) {
	store, ok := repo.(storage.AnnouncementStore)
	if !ok {
		writeError(w, "announcements are not supported by this store", 500)
		return nil, false
	}
	return store, true
}

func validAnnouncement(title, body string) bool {
	return title != "" && body != "" && len(title) <= announcementTitleLimit && len(body) <= announcementBodyLimit
}

// GET /api/announcements returns published announcements for the signed-in
// learner, newest first, with a read flag and the unread count.
func (s *Server) handleAnnouncements(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	store, ok := announcementStore(w, s.repo)
	if !ok {
		return
	}
	items, err := store.ListPublishedAnnouncements(studentID)
	if err != nil {
		writeError(w, "failed to load announcements", 500)
		return
	}
	out := make([]announcementJSON, 0, len(items))
	unread := 0
	for _, item := range items {
		if !item.Read {
			unread++
		}
		out = append(out, announcementDTO(item))
	}
	writeJSON(w, map[string]interface{}{"announcements": out, "unread": unread})
}

// POST /api/announcements/read marks every currently published announcement
// read for the caller. The operation is idempotent.
func (s *Server) handleAnnouncementsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	store, ok := announcementStore(w, s.repo)
	if !ok {
		return
	}
	now := time.Now().UTC()
	if err := store.MarkAnnouncementsRead(studentID, now); err != nil {
		writeError(w, "failed to mark announcements read", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "read_at": now.Format(time.RFC3339)})
}

// POST /api/announcements/clear hides the caller's current announcement feed.
// The global announcement rows and their read state stay intact.
func (s *Server) handleAnnouncementsClear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	store, ok := announcementStore(w, s.repo)
	if !ok {
		return
	}
	now := time.Now().UTC()
	if err := store.ClearAnnouncements(studentID, now); err != nil {
		writeError(w, "failed to clear announcements", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "cleared_at": now.Format(time.RFC3339)})
}

// GET /api/admin/announcements lists drafts and published announcements.
// POST creates one. The route floor is ContentRead; writes re-check
// ContentManage so a moderator can inspect but not publish.
func (s *Server) handleAdminAnnouncements(w http.ResponseWriter, r *http.Request) {
	store, ok := announcementStore(w, s.repo)
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodGet:
		items, err := store.ListAnnouncements()
		if err != nil {
			writeError(w, "failed to load announcements", 500)
			return
		}
		out := make([]announcementJSON, 0, len(items))
		for _, item := range items {
			out = append(out, announcementDTO(item))
		}
		writeJSON(w, map[string]interface{}{"announcements": out})
	case http.MethodPost:
		if !hasPermission(r, admin.ContentManage) {
			writeError(w, "insufficient permission", 403)
			return
		}
		var req struct {
			Title     string `json:"title"`
			Body      string `json:"body"`
			Published bool   `json:"published"`
		}
		if err := decodeJSON(w, r, &req); err != nil {
			writeError(w, "invalid request", 400)
			return
		}
		title := strings.TrimSpace(req.Title)
		body := strings.TrimSpace(req.Body)
		if !validAnnouncement(title, body) {
			writeError(w, "title and body are required and must be within the limits", 400)
			return
		}
		a := actorFrom(r)
		now := time.Now().UTC()
		item, err := store.CreateAnnouncement(storage.Announcement{
			Title:     title,
			Body:      body,
			Published: req.Published,
			CreatedAt: now,
			CreatedBy: a.ID,
		})
		if err != nil {
			writeError(w, "failed to create announcement", 500)
			return
		}
		afterJSON, _ := json.Marshal(map[string]interface{}{
			"title": item.Title, "body": item.Body, "published": item.Published,
		})
		_, _ = s.repo.RecordAdminAudit(storage.AdminAuditRecord{
			ActorID: a.ID, Action: storage.ActionAnnouncementCreate,
			EntityType: "announcement", EntityID: item.ID,
			AfterJSON: string(afterJSON), CreatedAt: now.Format(time.RFC3339),
		})
		writeJSON(w, map[string]interface{}{"announcement": announcementDTO(*item)})
	default:
		writeError(w, "method not allowed", 405)
	}
}

// PATCH /api/admin/announcements/{id} edits, publishes or unpublishes one
// announcement. Omitted fields keep their current value.
func (s *Server) handleAdminAnnouncement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		writeError(w, "method not allowed", 405)
		return
	}
	if !hasPermission(r, admin.ContentManage) {
		writeError(w, "insufficient permission", 403)
		return
	}
	id := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/admin/announcements/"), "/")
	if id == "" || strings.Contains(id, "/") {
		writeError(w, "invalid announcement id", 400)
		return
	}
	store, ok := announcementStore(w, s.repo)
	if !ok {
		return
	}
	before, err := store.GetAnnouncement(id)
	if err != nil {
		writeError(w, "failed to load announcement", 500)
		return
	}
	if before == nil {
		writeError(w, "announcement not found", 404)
		return
	}
	var req struct {
		Title     *string `json:"title"`
		Body      *string `json:"body"`
		Published *bool   `json:"published"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	title, body, published := before.Title, before.Body, before.Published
	if req.Title != nil {
		title = strings.TrimSpace(*req.Title)
	}
	if req.Body != nil {
		body = strings.TrimSpace(*req.Body)
	}
	if req.Published != nil {
		published = *req.Published
	}
	if !validAnnouncement(title, body) {
		writeError(w, "title and body are required and must be within the limits", 400)
		return
	}
	now := time.Now().UTC()
	item, err := store.UpdateAnnouncement(id, title, body, published, now)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, "announcement not found", 404)
		return
	}
	if err != nil {
		writeError(w, "failed to update announcement", 500)
		return
	}
	a := actorFrom(r)
	beforeJSON, _ := json.Marshal(map[string]interface{}{
		"title": before.Title, "body": before.Body, "published": before.Published,
	})
	afterJSON, _ := json.Marshal(map[string]interface{}{
		"title": item.Title, "body": item.Body, "published": item.Published,
	})
	_, _ = s.repo.RecordAdminAudit(storage.AdminAuditRecord{
		ActorID: a.ID, Action: storage.ActionAnnouncementUpdate,
		EntityType: "announcement", EntityID: item.ID,
		BeforeJSON: string(beforeJSON), AfterJSON: string(afterJSON),
		CreatedAt: now.Format(time.RFC3339),
	})
	writeJSON(w, map[string]interface{}{"announcement": announcementDTO(*item)})
}
