package storage

import (
	"testing"
	"time"
)

// The import is the boundary between production and a local workbench, so the
// two properties that matter are that it keeps the evidence intact and that
// running the same export twice does not double the queue.
func TestImportReport_IsIdempotentAndKeepsTheEvidence(t *testing.T) {
	store := newTestStore(t)
	filed := time.Date(2026, 10, 20, 9, 30, 0, 0, time.UTC)

	r := QuestionReport{
		ReporterID:  "rep-1",
		ConceptID:   "arith.div.remainder",
		Kind:        "question",
		Question:    "What is the remainder when 17 is divided by 5?",
		Expected:    "2",
		AnswerGiven: "3",
		Reason:      "wrong_expected",
		Detail:      "The expected answer looks wrong.",
		Status:      "open",
		CreatedAt:   filed,
	}

	inserted, err := store.ImportReport("73", r)
	if err != nil {
		t.Fatalf("first import: %v", err)
	}
	if !inserted {
		t.Fatal("first import reported the report as already present")
	}

	// A second run of the same export must not add a row, even if the payload
	// carries a later local edit: the file is a copy of one report, not a feed.
	r.Question = "edited after export"
	inserted, err = store.ImportReport("73", r)
	if err != nil {
		t.Fatalf("second import: %v", err)
	}
	if inserted {
		t.Error("a repeated import inserted a duplicate")
	}

	reports, err := store.ListReports("open", 10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(reports) != 1 {
		t.Fatalf("queue has %d reports, want 1", len(reports))
	}
	got := reports[0]
	if got.ConceptID != "arith.div.remainder" || got.Expected != "2" || got.AnswerGiven != "3" {
		t.Errorf("evidence not preserved: %+v", got)
	}
	if !got.CreatedAt.UTC().Equal(filed) {
		t.Errorf("created_at = %s, want %s", got.CreatedAt.UTC(), filed)
	}
	if got.Question != "What is the remainder when 17 is divided by 5?" {
		t.Errorf("the second payload overwrote the first: %q", got.Question)
	}
}

func TestImportReport_RequiresAnOriginID(t *testing.T) {
	store := newTestStore(t)
	if _, err := store.ImportReport("   ", QuestionReport{ConceptID: "a"}); err == nil {
		t.Fatal("import accepted an empty origin id")
	}
}

// Locally filed reports have no origin id, and the uniqueness rule must not
// treat "none" as a value: two reports from the same concept are ordinary.
func TestImportReport_DoesNotConstrainLocallyFiledReports(t *testing.T) {
	store := newTestStore(t)
	r := QuestionReport{ConceptID: "arith.div.remainder", Kind: "question", Status: "open", CreatedAt: time.Now().UTC()}
	if _, err := store.CreateReport(r); err != nil {
		t.Fatalf("create first: %v", err)
	}
	if _, err := store.CreateReport(r); err != nil {
		t.Fatalf("create second: %v", err)
	}
	reports, err := store.ListReports("open", 10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(reports) != 2 {
		t.Fatalf("queue has %d reports, want 2", len(reports))
	}
}

// TestImportReport_KeepsTheOriginIDReadable is about the number a person quotes.
// The local row id is not the production id — SQLite burns ids on rolled-back
// inserts, so a local #94 and the production #94 can be different reports. The
// origin column must therefore survive the import and be findable both in a list
// and by direct lookup, which is what makes "look at production #94" work here.
func TestImportReport_KeepsTheOriginIDReadable(t *testing.T) {
	store := newTestStore(t)
	if _, err := store.ImportReport("94", QuestionReport{
		ConceptID: "precalc.series.geom_infinite",
		Kind:      "question",
		Status:    "open",
	}); err != nil {
		t.Fatalf("import: %v", err)
	}
	rows, err := store.ListReports("open", 10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(rows) != 1 || rows[0].OriginReportID != "94" {
		t.Fatalf("list did not carry the origin id: %+v", rows)
	}
	found, err := store.GetReportByOrigin("94")
	if err != nil || found == nil {
		t.Fatalf("lookup by origin: %v found=%v", err, found)
	}
	if found.ConceptID != "precalc.series.geom_infinite" {
		t.Errorf("wrong report for origin 94: %+v", found)
	}
	if missing, err := store.GetReportByOrigin("no-such-origin"); err != nil || missing != nil {
		t.Errorf("unknown origin = %+v, %v; want nil, nil", missing, err)
	}
	// A locally filed report has no origin and must not be forced to invent one.
	local, err := store.CreateReport(QuestionReport{Kind: "message", Detail: "local", Status: "open"})
	if err != nil {
		t.Fatalf("create local: %v", err)
	}
	row, err := store.GetReport(local)
	if err != nil || row == nil {
		t.Fatalf("get local: %v %v", err, row)
	}
	if row.OriginReportID != "" {
		t.Errorf("locally filed report claims origin %q", row.OriginReportID)
	}
}
