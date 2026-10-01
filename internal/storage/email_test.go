package storage

import (
	"testing"
)

func TestEmailUniqueIndex_Created(t *testing.T) {
	s := newTestStore(t)
	var name string
	if err := s.db.QueryRow("SELECT name FROM sqlite_master WHERE type='index' AND name='uidx_students_email'").Scan(&name); err != nil {
		t.Fatalf("expected uidx_students_email index: %v", err)
	}
}

func TestSetEmail_DuplicateRejected(t *testing.T) {
	s := newTestStore(t)
	a, err := s.CreateUser("A", "adda", "hash")
	if err != nil {
		t.Fatal(err)
	}
	b, err := s.CreateUser("B", "bddb", "hash")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.SetEmail(a.ID, "dupe@example.com"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetEmail(b.ID, "dupe@example.com"); !IsUniqueViolation(err) {
		t.Fatalf("expected unique violation, got %v", err)
	}
}

func TestSetEmail_BlankRowsAllowed(t *testing.T) {
	s := newTestStore(t)
	a, err := s.CreateUser("A", "blanka", "hash")
	if err != nil {
		t.Fatal(err)
	}
	b, err := s.CreateUser("B", "blankb", "hash")
	if err != nil {
		t.Fatal(err)
	}
	// Legacy '' rows must stay legal under the partial index.
	if err := s.SetEmail(a.ID, ""); err != nil {
		t.Fatal(err)
	}
	if err := s.SetEmail(b.ID, ""); err != nil {
		t.Fatal(err)
	}
}

func TestFindByEmail_CaseInsensitive(t *testing.T) {
	s := newTestStore(t)
	a, err := s.CreateUser("A", "casea", "hash")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.SetEmail(a.ID, "User@Example.COM"); err != nil {
		t.Fatal(err)
	}
	got, err := s.FindByEmail("user@example.com")
	if err != nil || got == nil || got.ID != a.ID {
		t.Fatalf("expected case-insensitive match, got %+v %v", got, err)
	}
}
