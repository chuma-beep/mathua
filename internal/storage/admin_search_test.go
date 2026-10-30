package storage

import (
	"testing"
	"time"
)

func TestAdminSearchFindsReportsAndAudit(t *testing.T) {
	repo := newTestStore(t)
	search, ok := repo.(ReportSearchStore)
	if !ok {
		t.Fatalf("%T does not implement ReportSearchStore", repo)
	}
	audit, ok := repo.(AuditSearchStore)
	if !ok {
		t.Fatalf("%T does not implement AuditSearchStore", repo)
	}

	if _, err := repo.CreateReport(QuestionReport{
		ReporterID: "lapotimu", ConceptID: "ratio.basics.concept", Kind: "question",
		Question: "Write the ratio of 5 to 2 in simplest form.", Detail: "unclear", Status: "open", Reason: "unclear",
	}); err != nil {
		t.Fatalf("create ratio report: %v", err)
	}
	if _, err := repo.CreateReport(QuestionReport{
		ReporterID: "someone", ConceptID: "alg.quad.incomplete", Kind: "question",
		Question: "Solve x^2 = 256.", Status: "open", Reason: "wrong_answer",
	}); err != nil {
		t.Fatalf("create algebra report: %v", err)
	}
	if _, err := repo.RecordAdminAudit(AdminAuditRecord{
		ActorID: "owner-1", Action: "ANNOUNCEMENT_CREATE", EntityType: "announcement", EntityID: "a-1",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}); err != nil {
		t.Fatalf("record audit: %v", err)
	}

	reports, total, err := search.SearchReports("open", "ratio", 10, 0)
	if err != nil {
		t.Fatalf("search reports: %v", err)
	}
	if total != 1 || len(reports) != 1 || reports[0].ConceptID != "ratio.basics.concept" {
		t.Fatalf("ratio search = %d total, %d rows", total, len(reports))
	}

	reports, total, err = search.SearchReports("open", "lapotimu", 10, 0)
	if err != nil || total != 1 || len(reports) != 1 {
		t.Fatalf("reporter search = %d total, %d rows, %v", total, len(reports), err)
	}

	reports, total, err = search.SearchReports("open", "no-such-report", 10, 0)
	if err != nil || total != 0 || len(reports) != 0 {
		t.Fatalf("empty search = %d total, %d rows, %v", total, len(reports), err)
	}

	events, total, err := audit.SearchAdminAudit("announcement", 10, 0)
	if err != nil || total != 1 || len(events) != 1 {
		t.Fatalf("audit search = %d total, %d rows, %v", total, len(events), err)
	}
	events, total, err = audit.SearchAdminAudit("no-such-event", 10, 0)
	if err != nil || total != 0 || len(events) != 0 {
		t.Fatalf("empty audit search = %d total, %d rows, %v", total, len(events), err)
	}
}
