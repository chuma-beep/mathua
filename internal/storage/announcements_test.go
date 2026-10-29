package storage

import (
	"testing"
	"time"
)

func announcementStore(t *testing.T, repo Repository) AnnouncementStore {
	t.Helper()
	store, ok := repo.(AnnouncementStore)
	if !ok {
		t.Fatalf("repository %T does not implement AnnouncementStore", repo)
	}
	return store
}

func TestAnnouncementsPublishReadAndUnpublish(t *testing.T) {
	store := newTestStore(t)
	ann := announcementStore(t, store)
	learner := seedStudent(t, store)
	other := seedStudent(t, store)

	created, err := ann.CreateAnnouncement(Announcement{
		Title: "Beta update",
		Body:  "Timed quizzes are live.",
	})
	if err != nil {
		t.Fatalf("create draft: %v", err)
	}
	if created.Published || !created.PublishedAt.IsZero() {
		t.Fatalf("new draft = %+v, want unpublished with no publish time", created)
	}
	if got, _ := ann.ListPublishedAnnouncements(learner.ID); len(got) != 0 {
		t.Fatalf("draft visible to learner: %+v", got)
	}
	if got, _ := ann.ListAnnouncements(); len(got) != 1 {
		t.Fatalf("admin announcements = %d, want 1", len(got))
	}

	publishedAt := time.Now().UTC()
	published, err := ann.UpdateAnnouncement(created.ID, created.Title, created.Body, true, publishedAt)
	if err != nil {
		t.Fatalf("publish: %v", err)
	}
	if !published.Published || published.PublishedAt.IsZero() {
		t.Fatalf("published = %+v", published)
	}
	if !published.PublishedAt.Equal(publishedAt) {
		t.Fatalf("published at = %v, want %v", published.PublishedAt, publishedAt)
	}

	items, err := ann.ListPublishedAnnouncements(learner.ID)
	if err != nil {
		t.Fatalf("list learner announcements: %v", err)
	}
	if len(items) != 1 || items[0].Read {
		t.Fatalf("learner items = %+v, want one unread", items)
	}
	if err := ann.MarkAnnouncementsRead(learner.ID, time.Now().UTC()); err != nil {
		t.Fatalf("mark read: %v", err)
	}
	items, _ = ann.ListPublishedAnnouncements(learner.ID)
	if len(items) != 1 || !items[0].Read {
		t.Fatalf("after mark read = %+v, want read", items)
	}
	if otherItems, _ := ann.ListPublishedAnnouncements(other.ID); len(otherItems) != 1 || otherItems[0].Read {
		t.Fatalf("another learner inherited read state: %+v", otherItems)
	}

	// Republishing must not move the original publication time.
	later := publishedAt.Add(time.Hour)
	updated, err := ann.UpdateAnnouncement(created.ID, "Beta update revised", created.Body, true, later)
	if err != nil {
		t.Fatalf("update published: %v", err)
	}
	if !updated.PublishedAt.Equal(publishedAt) {
		t.Fatalf("published at moved on edit: got %v want %v", updated.PublishedAt, publishedAt)
	}

	if _, err := ann.UpdateAnnouncement(created.ID, created.Title, created.Body, false, later); err != nil {
		t.Fatalf("unpublish: %v", err)
	}
	if items, _ := ann.ListPublishedAnnouncements(learner.ID); len(items) != 0 {
		t.Fatalf("unpublished announcement still visible: %+v", items)
	}
	if all, _ := ann.ListAnnouncements(); len(all) != 1 {
		t.Fatalf("admin lost unpublished announcement: %+v", all)
	}

	if _, err := ann.UpdateAnnouncement(created.ID, created.Title, created.Body, true, later); err != nil {
		t.Fatalf("republish: %v", err)
	}
	items, _ = ann.ListPublishedAnnouncements(learner.ID)
	if len(items) != 1 || !items[0].Read {
		t.Fatalf("republished read state = %+v, want preserved read", items)
	}
}

func TestAnnouncementsSurviveResetAndAccountDeletion(t *testing.T) {
	store := newTestStore(t)
	ann := announcementStore(t, store)
	learner := seedStudent(t, store)
	other := seedStudent(t, store)

	created, err := ann.CreateAnnouncement(Announcement{
		Title:     "Account note",
		Body:      "This survives account changes.",
		Published: true,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := ann.MarkAnnouncementsRead(learner.ID, time.Now().UTC()); err != nil {
		t.Fatalf("mark read: %v", err)
	}
	if err := store.ResetProgress(learner.ID); err != nil {
		t.Fatalf("reset progress: %v", err)
	}
	items, _ := ann.ListPublishedAnnouncements(learner.ID)
	if len(items) != 1 || !items[0].Read {
		t.Fatalf("reset removed announcement or read state: %+v", items)
	}

	if err := store.DeleteAccount(learner.ID); err != nil {
		t.Fatalf("delete account: %v", err)
	}
	if all, _ := ann.ListAnnouncements(); len(all) != 1 || all[0].ID != created.ID {
		t.Fatalf("account deletion removed the global announcement: %+v", all)
	}
	items, _ = ann.ListPublishedAnnouncements(other.ID)
	if len(items) != 1 || items[0].Read {
		t.Fatalf("other learner read state = %+v, want unread", items)
	}
}

func TestAnnouncementsClearHidesCurrentFeed(t *testing.T) {
	store := newTestStore(t)
	ann := announcementStore(t, store)
	learner := seedStudent(t, store)
	other := seedStudent(t, store)
	now := time.Now().UTC()

	first, err := ann.CreateAnnouncement(Announcement{
		Title:       "First update",
		Body:        "This should be cleared.",
		Published:   true,
		CreatedAt:   now,
		PublishedAt: now,
	})
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	if items, _ := ann.ListPublishedAnnouncements(learner.ID); len(items) != 1 {
		t.Fatalf("first feed = %+v, want one item", items)
	}
	if err := ann.ClearAnnouncements(learner.ID, now.Add(time.Second)); err != nil {
		t.Fatalf("clear announcements: %v", err)
	}
	if items, _ := ann.ListPublishedAnnouncements(learner.ID); len(items) != 0 {
		t.Fatalf("cleared feed = %+v, want empty", items)
	}
	if items, _ := ann.ListPublishedAnnouncements(other.ID); len(items) != 1 || items[0].ID != first.ID {
		t.Fatalf("another learner lost the announcement: %+v", items)
	}

	second, err := ann.CreateAnnouncement(Announcement{
		Title:       "Second update",
		Body:        "This arrives after the clear.",
		Published:   true,
		CreatedAt:   now.Add(2 * time.Second),
		PublishedAt: now.Add(2 * time.Second),
	})
	if err != nil {
		t.Fatalf("create second: %v", err)
	}
	items, err := ann.ListPublishedAnnouncements(learner.ID)
	if err != nil {
		t.Fatalf("list after second: %v", err)
	}
	if len(items) != 1 || items[0].ID != second.ID {
		t.Fatalf("new announcement after clear = %+v, want only second", items)
	}
}

func TestAnnouncementUpdateMissingRow(t *testing.T) {
	store := newTestStore(t)
	ann := announcementStore(t, store)
	if _, err := ann.UpdateAnnouncement("missing", "Title", "Body", true, time.Now().UTC()); err == nil {
		t.Fatal("update missing announcement returned nil error")
	}
}
